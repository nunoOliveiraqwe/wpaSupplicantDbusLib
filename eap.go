package wpaSuppDBusLib

type eapProvider interface {
}

func init() {

}

type md5 struct {
	username string
	password string
}

type tls struct {
}

type ttls struct {
}

type peap struct {
}