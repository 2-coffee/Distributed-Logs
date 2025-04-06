package server

import (
	"bytes"
	"testing"
)

func TestSplitLastMessage(t *testing.T) {
	res := []byte("100\n101\n10")
	testTruncated, testRest := []byte("100\n101\n"), []byte("10")

	Truncated1, Rest1, err := splitLastMessage(res)
	if err != nil {
		t.Errorf("splitLastMessage(%q): got error %v", string(res), err)
	}
	if !bytes.Equal(testTruncated, Truncated1) || !bytes.Equal(testRest, Rest1) {
		t.Errorf("splitLastMessage(%q): got %q, %q; wanted %q, %q", string(res), string(Truncated1), string(Rest1), string(testTruncated), string(testRest))
	}
}

// could not get full message
func TestSplitLastMessageSmallBuffer(t *testing.T) {
	res := []byte("10010110")
	_, _, err := splitLastMessage(res)
	if err == nil {
		t.Errorf("splitLastMessage(%q): did not get an error, want an error here.", string(res))
	}
}
