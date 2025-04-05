package client

import "bytes"

const defaultScratchsize = 64 * 1024

type Simple struct {
	addrs []string
	buf   bytes.Buffer
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
		scratch = make([]byte, defaultScratchsize) // making a buffer incase we are waiting
	}
	n, err := s.buf.Read(scratch)
	if err != nil { // error in Read
		return nil, err
	}
	// return data and nil for error
	return scratch[0:n], nil
}
