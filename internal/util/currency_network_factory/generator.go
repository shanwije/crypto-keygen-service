package currency_network_factory

type KeyGenerator interface {
	GenerateKeyPair(userID int) (string, string, string, error)
}
