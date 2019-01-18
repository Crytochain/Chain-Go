// Copyright 2016 The MOAC-core Authors
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
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"

	"Chain3Go/lib/common"
	"Chain3Go/lib/crypto"
	"Chain3Go/lib/log"
	//"Chain3Go/lib/params"
)

//Added ErrUnproctedTX to exclude
var (
	ErrInvalidChainId = errors.New("invalid chain id for signer")
	ErrUnproctedTX    = errors.New("unprotected transaction from signer")

	errAbstractSigner     = errors.New("abstract signer")
	abstractSignerAddress = common.HexToAddress("ffffffffffffffffffffffffffffffffffffffff")
)

// sigCache is used to cache the derived sender and contains
// the signer used to derive it.
type sigCache struct {
	signer Signer
	from   common.Address
}

// ChainConfig is the core config which determines the blockchain settings.
//
// ChainConfig is stored in the database on a per block basis. This means
// that any network, identified by its genesis block, can have its own
// set of configuration options.
// default is to omitempty Accounts
// EIP158
//
type ChainConfig struct {
	ChainId *big.Int `json:"chainId"` // Chain id identifies the current chain and is used for replay protection

	PanguBlock         *big.Int `json:"panguBlock,omitempty"`         // Pangu switch block (nil = no fork, 0 = already pangu)
	RemoveEmptyAccount bool     `json:"removeEmptyAccount,omitempty"` //Replace EIP158 check and should be set to true

	// DAOForkBlock   *big.Int `json:"daoForkBlock,omitempty"`   // TheDAO hard-fork switch block (nil = no fork)
	// DAOForkSupport bool     `json:"daoForkSupport,omitempty"` // Whether the nodes supports or opposes the DAO hard-fork

	// // EIP150 implements the Gas price changes (https://github.com/ethereum/EIPs/issues/150)
	// EIP150Block *big.Int    `json:"eip150Block,omitempty"` // EIP150 HF block (nil = no fork)
	// EIP150Hash  common.Hash `json:"eip150Hash,omitempty"`  // EIP150 HF hash (fast sync aid)

	// EIP155Block *big.Int `json:"eip155Block,omitempty"` // EIP155 HF block
	// EIP158Block *big.Int `json:"eip158Block,omitempty"` // EIP158 HF block

	// ByzantiumBlock *big.Int `json:"byzantiumBlock,omitempty"` // Byzantium switch block (nil = no fork, 0 = alraedy on pangu)

	// Various consensus engines
	Ethash *EthashConfig `json:"ethash,omitempty"`
	// Clique *CliqueConfig `json:"clique,omitempty"`
}

// EthashConfig is the consensus engine configs for proof-of-work based sealing.
type EthashConfig struct{}

// String implements the stringer interface, returning the consensus engine details.
func (c *EthashConfig) String() string {
	return "ethash"
}

// MakeSigner returns a Signer based on the given chain config and block number.
// PANGU 0.8
// May use other signer if the protocol changes
// Following interfaces are following GETH 1.8

func MakeSigner(config *ChainConfig, blockNumber *big.Int) Signer {
	// var signer Signer
	// switch {
	// case config.IsEIP155(blockNumber):
	// 	signer = NewEIP155Signer(config.ChainId)
	// case config.IsPangu(blockNumber):
	// 	signer = PanguSigner{}
	// default:
	// 	signer = PanguSigner{}
	// }
	// return signer
	return NewPanguSigner(config.ChainId)
}

// SignTx signs the transaction using the given signer and private key
func SignTx(tx *Transaction, s Signer, prv *ecdsa.PrivateKey) (*Transaction, error) {
	// log.Info("[core/types/transaction_signing.go->SignTx]")
	h := s.Hash(tx)
	sig, err := crypto.Sign(h[:], prv)
	if err != nil {
		return nil, err
	}
	// Note this part changes from
	// s.WithSignature(tx, sig)
	// to
	// tx.WithSignature(s, sig)
	return tx.WithSignature(s, sig)
}

// Sender derives the sender from the tx using the signer derivation
// functions.

