package wallet

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/timkippdev/cryptochain/pkg/blockchain"
)

func GenerateSignature(privateKey *ecdsa.PrivateKey, data interface{}) string {
	hash := blockchain.Hash(fmt.Sprintf("%v", data))
	signature, _ := ecdsa.SignASN1(rand.Reader, privateKey, []byte(hash))
	return hex.EncodeToString(signature)
}

func VerifySignature(publicKey *ecdsa.PublicKey, data interface{}, signature string) bool {
	decodedSignature, _ := hex.DecodeString(signature)
	hash := blockchain.Hash(fmt.Sprintf("%v", data))
	return ecdsa.VerifyASN1(publicKey, []byte(hash), decodedSignature)
}
