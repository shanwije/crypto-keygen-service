package currency_network_factory

type KeyGenerator interface {
	GenerateKeyPair() (string, string, string, error)
}
