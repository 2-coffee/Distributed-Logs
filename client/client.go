package client

import (
	"bytes"
	"errors"
)

const defaultScratchsize = 64 * 1024

var errBufferTooSmall = errors.New("buffer is too small to fit one message")

type Simple struct {
	addrs         []string
	buf           bytes.Buffer
	restOfLastbuf bytes.Buffer
}

// Creating a new client
func NewSimple(addrs []string) *Simple {
	return &Simple{
		addrs: addrs,
	}

}

// Sends msgs to server
func (s *Simple) Send(msgs []byte) error {
	_, err := s.buf.Write(msgs)
	return err
}

// Receive data from server
// wait for new msgs or return an error in case something goes wrong
func (s *Simple) Receive(scratch []byte) ([]byte, error) {
	if scratch == nil {
		scratch = make([]byte, defaultScratchsize) // making a buffer
	}
	offset := 0

	if s.restOfLastbuf.Len() > 0 {
		if s.restOfLastbuf.Len() >= len(scratch) {
			return nil, errBufferTooSmall
		}
		n, err := s.restOfLastbuf.Read(scratch) // read the last message that got cut off into buffer
		if err != nil {
			return nil, err
		}
		offset += n

	}

	n, err := s.buf.Read(scratch[offset:]) // read into buffer
	if err != nil {                        // error in Read
		return nil, err
	}
	// fmt.Println("\n\n")
	// fmt.Println(scratch[0:n])

	// one of the problems I ran into is that read would get cut off
	truncated, rest, err := splitLastMessage(scratch[0 : n+offset])

	s.restOfLastbuf.Reset()
	s.restOfLastbuf.Write(rest)
	if err != nil {
		return nil, err
	}

	return truncated, nil
}

func splitLastMessage(res []byte) (truncated []byte, rest []byte, err error) {
	n := len(res)
	if len(res) == 0 {
		return res, nil, nil
	}
	if res[n-1] == '\n' { // all logs read
		return res, nil, nil
	}
	lastLog := bytes.LastIndexByte(res, '\n')
	if lastLog < 0 {
		return nil, nil, errBufferTooSmall
	}
	lastLog += 1
	return res[0:lastLog], res[lastLog:], nil
}
