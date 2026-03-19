package util

// IsBSUID returns true if the given string is a valid BSUID (business-scoped user ID).
func IsBSUID(v string) bool {
	if len(v) < 4 {
		return false
	}

	// check if the first two characters are between A-z
	if v[0] < 'A' || v[0] > 'z' || v[1] < 'A' || v[1] > 'z' {
		return false
	}

	// check if the third character is a dot
	if v[2] != '.' {
		return false
	}

	return true
}
