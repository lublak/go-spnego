//go:build windows
// +build windows

package sspi

import (
	"net/http"

	"github.com/alexbrainman/sspi"
	"github.com/alexbrainman/sspi/ntlm"
	"github.com/lublak/go-spnego/internal"
	"github.com/lublak/go-spnego/option"
)

type ntlmRoundTripper struct {
	r       http.RoundTripper
	options option.AuthOptions
}

func (t *ntlmRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var cred *sspi.Credentials
	var err error

	if t.options.User == nil {
		cred, err = ntlm.AcquireCurrentUserCredentials()
		if err != nil {
			return nil, err
		}
	} else {
		cred, err = ntlm.AcquireUserCredentials(t.options.User.Domain, t.options.User.Name, t.options.User.Password)
		if err != nil {
			return nil, err
		}
	}

	defer cred.Release()

	client, token, err := ntlm.NewClientContext(cred)
	if err != nil {
		return nil, err
	}
	defer client.Release()

	reqBody, pos, err := internal.SetBodyAsSeekCloser(req)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", internal.EncodeNtlmToken(token))

	res, err := t.r.RoundTrip(req)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusUnauthorized {
		return res, nil
	}

	ntlmToken, err := internal.FindAndDecodeNtlmToken(res.Header.Values("Www-Authenticate"))

	if err != nil {
		internal.DiscardResponseBody(res)
		return nil, err
	}

	authenticate, err := client.Update(ntlmToken)

	if err != nil {
		internal.DiscardResponseBody(res)
		return nil, err
	}

	req.Header.Set("Authorization", internal.EncodeNtlmToken(authenticate))

	internal.ResetRoundTrip(reqBody, pos, res)

	return t.r.RoundTrip(req)
}

func NewNtlmRoundTripper(base http.RoundTripper, options option.AuthOptions) http.RoundTripper {
	if base == nil {
		base = http.DefaultTransport
	}
	return &ntlmRoundTripper{
		r:       base,
		options: options,
	}
}
