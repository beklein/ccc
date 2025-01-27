# Context Carbon Copy (CCC)

Context Carbon Copy (CCC) is a tool for copying text from multiple files and directories into your clipboard to be then used as context for AI models when using the web interface.

In contrast to other tools, this one uses a `.ccc` file, a config file that defines files that should be used as context.

## Installation

- You can download the binary from the [releases page](https://github.com/beklein/ccc/releases) and add to your `$PATH`.
- Build from source using Go Modules
  - Run `go install github.com/beklein/ccc@latest`.

## Usage

Create a `.ccc` file in your project root.
Each line should be a file path, a directory path, or a glob pattern to include into the context.

Examples:

```bash
# This is a comment
some_file.xyz
src/*.go
README.md
```

Run the tool:

```bash
ccc
```

By default this will do the following:
- reads each line form the `.ccc` file in the current directory
- gathers all matched file contents
- copies them to your system clipboard

### Flags

- `-o, --output`
  - Print output to stdout instead of copying to the clipboard.
    - `ccc -o > out.txt`
- `-c, --config`
  - Specify another `.ccc` file.
    - `ccc -config .ccc.example`