package wpaSuppDBusLib

import (
	"errors"
	"fmt"
	"strings"
)

type PEAPVersion int

const (
	PEAPVersion0 PEAPVersion = 0
	PEAPVersion1 PEAPVersion = 1
)

var peapVersionSlice = []PEAPVersion{PEAPVersion0, PEAPVersion1}
var allowedInnerAuthTypes = []innerAuthType{InnerAuthMsChapV2, InnerAuthMD5, InnerAuthGTC}

type peapMethod struct {
	anonymousIdentity string        `json:"anonymousIdentity,omitempty"`
	identity          string        `json:"identity"`
	password          string        `json:"password"`
	peapVersion       PEAPVersion   `json:"peaplabel,omitempty"`
	caCertPath        string        `json:"ca_cert,omitempty"`
	innerAuth         innerAuthType `json:"phase2"`
}

func (p *peapMethod) GetEAPName() string {
	return "PEAP"
}

func (p *peapMethod) ToConfigString() string {
	builder := strings.Builder{}
	if p.anonymousIdentity != "" {
		builder.WriteString(fmt.Sprintf("  anonymous_identity=\"%s\"\n", p.anonymousIdentity))
	}
	if p.identity != "" {
		builder.WriteString(fmt.Sprintf("  identity=\"%s\"\n", p.identity))
	}
	if p.password != "" {
		builder.WriteString(fmt.Sprintf("  password=\"%s\"\n", p.password))
	}
	if p.peapVersion != -1 {
		builder.WriteString(fmt.Sprintf("  phase1=\"peaplabel=%d\"", p.peapVersion))
	}
	if p.caCertPath != "" {
		builder.WriteString(fmt.Sprintf("  ca_cert=\"%s\"\n", p.caCertPath))
	}
	if p.innerAuth != "" {
		builder.WriteString(fmt.Sprintf("  phase2=\"auth=%s\"\n", p.innerAuth))
	}
	return builder.String()
}

type PEAPBuilder struct {
	anonymousIdentity string        `json:"anonymousIdentity,omitempty"`
	identity          string        `json:"identity"`
	password          string        `json:"password"`
	peapVersion       PEAPVersion   `json:"peaplabel,omitempty"`
	caCertPath        string        `json:"ca_cert,omitempty"`
	innerAuth         innerAuthType `json:"phase2"`
}

func NewPEAPBuilder() PEAPBuilder {
	return PEAPBuilder{
		peapVersion: -1,
	}
}

func (b *PEAPBuilder) WithAnonymousIdentity(anonIdentity string) *PEAPBuilder {
	b.anonymousIdentity = anonIdentity
	return b
}

func (b *PEAPBuilder) WithIdentity(identity string) *PEAPBuilder {
	b.identity = identity
	return b
}

func (b *PEAPBuilder) WithPassword(password string) *PEAPBuilder {
	b.password = password
	return b
}

func (b *PEAPBuilder) WithPEAPVersion(peapVersion PEAPVersion) *PEAPBuilder {
	b.peapVersion = peapVersion
	return b
}

func (b *PEAPBuilder) WithCaCertPath(caCertPath string) *PEAPBuilder {
	b.caCertPath = caCertPath
	return b
}

func (b *PEAPBuilder) WithInnerAuthType(innerAuthType innerAuthType) *PEAPBuilder {
	b.innerAuth = innerAuthType
	return b
}

func (b *PEAPBuilder) Build() (eapMethod, error) {
	err := b.validate()
	if err != nil {
		return nil, err
	}
	eap := peapMethod{
		anonymousIdentity: b.anonymousIdentity,
		identity:          b.identity,
		password:          b.password,
		peapVersion:       b.peapVersion,
		caCertPath:        b.caCertPath,
		innerAuth:         b.innerAuth,
	}
	return &eap, nil
}

func (b *PEAPBuilder) validate() error {
	if b.identity == "" {
		return errors.New("invalid identity")
	}
	if b.password == "" {
		return errors.New("invalid password")
	}
	if b.innerAuth == "" {
		return errors.New("invalid inner auth (empty)")
	}
	if b.peapVersion != -1 && !contains(peapVersionSlice, b.peapVersion) {
		return errors.New("invalid peap version value")
	}
	if !contains(allowedInnerAuthTypes, b.innerAuth) {
		return errors.New("invalid inner auth (wrong value)")
	}
	return nil
}
