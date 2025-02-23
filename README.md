# writedisk

writedisk creates a specified number of files with a given size, filled with non-zero data. This
is intended for use as a diagnostic tool when testing disk speed, possible disk error
conditions (particularly on SSDs) and testing APFS cloned file functionality (for example,
using Hyperspace Mac app)

## Usage

The command writedisk can be followed by options. If no options are specified, an error
is reported. The --path value must be specified, all other options have a default value.
Options are specified using the Unix standard of "--" before long option names, or "-"
before short option names.

## Options

| Option            | Parameter   | Description |
|:------------------|:------------|:------------|
| --count, -c       |   number    | Number of files to create (default: 1) |
| --help, -h        |             | Display help output and exit |
| --path, -p        |   file path | Output path for where files are created (required)|
| --size, -s        |   number    | Size of each file in bytes (default is 10MB) |
| --therads, -t     |   number    | Number of threads to use (default is 2xCPU cores) |
| --verbose, -v     |             | If present, enable additonal logging |

The `--size` value can be expressed with a scale suffix; i.e. ten megabytes can be written as "10mb".
The scales "mb", "gb" and "kb" are allowed. A kilobyte is 1024 bytes.
