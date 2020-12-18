//go:generate sh generate.sh

//Package tld has the same API as net/url except
//tld.URL contains extra fields: Subdomain, Domain, TLD and Port.
package tld

import (
	"net/url"
	"strings"
	"github.com/weppos/publicsuffix-go/publicsuffix"
)

//URL embeds net/url and adds extra fields ontop
type URL struct {
	Subdomain, Domain, TLD, Port string
	ICANN                        bool
	*url.URL
}

//Parse mirrors net/url.Parse except instead it returns
//a tld.URL, which contains extra fields.
func Parse(s string) (*URL, error) {
	url, err := url.Parse(s)
	if err != nil {
		return nil, err
	}
	if url.Host == "" {
		return &URL{URL: url}, nil
	}
	_, port := domainPort(url.Host)

	pubDomain, err1 := publicsuffix.DomainFromListWithOptions(publicsuffix.DefaultList, url.Host, &publicsuffix.FindOptions{IgnorePrivate: true})
	pubParse, _ := publicsuffix.ParseFromListWithOptions(publicsuffix.DefaultList, url.Host, &publicsuffix.FindOptions{IgnorePrivate: true})

	// error parser rules
	sub := ""
	domName := ""
	tld := ""
	icann := false

	if err1 == nil {
		sub = url.Host
		domName = pubDomain
		tld = pubParse.TLD
		icann = true
	}

	return &URL{
		Subdomain: sub,
		Domain:    domName,
		TLD:       tld,
		Port:      port,
		ICANN:     icann,
		URL:       url,
	}, nil

}

func domainPort(host string) (string, string) {
	for i := len(host) - 1; i >= 0; i-- {
		if host[i] == ':' {
			return host[:i], host[i+1:]
		} else if host[i] < '0' || host[i] > '9' {
			return host, ""
		}
	}
	//will only land here if the string is all digits,
	//net/url should prevent that from happening
	return host, ""
}
