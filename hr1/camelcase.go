package main

import (
	"fmt"
)

const (
	// A a const to avoid hardcoded chars
	A = 'A'
	a = 'a'

	// Z z const to avoid hardcoded chars
	Z = 'Z'
	z = 'z'
)

// #1
// Complete the caesarCipher function below.
func caesarCipher(r rune, k int) rune {
	if r >= A && r <= Z {
		return rotate(r, A, k)
	}
	if r >= a && r <= z {
		return rotate(r, a, k)
	}
	return r
}

func rotate(r rune, base, k int) rune {
	tmp := int(r) - base
	tmp = (tmp + k) % 26
	return rune(tmp + base)
}

// #2 count the camelcase of a string
// Complete the camelcase function below.
//func camelcase(s string) int32 {
//count := 1
//for _, letter := range s {
//if letter <= 90 && letter >= 65 {
//count++
//}
//}
//return count
//}

func main() {
	var n, k int
	var s string

	fmt.Scanf("%d\n", &n)
	fmt.Scanf("%s\n", &s)
	fmt.Scanf("%d\n", &k)

	var ret []rune
	for _, ch := range s {
		ret = append(ret, caesarCipher(ch, k))
	}
	fmt.Printf("%s\n", string(ret))
}
