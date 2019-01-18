// Copyright 2015 The MOAC-core Authors
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

package common

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"math/rand"
	"reflect"
	"strconv"
	"strings"

	"Chain3Go/lib/common/base58util"
	"Chain3Go/lib/common/hexutil"
	"Chain3Go/lib/crypto/sha3"
	"Chain3Go/lib/log"
	"Chain3Go/lib/rlp"
)

const (
	HashLength    = 32
	AddressLength = 20

	MoacAddressLength = 32
	MoacVersion       = 0
)

var (
	hashT    = reflect.TypeOf(Hash{})
	addressT = reflect.TypeOf(Address{})
	//moacaddressT = reflect.TypeOf(MoacAddress{})
)

// Hash represents the 32 byte Keccak256 hash of arbitrary data.
type Hash [HashLength]byte

func RlpHash(x interface{}) (h Hash) {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}

func BytesToHash(b []byte) Hash {
	var h Hash
	h.SetBytes(b)
	return h
}
func StringToHash(s string) Hash { return BytesToHash([]byte(s)) }
func BigToHash(b *big.Int) Hash  { return BytesToHash(b.Bytes()) }
func HexToHash(s string) Hash    { return BytesToHash(FromHex(s)) }

// Get the string representation of the underlying hash
func (h Hash) Str() string   { return string(h[:]) }
func (h Hash) Bytes() []byte { return h[:] }
func (h Hash) Big() *big.Int { return new(big.Int).SetBytes(h[:]) }
func (h Hash) Hex() string   { return hexutil.Encode(h[:]) }

// TerminalString implements log.TerminalStringer, formatting a string for console
// output during logging.
func (h Hash) TerminalString() string {
	return h.String()
	//return fmt.Sprintf("%x…%x", h[:3], h[29:])
}

// String implements the stringer interface and is used also by the logger when
// doing full logging into a file.
func (h Hash) String() string {
	return h.Hex()
}

// Format implements fmt.Formatter, forcing the byte slice to be formatted as is,
// without going through the stringer interface used for logging.
func (h Hash) Format(s fmt.State, c rune) {
	fmt.Fprintf(s, "%"+string(c), h[:])
}

// UnmarshalText parses a hash in hex syntax.
func (h *Hash) UnmarshalText(input []byte) error {
	log.Debug("[common/types.go->Hash.UnmarshalText] input=" + string(input))
	return hexutil.UnmarshalFixedText("Hash", input, h[:])
}

// UnmarshalJSON parses a hash in hex syntax.
func (h *Hash) UnmarshalJSON(input []byte) error {
	// log.Debug("[common/types.go->Hash.UnmarshalJSON] input=" + string(input))
	return hexutil.UnmarshalFixedJSON(hashT, input, h[:])
}

// MarshalText returns the hex representation of h.
func (h Hash) MarshalText() ([]byte, error) {
	return hexutil.Bytes(h[:]).MarshalText()
}

// Sets the hash to the value of b. If b is larger than len(h), 'b' will be cropped (from the left).
func (h *Hash) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-HashLength:]
	}

	copy(h[HashLength-len(b):], b)
}

// Set string `s` to h. If s is larger than len(h) s will be cropped (from left) to fit.
func (h *Hash) SetString(s string) { h.SetBytes([]byte(s)) }

// Sets h to other
func (h *Hash) Set(other Hash) {
	for i, v := range other {
		h[i] = v
	}
}

// Generate implements testing/quick.Generator.
func (h Hash) Generate(rand *rand.Rand, size int) reflect.Value {
	m := rand.Intn(len(h))
	for i := len(h) - 1; i > m; i-- {
		h[i] = byte(rand.Uint32())
	}
	return reflect.ValueOf(h)
}

func EmptyHash(h Hash) bool {
	return h == Hash{}
}

// UnprefixedHash allows marshaling a Hash without 0x prefix.
type UnprefixedHash Hash

// UnmarshalText decodes the hash from hex. The 0x prefix is optional.
func (h *UnprefixedHash) UnmarshalText(input []byte) error {
	log.Debugf("[common/types.go->UnprefixedHash.UnmarshalText] input" + string(input))
	return hexutil.UnmarshalFixedUnprefixedText("UnprefixedHash", input, h[:])
}

// MarshalText encodes the hash as hex.
func (h UnprefixedHash) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(h[:])), nil
}

/////////// Address

// Address represents the 20 byte address of an MoacNode account.
type Address [AddressLength]byte

func BytesToAddress(b []byte) Address {
	var a Address
	a.SetBytes(b)
	return a
}

func StringToAddress(s string) Address { return BytesToAddress([]byte(s)) }
func BigToAddress(b *big.Int) Address  { return BytesToAddress(b.Bytes()) }
func HexToAddress(s string) Address    { return BytesToAddress(FromHex(s)) }

