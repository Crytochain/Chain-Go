// Copyright 2014 The MOAC-core Authors
// This file is part of the MOAC-core library.
//
// The MOAC-core library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The MOAC-core library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the MOAC-core library. If not, see <http://www.gnu.org/licenses/>.

package transaction

import (
	"container/heap"
	"errors"
	"fmt"
	"io"
	"math/big"
	"sync/atomic"

	"Chain3Go/lib/common"
	"Chain3Go/lib/common/hexutil"

	"Chain3Go/lib/rlp"
)

//go:generate gencodec -type txdata -field-override txdataMarshaling -out gen_tx_json.go

var (
	ErrInvalidSig = errors.New("invalid transaction v, r, s values")
	errNoSigner   = errors.New("missing signing methods")
)

// deriveSigner makes a *best* guess about which signer to use.
// For MOAC pangu release, just use one
// may change in the future.
// REturn a signer with chainID info
func deriveSigner(V *big.Int) Signer {
	// if V.Sign() != 0 && isProtectedV(V) {
	// 	return NewEIP155Signer(deriveChainId(V))
	// } else {
	return NewPanguSigner(deriveChainId(V))
	// }
}

type Transaction struct {
	TxData txdata

	// caches
	hash    atomic.Value
	size    atomic.Value //size of the txdata
	from    atomic.Value
	execblk big.Int
}

/*
 * 2018/03/16 keep ShardingFlag only
 */
type txdata struct {
	AccountNonce   uint64          `json:"nonce"    gencodec:"required"`
	SystemContract uint64          `json:"syscnt" gencodec:"required"`
	Price          *big.Int        `json:"gasPrice" gencodec:"required"`
	GasLimit       *big.Int        `json:"gas"      gencodec:"required"`
	Recipient      *common.Address `json:"to"       rlp:"nil"` // nil means contract creation
	Amount         *big.Int        `json:"value"    gencodec:"required"`
	Payload        []byte          `json:"input"    gencodec:"required"`
	// QueryFlag      uint64          `json:"queryFlag" gencodec:"required"`
	ShardingFlag uint64 `json:"shardingFlag" gencodec:"required"`
	//SubChainProtocol *common.Address `json:"SubChainProtocol"`
	// ShardMinMember   uint64         `json:"ShardMinMember" gencodec:"required"`
	// ScsConsensusAddr common.Address `json:"ScsConsensusAddr" gencodec:"required"`

	// Signature values
	V *big.Int `json:"v" gencodec:"required"`
	R *big.Int `json:"r" gencodec:"required"`
	S *big.Int `json:"s" gencodec:"required"`

	// This is only used when marshaling to JSON.
	Hash *common.Hash `json:"hash" rlp:"-"`
}

func (tdata *txdata) GetShardingFlag() uint64 { return tdata.ShardingFlag }
func (tdata *txdata) GetSystemFlag() uint64   { return tdata.SystemContract }
func (tdata *txdata) GetAmount() *big.Int     { return tdata.Amount }

func (tdata *txdata) SetShardingFlag(inflag uint64) {
	tdata.ShardingFlag = inflag
}

/*
 * Add the ControlFlag and ScsConsensusAddr
 * seems not used
 */
type txdataMarshaling struct {
	AccountNonce hexutil.Uint64
	ControlFlag  hexutil.Uint64
	Price        *hexutil.Big
	GasLimit     *hexutil.Big
	Amount       *hexutil.Big
	Payload      hexutil.Bytes
	// ScsConsensusAddr common.Address
	// QueryFlag    hexutil.Uint64
	// ShardingFlag hexutil.Uint64
	V *hexutil.Big
	R *hexutil.Big
	S *hexutil.Big
}

func NewTransaction(nonce uint64, to common.Address, amount, gasLimit, gasPrice *big.Int, shardingFlag uint64, data []byte) *Transaction {
	return newTransaction(nonce, &to, amount, gasLimit, gasPrice, shardingFlag, data)
}

func NewContractCreation(nonce uint64, amount, gasLimit, gasPrice *big.Int, shardingFlag uint64, data []byte) *Transaction {
	return newTransaction(nonce, nil, amount, gasLimit, gasPrice, shardingFlag, data)
}

