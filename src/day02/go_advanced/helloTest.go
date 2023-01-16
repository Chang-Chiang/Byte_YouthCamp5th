package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func HelloTom() string {
	return "Jerry"
}

func TestHelloTom(t *testing.T) {
	output := HelloTom()
	expectOutput := "Tom"
	// if output != expectOutput {
	// 	t.Errorf("Expected %s do not match actual %s", expectOutput, output)
	// }
	assert.Equal(t, expectOutput, output)
}

func main() {
	TestHelloTom()
}