// Sender returns the address derived from the signature (V, R, S) using secp256k1
// elliptic curve and an error if it failed deriving or upon an incorrect
// signature.
//
// Sender may cache the address, allowing it to be used regardless of
// signing method. The cache is invalidated if the cached signer does
// not match the signer used in the current call.
//
func Sender(signer Signer, tx *Transaction) (common.Address, error) {
	//if system contract, return address {100}
	if tx.TxData.GetSystemFlag() > 0 {
		return common.BytesToAddress([]byte{100}), nil
	}
	// log.Info("[core/types/transaction_signing.go->Sender:non system call]")
	if sc := tx.from.Load(); sc != nil {
		sigCache := sc.(sigCache)
		// If the signer used to derive from in a previous
		// call is not the same as used current, invalidate
		// the cache.
		if sigCache.signer.Equal(signer) {
			return sigCache.from, nil
		}
	}

	addr, err := signer.Sender(tx)
	if err != nil {
		log.Info("[core/types/transaction_signing.go->Sender:get addr from Sender error")
		return common.Address{}, err
	}
	tx.from.Store(sigCache{signer: signer, from: addr})
	return addr, nil
}

//Changed the interface to GETH 1.8
// type Signer interface {
// 	// Hash returns the rlp encoded hash for signatures
// 	Hash(tx *Transaction) common.Hash
// 	// PubilcKey returns the public key derived from the signature
// 	PublicKey(tx *Transaction) ([]byte, error)
// 	// WithSignature returns a copy of the transaction with the given signature.
// 	// The signature must be encoded in [R || S || V] format where V is 0 or 1.
// 	WithSignature(tx *Transaction, sig []byte) (*Transaction, error)
// 	// Checks for equality on the signers
// 	Equal(Signer) bool
// }
// Signer encapsulates transaction signature handling. Note that this interface is not a
// stable API and may change at any time to accommodate new protocol rules.
type Signer interface {
	// Sender returns the sender address of the transaction.
	Sender(tx *Transaction) (common.Address, error)
	// SignatureValues returns the raw R, S, V values corresponding to the
	// given signature.
	SignatureValues(tx *Transaction, sig []byte) (r, s, v *big.Int, err error)
	// Hash returns the hash to be signed.
	Hash(tx *Transaction) common.Hash
	// Equal returns true if the given signer is the same as the receiver.
	Equal(Signer) bool
}

// EIP155Transaction implements TransactionInterface using the
// EIP155 rules
//type EIP155Signer struct {
type PanguSigner struct {
	chainId, chainIdMul *big.Int
}

// func NewEIP155Signer(chainId *big.Int) EIP155Signer {
//Following the EIP155 rules
//
func NewPanguSigner(inchainID *big.Int) PanguSigner {
	if inchainID == nil {
		inchainID = new(big.Int)
	}
	return PanguSigner{
		chainId:    inchainID,
		chainIdMul: new(big.Int).Mul(inchainID, big.NewInt(2)),
	}
}

func (ps PanguSigner) Equal(s2 Signer) bool {
	pangu, ok := s2.(PanguSigner)
	return ok && pangu.chainId.Cmp(ps.chainId) == 0
}

/*
 * Compute the public key based on the input transaction
 * signature
 *
 */

// WithSignature returns a new transaction with the given signature. This signature
// needs to be in the [R || S || V] format where V is 0 or 1.
// Replaced by the SignatureValues
func (ps PanguSigner) SignatureValues(tx *Transaction, sig []byte) (R, S, V *big.Int, err error) {
	// R, S, V, err = PanguSigner{}.SignatureValues(tx, sig)
	if len(sig) != 65 {
		panic(fmt.Sprintf("wrong size for signature: got %d, want 65", len(sig)))
	}
	R = new(big.Int).SetBytes(sig[:32])
	S = new(big.Int).SetBytes(sig[32:64])
	V = new(big.Int).SetBytes([]byte{sig[64] + 27})

	log.Debugf("[core/types/transaction_signing.go->PANGU signer] chainID: %v\n", ps.chainId)
	if ps.chainId != nil {
		fmt.Printf("sign %v\n", ps.chainId.Sign())
		if ps.chainId.Sign() != 0 {
			V = big.NewInt(int64(sig[64] + 35))
			V.Add(V, ps.chainIdMul)
		}
	}

	return R, S, V, nil
}

// Hash returns the hash to be signed by the sender.
// It does not uniquely identify the transaction.
// Need to be updated with the transaction data structure.
// 2018/03/15 Updated the transaction with SystemContract
// and ShardingFlag

