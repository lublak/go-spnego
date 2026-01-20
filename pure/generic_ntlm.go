package pure

import (
	"errors"
	"net/http"
	"os"

	"github.com/Azure/go-ntlmssp"
	"github.com/lublak/go-spnego/internal"
	"github.com/lublak/go-spnego/options"
)

type genericNtlmRoundTripper struct {
	r              http.RoundTripper
	options        options.Options
	negotiateToken bool
}

func (t *genericNtlmRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.options.User == nil {
		return nil, errors.New("anonymous authentication not supported for ntlm")
	}
	workstation, _ := os.Hostname()
	domain, _ := internal.GetDomain()
	negotiate, err := ntlmssp.NewNegotiateMessage(domain, workstation)
	if err != nil {
		return nil, err
	}

	reqBody, pos, err := internal.SetBodyAsSeekCloser(req)
	if err != nil {
		return nil, err
	}

	if t.negotiateToken {
		req.Header.Set("Authorization", internal.EncodeNegotiateToken(negotiate))
	} else {
		req.Header.Set("Authorization", internal.EncodeNtlmToken(negotiate))
	}

	res, err := t.r.RoundTrip(req)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusUnauthorized {
		return res, nil
	}

	var ntlmToken []byte

	if t.negotiateToken {
		ntlmToken, err = internal.FindAndDecodeNegotiateToken(res.Header.Values("Www-Authenticate"))
	} else {
		ntlmToken, err = internal.FindAndDecodeNtlmToken(res.Header.Values("Www-Authenticate"))
	}

	if err != nil {
		internal.DiscardResponseBody(res)
		return nil, err
	}

	auth, err := ntlmssp.NewAuthenticateMessage(ntlmToken, internal.JoinDomainAndName(t.options.User.Domain, t.options.User.Name), t.options.User.Password, &ntlmssp.AuthenticateMessageOptions{
		WorkstationName: workstation,
		PasswordHashed:  false,
	})

	if err != nil {
		internal.DiscardResponseBody(res)
		return nil, err
	}

	if t.negotiateToken {
		req.Header.Set("Authorization", internal.EncodeNegotiateToken(auth))
	} else {
		req.Header.Set("Authorization", internal.EncodeNtlmToken(auth))
	}

	internal.ResetRoundTrip(reqBody, pos, res)

	return t.r.RoundTrip(req)
}

func newGenericNtlmRoundTripper(base http.RoundTripper, spnegoOptions options.Options, negotiateToken bool) http.RoundTripper {
	if base == nil {
		base = http.DefaultTransport
	}
	return &genericNtlmRoundTripper{
		r:              base,
		options:        spnegoOptions,
		negotiateToken: negotiateToken,
	}
}
