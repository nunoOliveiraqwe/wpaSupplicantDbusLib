package wpaSuppDBusLib

import (
	"errors"
	"fmt"
)

type md5EapMethod struct {
	username string `json:"username"`
	password string `json:"password"`
}

func (m *md5EapMethod) GetEAPName() string {
	return "MD5"
}

func (m *md5EapMethod) ToConfigString() string {
	return fmt.Sprintf("  username=\"%s\"\n  password=\"%s\"\n", m.username, m.password)
}

type MD5EAPBuilder struct {
	username string `json:"username"`
	password string `json:"password"`
}

func NewMd5EApBuilder() MD5EAPBuilder {
	return MD5EAPBuilder{}
}

func (b *MD5EAPBuilder) WithUsername(username string) *MD5EAPBuilder {
	b.username = username
	return b
}

func (b *MD5EAPBuilder) WithPassword(password string) *MD5EAPBuilder {
	b.password = password
	return b
}

func (b *MD5EAPBuilder) Build() (eapMethod, error) {
	err := b.validate()
	if err != nil {
		return nil, err
	}
	md5 := md5EapMethod{
		username: b.username,
		password: b.password,
	}
	return &md5, nil
}

func (b *MD5EAPBuilder) validate() error {
	if b.username == "" {
		return errors.New("invalid username")
	}
	if b.password == "" {
		return errors.New("invalid password")
	}
	return nil
}
