package account

import "regexp"

var (
	UsernameRegex *regexp.Regexp = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]{2,49}$`)
)
