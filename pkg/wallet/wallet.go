package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
)

const (
	startingWalletBalance float64 = 1000
)

type Wallet struct {
	balance    float64
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

func NewWallet() *Wallet {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	return &Wallet{
		balance:    startingWalletBalance,
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
	}
}

func (w *Wallet) CreateTransaction(recipient *ecdsa.PublicKey, amount float64) (*Transaction, error) {
	if amount > w.balance {
		return nil, errors.New("amount exceeds balance")
	}
	return NewTransaction(w, recipient, amount), nil
}

func (w *Wallet) GetBalance() float64 {
	return w.balance
}

func (w *Wallet) Sign(data interface{}) string {
	return GenerateSignature(w.privateKey, data)
}
