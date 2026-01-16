package utils

func MergeSlices(a []string, b []string) []string {
	return append(a, b...)
}

func RemoveElement(s []string, element string) []string {
	var result []string
	for _, str := range s {
		if str != element {
			result = append(result, str)
		}
	}
	return result
}

func RemoveDuplicates(s []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, str := range s {
		if !seen[str] {
			seen[str] = true
			result = append(result, str)
		}
	}
	return result
}
