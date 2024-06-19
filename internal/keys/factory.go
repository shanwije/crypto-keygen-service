package keys

import "errors"

func GetKeyGenerator(network string) (KeyGenerator, error) {
    switch network {
    default:
        return nil, errors.ErrUnsupported
    }
}
