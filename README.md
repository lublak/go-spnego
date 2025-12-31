# go-spnego

Go package that provides a roundtripper spnego implementations with fallback to either basic auth if enabled or ntlm.

It wraps following implementations:

| Lib                                  | Protocols     | Api  |
| ------------------------------------ | ------------- | ---- |
| https://github.com/alexbrainman/sspi | Kerberos/NTLM | SSPI |
| https://github.com/jcmturner/gokrb5  | Kerberos      | Pure |
| https://github.com/Azure/go-ntlmssp  | NTLM          | Pure |

Currently not ready to use, completly untested and missing complete own spnego implementation for pure implementation (currently only supports keberos here).