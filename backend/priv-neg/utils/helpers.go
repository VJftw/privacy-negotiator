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

// IsSubset - Checks if a is a subset of b
func IsSubset(a []string, b []string) bool {
	for _, vA := range a {
		inB := false
		for _, vB := range b {
			if vA == vB {
				inB = true
				break
			}
		}

		if !inB {
			return false
		}
	}

	return true
}

// ArrayUnion - Returns the union of 2 given arrays
func ArrayUnion(a []string, b []string) []string {
	c := []string{}

	for _, aV := range a {
		for _, bV := range b {
			if aV == bV && !IsIn(aV, c) {
				c = append(c, aV)
			}
		}
	}
	return c
}
