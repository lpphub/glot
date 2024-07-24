package utils

import (
	"fmt"
	"testing"
)

func TestGenerateToken(t *testing.T) {
	key := "12345"
	j1, _ := GenerateToken("aaa", key)
	fmt.Println(j1)

	r1, _ := ParseToken(j1, key)
	fmt.Println(r1)
}
