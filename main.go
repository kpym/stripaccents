// strip accents is a small utility to remove accents from a string or a file.
// Usage:
//
//		stripaccents -s "Café" => "Cafe"
//		stripaccents "Café" => "Cafe"
//		stripaccents -i in.txt -o out.txt
//		stripaccents in.txt
//		stripaccents -i in.txt > out.txt
//	 	cat in.txt | stripaccents > out.txt
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// check is a helper function to check for errors
func check(e error, msg string) {
	if e != nil {
		fmt.Fprintln(os.Stderr, msg)
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}
}

// define the usage message
func init() {
	flag.Usage = func() {
		// get the executable name without the path
		exe := filepath.Base(os.Args[0])
		// remove the .exe extension on Windows
		if filepath.Ext(exe) == ".exe" {
			exe = exe[:len(exe)-4]
		}
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options]\n", exe)
		flag.PrintDefaults()
		// Examples
		examples := `
Examples:
	%s -s "Café"
	%s "Café"
	%s -i in.txt -o out.txt
	%s in.txt
	%s -i in.txt > out.txt
	cat in.txt | %s > out.txt
`
		fmt.Fprintln(flag.CommandLine.Output(), fmt.Sprintf(examples, exe, exe, exe, exe, exe, exe))
	}
}

// the command line parameters
var params = struct {
	in, out, str string
}{}

// parses the command line parameters
func getParams() {
	flag.StringVar(&params.in, "i", "", "input file (incompatible with -s)")
	flag.StringVar(&params.out, "o", "", "output file")
	flag.StringVar(&params.str, "s", "", "string to process (incompatible with -i)")
	flag.Parse()
	// check if both -s and -i are provided
	if params.in != "" && params.str != "" {
		fmt.Fprintln(os.Stderr, "Error: -i and -s cannot be used together")
		flag.Usage()
		os.Exit(1)
	}
	// check if both -s and -i are not provided check the first positional argument
	if params.str == "" && params.in == "" && flag.NArg() > 0 {
		str := flag.Arg(0)
		// check if str is a file
		if _, err := os.Stat(str); err == nil {
			params.in = str
		} else {
			params.str = str
		}
	}
}

func sameFile(file1, file2 string) bool {
	if file1 == "" || file2 == "" {
		return false
	}
	if file1 == file2 {
		return true
	}
	fi1, err := os.Stat(file1)
	if err != nil {
		return false
	}
	fi2, err := os.Stat(file2)
	if err != nil {
		return false
	}
	return os.SameFile(fi1, fi2)
}

func main() {
	getParams()

	var in io.Reader

	if params.str != "" {
		// convert the string to a reader
		in = bytes.NewBufferString(params.str)
	} else if params.in != "" {
		if sameFile(params.in, params.out) {
			// read the file into a buffer
			data, err := os.ReadFile(params.in)
			check(err, "Error reading input file")
			in = bytes.NewBuffer(data)
		} else {
			f, err := os.Open(params.in)
			check(err, "Error opening input file")
			defer f.Close()
			in = f
		}
	} else {
		// check if data is being piped
		fi, _ := os.Stdin.Stat()
		if (fi.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
			// no input is provided
			flag.Usage()
			os.Exit(1)
		}
		in = os.Stdin
	}

	var out io.Writer
	if params.out != "" {
		f, err := os.Create(params.out)
		check(err, "Error creating output file")
		defer f.Close()
		out = f
	} else {
		out = os.Stdout
	}

	stripAccents(in, out)
}

func stripAccents(in io.Reader, out io.Writer) {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	r := transform.NewReader(in, t)
	_, err := io.Copy(out, r)
	check(err, "Error copying data")
}
