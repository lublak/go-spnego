package pure

import (
	"fmt"
	"net/http"
	"os"
	"os/user"
	"strings"

	"github.com/jcmturner/gokrb5/v8/client"
	"github.com/jcmturner/gokrb5/v8/credentials"
	"github.com/jcmturner/gokrb5/v8/spnego"
	"github.com/lublak/go-spnego/internal"
)

type negotiateRoundTripper struct {
	r http.RoundTripper
}

func (t *negotiateRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	config, err := kerberosConfig()
	if err != nil {
		return nil, err
	}

	var ccpath string

	ccname := os.Getenv("KRB5CCNAME")
	if strings.HasPrefix(ccname, "FILE:") {
		ccpath = ccname[len("FILE:"):]
	} else {
		u, err := user.Current()
		if err != nil {
			return nil, err
		}

		ccpath = "/tmp/krb5cc_" + u.Uid
	}

	ccache, err := credentials.LoadCCache(ccpath)
	if err != nil {
		return nil, err
	}

	client, err := client.NewFromCCache(ccache, config, client.DisablePAFXFAST(true))

	if err != nil {
		return nil, err
	}

	spn := internal.GetSpnFromRequest(req)

	s := spnego.SPNEGOClient(client, spn)
	err = s.AcquireCred()
	if err != nil {
		return nil, fmt.Errorf("could not acquire client credential: %v", err)
	}
	st, err := s.InitSecContext()
	if err != nil {
		return nil, fmt.Errorf("could not initialize context: %v", err)
	}
	clientToken, err := st.Marshal()
	if err != nil {
		return nil, fmt.Errorf("could not marshal SPNEGO. %v", err)
	}

	req.Header.Set("Authorization", internal.EncodeNegotiateToken(clientToken))

	res, err := t.r.RoundTrip(req)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return res, nil
	}

	_, err = internal.FindAndDecodeNegotiateToken(res.Header.Values("WWW-Authenticate"))
	if err != nil {
		internal.DiscardResponseBody(res)
		return nil, err
	}

	return res, nil
}

func NewNegotiateRoundTripper(base http.RoundTripper) http.RoundTripper {
	if base == nil {
		base = http.DefaultTransport
	}
	return &negotiateRoundTripper{
		r: base,
	}
}
