# Textblitz
Jam: Fast & Scalable Text Indexer

## Handling File Names with Spaces

When using the command-line interface of Textblitz, if your file names contain spaces, it's important to enclose them in quotes. This ensures that the entire file name is treated as a single argument, rather than being split into multiple arguments. For example:

```bash
go run main.go -c index -i "OpenStax - Physics.pdf" -o sample.idx
```

In this example, the quotes around "OpenStax - Physics.pdf" ensure that the file name is correctly interpreted, even though it contains spaces. Without quotes, the command-line parser would treat each word as a separate argument, leading to errors.

### Why Use Quotes?

1. **Argument Parsing**: Command-line parsers split input based on spaces. Quotes ensure that file names with spaces are treated as a single argument.
2. **Correct Interpretation**: Quotes preserve spaces and special characters within file names, allowing the application to process them correctly.
3. **Preventing Errors**: Using quotes prevents errors like "file not found" or "invalid argument" by ensuring correct argument parsing.

Always use quotes for file names with spaces to prevent such issues and ensure smooth operation of the tool.

### Doc Formats

- using  sajari/docconv package

- https://github.com/sajari/docconv

## Add popler utils for pdftotext to work
Error during indexing: failed to chunk file: failed to convert document to text: error converting data: exec: "pdftotext": executable file not found in %PATH%    

#### solution:
sudo apt update
sudo apt install poppler-utils -y

