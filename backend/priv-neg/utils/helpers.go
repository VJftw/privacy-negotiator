package utils

// IsIn - Returns true if a string exists in a given slice of strings, false otherwise.
func IsIn(needle string, haystack []string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}
