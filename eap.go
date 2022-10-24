package wpaSuppDBusLib

type eapMethod interface {
	ToConfigString() string
	GetEAPName() string
}

type eapBuilder interface {
	Build() (eapMethod, error)
}
