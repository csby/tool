package main

import (
	"fmt"
	"github.com/csby/security/aes"
	"github.com/csby/security/encoding"
)

type ConfigPassword struct {
	value *string
}

func (*ConfigPassword) CanSet() bool {
	return true
}

func (s *ConfigPassword) Get() interface{} {
	if s.value == nil {
		return ""
	}
	if len(*s.value) < 1 {
		return ""
	}

	base64 := &encoding.Base64{}
	data, err := base64.DecodeFromString(*s.value)
	if err != nil {
		return ""
	}
	aesEncoder := &aes.Aes{
		Key:       "Pwd#Crt@2019",
		Algorithm: "AES-128-CBC",
	}
	decData, err := aesEncoder.Decrypt(data)
	if err != nil {
		return ""
	}

	return string(decData)
}

func (s *ConfigPassword) Set(v interface{}) error {
	if s.value == nil {
		return fmt.Errorf("invalid value: nil")
	}
	value := fmt.Sprint(v)
	if len(value) < 1 {
		*s.value = ""
		return nil
	}
	aesEncoder := &aes.Aes{
		Key:       "Pwd#Crt@2019",
		Algorithm: "AES-128-CBC",
	}
	data, err := aesEncoder.Encrypt([]byte(value))
	if err != nil {
		*s.value = ""
		return nil
	}

	base64 := &encoding.Base64{}
	*s.value = base64.EncodeToString(data)
	return nil
}

func (*ConfigPassword) Zero() interface{} {
	return ""
}
