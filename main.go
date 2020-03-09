package main

import (
	"flag"
	"fmt"
	"io"
	"os"
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
	outputFilename string
}

var config Config

func init() {
	flag.BoolVar(&config.help, "h", false, "display help")
	flag.StringVar(&config.dataName, "n", "data", "name of the created array")
	flag.StringVar(&config.dataType, "t", "uint8_t", "type of the created array")
	flag.StringVar(&config.prefixFilename, "p", "", "name of file to be insterted at start of output")
	flag.StringVar(&config.outputFilename, "o", "", "name of output file. Output written to stdout if omitted")

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "file2src: ")
		fmt.Fprintln(os.Stderr, "  Usage: file2src [options] <input> ")
		flag.PrintDefaults()
	}
}

func processFile(size int64, input io.Reader, output io.Writer, conf *Config) error {

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
	fmt.Fprintln(output, "};")

	return nil
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

	infile, err := os.Open(input)
	info, err := infile.Stat()
	exitOnError(err, "Cant Open Inputfile")
	exitOnError(processFile(info.Size(), infile, outfile, &config), "Error reading file")

	if outfile != os.Stdout {
		outfile.Close()
	}
}
