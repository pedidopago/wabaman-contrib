package util

import "strings"

// IsBSUID reports whether v looks like a business-scoped user ID (BSUID) or a
// parent BSUID: a two-letter ISO country code, a period, then the identifier
// (e.g. "US.13491208655302741918" or "US.ENT.11815799212886844830"). Phone
// numbers (all digits) and emails are rejected.
func IsBSUID(v string) bool {
	if len(v) < 4 {
		return false
	}

	// The first two characters must be ASCII letters (ISO country code).
	if !isASCIILetter(v[0]) || !isASCIILetter(v[1]) {
		return false
	}

	// The third character is the separating dot.
	if v[2] != '.' {
		return false
	}

	// An '@' means it's an email, not a BSUID (e.g. "jo.doe@example.com" would
	// otherwise pass the country-code + dot shape check).
	if strings.ContainsRune(v, '@') {
		return false
	}

	return true
}

func isASCIILetter(b byte) bool {
	return (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z')
}
