package option

type ApiType string

const (
	PURE ApiType = "pure"
	SSPI ApiType = "sspi"
)

type AuthOptions struct {
	User             *User
	KerberosFilePath string
}

func Default() *AuthOptions {
	return &AuthOptions{}
}
