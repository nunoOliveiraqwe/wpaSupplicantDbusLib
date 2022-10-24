package wpaSuppDBusLib

type innerAuthType string

const (
	InnerAuthPAP      innerAuthType = "PAP"
	InnerAuthMsChap   innerAuthType = "MSCHAP"
	InnerAuthMsChapV2 innerAuthType = "MSCHAPV2"
	InnerAuthChap     innerAuthType = "CHAP"
	InnerAuthMD5      innerAuthType = "MD5"
	InnerAuthGTC      innerAuthType = "GTC"
)
