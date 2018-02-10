package libs

import (
	"hash"
)

func EVPKDF(password, salt []byte, iter, keySize int, h func() hash.Hash) []byte {
	var derivedKey []byte
	var block []byte
	hasher := h()

	for len(derivedKey) < keySize*4 {
		if block != nil {
			hasher.Write(block)
		}
		hasher.Write(password)
		hasher.Write(salt)
		block = hasher.Sum(nil)

		hasher.Reset()
		for i := 1; i < iter; i++ {
			hasher.Write(block)
			block = hasher.Sum(block)
			hasher.Reset()
		}

		derivedKey = append(derivedKey, block...)
	}

	return derivedKey
}
