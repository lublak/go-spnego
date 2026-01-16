package sspi

import (
	"errors"
	"net/http"

	"github.com/alexbrainman/sspi"
	"github.com/alexbrainman/sspi/negotiate"
	"github.com/lublak/go-spnego/internal"
	spnego_options "github.com/lublak/go-spnego/options"
)

type negotiateRoundTripper struct {
	r       http.RoundTripper
	options spnego_options.Options
}

func (t *negotiateRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var cred *sspi.Credentials
	var err error

	if t.options.User == nil {
		cred, err = negotiate.AcquireCurrentUserCredentials()
		if err != nil {
			return nil, err
		}
	} else {
		cred, err = negotiate.AcquireUserCredentials(t.options.User.Domain, t.options.User.Name, t.options.User.Password)
		if err != nil {
			return nil, err
		}
	}

	defer cred.Release()

	client, token, err := negotiate.NewClientContext(cred, internal.GetSpnFromRequest(req))
	if err != nil {
		return nil, err
	}
	defer client.Release()

	req.Header.Set("Authorization", internal.EncodeNegotiateToken(token))

	res, err := t.r.RoundTrip(req)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return res, nil
	}

	negotiateToken, err := internal.FindAndDecodeNegotiateToken(res.Header.Values("WWW-Authenticate"))
	if err != nil {
		internal.DiscardResponseBody(res)
		return nil, err
	}

	authCompleted, _, err := client.Update(negotiateToken)

	if err != nil {
		internal.DiscardResponseBody(res)
		return nil, err
	}

	if !authCompleted {
		internal.DiscardResponseBody(res)
		return nil, errors.New("client authentication not completed")
	}

	return res, nil
}

func NewNegotiateRoundTripper(base http.RoundTripper, options spnego_options.Options) http.RoundTripper {
	if base == nil {
		base = http.DefaultTransport
	}
	return &negotiateRoundTripper{
		r:       base,
		options: options,
	}
}
