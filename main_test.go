package main

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	unittest = true
	for _, test := range []struct {
		Args   []string
		Output string
	}{
		{
			Args:   []string{"./reloaded", "test/inHex.txt", "test/temp.txt"},
			Output: "Hex test 30",
		},
		{
			Args:   []string{"./reloaded", "test/inBin.txt", "test/temp.txt"},
			Output: "Bin test 5",
		},
		{
			Args:   []string{"./reloaded", "test/inUp.txt", "test/temp.txt"},
			Output: "Up TEST",
		},
		{
			Args:   []string{"./reloaded", "test/inLow.txt", "test/temp.txt"},
			Output: "Low test",
		},
		{
			Args:   []string{"./reloaded", "test/inCap.txt", "test/temp.txt"},
			Output: "Cap Test",
		},
		{
			Args:   []string{"./reloaded", "test/inAn.txt", "test/temp.txt"},
			Output: "An hour test",
		},
		{
			Args:   []string{"./reloaded", "test/inAll.txt", "test/temp.txt"},
			Output: "The.!,? 'punctuation ' test'",
		},
	} {
		t.Run("", func(t *testing.T) {
			os.Args = test.Args
			out = bytes.NewBuffer(nil)
			main()

			if actual := out.(*bytes.Buffer).String(); actual != test.Output {
				fmt.Println(actual, test.Output)
				t.Errorf("expected>>%s<<, but got>>%s<<", test.Output, actual)
			}
		})
	}
}
