package main

import (
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

func processFile(size int64, input io.Reader, output io.Writer, conf *config) error {

	fmt.Fprintf(output, "const size_t %s_sz = %v;\n", conf.dataName, size)
	fmt.Fprintf(output, "uint8_t %s[] = {\n", conf.dataName)

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

type config struct {
	dataName string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("picomerge: merges pico8 p8 cart with its included source files to produce a single file p8 cart.")
		fmt.Println("  Usage: picomerge <input.p8> [output.p8]")
		fmt.Println("  if no output file is specified, then output is printed to console")
		os.Exit(-1)
	}

	input := os.Args[1]
	outfile := os.Stdout

	if len(os.Args) > 2 {
		output := os.Args[2]
		var err error
		outfile, err = os.Create(output)
		exitOnError(err, "Cannot create output file")
	}

	conf := config{}
	conf.dataName = "data"

	infile, err := os.Open(input)
	info, err := infile.Stat()
	exitOnError(err, "Cant Open Inputfile")
	exitOnError(processFile(info.Size(), infile, outfile, &conf), "Error reading file")

	if outfile != os.Stdout {
		outfile.Close()
	}
}
