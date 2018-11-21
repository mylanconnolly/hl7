package main

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
	i := 0

	for {
		i++
		msg, err := reader.ReadMessage()

		if err != nil {
			break
		}
		for {
			_, err = msg.ReadSegment()

			if err != nil {
				break
			}
			fmt.Println("Found segment!")
		}
	}
}
