package wal

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"time"
)

const (
	WAL_LOG_FILENAME = "wal.log"
	DELETE_FLAG = 1
	SET_FLAT = 0
)

// log file schema
type WALRecord struct {
	EntrySize uint32   // 4B
	EntryType [1]byte  // 0 - set, 1 - delete
	KeyLen    uint32   // 4B
	Key       [32]byte // 32B key
	ValueLen  uint32   // 4B
	Value     []byte   // variable
	Checksum  uint32   // 4B
	Timestamp uint64   // time
}

func NewWALRecord(entryType string, key string, value string) *WALRecord {
	wr := &WALRecord{}

	actualKeyLen := len(key)
	wr.KeyLen = uint32(actualKeyLen)

	keyBytes := [32]byte{}
	copy(keyBytes[:], key)
	wr.Key = keyBytes

	wr.Value = []byte(value)

	valueLen := len(wr.Value)
	wr.ValueLen = uint32(valueLen)

	wr.Timestamp = uint64(time.Now().Unix())

	if entryType == "set" {
		wr.EntryType[0] = SET_FLAT
	} else if entryType == "delete" {
		wr.EntryType[0] = DELETE_FLAG
	}

	wr.Checksum = crc32.ChecksumIEEE(append(wr.Key[:], wr.Value...))

	entrySize := 4 + 1 + 4 + 32 + 4 + len(wr.Value) + 4 + 8 // sequentially from entrysize --> timestamp
	wr.EntrySize = uint32(entrySize)

	return wr
}

func (wr *WALRecord) WriteWALRecordToFile(index int) bool {

	data := wr.ToBytes()
	// fmt.Println("SAVING DATA LEN ", len(data))

	file, err := os.OpenFile(WAL_LOG_FILENAME, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return false
	}

	defer file.Close()

	_, err = file.Write(data)

	return err == nil

}

func (wr *WALRecord) ToBytes() []byte {
	// Calculate total size: 4 + 1 + 4 + 32 + 4 + valueLen + 4 + 8
	totalSize := 4 + 1 + 4 + 32 + 4 + len(wr.Value) + 4 + 8 // sequentially from entrysize --> timestamp
	data := make([]byte, totalSize)

	offset := 0

	// EntrySize
	binary.LittleEndian.PutUint32(data[offset:offset+4], wr.EntrySize)
	offset += 4

	// EntryType (1 byte)
	copy(data[offset:offset+1], wr.EntryType[:])
	offset += 1

	// KeyLen (4 bytes)
	binary.LittleEndian.PutUint32(data[offset:offset+4], wr.KeyLen)
	offset += 4

	// Key (32 bytes)
	copy(data[offset:offset+32], wr.Key[:])
	offset += 32

	// ValueLen (4 bytes)
	binary.LittleEndian.PutUint32(data[offset:offset+4], wr.ValueLen)
	offset += 4

	// Value (variable length)
	copy(data[offset:offset+len(wr.Value)], wr.Value)
	offset += len(wr.Value)

	// Checksum (4 bytes)
	binary.LittleEndian.PutUint32(data[offset:offset+4], wr.Checksum)
	offset += 4

	// Timestamp (8 bytes)
	binary.LittleEndian.PutUint64(data[offset:offset+8], wr.Timestamp)

	return data
}

func Decode(data []byte) *WALRecord {

	if len(data) < 53 { // min size = 1+4+32+4+0+4+8
		return nil
	}
	wr := &WALRecord{}
	offset := 0

	wr.EntrySize = binary.LittleEndian.Uint32(data[offset : offset+4])
	offset += 4

	// EntryType (1 byte)
	copy(wr.EntryType[:], data[offset:offset+1])
	offset += 1

	// KeyLen (4 bytes)
	wr.KeyLen = binary.LittleEndian.Uint32(data[offset : offset+4])
	offset += 4

	// Key (32 bytes)
	copy(wr.Key[:], data[offset:offset+32])
	offset += 32

	// ValueLen (4 bytes)
	wr.ValueLen = binary.LittleEndian.Uint32(data[offset : offset+4])
	offset += 4

	// Value (variable length)
	if len(data) < offset+int(wr.ValueLen)+8 { // Check if enough data remains
		return nil
	}
	wr.Value = make([]byte, wr.ValueLen)
	copy(wr.Value, data[offset:offset+int(wr.ValueLen)])
	offset += int(wr.ValueLen)

	// Checksum (4 bytes)
	wr.Checksum = binary.LittleEndian.Uint32(data[offset : offset+4])
	offset += 4

	// Timestamp (8 bytes)
	wr.Timestamp = binary.LittleEndian.Uint64(data[offset : offset+8])

	return wr

}

func GetAllWALRecords() []*WALRecord {

	file, err := os.Open(WAL_LOG_FILENAME)
	if err != nil {
		return nil
	}

	defer file.Close()

	var wals []*WALRecord 

	reader := bufio.NewReader(file)
	for {
		entrySizeBytes := make([]byte, 4)
		_, err := io.ReadFull(reader, entrySizeBytes)
		if err == io.EOF {
			break // reached end of file — normal
		}
		if err != nil {
			// something went wrong (e.g. partial write)
			fmt.Println("failed to read entry size:", err)
			break
		}
		entrySize := binary.LittleEndian.Uint32(entrySizeBytes)

		entryBuf := make([]byte, entrySize-4)
		_, err = io.ReadFull(reader, entryBuf)
		if err == io.EOF {
			fmt.Println("partial WAL entry detected at end — ignoring")
			break
		}
		if err != nil {
			fmt.Println("failed to read full WAL entry:", err)
			break
		}
		walBuf := append(entrySizeBytes, entryBuf...)
		record := Decode(walBuf)

		// append to wals array
		wals = append(wals, record)
	}

	return wals

}

// func RecoverFromLogs() {

// 	file, err := os.Open(WAL_LOG_FILENAME)
// 	if err != nil {
// 		return
// 	}

// 	defer file.Close()

// 	reader := bufio.NewReader(file)
// 	for {
// 		entrySizeBytes := make([]byte, 4)
// 		_, err := io.ReadFull(reader, entrySizeBytes)
// 		if err == io.EOF {
// 			break // reached end of file — normal
// 		}
// 		if err != nil {
// 			// something went wrong (e.g. partial write)
// 			fmt.Println("failed to read entry size:", err)
// 			break
// 		}
// 		entrySize := binary.LittleEndian.Uint32(entrySizeBytes)

// 		entryBuf := make([]byte, entrySize-4)
// 		_, err = io.ReadFull(reader, entryBuf)
// 		if err == io.EOF {
// 			fmt.Println("partial WAL entry detected at end — ignoring")
// 			break
// 		}
// 		if err != nil {
// 			fmt.Println("failed to read full WAL entry:", err)
// 			break
// 		}
// 		walBuf := append(entrySizeBytes, entryBuf...)
// 		record := Decode(walBuf)

// 		switch string(record.EntryType[:]) {
// 		case "set":
// 			// Set(string(record.Key[:]), string(record.Value))
// 			continue
// 		case "delete":
// 			fmt.Println("To be implemented")
// 		default:
// 			continue
// 		}

// 	}

// }
