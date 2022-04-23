package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

func main() {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	fmt.Printf("%+v\n\n", privateKey)

	hash := sha256.Sum256([]byte("test"))

	publicKey := privateKey.PublicKey
	fmt.Printf("%+v\n", privateKey)

	fmt.Printf("%s\n", string(elliptic.Marshal(secp256k1.S256(), publicKey.X, publicKey.Y)))
	fmt.Printf("%s\n", hex.EncodeToString(elliptic.Marshal(secp256k1.S256(), publicKey.X, publicKey.Y)))

	sig, _ := ecdsa.SignASN1(rand.Reader, privateKey, hash[:])
	// fmt.Printf("%x\n", sig)

	valid := ecdsa.VerifyASN1(&privateKey.PublicKey, hash[:], sig)
	fmt.Println("valid:", valid)
}
