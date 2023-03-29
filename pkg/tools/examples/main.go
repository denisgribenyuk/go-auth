package main

import (
	"fmt"

	"gitlab.assistagro.com/back/back.auth.go/pkg/tools"
)

func main() {
	a := struct {
		FieldStr1   string
		FieldInt1   int64
		FieldStr2   string
		StructField struct {
			Subfield1 string
			Subfield2 string
		}
	}{
		FieldStr1: "   Hello   ",
		FieldInt1: 123,
		FieldStr2: " World",
		StructField: struct {
			Subfield1 string
			Subfield2 string
		}{
			Subfield1: "Sub1   ",
			Subfield2: "   Sub2",
		},
	}
	fmt.Printf("Before: %q\n", a)
	tools.TrimSpaces(&a)
	fmt.Printf("After: %q\n", a)

	b := "   abc  "
	fmt.Printf("Before: %q\n", b)
	tools.TrimSpaces(&b)
	fmt.Printf("After: %q\n", b)

	c := []string{" test1 ", " test2 ", "    test3"}
	fmt.Printf("Before: %q\n", c)
	tools.TrimSpaces(&c)
	fmt.Printf("After: %q\n", c)
}
