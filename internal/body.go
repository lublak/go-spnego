package internal

import (
	"bytes"
	"io"
	"net/http"
)

func SetBodyAsSeekCloser(req *http.Request) (io.ReadSeekCloser, int64, error) {
	if req.Body == nil {
		return nil, 0, nil
	}

	if body, ok := req.Body.(io.ReadSeeker); ok {
		pos, err := body.Seek(0, io.SeekCurrent)
		if err == nil {
			if body, ok := body.(io.ReadSeekCloser); ok {
				return body, pos, nil
			}
			sc := NopSeekerCloser(body)
			req.Body = sc
			return sc, pos, nil
		}
	}
	body := bytes.Buffer{}
	_, err := body.ReadFrom(req.Body)
	if err != nil {
		return nil, 0, err
	}
	req.Body.Close()
	sc := NopSeekerCloser(bytes.NewReader(body.Bytes()))
	req.Body = sc
	return sc, 0, nil
}

func DiscardResponseBody(res *http.Response) {
	_, _ = io.Copy(io.Discard, res.Body)
	_ = res.Body.Close()
}

func ResetRoundTrip(reqBody io.ReadSeekCloser, pos int64, res *http.Response) {
	if reqBody != nil {
		reqBody.Seek(pos, io.SeekStart)
	}
	DiscardResponseBody(res)
}
