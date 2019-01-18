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

package types

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"

	"Chain3Go/lib/common"
	"Chain3Go/lib/crypto"
	"Chain3Go/lib/log"
	"Chain3Go/lib/params"
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

// MakeSigner returns a Signer based on the given chain config and block number.
// PANGU 0.8
// May use other signer if the protocol changes
// Following interfaces are following GETH 1.8
//

func MakeSigner(config *params.ChainConfig, blockNumber *big.Int) Signer {
	// var signer Signer
	// Reserved for future changes
	// switch {
	// case config.IsEIP155(blockNumber):
	// 	signer = NewEIP155Signer(config.ChainId)
	// case config.IsPangu(blockNumber):
	// 	signer = PanguSigner{}
	// default:
	// 	signer = PanguSigner{}
	// }
	// return signer
	// log.Debugf("[core/types/New .go->MakeSigner:ChainId] %v", config.ChainId)
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
	// fmt.Printf("Assign new PanguSigner with %v\n", inchainID)
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
// 2018/04/23
// Add Via field

func (ps PanguSigner) Hash(tx *Transaction) common.Hash {
	log.Debugf("[core/types/transaction_signing.go->PanguSigner.Hash]:%v", ps.chainId)
	return common.RlpHash([]interface{}{
		tx.TxData.AccountNonce,
		tx.TxData.SystemContract,
		tx.TxData.Price,
		tx.TxData.GasLimit,
		tx.TxData.Recipient,
		tx.TxData.Amount,
		tx.TxData.Payload,
		tx.TxData.ShardingFlag,
		tx.TxData.Via,
		ps.chainId, uint(0), uint(0),
	})
}

/*
 * Derived the sender info from the TX input
 * panguSigner 0.8.2
 *
 */
func (ps PanguSigner) Sender(tx *Transaction) (common.Address, error) {
	if !tx.Protected() {
		//Report error if the input signature is not protected
		return common.Address{}, ErrUnproctedTX
	}
	// fmt.Printf("Sender TX chainID in protected TX: %v, ps id %v\n", tx.ChainId(), ps.chainId)
	log.Debugf("[core/types/transaction_signing.go->PanguSigner.Sender:%v, ps %v\n", tx.ChainId(), ps.chainId)
	if tx.ChainId().Cmp(ps.chainId) != 0 {
		// log.Debugf("[core/types/transaction_signing.go->PanguSigner.Sender:unmatched chain ID%v\n", tx.ChainId())
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
	//Return the address
	var addr common.Address
	copy(addr[:], crypto.Keccak256(pub[1:])[12:])
	return addr, nil
}

// deriveChainId derives the chain id from the given v parameter
func deriveChainId(v *big.Int) *big.Int {
	if v.BitLen() <= 64 {
		//If v value is a UINT64 number
		v := v.Uint64()
		// No chainID is included in the V
		if v == 27 || v == 28 {
			//This should not happen in MOAC network

			return new(big.Int)
		}
		//EIP155 compute
		return new(big.Int).SetUint64((v - 35) / 2)
	}
	//If v is really large
	v = new(big.Int).Sub(v, big.NewInt(35))
	return v.Div(v, big.NewInt(2))
}
