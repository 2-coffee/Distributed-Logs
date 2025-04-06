package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/2-coffee/Distributed-Logs/client"
)

const maxN = 10000000
const maxBufferSize = 1 << 20

func main() {
	s := client.NewSimple([]string{"http://localhost:8080"})
	// don't want to keep all of the data in memory
	want, err := send(s)
	if err != nil {
		log.Fatalf("Send error: %v send", err)
	}

	got, err := receive(s)
	if err != nil {
		log.Fatalf("Receive error: %v receive", err)
	}

	if want != got {
		log.Fatalf("The expected sum %d is not equal to the actual sum %d", want, got)
	}

	log.Printf("The test passed")
}

// checking client received msgs
func receive(s *client.Simple) (sum int64, err error) {
	buf := make([]byte, maxBufferSize)

	for {
		res, err := s.Receive(buf)
		if err == io.EOF { // end of file
			return sum, nil
		} else if err != nil {
			return 0, err
		}
		// each log has a new line at the end
		ints := strings.Split(string(res), "\n")
		for _, str := range ints {
			if str == "" {
				continue
			}
			i, err := strconv.Atoi(str)
			if err != nil {
				return 0, err
			}
			sum += int64(i)
		}
	}
}

// send to client and also return how many was sent
func send(s *client.Simple) (sum int64, err error) {
	var b bytes.Buffer
	for i := 0; i <= maxN; i++ {
		sum += int64(i)

		fmt.Fprintf(&b, "%d\n", i)

		if b.Len() >= maxBufferSize {
			if err := s.Send(b.Bytes()); err != nil {
				return 0, err
			}
			b.Reset()
		}
	}
	if b.Len() != 0 {
		if err := s.Send(b.Bytes()); err != nil {
			return 0, err
		}
	}
	return sum, nil
}
