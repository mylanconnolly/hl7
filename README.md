# HL7

This is a basic HL7 parser written in Go.

## Installation

Installation is simple:

```bash
$ go get github.com/mylanconnolly/hl7
```

## Usage

Usage is meant to mimick the semantics of a reader. For example:

```go
package main

func main() {
  // Error handling elided for brevity
  file, _ := os.Open("/path/to/file.hl7")
  reader := hl7.NewReader(file)

  for {
    msg, err := reader.ReadMessage()

    // The only error you'll encounter here is io.EOF, which means there are no
    // messages.
    if err != nil {
      break
    }
    for {
      segment, err := msg.ReadSegment()

      // Like reading messages, the only error you'll encounter here is io.EOF,
      // for the same reason (no more segments).
      if err != nil {
        break
      }
      // Do something with the segment data
    }
  }
}
```

Further usage information can be found on
[godocs.org](https://godoc.org/github.com/mylanconnolly/hl7).

## TODO

I would like to add the following functionality, but it's not on the immediate
schedule:

- [ ] A way to handle unmarshalling using Go semantics (struct tags, etc.).
- [ ] A way to marshal data back into HL7.
- [ ] Some validation of the input data (this isn't likely; it means this
      program will need to know a lot about HL7 and I might not have time to
      implement it correctly).
