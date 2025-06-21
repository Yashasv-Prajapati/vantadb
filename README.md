# VantaDB

This is my personal project to learn about databases, and different aspects of it. I'm trying to build a key-value database from scratch, and I will be documenting my progress here. It even uses a custom binary file to store the entire metadata on disk. I have also implemented a Write-Ahead Logging (WAL) mechanism for crash recovery. It stores the data in a custom format which is used to recover the database in case of deletion.

# Features

- Key-Value Store
- Custom Binary File Format
- WAL (Write-Ahead Logging) for crash recovery (currently only implemented if complete database is deleted)
- REPL support for interactive commands
- REST API for programmatic access
- Built from scratch in Golang
- Lots of learning about file systems, data structures and database internals.

# Future Plans

- Implement WAL recovery for partial updates using timestamp
- Add more data structures like B-Trees, LSM Trees, etc.
- Implement a query language
- Add support for transactions
- Improve performance and scalability

# How to Run

To run the project, you need to have Go installed on your machine. You can clone the repository and run the following commands:

```bash
git clone https://github.com/Yashasv-Prajapati/vantadb.git
cd vantadb
go mod tidy
go install .
```

Then you can
run the server using the following command:

```bash
vantadb init .vdsk
vantadb serve --port 8080 -f .vdsk
```

The above command will start the server on port 8080 and use the `.vdsk` file to store the database. If the file does not exist, it will be created. It will also start a REPL shell for interactive commands.

# Contributing

If you want to contribute to the project, feel free to open an issue or a pull request. I welcome any contributions, whether it's bug fixes, new features, or documentation improvements.

## How to begin contribution

After forking the repository, you can start by exploring the codebase. Here are some steps to get started:

1. Start with the `main.go` file in the root directory. This file contains the entry point of the application and sets up the server.
2. Explore the `cmd` directory for command-line interface (CLI) commands.
3. Check the `internal` directory for the core logic of the database, including the key-value store implementation, file handling, and WAL mechanism.

## How to setup for contribution

1. Clone the repository:
   ```bash
   git clone
    cd vantadb
    go mod tidy
   ```
2. Try running the server using(also starts the REPL shell)
   ```bash
   go run main.go serve --port 8080 -f .vdsk
   ```
