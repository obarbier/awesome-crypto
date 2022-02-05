package domain

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

var _ Key = (*RSA)(nil)

func NewRsaKey(keySize int) (*RSA, error) {
	// generate key
	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return nil, fmt.Errorf("Cannot generate RSA key\n")
	}
	publicKey := privateKey.Public()

	resp := &RSA{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		KetType:    "rsa",
	}

	return resp, nil
}

type RSA struct {
	KetType    string
	PrivateKey *rsa.PrivateKey
	PublicKey  crypto.PublicKey
}

func (R *RSA) Type() string {
	return R.KetType
}

func (R *RSA) UnMarshal() ([]byte, error) {
	return rsaUnMarshal(R.PrivateKey)
}

func (R *RSA) Marshal(privatePEM string) error {
	key, err := rsaMarshal(privatePEM)
	if err != nil {
		return err
	}
	*R = RSA{
		PrivateKey: key,
		PublicKey:  key.Public(),
		KetType:    "rsa",
	}
	return nil
}

func rsaUnMarshal(privateKey *rsa.PrivateKey) ([]byte, error) {
	privateBytes := x509.MarshalPKCS1PrivateKey(privateKey)

	privateBytesPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privateBytes,
		},
	)
	return privateBytesPem, nil
}

func rsaMarshal(privatePEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privatePEM))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil

}
