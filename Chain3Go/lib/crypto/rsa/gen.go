package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"os"
)

// func main() {
// 	var bits int
// 	flag.IntVar(&bits, "b", 2048, "key length, default 1024.")
// 	if err := GenRsaKey(bits); err != nil {
// 		log.Fatal("create pem failed.")
// 	}
// 	log.Println("create pem success")
// }

func GenRsaKey(bits int) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "privatekey",
		Bytes: derStream,
	}
	file, err := os.Create("private.pem")
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}

	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	block = &pem.Block{
		Type:  "publickey",
		Bytes: derPkix,
	}
	file, err = os.Create("public.pem")
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	return nil
}


func GenExtRsaKey(bits int) ([]byte, []byte, error) {
	extPrivateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return []byte(""), []byte(""), err
	}
	extPrivateKeyBytes := x509.MarshalPKCS1PrivateKey(extPrivateKey)
	extPublicKey := &extPrivateKey.PublicKey
	extPublicKeyBytes, err := x509.MarshalPKIXPublicKey(extPublicKey)
	return extPrivateKeyBytes, extPublicKeyBytes, err
}