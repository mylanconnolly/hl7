package main

// This example program accepts a single flag (-file) with an argument of the
// filename. It then opens that file and parses all of the HL7 messages within
// it. There can be one or multiple messages per file.
//
// Note that the functionality of the library allows you to use any io.Reader
// when building an hl7.Reader, so this could just as easily be a TCP connection
// or bytes.Buffer or something else (the tests use bytes.Buffers).

import (
	"flag"
	"fmt"
	"os"

	"github.com/mylanconnolly/hl7"
)

var filename = flag.String("file", "", "The file to parse")

func main() {
	flag.Parse()

	if filename == nil || *filename == "" {
		fmt.Fprintln(os.Stderr, "Must specify the filename using the -file flag!")
		os.Exit(1)
	}
	file, err := os.Open(*filename)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error encountered opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	reader := hl7.NewReader(file)

	if err = reader.EachMessage(func(msg *hl7.Message) error {
		fmt.Println("Found a message!")
		return nil
	}); err != nil {
		fmt.Fprintln(os.Stderr, "Error encountered:", err)
	}
}
