// fivesmallest implements solution to "five smallest" problem.
package main

import "fmt"

type set map[int]struct{}

// Smallest returns the n smallest ints in the set.
func Smallest(data set, n int) []int {
	r := make([]int, n)
	for i := 0; i < n; i++ {
		r[i] = i
	}
	return r
}

func main() {
	d := set{
		3: struct{}{},
		1: struct{}{},
		5: struct{}{},
		7: struct{}{},
	}
	s3 := Smallest(d, 3)
	fmt.Printf("Smallest 3 of %v: %v\n", d, s3)
}
