package sdk

import "os"

func diff(s1 []string, s2 []string) ([]string, []string) {
	o1 := itemIn(s1, s2)
	o2 := itemIn(s2, s1)

	return o1, o2
}

func itemIn(s1 []string, s2 []string) []string {
	o := make([]string, 0)
	m := toMap(s1)

	for _, x := range s2 {
		if _, found := m[x]; !found {
			o = append(o, x)
		}
	}
	return o
}

func toMap(s []string) map[string]struct{} {
	mb := make(map[string]struct{}, len(s))

	for _, x := range s {
		mb[x] = struct{}{}
	}

	return mb
}

func writeIn(path string, content string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(content + "\n"); err != nil {
		return err
	}
	return nil
}
