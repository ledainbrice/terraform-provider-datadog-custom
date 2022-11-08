package sdk

import (
	"fmt"
	"testing"
)

func TestDiff(t *testing.T) {
	s1 := []string{"a", "b"}
	s2 := []string{"b", "c"}
	o1, o2 := diff(s1, s2)
	// if err != nil {
	// 	print(err)
	// 	t.Errorf("Func fail")
	// }
	fmt.Println(o1)
	fmt.Println(o2)
}
