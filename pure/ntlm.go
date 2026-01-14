package pure

import (
	"net/http"

	"github.com/lublak/go-spnego/option"
)

func NewNtlmRoundTripper(base http.RoundTripper, options option.AuthOptions) http.RoundTripper {
	return newGenericNtlmRoundTripper(base, options, false)
}
