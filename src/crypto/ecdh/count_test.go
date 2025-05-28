package ecdh

import (
	"crypto/internal/fips140/edwards25519/field"
	"crypto/rand"
	"encoding/base64"
	"testing"
)

func generatePrivateKey() ([32]byte, error) {
	var key [32]byte
	rand.Read(key[:])
	// Modify random bytes using algorithm described at:
	// https://cr.yp.to/ecdh.html.
	key[0] &= 248
	key[31] &= 127
	key[31] |= 64

	return key, nil
}

func curve25519_scalarBaseMult(dst, scalar *[32]byte) {
	curve := X25519()
	priv, err := curve.NewPrivateKey(scalar[:])
	if err != nil {
		panic("curve25519: internal error: scalarBaseMult was not 32 bytes")
	}
	copy(dst[:], priv.PublicKey().Bytes())
}

func TestCountX25519(t *testing.T) {
	private, _ := generatePrivateKey()
	var public [32]byte

	// run once to init internals if any
	curve25519_scalarBaseMult(&public, &private)

	before := field.Multiplications

	curve25519_scalarBaseMult(&public, &private)

	after := field.Multiplications

	t.Logf("Private Key: %s", base64.StdEncoding.EncodeToString(private[:]))
	t.Logf("Public Key:  %s", base64.StdEncoding.EncodeToString(public[:]))
	t.Logf("Multiplications: %d", after-before)
}
