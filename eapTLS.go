package wpaSuppDBusLib

import (
	"errors"
	"fmt"
	"strings"
)

type tlsMethod struct {
	identity           string `json:"identity"`
	caCertPath         string `json:"ca_cert,omitempty"`
	clientCert         string `json:"client_cert"`
	privateKey         string `json:"private_key"`
	privateKeyPassword string `json:"private_key_passwd,omitempty"`
}

func (t *tlsMethod) GetEAPName() string {
	return "TLS"
}

func (t *tlsMethod) ToConfigString() string {
	builder := strings.Builder{}
	if t.identity != "" {
		builder.WriteString(fmt.Sprintf("  identity=\"%s\"\n", t.identity))
	}
	if t.caCertPath != "" {
		builder.WriteString(fmt.Sprintf("  ca_cert=\"%s\"\n", t.caCertPath))
	}
	if t.clientCert != "" {
		builder.WriteString(fmt.Sprintf("  client_cert=\"%s\"\n", t.clientCert))
	}
	if t.privateKey != "" {
		builder.WriteString(fmt.Sprintf("  private_key=\"%s\"\n", t.privateKey))
	}
	if t.privateKeyPassword != "" {
		builder.WriteString(fmt.Sprintf("  private_key_passwd=\"%s\"\n", t.privateKeyPassword))
	}
	return builder.String()
}

type TLSBuilder struct {
	identity           string `json:"identity"`
	caCertPath         string `json:"ca_cert,omitempty"`
	clientCert         string `json:"client_cert"`
	privateKey         string `json:"private_key"`
	privateKeyPassword string `json:"private_key_passwd,omitempty"`
}

func NewTLSBuilder() TLSBuilder {
	return TLSBuilder{}
}

func (t *TLSBuilder) WithIdentity(identity string) *TLSBuilder {
	t.identity = identity
	return t
}

func (t *TLSBuilder) WithPrivateKeyPassword(password string) *TLSBuilder {
	t.privateKeyPassword = password
	return t
}

func (t *TLSBuilder) WithCaCertPath(caCertPath string) *TLSBuilder {
	t.caCertPath = caCertPath
	return t
}

func (t *TLSBuilder) WithClientCertPath(clientCertPath string) *TLSBuilder {
	t.clientCert = clientCertPath
	return t
}

func (t *TLSBuilder) WithPrivateKeyPath(privateKeyPath string) *TLSBuilder {
	t.privateKey = privateKeyPath
	return t
}

func (t *TLSBuilder) Build() (eapMethod, error) {
	err := t.validate()
	if err != nil {
		return nil, err
	}
	tls := tlsMethod{
		identity:           t.identity,
		caCertPath:         t.caCertPath,
		clientCert:         t.clientCert,
		privateKey:         t.privateKey,
		privateKeyPassword: t.privateKeyPassword,
	}
	return &tls, nil
}

func (t *TLSBuilder) validate() error {
	if t.identity == "" {
		errors.New("invalid value for identity")
	}
	if t.clientCert == "" {
		errors.New("invalid value for client cert")
	}
	if t.privateKey == "" {
		errors.New("invalid value for private key")
	}
	return nil
}
