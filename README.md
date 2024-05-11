# stripaccents

A small tool to remove accents from a string or a file.

## Usage

```shell
> stripaccents -h
Usage: stripaccents [options]
  -i string
        input file (incompatible with -s)
  -o string
        output file
  -s string
        string to process (incompatible with -i)

Examples:
        stripaccents -s "Café"
        stripaccents "Café"
        stripaccents -i in.txt -o out.txt
        stripaccents in.txt
        stripaccents -i in.txt > out.txt
        cat in.txt | stripaccents > out.txt
```

## Installation

### From binary

Download the binary from the [releases](https://github.com/kpym/stripaccents/releases) page.

### From source

```shell
> go install github.com/kpym/stripaccents@latest
```

## License

[MIT](LICENSE)
