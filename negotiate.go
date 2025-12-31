package spnego

import (
	"io"
	"net/http"

	"github.com/lublak/go-spnego/internal"
	"github.com/lublak/go-spnego/option"
	"github.com/lublak/go-spnego/pure"
	"github.com/lublak/go-spnego/sspi"
)

type roundTripper struct {
	base           http.RoundTripper
	user           *option.User
	allowBasicAuth bool
	negotiate      http.RoundTripper
	ntlm           http.RoundTripper
}

func (t *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	reqBody, pos, err := internal.SetBodyAsSeekCloser(req)
	if err != nil {
		return nil, err
	}
	res, err := t.base.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusUnauthorized {
		return res, nil
	}

	wwwAuthenticateHeaders := res.Header.Values("WWW-Authenticate")

	if t.allowBasicAuth && internal.HasBasic(wwwAuthenticateHeaders) {
		if t.user == nil {
			return res, nil
		}
		if reqBody != nil {
			reqBody.Seek(pos, io.SeekStart)
		}
		internal.DiscardResponseBody(res)
		req.SetBasicAuth(internal.JoinDomainAndName(t.user.Domain, t.user.Name), t.user.Password)
		return t.base.RoundTrip(req)
	}

	if t.negotiate != nil && internal.HasNegotiate(wwwAuthenticateHeaders) {
		if reqBody != nil {
			reqBody.Seek(pos, io.SeekStart)
		}
		internal.DiscardResponseBody(res)
		return t.negotiate.RoundTrip(req)
	}

	if t.ntlm != nil && internal.HasNtlm(wwwAuthenticateHeaders) {
		if reqBody != nil {
			reqBody.Seek(pos, io.SeekStart)
		}
		internal.DiscardResponseBody(res)
		return t.ntlm.RoundTrip(req)
	}

	return res, err
}

func NewRoundTripper(base http.RoundTripper, api option.ApiType, options option.AuthOptions, allowBasicAuth bool) http.RoundTripper {
	if base == nil {
		base = http.DefaultTransport
	}
	var negotiate http.RoundTripper
	var ntlm http.RoundTripper
	switch api {
	case option.PURE:
		negotiate = pure.NewNegotiateRoundTripper(base)
		ntlm = pure.NewNtlmRoundTripper(base, options)
	case option.SSPI:
		negotiate = sspi.NewNegotiateRoundTripper(base, options)
		ntlm = sspi.NewNtlmRoundTripper(base, options)
	}
	if ntlm != nil && negotiate != nil {
		return nil
	}
	return &roundTripper{
		base:           base,
		user:           options.User,
		negotiate:      negotiate,
		ntlm:           ntlm,
		allowBasicAuth: allowBasicAuth,
	}
}
