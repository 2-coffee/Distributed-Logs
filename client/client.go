package client

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
)

const defaultScratchsize = 64 * 1024

var errBufferTooSmall = errors.New("buffer is too small to fit one message")

type Simple struct {
	addrs  []string
	cl     *http.Client // from Go http
	offset uint64
}

// Creating a new client
func NewSimple(addrs []string) *Simple {
	return &Simple{
		addrs: addrs,
		cl:    &http.Client{},
	}

}

// Sends msgs to server
func (s *Simple) Send(msgs []byte) error {
	resp, err := s.cl.Post(s.addrs[0]+"/write", "application/octer-stream", bytes.NewReader(msgs))

	if err != nil {
		return err
	}

	defer resp.Body.Close() // schedules to make sure we close incoming stream

	if resp.StatusCode != http.StatusOK {
		var b bytes.Buffer
		io.Copy(&b, resp.Body)                                                           // temporary to make every write succeed
		return fmt.Errorf("sending data: http code %d, %s", resp.StatusCode, b.String()) // shows the rest of the data that was not sent
	}
	io.Copy(io.Discard, resp.Body)
	return nil
}

// Receive data from server
// wait for new msgs or return an error in case something goes wrong
func (s *Simple) Receive(scratch []byte) ([]byte, error) {
	if scratch == nil {
		scratch = make([]byte, defaultScratchsize) // making a buffer
	}
	addrIndx := rand.Intn(len(s.addrs))
	addr := s.addrs[addrIndx]
	readURL := fmt.Sprintf("%s/read?off=%d&maxSize=%d", addr, s.offset, len(scratch)) // http request
	resp, err := s.cl.Get(readURL)                                                    // get request

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close() // schedules to make sure we close incoming stream

	if resp.StatusCode != http.StatusOK {
		var b bytes.Buffer
		io.Copy(&b, resp.Body)                                                                // temporary to make every write succeed
		return nil, fmt.Errorf("sending data: http code %d, %s", resp.StatusCode, b.String()) // shows the rest of the data that was not sent
	}

	b := bytes.NewBuffer(scratch[0:0])
	_, err = io.Copy(b, resp.Body) // writing

	if err != nil {
		return nil, err
	}
	// 0 bytes reads, no errors encountered. Copy will return nil. This is EOF by convention
	if b.Len() == 0 {
		if err := s.ackCurrentChunk(addr); err != nil {
			return nil, err
		}
		return nil, io.EOF // without this error, it will keep looking for more lines
	}

	s.offset += uint64(b.Len())
	return b.Bytes(), nil
}

func (s *Simple) ackCurrentChunk(addr string) error {
	resp, err := s.cl.Get(addr + "/ack") // get request

	if err != nil {
		return err
	}

	defer resp.Body.Close() // schedules to make sure we close incoming stream
	if resp.StatusCode != http.StatusOK {
		var b bytes.Buffer
		io.Copy(&b, resp.Body)                                                           // temporary to make every write succeed
		return fmt.Errorf("sending data: http code %d, %s", resp.StatusCode, b.String()) // shows the rest of the data that was not sent
	}

	io.Copy(io.Discard, resp.Body) // don't want anything

	return nil
}
