package currency_factory

type KeyGenerator interface {
	GenerateKeyPair() (string, string, string, error)
}
