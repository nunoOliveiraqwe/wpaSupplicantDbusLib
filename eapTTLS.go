package wpaSuppDBusLib

import (
	"errors"
	"fmt"
	"strings"
)

var allowedTTLSInnerAuthTypes = []innerAuthType{InnerAuthPAP, InnerAuthMsChap, InnerAuthMsChapV2, InnerAuthChap, InnerAuthMD5, InnerAuthGTC}

type ttlsMethod struct {
	anonymousIdentity string        `json:"anonymousIdentity,omitempty"`
	identity          string        `json:"identity"`
	caCertPath        string        `json:"ca_cert,omitempty"`
	password          string        `json:"password,omitempty"`
	innerAuth         innerAuthType `json:"phase2"`
}

func (t *ttlsMethod) GetEAPName() string {
	return "TTLS"
}

func (t *ttlsMethod) ToConfigString() string {
	builder := strings.Builder{}
	if t.anonymousIdentity != "" {
		builder.WriteString(fmt.Sprintf("   anonymous_identity=\"%s\"\n", t.anonymousIdentity))
	}
	if t.identity != "" {
		builder.WriteString(fmt.Sprintf("  identity=\"%s\"\n", t.identity))
	}
	if t.caCertPath != "" {
		builder.WriteString(fmt.Sprintf("  ca_cert=\"%s\"\n", t.caCertPath))
	}
	if t.password != "" {
		builder.WriteString(fmt.Sprintf("  password=\"%s\"\n", t.password))
	}
	if t.innerAuth != "" {
		builder.WriteString(fmt.Sprintf("  phase2=\"auth=%s\"\n", t.innerAuth))
	}
	return builder.String()
}

type TTLSBuilder struct {
	anonymousIdentity string        `json:"anonymousIdentity,omitempty"`
	identity          string        `json:"identity"`
	caCertPath        string        `json:"ca_cert,omitempty"`
	password          string        `json:"password"`
	innerAuth         innerAuthType `json:"phase2"`
}

func NewTTLSBuilder() TTLSBuilder {
	return TTLSBuilder{}
}

func (t *TTLSBuilder) WithAnonymousIdentity(anonIdentity string) *TTLSBuilder {
	t.anonymousIdentity = anonIdentity
	return t
}

func (t *TTLSBuilder) WithIdentity(identity string) *TTLSBuilder {
	t.identity = identity
	return t
}

func (t *TTLSBuilder) WithCaCertPath(caCertPath string) *TTLSBuilder {
	t.caCertPath = caCertPath
	return t
}

func (t *TTLSBuilder) WithPassword(password string) *TTLSBuilder {
	t.password = password
	return t
}

func (t *TTLSBuilder) WithInnerAuthType(innerAuthType innerAuthType) *TTLSBuilder {
	t.innerAuth = innerAuthType
	return t
}

func (t *TTLSBuilder) Build() (eapMethod, error) {
	err := t.validate()
	if err != nil {
		return nil, err
	}
	tls := ttlsMethod{
		anonymousIdentity: t.anonymousIdentity,
		identity:          t.identity,
		caCertPath:        t.caCertPath,
		password:          t.password,
		innerAuth:         t.innerAuth,
	}
	return &tls, nil
}

func (t *TTLSBuilder) validate() error {
	if t.identity == "" {
		return errors.New("invalid identity")
	}
	if t.password == "" {
		return errors.New("invalid password")
	}
	if t.innerAuth == "" {
		return errors.New("invalid inner auth (empty)")
	}
	if !contains(allowedTTLSInnerAuthTypes, t.innerAuth) {
		return errors.New("invalid inner auth (wrong value)")
	}
	return nil
}