func newTransaction(nonce uint64, to *common.Address, amount, gasLimit, gasPrice *big.Int, shardingFlag uint64, data []byte) *Transaction {
	if len(data) > 0 {
		data = common.CopyBytes(data)
	}

	d := txdata{
		AccountNonce:   nonce,
		Recipient:      to,
		Payload:        data,
		Amount:         new(big.Int),
		GasLimit:       new(big.Int),
		Price:          new(big.Int),
		V:              new(big.Int),
		R:              new(big.Int),
		S:              new(big.Int),
		SystemContract: 0,
		ShardingFlag:   shardingFlag,
	}

	//Set new shardingFlag
	if shardingFlag > 0 {
		d.SetShardingFlag(shardingFlag)
	}

	if amount != nil {
		d.Amount.Set(amount)
	}
	if gasLimit != nil {
		d.GasLimit.Set(gasLimit)
	}
	if gasPrice != nil {
		d.Price.Set(gasPrice)
	}

	return &Transaction{TxData: d}
}

// ChainId returns which chain id this transaction was signed for (if at all)
func (tx *Transaction) ChainId() *big.Int {
	return deriveChainId(tx.TxData.V)
}

// Protected returns whether the transaction is protected from replay protection.
func (tx *Transaction) Protected() bool {
	return isProtectedV(tx.TxData.V)
}

func isProtectedV(V *big.Int) bool {
	if V.BitLen() <= 8 {
		v := V.Uint64()
		return v != 27 && v != 28
	}
	// anything not 27 or 28 are considered unprotected
	return true
}

// DecodeRLP implements rlp.Encoder
func (tx *Transaction) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, &tx.TxData)
}

// DecodeRLP implements rlp.Decoder
// Add the check on the data size

func (tx *Transaction) DecodeRLP(s *rlp.Stream) error {
	// fmt.Printf("In transaction.go:DecodeRLP: 190")

	_, size, _ := s.Kind()
	// fmt.Printf("stream size: %v\n", size)
	err := s.Decode(&tx.TxData)
	if err == nil {
		tx.size.Store(common.StorageSize(rlp.ListSize(size)))
	}

	return err
}

func (tx *Transaction) Data() []byte             { return common.CopyBytes(tx.TxData.Payload) }
func (tx *Transaction) Gas() *big.Int            { return new(big.Int).Set(tx.TxData.GasLimit) }
func (tx *Transaction) GasPrice() *big.Int       { return new(big.Int).Set(tx.TxData.Price) }
func (tx *Transaction) Value() *big.Int          { return new(big.Int).Set(tx.TxData.Amount) }
func (tx *Transaction) Nonce() uint64            { return tx.TxData.AccountNonce }
func (tx *Transaction) CheckNonce() bool         { return true }
func (tx *Transaction) UpdateNonce(nonce uint64) { tx.TxData.AccountNonce = nonce }

//functions to pass flag values
func (tx *Transaction) GetShardingFlag() uint64 { return tx.TxData.GetShardingFlag() }
func (tx *Transaction) GetSystemFlag() uint64   { return tx.TxData.GetSystemFlag() }
func (tx *Transaction) ShardingFlag() uint64    { return tx.TxData.GetShardingFlag() }
func (tx *Transaction) SystemFlag() uint64      { return tx.TxData.GetSystemFlag() }

//Process the controlFlag in tx.TxData
func (tx *Transaction) SetShardingFlag(inflag uint64) { tx.TxData.ShardingFlag = inflag }
func (tx *Transaction) SetSystemFlag(inflag uint64)   { tx.TxData.SystemContract = inflag }

// To returns the recipient address of the transaction.
// It returns nil if the transaction is a contract creation.
func (tx *Transaction) To() *common.Address {
	if tx.TxData.Recipient == nil {
		return nil
	} else {
		to := *tx.TxData.Recipient
		return &to
	}
}

// To returns the sender address of the transaction.
//func (tx *Transaction) From() *common.Address {
//	return tx.from.Load().(*common.Address)
//}

// Hash hashes the RLP encoding of tx.
// It uniquely identifies the transaction.
func (tx *Transaction) Hash() common.Hash {
	if hash := tx.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	v := rlpHash(tx)
	tx.hash.Store(v)
	return v
}

