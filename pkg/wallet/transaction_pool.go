package wallet

type TransactionPool struct {
	transactionMap map[string]*Transaction
}

func NewTransactionPool() *TransactionPool {
	return &TransactionPool{
		transactionMap: make(map[string]*Transaction),
	}
}

func (tp *TransactionPool) AddTransaction(transaction *Transaction) {
	tp.transactionMap[transaction.ID] = transaction
}
