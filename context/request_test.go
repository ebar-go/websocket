package context

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRequest(t *testing.T) {
	source := `{"uri":"/index","body":{"name":"hongker","extends":{"age":27}}}`
	req, err := NewRequest([]byte(source))
	assert.Nil(t, err)
	fmt.Println(req)

}