// SigHash returns the hash to be signed by the sender.
// It does not uniquely identify the transaction.
func (tx *Transaction) SigHash(signer Signer) common.Hash {
	return signer.Hash(tx)
}

func (tx *Transaction) Size() common.StorageSize {
	if size := tx.size.Load(); size != nil {
		return size.(common.StorageSize)
	}
	c := writeCounter(0)
	rlp.Encode(&c, &tx.TxData)
	tx.size.Store(common.StorageSize(c))
	return common.StorageSize(c)
}

// AsMessage returns the transaction as a core.Message.
//
// AsMessage requires a signer to derive the sender.
//
// XXX Rename message to something less arbitrary?
func (tx *Transaction) AsMessage(s Signer) (Message, error) {
	var (
		// syncFlag  bool
		autoFlush bool
		blk       *big.Int
	)

	// if tx.TxData.GetQueryFlag() > 0 {
	// 	syncFlag = false
	// 	autoFlush = true
	// 	blk = big.NewInt(int64(tx.TxData.GetQueryFlag()))
	// } else {
	// syncFlag = true
	autoFlush = true
	blk = big.NewInt(0)
	// }

	msg := Message{
		nonce:      tx.TxData.AccountNonce,
		price:      new(big.Int).Set(tx.TxData.Price),
		gasLimit:   new(big.Int).Set(tx.TxData.GasLimit),
		to:         tx.TxData.Recipient,
		amount:     tx.TxData.Amount,
		data:       tx.TxData.Payload,
		checkNonce: true,
		system:     tx.TxData.SystemContract,
		// syncFlag:        syncFlag,
		autoFlush:       autoFlush,
		waitBlockNumber: blk,
		shardFlag:       tx.TxData.ShardingFlag,
	}

	var err error
	msg.from, err = Sender(s, tx)
	return msg, err
}

// WithSignature returns a new transaction with the given signature.
// This signature needs to be formatted as described in the yellow paper (v+27).
// updated to GETH 1.7 version
func (tx *Transaction) WithSignature(signer Signer, sig []byte) (*Transaction, error) {
	// return signer.WithSignature(tx, sig)
	r, s, v, err := signer.SignatureValues(tx, sig)
	if err != nil {
		return nil, err
	}
	cpy := &Transaction{TxData: tx.TxData}
	cpy.TxData.R, cpy.TxData.S, cpy.TxData.V = r, s, v
	return cpy, nil
}

// Cost returns amount + gasprice * gaslimit.
func (tx *Transaction) Cost() *big.Int {
	total := new(big.Int).Mul(tx.TxData.Price, tx.TxData.GasLimit)
	// fmt.Printf("Gas cost %v, limit %v\n", tx.TxData.Price, tx.TxData.GasLimit)
	// fmt.Printf("Gas cost %v, Send value %v\n", total, tx.TxData.Amount)

	total.Add(total, tx.TxData.Amount)
	return total
}

func (tx *Transaction) RawSignatureValues() (*big.Int, *big.Int, *big.Int) {
	return tx.TxData.V, tx.TxData.R, tx.TxData.S
}

/*
 * Return the sender from the sig as string
 */
func (tx *Transaction) GetSender() string {
	var from string
	if tx.TxData.V != nil {
		// make a best guess about the signer and use that to derive
		// the sender.
		signer := deriveSigner(tx.TxData.V)
		if f, err := Sender(signer, tx); err != nil { // derive but don't cache
			from = "[invalid sender: invalid sig]"
		} else {
			from = fmt.Sprintf("%x", f[:])
		}
	} else {
		from = "[invalid sender: nil V field]"
	}
	return from
}

/*
 * Return the sender as address
 */
func (tx *Transaction) From() common.Address {
	if tx.TxData.V != nil {
		// make a best guess about the signer and use that to derive
		// the sender.
		signer := deriveSigner(tx.TxData.V)
		if from, err := Sender(signer, tx); err == nil {
			return from
		}
	}
	return common.Address{}
}

/*
 * Return the transaction info as string
 */
