package options

type Api string

const (
	PURE Api = "pure"
	SSPI Api = "sspi"
)

type Options struct {
	AllowBasicAuth      bool
	User                *User
	UserOnlyForFallback bool
	Kerberos            *Kerberos
}

func Default() *Options {
	return &Options{}
}
