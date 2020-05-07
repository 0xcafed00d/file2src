picomerge
# file2src 
encodes a file as an array of bytes in a C/C++ source file, to allow data files to be compiled in a executable
  Usage: file2src [options] <input file name> 
  -P string
        text to be inserted at start of output
  -h    display help
  -l string
        language of src output (suppored: c, go) (default "c")
  -n string
        name of the created array (default "data")
  -o string
        name of output file. Output written to stdout if omitted
  -p string
        name of file to be insterted at start of output
  -t string
        type of the created array
  -z    place a Zero byte at the end of the data.
        Extends length of data by 1 byte.
        (useful for null terminating string data)

