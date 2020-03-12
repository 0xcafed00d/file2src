package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func exitOnError(e error, msg string) {
	if e != nil {
		abend(msg + " : " + e.Error())
	}
}

func abend(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(-1)
}

// Config settings from invocation flags
type Config struct {
	help           bool
	dataName       string
	dataType       string
	prefixFilename string
	prefixText     string
	outputFilename string
	nullTerminate  bool
}

var config Config

func init() {
	flag.BoolVar(&config.help, "h", false, "display help")
	flag.StringVar(&config.dataName, "n", "data", "name of the created array")
	flag.StringVar(&config.dataType, "t", "unsigned char", "type of the created array")
	flag.StringVar(&config.prefixFilename, "p", "", "name of file to be insterted at start of output")
	flag.StringVar(&config.prefixText, "P", "", "text to be inserted at start of output")
	flag.StringVar(&config.outputFilename, "o", "", "name of output file. Output written to stdout if omitted")
	flag.BoolVar(&config.nullTerminate, "z", false, "place a Zero byte at the end of the data.\nExtends length of data by 1 byte.\n(useful for null terminating string data)")

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "file2src: encodes a file as an array of bytes in a C/C++ source file, to allow data files to be compiled in a executable")
		fmt.Fprintln(os.Stderr, "  Usage: file2src [options] <input file name> ")
		flag.PrintDefaults()
	}
}

func processFile(size int64, input io.Reader, output io.Writer, conf *Config) error {

	if config.nullTerminate {
		size++
	}

	fmt.Fprintf(output, "const size_t %s_sz = %v;\n", conf.dataName, size)
	fmt.Fprintf(output, "%s %s[%s_sz] = {\n", conf.dataType, conf.dataName, conf.dataName)

	byteCount := int64(0)

	for {
		buffer := make([]byte, 16)
		n, err := input.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}

		if n == 0 && err == io.EOF {
			break
		}

		fmt.Fprint(output, "\t")

		for i := 0; i < n; i++ {
			fmt.Fprintf(output, "0x%02x", buffer[i])
			byteCount++

			if byteCount < size {
				fmt.Fprint(output, ",")
			}
		}
		fmt.Fprint(output, "\n")
	}
	if config.nullTerminate {
		fmt.Fprintln(output, "\t0")
	}

	fmt.Fprintln(output, "};")

	return nil
}

func unescapeString(str string) string {
	out := []rune{'"'}

	for _, r := range str {
		if r == '"' {
			out = append(out, '\\')
		}
		out = append(out, r)
	}
	out = append(out, '"')

	unescapedString, err := strconv.Unquote(string(out))
	exitOnError(err, "failed to unquote: "+string(out))

	return unescapedString
}

func main() {
	flag.Parse()

	if len(flag.Args()) == 0 || config.help {
		flag.Usage()
		os.Exit(1)
	}

	input := flag.Args()[0]
	outfile := os.Stdout

	if config.outputFilename != "" {
		var err error
		outfile, err = os.Create(config.outputFilename)
		exitOnError(err, "Cannot create output file")
	}

	if config.prefixFilename != "" {
		prefixfile, err := os.Open(config.prefixFilename)
		exitOnError(err, "Cant Open Prefix file")
		defer prefixfile.Close()
		_, err = io.Copy(outfile, prefixfile)
		exitOnError(err, "Failed Writing Prefix file")
	}

	if config.prefixText != "" {
		txt := unescapeString(config.prefixText)
		prefixtext := strings.NewReader(txt)
		_, err := io.Copy(outfile, prefixtext)
		exitOnError(err, "Failed Writing Prefix text")
	}

	infile, err := os.Open(input)
	info, err := infile.Stat()
	exitOnError(err, "Cant Open Inputfile")
	exitOnError(processFile(info.Size(), infile, outfile, &config), "Error reading file")

	if outfile != os.Stdout {
		outfile.Close()
	}
}
