//go:build unit

package pocket

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCalAdd(t *testing.T) {
	a := 0.1
	b := 0.2

	result := Add(a, b)

	assert.Equal(t, result, 0.3)

}

func TestCalSub(t *testing.T) {
	a := 0.1
	b := 0.2

	result := Sub(b, a)

	assert.Equal(t, result, 0.1)
}
