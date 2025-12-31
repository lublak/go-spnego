package internal

import "io"

func NopSeekerCloser(r io.ReadSeeker) io.ReadSeekCloser {
	if _, ok := r.(io.WriterTo); ok {
		return nopCloserWriterTo{r}
	}
	return nopSeekerCloser{r}
}

type nopSeekerCloser struct {
	io.ReadSeeker
}

func (nopSeekerCloser) Close() error { return nil }

type nopCloserWriterTo struct {
	io.ReadSeeker
}

func (nopCloserWriterTo) Close() error { return nil }

func (c nopCloserWriterTo) WriteTo(w io.Writer) (n int64, err error) {
	return c.ReadSeeker.(io.WriterTo).WriteTo(w)
}
