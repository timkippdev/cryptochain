package wallet

import (
	"crypto/ecdsa"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewWallet(t *testing.T) {
	w := NewWallet()
	assert.Equal(t, startingWalletBalance, w.balance)
	assert.NotEmpty(t, w.publicKey)

	assert.Equal(t, startingWalletBalance, w.GetBalance())
}

func Test_CreateTransaction_AmountExceedsBalance(t *testing.T) {
	_, err := NewWallet().CreateTransaction(&ecdsa.PublicKey{}, 999999)
	assert.Error(t, err)
}

func Test_CreateTransaction_VerifyInputAndOutput(t *testing.T) {
	w := NewWallet()
	amount := float64(50)
	recipient := &ecdsa.PublicKey{}

	transaction, err := w.CreateTransaction(recipient, amount)
	assert.NoError(t, err)
	assert.Equal(t, w.publicKey, transaction.Input.Address)
	assert.Equal(t, w.publicKey, transaction.Input.Address)
	assert.Equal(t, amount, transaction.OutputMap[recipient])
}

func Test_VerifySignature_Valid(t *testing.T) {
	w := NewWallet()
	data := "test"
	assert.True(t, VerifySignature(w.publicKey, data, w.Sign(data)))
}

func Test_VerifySignature_Invalid(t *testing.T) {
	w := NewWallet()
	data := "test"
	assert.False(t, VerifySignature(w.publicKey, data, NewWallet().Sign(data)))
}
