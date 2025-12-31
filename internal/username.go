package internal

func JoinDomainAndName(domain string, name string) string {
	if len(domain) > 0 {
		name = domain + "\\" + name
	}

	return name
}
