# Textblitz

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.16+-00ADD8?style=flat&logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License">
  <img src="https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-lightgrey" alt="Platform">
</p>

## Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [Architecture](#architecture)
- [Installation](#installation)
- [Usage](#usage)
  - [Indexing Files](#indexing-files)
  - [Looking Up by SimHash](#looking-up-by-simhash)
  - [Handling File Names with Spaces](#handling-file-names-with-spaces)
- [Error Handling](#error-handling)
- [Performance Considerations](#performance-considerations)
- [Contributing](#contributing)
- [License](#license)

## 🚀 Introduction

Textblitz is a fast and scalable text indexing system written in Go, designed to efficiently search and retrieve data from large text files. It tackles the common challenge of quickly searching through extensive text corpora by implementing a SimHash-based chunking and indexing strategy.

The system breaks down large files into manageable chunks, computes similarity hashes for each chunk, and builds an optimized in-memory index. This approach enables lightning-fast retrieval of content based on hash values, making it ideal for applications requiring quick text lookups.
## Features

- **Efficient Chunking**: Splits text files into configurable fixed-size chunks
- **SimHash Fingerprinting**: Generates hash signatures that group similar text chunks together
- **Fast Lookup**: Provides immediate retrieval of text chunks based on hash values
- **Multi-threaded Processing**: Utilizes Go's concurrency for parallel processing
- **Memory Efficient**: Optimized for handling large files with minimal memory footprint
- **Simple CLI**: Easy-to-use command-line interface for indexing and lookup operations
## Architecture

Textblitz follows a pipeline architecture for processing text files:

```mermaid
graph TB
    Input[Text File] --> Chunker[Chunk Splitter]
    Chunker --> WorkerPool{Worker Pool}
    WorkerPool --> Worker1[Worker 1]
    WorkerPool --> Worker2[Worker 2]
    WorkerPool --> WorkerN[Worker N]
    Worker1 --> HashGen[SimHash Generator]
    Worker2 --> HashGen
    WorkerN --> HashGen
    HashGen --> IndexBuilder[Index Builder]
    IndexBuilder --> IndexFile[(Index File)]
    
    LookupCmd[Lookup Command] --> SearchIndex[Search Index]
    SearchIndex --> RetrieveChunk[Retrieve Chunk]
    IndexFile -.-> SearchIndex
    
    classDef input fill:#d1f0d1,stroke:#53a653,stroke-width:2px,color:#1a3a1a
    classDef process fill:#d1e8f0,stroke:#4a6da7,stroke-width:2px,color:#1a3a5a
    classDef worker fill:#ffd8b6,stroke:#e67e22,stroke-width:2px,color:#5a3a1a
    classDef storage fill:#e6d8e6,stroke:#9b59b6,stroke-width:2px,color:#3a1a3a
    
    class Input,LookupCmd input
    class Chunker,HashGen,IndexBuilder,SearchIndex,RetrieveChunk process
    class WorkerPool,Worker1,Worker2,WorkerN worker
    class IndexFile storage
```

<details>
<summary>Alternative Static Diagram (if Mermaid doesn't render)</summary>

```
┌─────────────┐     ┌───────────────┐     ┌─────────────┐     ┌─────────────┐     ┌────────────┐
│  Text File  │────▶│ Chunk Splitter│────▶│ Worker Pool │────▶│   SimHash   │────▶│    Index   │
└─────────────┘     └───────────────┘     │   ┌─────┐   │     │  Generator  │     │   Builder  │
                                           │   │  W1 │   │     └─────────────┘     └─────┬──────┘
                                           │   └─────┘   │                               │
                                           │   ┌─────┐   │                               ▼
                                           │   │  W2 │   │                          ┌────────────┐
                                           │   └─────┘   │                          │ Index File │
                                           │   ┌─────┐   │                          └──────┬─────┘
                                           │   │  WN │   │                                 │
                                           └───┴─────┴───┘                                 │
                                                                                          │
┌─────────────┐     ┌───────────────┐     ┌─────────────┐                                 │
│    Lookup   │────▶│ Search Index  │◀----┘             │                                 │
│   Command   │     └───────┬───────┘                    │                                 │
└─────────────┘             │                            │                                 │
                            ▼                            │                                 │
                     ┌─────────────┐                     │                                 │
                     │   Retrieve  │◀────────────────────┘                                 │
                     │    Chunk    │                                                       │
                     └─────────────┘                                                       │
```
</details>

The diagram above illustrates the data flow through the Textblitz system:

1. **Input Handling**: Parses text files and command-line arguments
2. **Chunk Splitting**: Divides text into fixed-size chunks (configurable)
3. **Worker Pool**: Distributes processing across multiple goroutines
4. **SimHash Generation**: Computes similarity hashes for each chunk
5. **Index Construction**: Maps hash values to byte offsets in the original file
6. **Lookup System**: Retrieves chunks based on their SimHash values
## 💻 Installation

### Prerequisites
- Go 1.16 or higher
- Git (for cloning the repository)

### Building from Source

```bash
# Clone the repository
git clone https://github.com/bravian1/Textblitz.git
cd Textblitz

# Build the executable
go build -o textindex

# Verify installation
./textindex --help
```
## 📝 Usage

Textblitz provides two primary commands: indexing and lookup.

### Indexing Files

Process a text file by splitting it into chunks, computing SimHash values, and creating an index:

```bash
textindex -c index -i <input_file.txt> -s <chunk_size> -o <index_file.idx> -w <workers>
```

**Arguments:**
- `-c index`: Specifies the indexing command
- `-i <input_file.txt>`: Path to the input text file
- `-s <chunk_size>`: Size of each chunk in bytes (default: 4096)
- `-o <index_file.idx>`: Path to save the generated index
- `-w <workers>`: Number of worker goroutines for parallel processing (default: 4)

**Example:**

```bash
textindex -c index -i large_text.txt -s 4096 -o index.idx -w 8
```
### Looking Up by SimHash

Find a chunk in the indexed file based on its SimHash value:

```bash
textindex -c lookup -i <index_file.idx> -h <simhash_value>
```

**Arguments:**
- `-c lookup`: Specifies the lookup command
- `-i <index_file.idx>`: Path to the previously generated index file
- `-h <simhash_value>`: SimHash value to search for

**Example:**

```bash
textindex -c lookup -i index.idx -h 3e4f1b2c98a61
```
## Handling File Names with Spaces

When using the command-line interface of Textblitz, if your file names contain spaces, it's important to enclose them in quotes. This ensures that the entire file name is treated as a single argument, rather than being split into multiple arguments. For example:

```bash
go run main.go -c index -i "OpenStax - Physics.pdf" -o sample.idx
```

In this example, the quotes around "OpenStax - Physics.pdf" ensure that the file name is correctly interpreted, even though it contains spaces. Without quotes, the command-line parser would treat each word as a separate argument, leading to errors.
## ⚠️ Error Handling

Textblitz implements comprehensive error handling to ensure reliability across various operations. This section provides detailed guidance for troubleshooting common issues.

### Command-Line Errors

| Error | Description | Solution |
|-------|-------------|----------|
| **Missing Command** | `Error: Missing Command (-c 'index' or 'lookup')` | Specify the command using `-c index` or `-c lookup` |
| **Missing Input File** | `Error: Input file (-i) is required` | Provide the input file path with `-i <filename>` |
| **Missing Output File** | `Error: Output file (-o) is required for indexing` | Specify the output path with `-o <filename>` |
| **Missing SimHash** | `Error: SimHash value (-h) is required for lookup` | Include the hash value with `-h <simhash>` |
| **Invalid Parameter** | `Error: Invalid parameter value` | Check parameter format and allowed values |

### File Operation Errors

| Error | Description | Solution |
|-------|-------------|----------|
| **File Not Found** | `Error: open <filename>: no such file or directory` | Verify file path and existence |
| **Permission Denied** | `Error: permission denied` | Check file permissions (read for input, write for output) |
| **Index File Corruption** | `Error: invalid index file format` | Regenerate the index file |
| **File Too Large** | `Error: file size exceeds maximum allowed` | Process the file in smaller chunks or increase memory allocation |

### Processing Errors

| Error | Description | Solution |
|-------|-------------|----------|
| **Memory Allocation** | `Error: cannot allocate memory` | Reduce worker count (`-w`) or chunk size (`-s`) |
| **Worker Pool Failure** | `Error: worker pool initialization failed` | Decrease worker count or restart application |
| **Processing Timeout** | `Error: processing timed out` | Increase chunk size or reduce file size |
| **Hash Collision** | `Error: hash collision detected` | Adjust chunk size to reduce collision probability |

### Debugging Tips

```bash
# Enable verbose logging
textindex -c index -i input.txt -o output.idx --verbose

# Generate debug information
textindex -c index -i input.txt -o output.idx --debug

# Validate an index file
textindex --validate-index -i index.idx
```

### Error Exit Codes

| Code | Meaning |
|------|--------|
| 1 | General error |
| 2 | Command-line argument error |
| 3 | File operation error |
| 4 | Processing error |
| 5 | Memory allocation error |

### Best Practices

- **Incremental Processing**: For very large files, process in batches
- **Resource Monitoring**: Watch memory usage during indexing of large files
- **Validation**: Always validate generated index files before deployment
- **Backup**: Keep backups of original files before processing
- **Error Logs**: Save error outputs for troubleshooting
## 🤝 Contributing

Contributions to Textblitz are welcome! Here's how you can help:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Please ensure your code follows the project's style guidelines and includes appropriate tests.

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.
