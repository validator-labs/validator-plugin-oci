package strings

// DeDupeStrSlice deduplicates a slices of strings
func DeDupeStrSlice(ss []string) []string {
	found := make(map[string]bool)
	l := []string{}
	for _, s := range ss {
		if _, ok := found[s]; !ok {
			found[s] = true
			l = append(l, s)
		}
	}
	return l
}
