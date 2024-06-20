package network_factory

type KeyPairAndAddress struct {
	PublicKey  string
	PrivateKey string
	Address    string
}

type KeyGenerator interface {
	GenerateKeyPairAndAddress(userID int) (KeyPairAndAddress, error)
}
