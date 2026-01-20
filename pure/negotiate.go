package pure

import (
	"fmt"
	"net/http"
	"os"

	"github.com/jcmturner/gokrb5/v8/client"
	"github.com/jcmturner/gokrb5/v8/config"
	"github.com/jcmturner/gokrb5/v8/credentials"
	"github.com/jcmturner/gokrb5/v8/spnego"
	"github.com/lublak/go-spnego/internal"
	spnego_options "github.com/lublak/go-spnego/options"
)

type negotiateRoundTripper struct {
	r       http.RoundTripper
	ntlm    http.RoundTripper
	options spnego_options.Options
}

func (t *negotiateRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var config *config.Config
	var err error

	if t.options.Kerberos != nil && len(t.options.Kerberos.ConfigFilePath) > 0 {
		config, err = kerberosConfigFromPath(t.options.Kerberos.ConfigFilePath)
	} else {
		config, err = defaultKerberosConfig()
		if os.IsNotExist(err) {
			// fallback to ntlm if kerberos config not exists

			res, roundTripErr := t.ntlm.RoundTrip(req)

			if roundTripErr != nil {
				return nil, fmt.Errorf("kerberos config missing: %v, fallback ntlm error: %v", err, roundTripErr)
			}

			return res, nil
		}
	}

	if err != nil {
		return nil, err
	}

	var c *client.Client

	if t.options.User == nil || !t.options.UserOnlyForFallback {
		var ccache *credentials.CCache

		if t.options.Kerberos != nil && len(t.options.Kerberos.CCName) > 0 {
			ccache, err = kerberosCCacheFromName(t.options.Kerberos.ConfigFilePath)
		} else {
			ccache, err = defaultKerberosCCache()

		}

		if err != nil {
			return nil, err
		}

		c, err = client.NewFromCCache(ccache, config, client.DisablePAFXFAST(true))

		if err != nil {
			return nil, err
		}
	} else {
		c = client.NewWithPassword(t.options.User.Name, t.options.User.Domain, t.options.User.Password, config, client.DisablePAFXFAST(true))
	}

	spn := internal.GetSpnFromRequest(req)

	s := spnego.SPNEGOClient(c, spn)
	err = s.AcquireCred()
	if err != nil {
		return nil, fmt.Errorf("could not acquire client credential: %v", err)
	}

	st, err := s.InitSecContext()
	if err != nil {
		// fallback to ntlm if init doesn't work

		res, roundTripErr := t.ntlm.RoundTrip(req)

		if roundTripErr != nil {
			return nil, fmt.Errorf("could not initialize context: %v, fallback ntlm error: %v", err, roundTripErr)
		}

		return res, nil
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

func NewNegotiateRoundTripper(base http.RoundTripper, options spnego_options.Options) http.RoundTripper {
	if base == nil {
		base = http.DefaultTransport
	}
	return &negotiateRoundTripper{
		r:       base,
		ntlm:    newGenericNtlmRoundTripper(base, options, true),
		options: options,
	}
}
