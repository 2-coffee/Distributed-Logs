package server

import (
	"io"
	"os"
)

const readBlockSize = 1 << 20

type OnDisk struct {
	fp *os.File
}

// Creates a server that stores all of its data on disk.
func NewOnDisk(fp *os.File) *OnDisk {
	return &OnDisk{fp: fp}
}

// Accepts messages from clients then store them on disc.
func (s *OnDisk) Write(msgs []byte) error {
	_, err := s.fp.Write(msgs)
	return err
}

// Reading from disk.
func (s *OnDisk) Read(offset uint64, maxSize uint64, w io.Writer) (err error) {
	// When streaming results, we will need to know when to stop
	// also cut off the last message if it was not fully read
	// Need to limit log size as well incase it is bigger than our buffer
	buf := make([]byte, maxSize)

	n, err := s.fp.ReadAt(buf, int64(offset)) // reading into buffer

	if n == 0 { // for some reason ReadAt was returning EOF when it didn't finish reading
		if err == io.EOF { // done
			return nil
		} else if err != nil {
			return err
		}
	}

	// message is greater than limit
	// cutting off the rest of the message
	truncated, _, err := splitLastMessage(buf[0:n])
	if err != nil {
		return err
	}
	// keep reading from file and writing to memory
	if _, err := w.Write(truncated); err != nil {
		return err
	}
	return nil
}

// Marks the current chunk as completed and free up memory
func (s *OnDisk) Ack() error {
	err := s.fp.Truncate(0)
	if err != nil {
		return err
	}

	// writes will continue from previous offset
	_, err = s.fp.Seek(0, os.SEEK_SET)
	return err
}

// func splitLastMessage(res []byte) (truncated []byte, rest []byte, err error) {
// 	n := len(res)
// 	if len(res) == 0 {
// 		return res, nil, nil
// 	}
// 	if res[n-1] == '\n' { // all logs read
// 		return res, nil, nil
// 	}
// 	lastLog := bytes.LastIndexByte(res, '\n')
// 	if lastLog < 0 {
// 		return nil, nil, errBufferTooSmall
// 	}
// 	lastLog += 1
// 	return res[0:lastLog], res[lastLog:], nil
// }
