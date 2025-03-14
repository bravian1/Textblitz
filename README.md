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
##  Architecture

Textblitz follows a pipeline architecture for processing text files:

```mermaid
flowchart TB
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

The diagram above illustrates the data flow through the Textblitz system:

1. **Input Handling**: Parses text files and command-line arguments
2. **Chunk Splitting**: Divides text into fixed-size chunks (configurable)
3. **Worker Pool**: Distributes processing across multiple goroutines
4. **SimHash Generation**: Computes similarity hashes for each chunk
5. **Index Construction**: Maps hash values to byte offsets in the original file
6. **Lookup System**: Retrieves chunks based on their SimHash values
