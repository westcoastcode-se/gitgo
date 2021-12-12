package utils

import "unicode"

func testSpace(r rune) bool {
	return unicode.IsSpace(r) || r == '\u0005'
}

// SplitOnSpace the supplied string into two parts. The first part contains the command to
// be executed. The second part contains the arguments sent to the command. It will search for the first space and use
// that as the delimiter.
func SplitOnSpace(s string) []string {
	l := len(s)
	i := 0
	for i = 0; i < l; i++ {
		r := rune(s[i])
		if testSpace(r) {
			break
		}
	}
	if i == l {
		return []string{s}
	}
	result := []string{
		s[0:i],
		s[i+1:],
	}
	return result
}
