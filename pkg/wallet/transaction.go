package wallet

import (
	"crypto/ecdsa"
	"errors"

	"github.com/segmentio/ksuid"
	"github.com/timkippdev/cryptochain/pkg/timestamp"
)

type Transaction struct {
	ID        string                       `json:"id"`
	Input     TransactionInput             `json:"input"`
	OutputMap map[*ecdsa.PublicKey]float64 `json:"-"`
}

type TransactionInput struct {
	Address   *ecdsa.PublicKey    `json:"-"`
	Amount    float64             `json:"amount"`
	Signature string              `json:"signature"`
	Timestamp timestamp.Timestamp `json:"timestamp"`
}

func NewTransactionInput(senderWallet *Wallet, outputMap map[*ecdsa.PublicKey]float64) TransactionInput {
	return TransactionInput{
		Address:   senderWallet.publicKey,
		Amount:    senderWallet.balance,
		Signature: senderWallet.Sign(outputMap),
		Timestamp: timestamp.Now(),
	}
}

func NewTransaction(senderWallet *Wallet, recipient *ecdsa.PublicKey, amount float64) *Transaction {
	outputMap := map[*ecdsa.PublicKey]float64{
		recipient:              amount,
		senderWallet.publicKey: senderWallet.balance - amount,
	}
	input := NewTransactionInput(senderWallet, outputMap)
	return &Transaction{
		ID:        ksuid.New().String(),
		Input:     input,
		OutputMap: outputMap,
	}
}

func (t *Transaction) Update(senderWallet *Wallet, recipient *ecdsa.PublicKey, amount float64) error {
	if amount > t.OutputMap[senderWallet.publicKey] {
		return errors.New("amount exceeds balance")
	}

	if _, found := t.OutputMap[recipient]; !found {
		t.OutputMap[recipient] = amount
	} else {
		t.OutputMap[recipient] = t.OutputMap[recipient] + amount
	}

	t.OutputMap[senderWallet.publicKey] = t.OutputMap[senderWallet.publicKey] - amount
	t.Input = NewTransactionInput(senderWallet, t.OutputMap)
	return nil
}

func ValidateTransaction(transaction *Transaction) bool {
	outputAmountTotal := float64(0)
	for _, v := range transaction.OutputMap {
		outputAmountTotal += v
	}

	if transaction.Input.Amount != outputAmountTotal {
		return false
	}

	if !VerifySignature(transaction.Input.Address, transaction.OutputMap, transaction.Input.Signature) {
		return false
	}

	return true
}
