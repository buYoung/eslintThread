package module

import "strings"

type pathSorter []string

func (s pathSorter) Len() int {
	return len(s)
}

func (s pathSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s pathSorter) Less(i, j int) bool {
	s1 := strings.Replace(s[i], "/", "\x00", -1)
	s2 := strings.Replace(s[j], "/", "\x00", -1)
	return s1 < s2

}
