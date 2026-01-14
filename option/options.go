package option

type ApiType string

const (
	PURE ApiType = "pure"
	SSPI ApiType = "sspi"
	AUTO ApiType = "auto"
)

type AuthOptions struct {
	AllowBasicAuth      bool
	User                *User
	UserOnlyForFallback bool
	Kerberos            *Kerberos
}

func Default() *AuthOptions {
	return &AuthOptions{}
}
