# 🖨️ Creeder 

_README.md fully written by ChatGPT_ 🤖

Creeder is a command-line interface written in Go that reads all files in a source code directory, prints the file tree of that directory, and prints the contents of each file.

## Installation

To use the CLI tool, you must have Go installed on your computer. You can then install using the following command:

```bash
$ go get github.com/impzero/creeder
```

## Usage

To use the CLI Tool, run the following command:

```bash
$ creeder <path> [-f <filter>] [-i <ignore>]
```

where:

- `<path>` is the path to the source code directory
- `-f <filter>` is an optional comma-separated list of file extensions to include (default is all files)
- `-i <ignore>` is an optional comma-separated list of directories or files to ignore (default is none)

For example, to print the files tree and contents of all `.go` and `.txt` files in the `/path/to/dir` directory, ignoring the `node_modules` and `vendor` directories, run the following command:

```bash
$ creeder /path/to/dir -f go,txt -i node_modules,vendor
```

## Purpose

Build this tool for a couple of reasons:

- ☕ Lazy and uneducated about any others tools that would do this job
- ⚒️ Need the response produced by this tool in order to create embeddings for ChatGPT for another project
- 📚 Learn more about Go and how to build CLI tools
- 🤖 Experiment being a "professional" ChatGPT code prompter

## Contributing

If you have any suggestions or find any issues, please feel free to open an issue or pull request on GitHub.

## License

Creeder is licensed under the MIT license. See `LICENSE` for more information.
