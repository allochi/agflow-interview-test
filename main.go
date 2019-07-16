package main

import (
	"fmt"
	"./opts"
)

func main() {
	rs := opts.NewRuleSet()
	rs.AddDep("a", "b")
	rs.AddDep("b", "c")
	rs.AddDep("c", "a")
	rs.AddDep("d", "e")
	rs.AddConflict("c", "e")

	selected := opts.New(rs)
	fmt.Printf("%v\n", selected.StringSlice())

	selected.Toggle("a")
	fmt.Printf("%v\n", selected.StringSlice())

	rs.AddDep("f", "f")
	selected.Toggle("f")
	fmt.Printf("%v\n", selected.StringSlice())

	selected.Toggle("e")
	fmt.Printf("%v\n", selected.StringSlice())
}