func (tx *Transaction) String() string {
	var from, to string
	if tx.TxData.V != nil {
		// make a best guess about the signer and use that to derive
		// the sender.
		signer := deriveSigner(tx.TxData.V)
		if f, err := Sender(signer, tx); err != nil { // derive but don't cache
			from = "[invalid sender: invalid sig]"
		} else {
			from = fmt.Sprintf("%x", f[:])
		}
	} else {
		from = "[invalid sender: nil V field]"
	}

	if tx.TxData.Recipient == nil {
		to = "[contract creation]"
	} else {
		to = fmt.Sprintf("%x", tx.TxData.Recipient[:])
		// fmt.Printf("Transaction:string %x\n", tx.TxData.Recipient[:])
	}
	enc, _ := rlp.EncodeToBytes(&tx.TxData)
	return fmt.Sprintf(`
	TX(%x)
	Contract: %v
	From:     %s
	To:       %s
	Nonce:    %v
	GasPrice: %#x
	GasLimit  %#x
	Value:    %#x
	Data:     0x%x
	V:        %#x
	R:        %#x
	S:        %#x
	Hex:      %x
	SysCnt:	  %v
	ShardingFlag: %v
`,
		tx.Hash(),
		tx.TxData.Recipient == nil,
		from,
		to,
		tx.TxData.AccountNonce,
		tx.TxData.Price,
		tx.TxData.GasLimit,
		tx.TxData.Amount,
		tx.TxData.Payload,
		tx.TxData.V,
		tx.TxData.R,
		tx.TxData.S,
		enc,
		tx.TxData.GetSystemFlag(),
		tx.TxData.GetShardingFlag(),
	)
}

// Transaction slice type for basic sorting.
type Transactions []*Transaction

// Len returns the length of s
func (s Transactions) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s
func (s Transactions) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// GetRlp implements Rlpable and returns the i'th element of s in rlp
func (s Transactions) GetRlp(i int) []byte {
	enc, _ := rlp.EncodeToBytes(s[i])
	return enc
}

// Returns a new set t which is the difference between a to b
func TxDifference(a, b Transactions) (keep Transactions) {
	keep = make(Transactions, 0, len(a))

	remove := make(map[common.Hash]struct{})
	for _, tx := range b {
		remove[tx.Hash()] = struct{}{}
	}

	for _, tx := range a {
		if _, ok := remove[tx.Hash()]; !ok {
			keep = append(keep, tx)
		}
	}

	return keep
}

// TxByNonce implements the sort interface to allow sorting a list of transactions
// by their nonces. This is usually only useful for sorting transactions from a
// single account, otherwise a nonce comparison doesn't make much sense.
type TxByNonce Transactions

func (s TxByNonce) Len() int           { return len(s) }
func (s TxByNonce) Less(i, j int) bool { return s[i].TxData.AccountNonce < s[j].TxData.AccountNonce }
func (s TxByNonce) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// TxByPrice implements both the sort and the heap interface, making it useful
// for all at once sorting as well as individually adding and removing elements.
type TxByPrice Transactions

