# HL7

This is a basic HL7 parser written in Go.

## Installation

Installation is simple:

```bash
$ go get github.com/mylanconnolly/hl7
```

## Usage

Usage is meant to mimic the semantics of a reader. Additionally, the HL7 data
is read from an `io.Reader`, so it's relatively trivial to read HL7 data from
a variety of sources (a file, TCP connection, `bytes.Buffer`, etc.). Actual
parsing of the message happens when you fetch a segment. This laziness should
help with large messages, particularly when you find out halfway through the
message that you don't care about it anymore (maybe it doesn't have anything
to do with your business case).

For example usage, see the [example program](example/main.go).

Further usage information can be found
[here](https://pkg.go.dev/github.com/mylanconnolly/hl7).

## TODO

I would like to add the following functionality, but it's not on the immediate
schedule:

- [ ] A way to handle unmarshalling using Go semantics (struct tags, etc.).
- [ ] Some validation of the input data (this isn't likely; it means this
      program will need to know a lot about HL7 and I might not have time to
      implement it correctly).
