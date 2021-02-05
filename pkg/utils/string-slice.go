package utils

/*
Index returns the first index of the target string t,
or -1 if no match is found.
*/
func Index(vs []string, t string) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

/*
Include returns true if the target string t is in the slice.
*/
func Include(vs []string, t string) bool {
	return Index(vs, t) >= 0
}

/*
Removes return a new slice without any occurence of the
target string t
*/
func Removes(vs []string, t string) []string {
	var newSlice []string
	for _, v := range vs {
		if v != t {
			newSlice = append(newSlice, v)
		}
	}
	return newSlice
}

/*
StringSliceUnion returns the union of the two string slices.
Remember that union dont have repeated elements
*/
func StringSliceUnion(a, b []string) []string {
	check := make(map[string]int)
	d := append(a, b...)
	res := make([]string, 0)
	for _, val := range d {
		check[val] = 1
	}

	for letter := range check {
		res = append(res, letter)
	}
	return res
}