func (s TxByPrice) Len() int           { return len(s) }
func (s TxByPrice) Less(i, j int) bool { return s[i].TxData.Price.Cmp(s[j].TxData.Price) > 0 }
func (s TxByPrice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (s *TxByPrice) Push(x interface{}) {
	*s = append(*s, x.(*Transaction))
}

func (s *TxByPrice) Pop() interface{} {
	old := *s
	n := len(old)
	x := old[n-1]
	*s = old[0 : n-1]
	return x
}

// TransactionsByPriceAndNonce represents a set of transactions that can return
// transactions in a profit-maximising sorted order, while supporting removing
// entire batches of transactions for non-executable accounts.
type TransactionsByPriceAndNonce struct {
	txs    map[common.Address]Transactions // Per account nonce-sorted list of transactions
	heads  TxByPrice                       // Next transaction for each unique account (price heap)
	signer Signer                          // Signer for the set of transactions
}

// NewTransactionsByPriceAndNonce creates a transaction set that can retrieve
// price sorted transactions in a nonce-honouring way.
//
// Note, the input map is reowned so the caller should not interact any more with
// if after providing it to the constructor.
func NewTransactionsByPriceAndNonce(signer Signer, txs map[common.Address]Transactions) *TransactionsByPriceAndNonce {
	// Initialize a price based heap with the head transactions
	heads := make(TxByPrice, 0, len(txs))
	for _, accTxs := range txs {
		heads = append(heads, accTxs[0])
		// Ensure the sender address is from the signer
		acc, _ := Sender(signer, accTxs[0])
		txs[acc] = accTxs[1:]
	}
	heap.Init(&heads)

	// Assemble and return the transaction set
	return &TransactionsByPriceAndNonce{
		txs:    txs,
		heads:  heads,
		signer: signer,
	}
}

// Peek returns the next transaction by price.
func (t *TransactionsByPriceAndNonce) Peek() *Transaction {
	if len(t.heads) == 0 {
		return nil
	}
	return t.heads[0]
}

// Shift replaces the current best head with the next one from the same account.
func (t *TransactionsByPriceAndNonce) Shift() {
	acc, _ := Sender(t.signer, t.heads[0])
	if txs, ok := t.txs[acc]; ok && len(txs) > 0 {
		t.heads[0], t.txs[acc] = txs[0], txs[1:]
		heap.Fix(&t.heads, 0)
	} else {
		heap.Pop(&t.heads)
	}
}

// Pop removes the best transaction, *not* replacing it with the next one from
// the same account. This should be used when a transaction cannot be executed
// and hence all subsequent ones should be discarded from the same account.
func (t *TransactionsByPriceAndNonce) Pop() {
	heap.Pop(&t.heads)
}

// Message is a fully derived transaction and implements core.Message
//
// NOTE: In a future PR this will be removed.
type Message struct {
	to                      *common.Address
	from                    common.Address
	nonce                   uint64
	amount, price, gasLimit *big.Int
	data                    []byte
	checkNonce              bool
	system                  uint64
	//	syncFlag                bool
	autoFlush        bool
	waitBlockNumber  *big.Int
	scsConsensusAddr common.Address
	shardFlag        uint64
}

func NewMessage(from common.Address, to *common.Address, nonce uint64, amount, gasLimit, price *big.Int,
	data []byte, checkNonce bool, syncFlag bool, autoFlush bool, waitBlockNumber *big.Int,
	scsConsensusAddr common.Address, shardingFlag uint64) Message {
	return Message{
		from:       from,
		to:         to,
		nonce:      nonce,
		amount:     amount,
		price:      price,
		gasLimit:   gasLimit,
		data:       data,
		checkNonce: checkNonce,
		system:     0,
		// syncFlag:         syncFlag,
		autoFlush:        autoFlush,
		waitBlockNumber:  waitBlockNumber,
		scsConsensusAddr: scsConsensusAddr,
		shardFlag:        shardingFlag,
	}
}

func (m Message) From() common.Address { return m.from }
func (m Message) To() *common.Address  { return m.to }
func (m Message) GasPrice() *big.Int   { return m.price }
func (m Message) Value() *big.Int      { return m.amount }
func (m Message) Gas() *big.Int        { return m.gasLimit }
func (m Message) Nonce() uint64        { return m.nonce }
func (m Message) Data() []byte         { return m.data }
func (m Message) CheckNonce() bool     { return m.checkNonce }
func (m Message) GetSystem() uint64    { return m.system }

// func (m Message) SyncFlag() bool            { return m.syncFlag }
func (m Message) AutoFlush() bool           { return m.autoFlush }
func (m Message) WaitBlockNumber() *big.Int { return m.waitBlockNumber }

func (m Message) ScsConsensusAddr() common.Address { return m.scsConsensusAddr }
func (m Message) ShardFlag() uint64                { return m.shardFlag }
func (m *Message) SetShardingFlag(sharding uint64) {
	if sharding > 0 {
		// m.syncFlag = false
		m.autoFlush = false
		m.waitBlockNumber = big.NewInt(0)
	} else {
		// m.syncFlag = true
		m.autoFlush = false
		m.waitBlockNumber = big.NewInt(0)
	}
}

// func (m *Message) SetFlag(query uint64, sharding uint64) {
// 	if query > 0 {
// 		m.syncFlag = false
// 		m.autoFlush = true
// 		m.waitBlockNumber = big.NewInt(int64(query))
// 	} else if sharding > 0 {
// 		m.syncFlag = false
// 		m.autoFlush = false
// 		m.waitBlockNumber = big.NewInt(0)
// 	} else {
// 		m.syncFlag = true
// 		m.autoFlush = false
// 		m.waitBlockNumber = big.NewInt(0)
// 	}
// }
