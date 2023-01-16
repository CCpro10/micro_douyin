package util

import (
	"github.com/CCpro10/micro_douyin/conf"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	ConfPath = "./conf/conf.yaml"
)

func TestJWT(t *testing.T) {
	conf.Init(ConfPath)
	InitJWTVal()
	userId := int64(1 << 50)
	for i := 0; i < 10000; i++ {

		token, err := GenerateToken(userId)
		assert.Equal(t, nil, err)

		c, err := parseToken(token)
		assert.Equal(t, nil, err)
		assert.Equal(t, userId, c.UserId)

		id, e := ParseToken(token)
		assert.Equal(t, userId, id)
		assert.Equal(t, nil, e)
	}
}
