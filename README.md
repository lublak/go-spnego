# go-spnego

Go package that provides a roundtripper spnego implementations with fallback to either basic auth if enabled or ntlm.

It wraps following implementations:

| Lib                                  | Protocols     | Api                    |
| ------------------------------------ | ------------- | ---------------------- |
| https://github.com/alexbrainman/sspi | Kerberos/NTLM | SSPI (Only on Windows) |
| https://github.com/jcmturner/gokrb5  | Kerberos      | Pure                   |
| https://github.com/Azure/go-ntlmssp  | NTLM          | Pure                   |
| Currently not supported              | Kerberos/NTLM | GSSAPI                 |

This library currently not support GSSAPI and I have no interest in integrating this function directly into my site.
However, if there were a pure Golang library available, I would be happy to use it and integrate it here.
In general, though, “Pure” should be completely sufficient here.

## The Apis:

All apis can be configurated with an user. Otherwise, use the configurations specific to the API.

### Auto

Just uses automatically uses sspi on windows and pure on non windows targets.

### SSPI (Only on Windows)

Does not require any platform-specific configurations.
The Kerberos options are ignored here.
When used on targets other than Windows, the function returns nil.
If user is not configured, the currently logged-in user is used.

### Pure

If user is not configured, the Kerberos keytab file is used.
A keytab file can be created with ktpass on windows, ktutil on linux and osx.
If a fallback to NTLM is required (no Kerberos), a user configuration is required.
To use this user only as a fallback, there is the UserOnlyForFallback setting.

There are two methods for configuring this function.
Option 1 the commonly known environment variables:

| Env         | Description                                                                           | Default                                                                            |
| ----------- | ------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- |
| KRB5_CONFIG | the path of the krb5 configuration                                                    | linux: /etc/krb5.conf; mac: /opt/local/etc/krb5.conf; windows: c:\windows\krb5.ini |
| KRB5CCNAME  | the path or folder of the ccache (allowed: DIR:/path/to/folder or FILE:/path/to/file) | /tmp/krb5cc_ + Uid                                                                 |

Option 2 the options:

When the options are set, it overrides the reading of the environment variables.
