package main

import "testing"

func TestMain(t *testing.T) {
	a, b := 1, 2
	s := Add(a, b)
	if s != a+b {
		t.Errorf("Add(%d, %d) = %d", a, b, s)
	}
}