//Decode the input string using base 58 encoding
//if the input string has the right length.

func Base58ToAddress(s string) Address {
	if len(s) == MoacAddressLength {
		// fmt.Println("Input length match")

		res, outVersion, err := base58util.MoacDecode(s)

		if err != nil {
			fmt.Println("Decode address as MOAC address failed", err)

		} else if outVersion != MoacVersion {
			fmt.Printf("MOAC address version mismatch: %X .vs. %X\n", outVersion, MoacVersion)

		}
		//TODO, handle the error in the decoding process
		outAddr := BytesToAddress(res)
		// fmt.Println("MOAC address returned as:", outAddr.Hex())
		return outAddr
	} else if IsHexAddress(s) {
		// fmt.Println("Input address is HEX! Try to convert......")
		return HexToAddress(s)

	} else {
		//If not the right length, try other functions to
		//crate the address
		// fmt.Println("Not MOAC or HEX, try string!")
		return StringToAddress(s)
	}
}

// IsHexAddress verifies whether a string can represent a valid hex-encoded
// MoacNode address or not.
func IsHexAddress(s string) bool {
	if len(s) == 2+2*AddressLength && IsHex(s) {
		return true
	}
	if len(s) == 2*AddressLength && IsHex("0x"+s) {
		return true
	}
	return false
}

// Get the string representation of the underlying address
func (a Address) Str() string {
	//fmt.Println("Use Str() in types.go")
	return string(a[:])
}
func (a Address) Bytes() []byte { return a[:] }
func (a Address) Big() *big.Int { return new(big.Int).SetBytes(a[:]) }
func (a Address) Hash() Hash    { return BytesToHash(a[:]) }

func (a Address) GetBase58Str() string {
	// fmt.Println("types.go: Get Base58 in HEX address:", a.Hex())

	enstr := base58util.MoacEncode(a.Bytes(), MoacVersion)
	return enstr //base58util.MoacEncode(a.Bytes(), 'm')

}

// Hex returns an EIP55-compliant hex string representation of the address.
func (a Address) Hex() string {
	unchecksummed := hex.EncodeToString(a[:])
	sha := sha3.NewKeccak256()
	sha.Write([]byte(unchecksummed))
	hash := sha.Sum(nil)

	result := []byte(unchecksummed)
	for i := 0; i < len(result); i++ {
		hashByte := hash[i/2]
		if i%2 == 0 {
			hashByte = hashByte >> 4
		} else {
			hashByte &= 0xf
		}
		if result[i] > '9' && hashByte > 7 {
			result[i] -= 32
		}
	}
	// fmt.Println("Use interface HEX types.go")
	return "0x" + string(result)
}

// String implements the stringer interface and is used also by the logger.
func (a Address) String() string {
	//fmt.Println("Use interface String types.go")
	// return a.GetBase58Str()
	return a.Hex()
}

// Format implements fmt.Formatter, forcing the byte slice to be formatted as is,
// without going through the stringer interface used for logging.
func (a Address) Format(s fmt.State, c rune) {
	// fmt.Println("Use interface format types.go")
	fmt.Fprintf(s, "%"+string(c), a[:])
}

// Sets the address to the value of b. If b is larger than len(a) it will panic
func (a *Address) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-AddressLength:]
	}
	copy(a[AddressLength-len(b):], b)
}

// Set string `s` to a. If s is larger than len(a) it will panic
func (a *Address) SetString(s string) { a.SetBytes([]byte(s)) }

// Sets a to other
func (a *Address) Set(other Address) {
	for i, v := range other {
		a[i] = v
	}
}

// MarshalText returns the hex representation of a.
func (a Address) MarshalText() ([]byte, error) {
	return hexutil.Bytes(a[:]).MarshalText()
}

// UnmarshalText parses a hash in hex syntax.
func (a *Address) UnmarshalText(input []byte) error {
	log.Debug("[common/types.go->Address.UnmarshalText] input" + string(input))
	return hexutil.UnmarshalFixedText("Address", input, a[:])
}

// UnmarshalJSON parses a hash in hex syntax.
func (a *Address) UnmarshalJSON(input []byte) error {
	inputStr := string(input)
	if len(inputStr) > 0 && inputStr[0] == '"' {
		inputStr = inputStr[1:]
	}
	if len(inputStr) > 0 && inputStr[len(inputStr)-1] == '"' {
		inputStr = inputStr[:len(inputStr)-1]
	}
	log.Debugf("[common/types.go->Address.UnmarshalJSON] input=" + inputStr)
	if IsMoacAddress(inputStr) {
		log.Debugf("It is a Moac address.")
		a.Set(MoacToAddress(inputStr))
		return nil
	}
	return hexutil.UnmarshalFixedJSON(addressT, input, a[:])
}

