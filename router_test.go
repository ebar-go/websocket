package websocket

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRouter(t *testing.T) {
	router := NewRouter()
	user := router.Group("user")
	{
		user.Route("list", func(ctx Context) {

		})
		user.Route("create", nil)
		user.Route("update", nil)
		user.Route("delete", nil)
	}

	router.Route("index", nil)

	handler, exist := router.Get("/user/list")
	assert.True(t, exist)
	fmt.Println(handler)

	router.print()
}
