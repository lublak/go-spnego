package pure

import (
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/Azure/go-ntlmssp"
	"github.com/lublak/go-spnego/internal"
	"github.com/lublak/go-spnego/option"
)

type ntlmRoundTripper struct {
	r       http.RoundTripper
	options option.AuthOptions
}

func (t *ntlmRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.options.User == nil {
		return nil, errors.New("anonymous authentication not supported")
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

	req.Header.Set("Authorization", internal.EncodeNtlmToken(negotiate))

	res, err := t.r.RoundTrip(req)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusUnauthorized {
		return res, nil
	}

	ntlmToken, err := internal.FindAndDecodeNtlmToken(res.Header.Values("Www-Authenticate"))
	internal.DiscardResponseBody(res)
	if err != nil {
		return nil, err
	}

	auth, err := ntlmssp.NewAuthenticateMessage(ntlmToken, internal.JoinDomainAndName(t.options.User.Domain, t.options.User.Name), t.options.User.Password, &ntlmssp.AuthenticateMessageOptions{
		WorkstationName: workstation,
		PasswordHashed:  false,
	})

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", internal.EncodeNtlmToken(auth))

	if reqBody != nil {
		reqBody.Seek(pos, io.SeekStart)
	}

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
