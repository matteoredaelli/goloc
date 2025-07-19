# GoLoc

A fast and efficient command-line tool for analyzing source code files and directories. GoLoc counts lines of code, blank lines, and comments across multiple programming languages, providing detailed statistics about your codebase.

## Features

- **Multi-language support** - Recognizes and analyzes dozens of programming languages
- **Comprehensive metrics** - Counts lines of code, blank lines, comments, and total files
- **Fast performance** - Built in Go for speed and efficiency
- **Flexible input** - Analyze individual files or entire directory trees
- **Multiple output formats** - Table (default), CSV, and JSON output formats
- **File counting mode** - Fast file counting without line parsing (`-f` flag)
- **Cross-platform** - Works on Linux, macOS, and Windows

## Installation

### From Source

```bash
git clone https://github.com/yourusername/goloc.git
cd goloc
go build -o goloc
```

### Using Go Install

```bash
go install github.com/yourusername/goloc@latest
```

## Usage

### Basic Usage

```bash
# Analyze current directory
goloc

# Analyze specific files
goloc main.go utils.go

# Analyze specific directories
goloc ./src ./tests

# Analyze mixed files and directories
goloc main.go ./src ./docs
```

### Command Line Options

```bash
goloc [options] [file or dir] ...
```

**Options:**
- `-f` - Count files without parsing lines (faster for file counting only)
- `-l` - Show supported languages/extensions and exit
- `-o string` - Output format: table|csv|json (default: "table")
- `-u` - Count and show files with unknown extension
- `-h` - Show help message

### Examples

```bash
# Analyze current project
goloc

# Analyze multiple directories
goloc ./frontend ./backend ./shared

# Count files only (faster, no line parsing)
goloc -f ./large-project

# Include files with unknown extensions
goloc -u ./project

# Show supported languages
goloc -l

# Output as CSV for further analysis
goloc -o csv ./src

# Output as JSON
goloc -o json ./project

# Analyze specific files
goloc main.go config.go utils/*.go
```

## Output Format

GoLoc produces a clean tabular output showing statistics for each detected language:

```
+----------+-------+---------+-------+------+----------+--------+
|   LANG   | FILES | SKIPPED | LINES | CODE | COMMENTS | BLANKS |
+----------+-------+---------+-------+------+----------+--------+
| Go       |     6 |       0 |   853 |  626 |      104 |    123 |
| Json     |     4 |       0 |  2051 | 2051 |        0 |      0 |
| Markdown |     2 |       0 |   147 |  105 |       20 |     22 |
| Python   |     3 |       0 |    57 |   43 |        3 |     11 |
| Rust     |     2 |       0 |    15 |    6 |        6 |      3 |
| Sql      |     1 |       0 |    33 |   28 |        0 |      5 |
| Yaml     |     1 |       0 |     6 |    6 |        0 |      0 |
+----------+-------+---------+-------+------+----------+--------+
|    TOTAL |    19 |       0 |  3162 | 2865 |      133 |    164 |
+----------+-------+---------+-------+------+----------+--------+
```

**Column Definitions:**
- **LANG**: Programming language or file type
- **FILES**: Number of files for each language
- **SKIPPED**: Number of files skipped (e.g., due to errors or filters)
- **LINES**: Total lines including code, comments, and blanks
- **CODE**: Lines containing actual code
- **COMMENTS**: Lines containing comments (single-line and multi-line)
- **BLANKS**: Empty lines or lines with only whitespace

## Supported Languages

GoLoc supports analysis of the following programming languages and file types:

- **Programming Languages**: Go, Python, JavaScript, TypeScript, Java, C, C++, C#, Rust, Ruby, PHP, Swift, Kotlin, Scala, Clojure, Haskell, Erlang, Elixir, ...
- **Web Technologies**: HTML, CSS, SCSS, SASS, Less, Vue.js, React (JSX)
- **Scripts & Config**: Shell (Bash, Zsh), PowerShell, Dockerfile, YAML, TOML, JSON, XML
- **Documentation**: Markdown, reStructuredText, AsciiDoc
- **Databases**: SQL, GraphQL
- **Other**: Lua, R, MATLAB, Vim Script, and many more

*File types are detected based on file extensions. Files with unknown extensions can be optionally skipped or included.*


## Why choose GoLoc?**
- Simple, focused tool with multiple output formats (table, CSV, JSON)
- Fast file counting mode for large projects
- Easy to build and modify (Go source code)
- Lightweight with minimal dependencies
- Built-in language detection and comprehensive statistics
- Shows skipped files for transparency

## Configuration

GoLoc uses built-in language definitions and file extension mappings. Future versions may include:
- Custom configuration files
- Additional output formats (JSON, CSV)
- Exclude patterns and filters
- Custom language definitions

## Contributing

Contributions are welcome! Please feel free to:

- Report bugs or request features via [GitHub Issues](https://github.com/yourusername/goloc/issues)
- Submit pull requests for improvements
- Add support for additional languages
- Improve documentation

### Development Setup

```bash
git clone https://github.com/yourusername/goloc.git
cd goloc
go mod tidy
go run main.go
```

### Running Tests

```bash
go test ./...
```

## License

This project is licensed under the GPL v3+ 
See the [LICENSE](LICENSE) file for details.

## Acknowledgments

Inspired by other excellent code analysis tools:
- [tokei](https://github.com/XAMPPRocky/tokei) - Fast code statistics
- [scc](https://github.com/boyter/scc) - Sloc Cloc and Code counter

---

**Questions or Issues?** Please open an issue on [GitHub](https://github.com/yourusername/goloc/issues).
