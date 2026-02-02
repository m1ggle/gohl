# gohl - MySQL CLI Tool

A lightweight and interactive MySQL command-line interface tool built with Go.

## Features

- Interactive MySQL shell with autocomplete support
- Syntax highlighting for SQL keywords
- Tabular result display
- Support for common MySQL commands (SELECT, INSERT, UPDATE, DELETE, SHOW, etc.)
- Configuration file support
- Auto-completion for common MySQL keywords

## Installation

First, make sure you have Go installed on your machine. Then you can install `gohl` using:

```bash
go install github.com/yahve/gohl@latest
```

Or clone the repository and build manually:

```bash
git clone https://github.com/yahve/gohl.git
cd gohl
make build
```

## Usage

### Login to MySQL Server

```bash
# Interactive login
gohl login

# Or with command line parameters
gohl login --host=localhost --port=3306 --user=root --pass=password --db=database_name
```

### Execute Single Query

```bash
gohl query "SELECT * FROM users LIMIT 10"
```

### Command Line Options

Common options available for both `login` and `query` commands:

- `--host` or `-H`: MySQL host (default: "127.0.0.1")
- `--port` or `-P`: MySQL port (default: "3306")
- `--user` or `-u`: MySQL user (default: "root")
- `--pass` or `-p`: MySQL password (if empty will prompt)
- `--db` or `-d`: Database name
- `--config`: Config file path (default: $HOME/.gohl.json)

## Autocomplete Features

The interactive shell provides smart autocomplete for:

- SQL Keywords: SELECT, INSERT, UPDATE, DELETE, CREATE, DROP, ALTER, SHOW, DESCRIBE, etc.
- Partial matches: Type "SH" to get "SHOW", "SEL" to get "SELECT", etc.
- Database and table related keywords

## Configuration

The tool supports configuration via YAML file. By default, it looks for `conf/ghl.yaml` in the current directory.

Example configuration:

```yaml
# MySQL Database Configuration
database:
  dbname: mysql
  user: root
  password: your_password
  host: localhost
  port: 3306
```

## Build

To build the project from source:

```bash
make build          # Build binary
make test           # Run tests
make run            # Run directly
make clean          # Clean build artifacts
```

## Development

This project uses:

- [Cobra](https://github.com/spf13/cobra) for CLI framework
- [Viper](https://github.com/spf13/viper) for configuration management
- [go-prompt](https://github.com/c-bata/go-prompt) for interactive shell with autocomplete
- [tablewriter](https://github.com/olekukonko/tablewriter) for tabular data display
- [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql) for MySQL driver

## License

MIT License
