# stripaccents

A small tool to remove accents from a string or a file written in [go](https://golang.org/).

## Usage

```shell
> stripaccents -h
Usage: stripaccents [options]
  -i string
        input file
  -o string
        output file
  -s string
        string to process

Examples:
        stripaccents -s "Café"
        stripaccents -i in.txt -o out.txt
        stripaccents -i in.txt > out.txt
        cat in.txt | stripaccents > out.txt
```

## Installation

### From binary

Download the binary from the [releases](github.com/kpym/stripaccents/releases) page.

### From source

```shell
> go install github.com/kpym/stripaccents@latest
```

## License

[MIT](LICENSE)
