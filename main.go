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
	flag.StringVar(&params.in, "i", "", "input file")
	flag.StringVar(&params.out, "o", "", "output file")
	flag.StringVar(&params.str, "s", "", "string to process")
	flag.Parse()
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

func main() {
	getParams()

	if params.str != "" {
		fmt.Println(stripAccentsString(params.str))
		return
	}

	var in io.Reader
	if params.in != "" {
		f, err := os.Open(params.in)
		check(err, "Error opening input file")
		defer f.Close()
		in = f
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

	stripAccentsFile(in, out)
}

func stripAccentsString(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	s, _, _ = transform.String(t, s)
	return s
}

func stripAccentsFile(in io.Reader, out io.Writer) {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	r := transform.NewReader(in, t)
	_, err := io.Copy(out, r)
	check(err, "Error copying data")
}
