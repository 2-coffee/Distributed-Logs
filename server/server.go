package server

import (
	"bytes"
	"errors"
	"io"
)

var (
	errBufferTooSmall          = errors.New("buffer is too small to fit one message")
	errOffsetMustBeNonNegative = errors.New("offset value must be non-negative")
	errMaxSizeMustBePositive   = errors.New("maxsize must be a positive number")
)

type InMemory struct {
	buf []byte // main buffer
}

// Accepts messages from clients then store them on disc
func (s *InMemory) Write(msgs []byte) error {
	s.buf = append(s.buf, msgs...)
	return nil
}

// Copies data from the buffer that was written to the server
// then writes the data to the provided Writer starting with the offset provided.
// New offset is returned.
func (s *InMemory) Read(offset uint64, maxSize uint64, w io.Writer) (err error) {
	// get offset from client
	maxOffset := uint64(len(s.buf))

	if offset >= maxOffset {
		// at the end
		return nil
	} else if offset+maxSize >= maxOffset {
		// Client side wants to read until the end
		w.Write(s.buf[offset:])
		return nil
	}

	// fmt.Println("\n\n")
	// fmt.Println(scratch[0:n])

	// one of the problems I ran into is that read would get cut off
	truncated, _, err := splitLastMessage(s.buf[offset : offset+maxSize])

	if err != nil {
		return err
	}

	if _, err := w.Write(truncated); err != nil {
		return err
	}

	return nil
}

// Marks the current chunk as completed and free up memory
func (s *InMemory) Ack() error {
	// TODO: mark chunk as completed?
	s.buf = nil
	return nil
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
