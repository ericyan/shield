// Package stream implements stream cipher by wrapping a PRNG.
package stream

import "encoding/binary"

// A PRNG represents a pseudorandom number generator.
type PRNG interface {
	Uint64() uint64
}

// A Cipher represents a stream cipher with an instance of PRNG wrapped
// in. It implements the cipher.Stream interface.
type Cipher struct {
	keystream PRNG
}

// NewCipher returns a new Cipher which generates its keystream using
// the given PRNG.
func NewCipher(prng PRNG) *Cipher {
	return &Cipher{prng}
}

// XORKeyStream XORs each byte in the given slice with a byte from the
// cipher's keystream.
func (c *Cipher) XORKeyStream(dst, src []byte) {
	if len(src) == 0 {
		return
	}

	for i := 0; i < len(src)/8; i++ {
		v := binary.LittleEndian.Uint64(src[8*i:8*i+8]) ^ c.keystream.Uint64()
		binary.LittleEndian.PutUint64(dst[8*i:8*i+8], v)
	}

	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, c.keystream.Uint64())
	for i := (len(src) - len(src)%8); i < len(src); i++ {
		dst[i] = src[i] ^ buf[len(src)-i]
	}
}
