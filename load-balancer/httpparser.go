package main

import (
	"bytes"
)

var httpString = []byte("GET /test HTTP/1.1\r\nheader1: header1\r\nheader2: header2\r\nheader3: header3\r\nheader4: header4\r\nheader5: header5\r\nheader6: header6\r\n\r\ncontent\r\n")
var CRFL = []byte{'\r', '\n'}

type Header struct {
	name  []byte
	value []byte
}

type HttpParser struct {
	rawData []byte
	method  []byte
	url     []byte
	version []byte
	headers []Header
	content []byte
}

func NewHttpParser(r []byte, lazy bool) *HttpParser {
	parser := &HttpParser{
		rawData: r,
	}
	if !lazy {
		parser.Parse()
	}
	return parser
}
func (parser *HttpParser) Parse() {
	l, r := 0, 0
	for parser.rawData[r] != ' ' {
		r++
	}
	parser.method = parser.rawData[l:r]
	l, r = r+1, r+1
	for parser.rawData[r] != ' ' {
		r++
	}
	parser.url = parser.rawData[l:r]
	l, r = r+1, r+1
	for parser.rawData[r] != '/' {
		r++
	}
	l, r = r+1, r+1
	for !(parser.rawData[r] == '\r' && parser.rawData[r+1] == '\n') {
		r++
	}
	parser.version = parser.rawData[l:r]
	l, r = r+2, r+2
	// parser.headers = make([]Header, 0, 5)
	parser.headers = nil
	for !(parser.rawData[l] == '\r' && parser.rawData[l+1] == '\n') {
		header := Header{}
		for parser.rawData[r] != ':' {
			r++
		}
		header.name = parser.rawData[l:r]
		l, r = r+2, r+2
		for !(parser.rawData[r] == '\r' && parser.rawData[r+1] == '\n') {
			r++
		}
		header.value = parser.rawData[l:r]
		l, r = r+2, r+2
		if !bytes.Equal(header.name, []byte("Host")) {
			parser.headers = append(parser.headers, header)
		}
	}
	l, r = r+2, r+2
	parser.content = parser.rawData[l:]
}

func (httpparser *HttpParser) Serealize(writer *[]byte) {
	// starter line
	*writer = append(*writer, httpparser.method...)
	*writer = append(*writer, []byte{' '}...)
	*writer = append(*writer, httpparser.url...)
	*writer = append(*writer, []byte{' '}...)
	*writer = append(*writer, []byte("HTTP/")...)
	*writer = append(*writer, httpparser.version...)
	*writer = append(*writer, CRFL...)
	// headers
	for _, header := range httpparser.headers {
		*writer = append(*writer, header.name...)
		*writer = append(*writer, []byte(": ")...)
		*writer = append(*writer, header.value...)
		*writer = append(*writer, CRFL...)
	}
	*writer = append(*writer, CRFL...)
	// content
	*writer = append(*writer, httpparser.content...)
	*writer = append(*writer, CRFL...)
}
