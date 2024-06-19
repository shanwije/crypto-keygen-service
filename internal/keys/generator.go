package keys

type KeyGenerator interface {
	GenerateKeyPair(userID int) error
}
