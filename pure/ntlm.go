package pure

import (
	"net/http"

	"github.com/lublak/go-spnego/options"
)

func NewNtlmRoundTripper(base http.RoundTripper, spnegoOptions options.Options) http.RoundTripper {
	return newGenericNtlmRoundTripper(base, spnegoOptions, false)
}
