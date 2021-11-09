package main

import (
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"strings"
	"testing"
)

func TestGenerateSting(t *testing.T) {
	var dummy = mapset.NewSet("hello", "sun", "world", "space", "moon", "crypto", "sky", "ocean", "universe", "human")
	s := generateString()
	testArr := strings.Split(s, " ")
	for _, st := range testArr {
		if !dummy.Contains(st) {
			t.FailNow()
		}
	}
	fmt.Println(s)
}

func TestGreet(t *testing.T) {
	var st string
	st = greet()
	if st == "" {
		t.FailNow()
	}
	fmt.Println(st)
}
