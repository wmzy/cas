package cas

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewV2(t *testing.T) {
	a := assert.New(t)

	cas2 := NewV2("http://cas.example.com")
	a.Equal("http://cas.example.com", cas2.BaseUrl, "base url")
	a.Equal(2, cas2.version, "version should be 2")
}
