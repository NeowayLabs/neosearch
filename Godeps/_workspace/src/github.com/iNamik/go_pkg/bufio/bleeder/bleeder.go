// Bleeder provides a mechanism for seamlessly exhausting the buffered content of
// a bufio.Reader, then delegating further read requests to a separate reader.
//
// The original use case for this implementation was an efficient means of
// creating an 'expanding' bufio.Reader.  Just continuing to create new bufio.Readers
// using the previous bufio.Reader as the input could cause a lot of unnecessary memory
// usage as all the previous buffers would still be hanging around.
//
// Using Bleeder, you only store one extra bufio.Reader, and only until it has been
// exhausted, which will very often be on the first call to read(), unless the new
// reader is smaller than the old reader.
//
// The only caveat to this use case is that you need to know the original io.Reader
// as bufio.Reader does not provide access to it (as of this writing)
//
// *NOTE* Although the original use case for the function expects the specified
// io.Reader to the be the same one that feeds the specified bufio.Reader, there
// is no enforcement of the rule, so you could chain any io.Reader to the bufio.Reader.

package bleeder

import (
	"bufio"
	"io"
)

// New returns a new io.Reader that will exhaust (bleed) the specified
// bufio.Reader, then delegate any further requests to the specified io.Reader
//
// *NOTE* Although the original use case for the function expects the specified
// io.Reader to the be the same one that feeds the specified bufio.Reader, there
// is no enforcement of the rule, so you could chain any io.Reader to the bufio.Reader.
func New(br *bufio.Reader, r io.Reader) io.Reader {
	return &bleeder{br: br, r: r}
}

// bleeder stores internal state
type bleeder struct {
	br *bufio.Reader
	r  io.Reader
}

// Read bleeds the bufio.Reader then delegates to the io.Reader
func (b *bleeder) Read(p []byte) (n int, err error) {
	var head int = 0
	// If bufio.Reader is still active
	if b.br != nil {
		avail := b.br.Buffered() // Number of bytes available
		// If there are still any bytes available
		if avail > 0 {
			req := len(p) // Bytes requested
			// If requested bytes is less/equal available bytes
			// i.e. Our buffered reader can handle the entire request
			if req <= avail {
				// Read requested bytes
				n, err = b.br.Read(p[0:req])
				// If error or too few bytes read
				if err != nil || n < req {
					return
				}
				// If we've exhausted all available bytes
				if n == avail {
					b.br = nil // Free the bufio reader
				}
				return
				// Else requested bytes more than available bytes
				// i.e. The buffered reader can only partially handle the request
			} else {
				// Read available bytes
				n, err = b.br.Read(p[0:avail])
				// If error or too few bytes read
				if err != nil || n < avail {
					return
				}
				b.br = nil // Free the bufio reader
				head = n   // Start the standard reader here
			}
			// Else no more bytes available
		} else {
			b.br = nil // Free the bufio reader
		}
	}
	n, err = b.r.Read(p[head:])
	return n + head, err
}
