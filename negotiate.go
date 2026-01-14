package spnego

import (
	"net/http"
	"runtime"

	"github.com/lublak/go-spnego/internal"
	"github.com/lublak/go-spnego/option"
	"github.com/lublak/go-spnego/pure"
	"github.com/lublak/go-spnego/sspi"
)

type roundTripper struct {
	base      http.RoundTripper
	negotiate http.RoundTripper
	ntlm      http.RoundTripper
	options   option.AuthOptions
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

	if t.options.AllowBasicAuth && internal.HasBasic(wwwAuthenticateHeaders) {
		if t.options.User == nil {
			return res, nil
		}
		internal.ResetRoundTrip(reqBody, pos, res)
		req.SetBasicAuth(internal.JoinDomainAndName(t.options.User.Domain, t.options.User.Name), t.options.User.Password)
		return t.base.RoundTrip(req)
	}

	if t.negotiate != nil && internal.HasNegotiate(wwwAuthenticateHeaders) {
		internal.ResetRoundTrip(reqBody, pos, res)
		return t.negotiate.RoundTrip(req)
	}

	if t.ntlm != nil && internal.HasNtlm(wwwAuthenticateHeaders) {
		internal.ResetRoundTrip(reqBody, pos, res)
		return t.ntlm.RoundTrip(req)
	}

	return res, err
}

func NewRoundTripper(base http.RoundTripper, api option.ApiType, options option.AuthOptions) http.RoundTripper {
	if base == nil {
		base = http.DefaultTransport
	}
	var negotiate http.RoundTripper
	var ntlm http.RoundTripper
	switch api {
	case option.PURE:
		negotiate = pure.NewNegotiateRoundTripper(base, options)
		ntlm = pure.NewNtlmRoundTripper(base, options)
	case option.SSPI:
		negotiate = sspi.NewNegotiateRoundTripper(base, options)
		ntlm = sspi.NewNtlmRoundTripper(base, options)
	case option.AUTO:
		if runtime.GOOS == "windows" {
			negotiate = sspi.NewNegotiateRoundTripper(base, options)
			ntlm = sspi.NewNtlmRoundTripper(base, options)
		} else {
			negotiate = pure.NewNegotiateRoundTripper(base, options)
			ntlm = pure.NewNtlmRoundTripper(base, options)
		}
	}
	if ntlm != nil && negotiate != nil {
		return nil
	}
	return &roundTripper{
		base:      base,
		negotiate: negotiate,
		ntlm:      ntlm,
		options:   options,
	}
}
