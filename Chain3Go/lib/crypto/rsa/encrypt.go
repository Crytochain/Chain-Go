package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"flag"
	"io/ioutil"
	"os"

	"Chain3Go/lib/log"
)

var decrypted string
var privateKey, publicKey []byte

func init() {
	var err error
	_, err = os.Stat("public.pem")
	//fmt.Println(err)
	if err != nil {
		var bits int
		flag.IntVar(&bits, "b", 2048, "key length, default 1024.")
		if err := GenRsaKey(bits); err != nil {
			log.Debugf("create pem failed.")
		}
		log.Debugf("create pem success")
	}

	// flag.StringVar(&decrypted, "d", "", "encrypt data")
	// flag.Parse()
	publicKey, err = ioutil.ReadFile("public.pem")

	if err != nil {
		os.Exit(-1)
	}
	privateKey, err = ioutil.ReadFile("private.pem")
	if err != nil {
		os.Exit(-1)
	}

}

func GetPubkey() (pk []byte) {
	//log.Debugf("Scs pubkey: %v", publicKey)
	return publicKey
}

// func main() {
// 	var data []byte
// 	var err error
// 	data, err = RsaEncrypt([]byte("fyxichen"))
// 	if err != nil {
// 		panic(err)
// 	}
// 	origData, err := RsaDecrypt(data)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(string(origData))
// }

// Encrypt
func RsaEncrypt(origData []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

// Decrypt
func RsaDecrypt(ciphertext []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}

func RsaEncryptWithKey(origData []byte, aPublicKeyBytes []byte) ([]byte, error){
	block, _ := pem.Decode(aPublicKeyBytes)

	if block == nil {
		return nil, errors.New("public key error " + string(aPublicKeyBytes) )
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

func RsaDecryptWithKey(ciphertext []byte, aPrivateKeyBytes []byte) ([]byte, error) {
	block, _ := pem.Decode(aPrivateKeyBytes)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}