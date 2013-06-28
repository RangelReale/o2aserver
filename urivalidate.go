package o2aserver

import (
	"net/url"
	"strings"
	"log"
)

func ValidateUri(baseUri string, redirectUri string) bool {
	if baseUri == "" || redirectUri == "" {
		log.Print("urls cannot be blank.\n")
		return false
	}

	base, err := url.Parse(baseUri)
	if err != nil {
		log.Printf("Error: %s\n", err)
		return false
	}

	redirect, err := url.Parse(redirectUri)
	if err != nil {
		log.Printf("Error: %s\n", err)
		return false
	}

	// must not have fragment
	if base.Fragment != "" || redirect.Fragment != "" {
		log.Print("Error: url must not include fragment.\n")
		return false
	}

	if base.Scheme == redirect.Scheme && base.Host == redirect.Host && len(redirect.Path) >= len(base.Path) && strings.HasPrefix(redirect.Path, base.Path) {
		return true
	}

	log.Printf("urls don't validate: %s / %s\n", baseUri, redirectUri)
	return false
}
