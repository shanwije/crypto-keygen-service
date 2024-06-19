package keys

type KeyGenerator interface {
	GenerateKeyPair() (string, string, string, error)
}
