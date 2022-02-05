// MIT license · Gregor Noczinski · gregor [at] noczinski [dot] eu · 12/2021

package jsonlex

import (
	"fmt"
	"io"
)

// EnsureAtLeastSingleByteUnreadableReader will provide an instance of the
// provided io.Reader which can at least unread a single byte.
func EnsureAtLeastSingleByteUnreadableReader(r io.Reader) UnreadableReader {
	if ur, ok := r.(UnreadableReader); ok {
		return ur
	}
	return &singleByteUnreadableReader{
		delegate: r,
	}
}

type singleByteUnreadableReader struct {
	delegate io.Reader
	lastByte byte
	state    singleByteUnreadableReaderState
}

func (r *singleByteUnreadableReader) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}

	switch r.state {
	case singleByteUnreadableReaderStateBuffered, singleByteUnreadableReaderStateEmpty:
		n, err := r.delegate.Read(p)
		if n > 0 {
			r.lastByte = p[n-1]
			r.state = singleByteUnreadableReaderStateBuffered
		}
		return n, err

	case singleByteUnreadableReaderStateRewind:
		p[0] = r.lastByte
		r.state = singleByteUnreadableReaderStateEmpty
		if len(p) == 1 {
			return 1, nil
		}
		n, err := r.delegate.Read(p[1:])
		return n + 1, err

	default:
		panic(fmt.Sprintf("unknown state: %d", r.state))
	}
}

func (r *singleByteUnreadableReader) UnreadByte() error {
	switch r.state {
	case singleByteUnreadableReaderStateBuffered:
		r.state = singleByteUnreadableReaderStateRewind
		return nil

	case singleByteUnreadableReaderStateEmpty, singleByteUnreadableReaderStateRewind:
		return io.ErrShortBuffer

	default:
		panic(fmt.Sprintf("unknown state: %d", r.state))
	}
}

type singleByteUnreadableReaderState uint8

const (
	singleByteUnreadableReaderStateEmpty singleByteUnreadableReaderState = iota
	singleByteUnreadableReaderStateBuffered
	singleByteUnreadableReaderStateRewind
)