func (ps PanguSigner) Hash(tx *Transaction) common.Hash {
	// log.Info("[core/types/transaction_signing.go->EIP155Signer.Hash]")
	return rlpHash([]interface{}{
		tx.TxData.AccountNonce,
		tx.TxData.SystemContract,
		tx.TxData.Price,
		tx.TxData.GasLimit,
		tx.TxData.Recipient,
		tx.TxData.Amount,
		tx.TxData.Payload,
		tx.TxData.ShardingFlag,
		ps.chainId, uint(0), uint(0),
	})
}

/*
 * Derived the sender info from the TX input
 */
func (ps PanguSigner) Sender(tx *Transaction) (common.Address, error) {
	if !tx.Protected() {
		//Report error if the input signature is not protected
		return common.Address{}, ErrUnproctedTX
	}
	// fmt.Printf("Sender TX chainID in protected TX: %v, ps id %v\n", tx.ChainId(), ps.chainId)
	log.Debugf("[core/types/transaction_signing.go->PanguSigner.Sender:%v\n", tx.ChainId())
	if tx.ChainId().Cmp(ps.chainId) != 0 {
		log.Debugf("[core/types/transaction_signing.go->PanguSigner.Sender:unmatched chain ID%v\n", tx.ChainId())
		return common.Address{}, ErrInvalidChainId
	}

	//Get the value from input TX data.V to
	// fmt.Printf("tx.TxData.V: %v,  ps.chainIdMul %v\n", tx.TxData.V, ps.chainIdMul)
	V := new(big.Int).Sub(tx.TxData.V, ps.chainIdMul)
	// fmt.Printf("TX V after subtract 8: %v\n", V)
	var big8 = big.NewInt(8)
	V.Sub(V, big8)
	// fmt.Printf("TX V after subtract 8: %v\n", V)
	//Need to make sure the input is 27,
	//Get the Sender info
	return recoverPlain(ps.Hash(tx), tx.TxData.R, tx.TxData.S, V, true)
}

// New function used in GETH 1.7 and later
// Return the Address of the transaction Sender
//
func recoverPlain(sighash common.Hash, R, S, Vb *big.Int, pangu bool) (common.Address, error) {
	if Vb.BitLen() > 8 {
		return common.Address{}, ErrInvalidSig
	}
	//Compute the actual V value
	//For EIP155: v = CHAIN_ID * 2 + 35 or v = CHAIN_ID * 2 + 36
	//then when computing the hash of a transaction for purposes of signing or recovering,
	//instead of hashing only the first six elements (i.e. nonce, gasprice, startgas, to, value, data),
	//hash nine elements, with v replaced by CHAIN_ID, r = 0 and s = 0.
	//The currently existing signature scheme using v = 27 and v = 28,
	// remains valid and continues to operate under the same rules as it does now.
	//v is the `recovery id', a 1 byte value specifying the sign and niteness of the curve point; this
	//value is in the range of [27; 30], however we declare the upper two possibilities, representing innite values, invalid
	V := byte(Vb.Uint64() - 27)
	// fmt.Printf("V used to ValidateSignatureValues:%v\n", V)

	//When validate, V need to be either 27 or 28
	if !crypto.ValidateSignatureValues(V, R, S, pangu) {
		return common.Address{}, ErrInvalidSig
	}
	// encode the snature in uncompressed format
	r, s := R.Bytes(), S.Bytes()
	sig := make([]byte, 65)
	copy(sig[32-len(r):32], r)
	copy(sig[64-len(s):64], s)
	sig[64] = V

	// fmt.Printf("After ValidateSignatureValues:%v\n", V)
	// recover the public key from the snature
	pub, err := crypto.Ecrecover(sighash[:], sig)
	if err != nil {
		return common.Address{}, err
	}
	if len(pub) == 0 || pub[0] != 4 {
		return common.Address{}, errors.New("invalid public key")
	}
	// fmt.Printf("Pubkey:%v\n", pub)
	var addr common.Address
	copy(addr[:], crypto.Keccak256(pub[1:])[12:])
	return addr, nil
}

// deriveChainId derives the chain id from the given v parameter
func deriveChainId(v *big.Int) *big.Int {
	if v.BitLen() <= 64 {
		v := v.Uint64()
		if v == 27 || v == 28 {
			return new(big.Int)
		}
		return new(big.Int).SetUint64((v - 35) / 2)
	}
	v = new(big.Int).Sub(v, big.NewInt(35))
	return v.Div(v, big.NewInt(2))
}