/////////// UnprefixedAddress
// UnprefixedHash allows marshaling an Address without 0x prefix.
type UnprefixedAddress Address

// UnmarshalText decodes the address from hex. The 0x prefix is optional.
func (a *UnprefixedAddress) UnmarshalText(input []byte) error {
	log.Debug("[common/types.go->UnprefixedAddress.UnmarshalText] input" + string(input))
	return hexutil.UnmarshalFixedUnprefixedText("UnprefixedAddress", input, a[:])
}

// MarshalText encodes the address as hex.
func (a UnprefixedAddress) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(a[:])), nil
}

/////////// MOAC Address
// Updated with prefix byte(1) and checksum bytes (4)

// MOAC Address represents the base58 encoding of 20 byte address of the account.
// Change to PrefixAddress
//
type MoacAddress string

// IsMoacAddress verifies whether a string can represent a valid base58-encoded
// MOAC account address

//
// Check if the input string is a valid
// MOAC encoded address, contract address
// or secret

func IsMoacAddress(data string) bool {

	// log.Debugf("Input data length IsMoacAddress:%v", len(data))
	if len(data) <= MoacAddressLength {

		res, outVersion, err := base58util.MoacDecode(data)

		if err != nil {
			log.Debug("Decode address as MOAC address failed")
			return false
		} else if outVersion != MoacVersion {
			log.Debug("MOAC address version mismatch: %X .vs. %X\n", outVersion, MoacVersion)
			return false
		}
		outAddr := BytesToAddress(res)
		log.Debug("MOAC address returned as:" + outAddr.Hex())
		return IsHexAddress(outAddr.Hex())
	}
	//fmt.Println("Input length not match!!!")
	return false
}

func BytesToMoacAddress(b []byte) MoacAddress {

	val := base58util.MoacEncode(b, MoacVersion)
	return MoacAddress(val)
}

//Conver input HEX string to MOAC address
func HexToMoac(in_str string) MoacAddress {

	b := FromHex(in_str)
	val := base58util.MoacEncode(b, MoacVersion)
	// fmt.Printf("HEXTOMOAC: %s\n", val)
	return MoacAddress(val)
}

// Get the string representation of the underlying address
func (a MoacAddress) Str() string   { return string(a) }
func (a MoacAddress) Big() *big.Int { return new(big.Int).SetBytes(a.Bytes()) }

//func (a MoacAddress) Hash() Hash    { return BytesToHash(a[:]) }

//Convert the base58 based address to HEX address as string
func (a MoacAddress) Hex() string {
	//decode the moac address into byte array
	res, outVersion, err := base58util.MoacDecode(string(a))

	if err != nil {
		fmt.Println("Error: ", err)
		return ""

	} else if outVersion == MoacVersion {
		return hexutil.Encode(res)

	} else {
		return "0x" + hexutil.Encode(res)
	}

}

//
func (a MoacAddress) Bytes() []byte {
	//decode the moac address into byte array
	res, outVersion, err := base58util.MoacDecode(a.Str())

	if err == nil && outVersion == MoacVersion {
		return res
	}
	return nil
}

//Convert the Moac address to HEX address
func MoacToAddress(inadd string) Address {

	res, outVersion, err := base58util.MoacDecode(inadd)

	var result Address

	if err == nil && outVersion == MoacVersion {
		result = BytesToAddress([]byte(res))
	}

	//Set the err code if version does not fit the current version
	return result
}

func KeytoKey(key string) string {
	hash := sha3.NewKeccak256()
	var buf []byte
	b, _ := hex.DecodeString(key)
	hash.Write(b)
	buf = hash.Sum(buf)
	return Bytes2Hex(buf)
}

func IncreaseHexByNum(num int64, hex string) string {
	a := strings.Split(hex, "")

	for i := len(a) - 1; i >= 0; i-- {
		if num <= 0 {
			break
		} else {
			x, _ := strconv.ParseInt("0x"+a[i], 0, 64)
			x += num
			a[i] = strconv.FormatInt(x%16, 16)
			num = x / 16
		}
	}

	return strings.Join(a, "")
}

func IncreaseHexByOne(hex string) string {
	a := strings.Split(hex, "")

	for i := len(a) - 1; i >= 0; i-- {
		if a[i] == "f" {
			a[i] = "0"
		} else {
			x, _ := strconv.ParseInt("0x"+a[i], 0, 64)
			a[i] = strconv.FormatInt(x+1, 16)
			break
		}
	}

	return strings.Join(a, "")
}

func GetValue(dict map[string]string, key string) string {
	for k, v := range dict {
		if k == key {
			if len(v) > 2 {
				return v[2:]
			} else {
				return v
			}
		}
	}

	return ""
}
