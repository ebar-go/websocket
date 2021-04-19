package epoll

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreate(t *testing.T) {
	epoll, err := Create()
	assert.Nil(t, err)
	assert.NotEqual(t, 0, epoll.fd)
}

func TestEpollImpl_Add(t *testing.T) {

}
