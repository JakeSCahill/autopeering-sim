package peer

import (
	"crypto/ed25519"
	"crypto/sha256"
	"math/rand"
	"testing"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/blake2s"
)

func benchmarkIDHash(b *testing.B, hash func([]byte) [32]byte) {
	data := make([][]byte, b.N)
	for n := 0; n < b.N; n++ {
		data[n] = make([]byte, ed25519.PublicKeySize)
		rand.Read(data[n])
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = hash(data[n])
	}
}

func BenchmarkSHA256ID(b *testing.B) {
	benchmarkIDHash(b, sha256.Sum256)
}

func BenchmarkBLAKE2b256ID(b *testing.B) {
	benchmarkIDHash(b, blake2b.Sum256)
}

func BenchmarkBLAKE2s256ID(b *testing.B) {
	benchmarkIDHash(b, blake2s.Sum256)
}
