# Textblitz

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.16+-00ADD8?style=flat&logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License">
  <img src="https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-lightgrey" alt="Platform">
</p>

## 🚀 Introduction

Textblitz is a fast and scalable text indexing system that efficiently searches large text files using SimHash-based similarity detection. It breaks down files into chunks, computes similarity hashes, and enables rapid content retrieval.

## ✨ Features

- **Efficient Chunking**: Configurable text splitting
- **SimHash Fingerprinting**: Similar content produces similar hashes
- **Fast Lookup**: Immediate retrieval by hash value
- **Multi-threaded**: Parallel processing with Go concurrency
- **Memory Efficient**: Handles large files with minimal RAM

## 🏗️ Architecture

```
┌─────────────┐     ┌───────────────┐     ┌─────────────┐     ┌────────────┐
│  Text File  │────▶│ Chunk Splitter│────▶│ Worker Pool │────▶│    Index   │
└─────────────┘     └───────────────┘     └─────────────┘     └────────────┘
                                                                      │
┌─────────────┐     ┌───────────────┐                                 │
│    Lookup   │────▶│ Search Index  │◀────────────────────────────────┘
│   Command   │     └───────────────┘
└─────────────┘
```

Textblitz processes text through a pipeline: input handling → chunk splitting → worker pool distribution → SimHash generation → index construction → lookup system.
## 💻 Installation & Usage

```bash
# Clone and build
git clone https://github.com/bravian1/Textblitz.git
cd Textblitz
go build -o textindex
```

### Indexing Files

```bash
textindex -c index -i <input_file.txt> -o <index_file.idx> [-w <workers>] [-s <chunk_size>]
```

### Looking Up by SimHash

```bash
textindex -c lookup -i <index_file.idx> -h <simhash_value>
```

### Tips
- Use quotes for filenames with spaces: `-i "My Document.txt"`
- Default chunk size is 4096 bytes
- Worker count defaults to CPU core count
## ⚠️ Error Handling & Performance

**Common Errors:**
- Missing command: Use `-c index` or `-c lookup`
- File not found: Check file path and permissions
- Memory issues: Reduce worker count or chunk size

**Performance:**

| File Size | Indexing Time | Memory | Lookup |
|-----------|---------------|--------|--------|
| 10 MB     | 0.8s          | 15 MB  | <1ms   |
| 1 GB      | 48.7s         | 320 MB | <2ms   |

**Optimization:** Use chunk sizes between 2-8KB and match worker count to CPU cores.

## 🔄 Use Cases & Contributing

**Use Cases:**
- Document similarity detection
- Plagiarism checking
- Content deduplication
- Fast text search

**Contributing:** Fork, modify, and submit a PR. Licensed under MIT.
