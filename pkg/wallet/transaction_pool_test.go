package wallet

import (
	"crypto/ecdsa"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AddTransaction(t *testing.T) {
	w := NewWallet()
	transaction := NewTransaction(w, &ecdsa.PublicKey{}, float64(50))
	transactionPool := NewTransactionPool()

	transactionPool.AddTransaction(transaction)

	assert.Len(t, transactionPool.transactionMap, 1)
	assert.Equal(t, transaction, transactionPool.transactionMap[transaction.ID])
}
