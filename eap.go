package wpaSuppDBusLib

type eapMethod interface {
	ToConfigString() string
	GetEAPName() string
}

type eapBuilder interface {
	Build() (eapMethod, error)
}

func init() {

}

type tls struct {
}

type ttls struct {
}

type peap struct {
}
