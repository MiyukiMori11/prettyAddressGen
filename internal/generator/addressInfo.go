package generator

type addressInfo struct {
	address    string
	publicKey  string
	privateKey string
}

func (a *addressInfo) Address() string {
	return a.address
}

func (a *addressInfo) PublicKey() string {
	return a.publicKey
}

func (a *addressInfo) PrivateKey() string {
	return a.privateKey
}
