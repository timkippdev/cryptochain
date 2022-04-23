package wallet

import (
	"crypto/ecdsa"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewTransaction(t *testing.T) {
	senderWallet := NewWallet()
	senderWallet.balance = 1000
	recipient := &ecdsa.PublicKey{}
	amount := float64(50)

	transaction := NewTransaction(senderWallet, recipient, amount)
	assert.NotEmpty(t, transaction.ID)
	assert.Len(t, transaction.OutputMap, 2)
	assert.Equal(t, amount, transaction.OutputMap[recipient])
	assert.Equal(t, senderWallet.balance-amount, transaction.OutputMap[senderWallet.publicKey])
	assert.Greater(t, int64(transaction.Input.Timestamp), int64(0))
	assert.Equal(t, senderWallet.publicKey, transaction.Input.Address)
	assert.Equal(t, senderWallet.balance, transaction.Input.Amount)
	assert.True(t, VerifySignature(senderWallet.publicKey, transaction.OutputMap, transaction.Input.Signature))
}

func Test_Update_InvalidAmount(t *testing.T) {
	senderWallet := NewWallet()
	recipient := &ecdsa.PublicKey{}
	amount := float64(50)
	transaction := NewTransaction(senderWallet, recipient, amount)

	newRecipient := &ecdsa.PublicKey{}
	newAmount := float64(99999)

	err := transaction.Update(senderWallet, newRecipient, newAmount)
	assert.Error(t, err)
}

func Test_Update_NewAddress(t *testing.T) {
	senderWallet := NewWallet()
	recipient := &ecdsa.PublicKey{}
	amount := float64(50)
	transaction := NewTransaction(senderWallet, recipient, amount)

	senderOutputAmountBeforeUpdate := transaction.OutputMap[senderWallet.publicKey]
	signatureBeforeUpdate := transaction.Input.Signature

	newRecipient := &ecdsa.PublicKey{}
	newAmount := float64(50)

	err := transaction.Update(senderWallet, newRecipient, newAmount)
	assert.NoError(t, err)

	assert.Equal(t, newAmount, transaction.OutputMap[newRecipient])
	assert.Equal(t, senderOutputAmountBeforeUpdate-newAmount, transaction.OutputMap[senderWallet.publicKey])

	totalOutputAmount := float64(0)
	for _, v := range transaction.OutputMap {
		totalOutputAmount += v
	}
	assert.Equal(t, totalOutputAmount, transaction.Input.Amount)
	assert.NotEqual(t, signatureBeforeUpdate, transaction.Input.Signature)
}

func Test_Update_ExistingAddress(t *testing.T) {
	senderWallet := NewWallet()
	recipient := &ecdsa.PublicKey{}
	amount := float64(50)
	transaction := NewTransaction(senderWallet, recipient, amount)

	senderOutputAmountBeforeUpdate := transaction.OutputMap[senderWallet.publicKey]

	newAmount := float64(50)

	err := transaction.Update(senderWallet, recipient, newAmount)
	assert.NoError(t, err)

	assert.Equal(t, amount+newAmount, transaction.OutputMap[recipient])
	assert.Equal(t, senderOutputAmountBeforeUpdate-newAmount, transaction.OutputMap[senderWallet.publicKey])
}

func Test_ValidateTransaction_Valid(t *testing.T) {
	senderWallet := NewWallet()
	senderWallet.balance = 1000
	recipient := &ecdsa.PublicKey{}
	amount := float64(50)

	transaction := NewTransaction(senderWallet, recipient, amount)
	assert.True(t, ValidateTransaction(transaction))
}

func Test_ValidateTransaction_InvalidOutputMap(t *testing.T) {
	senderWallet := NewWallet()
	senderWallet.balance = 1000
	recipient := &ecdsa.PublicKey{}
	amount := float64(50)

	transaction := NewTransaction(senderWallet, recipient, amount)
	transaction.OutputMap[senderWallet.publicKey] = 999999
	assert.False(t, ValidateTransaction(transaction))
}

func Test_ValidateTransaction_InvalidInputSignature(t *testing.T) {
	senderWallet := NewWallet()
	senderWallet.balance = 1000
	recipient := &ecdsa.PublicKey{}
	amount := float64(50)

	transaction := NewTransaction(senderWallet, recipient, amount)
	transaction.Input.Signature = NewWallet().Sign("different-data")
	assert.False(t, ValidateTransaction(transaction))
}
