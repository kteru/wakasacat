package main

import (
	"os"
	"testing"
)

func Test_vfprintf(t *testing.T) {
	var pair = []struct {
		in1 bool
		in2 string
		out int
	}{
		{true, "test output\n", 12},
		{false, "test output\n", 0},
	}

	for _, p := range pair {
		option.verbose = p.in1
		result, _ := vfprintf(os.Stderr, p.in2)
		if result != p.out {
			t.Errorf("vfprintf() = %v, err | want %v", result, p.out)
		}
	}

}
