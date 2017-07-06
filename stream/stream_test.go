package stream

import "testing"

type fakePRNG struct{}

func (*fakePRNG) Uint64() uint64 {
	return uint64(0XDEADBEEF)
}

func TestStreamCipher(t *testing.T) {
	plaintext := []byte("hello, world")

	ciphertext := make([]byte, len(plaintext))
	c1 := NewCipher(new(fakePRNG))
	c1.XORKeyStream(ciphertext, plaintext)

	plaintext2 := make([]byte, len(plaintext))
	c2 := NewCipher(new(fakePRNG))
	c2.XORKeyStream(plaintext2, ciphertext)

	for i, v := range plaintext2 {
		if v != plaintext[i] {
			t.Fatalf("mismatch at byte %d:\nhave: % x\nwant: % x", i, plaintext2, plaintext)
		}
	}
}
