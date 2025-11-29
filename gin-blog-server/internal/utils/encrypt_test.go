package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncrypt(t *testing.T) {
	password := "123456"

	// Encrypt
	hashPassword, err := BcryptHash(password)
	assert.Nil(t, err)

	// Verify
	result := BcryptCheck(password, hashPassword)
	assert.True(t, result)
}

func TestMD5(t *testing.T) {
	assert.Equal(t, "e10adc3949ba59abbe56e057f20f883e", MD5("123456"))
}
