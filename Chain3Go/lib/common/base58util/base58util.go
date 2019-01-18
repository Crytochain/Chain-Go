// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Adapted from: https://github.com/btcsuite/btcutil

// Package base58util implements base58 encode/decode operations.

package base58util

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"
)

var bigRadix = big.NewInt(58)
var bigZero = big.NewInt(0)

const (
	// alphabet is the modified base58 alphabet used by Bitcoin.
	//alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	//moac
	base58alphabet = "mcsj1qinh2xue3ora4fy56g7b89tzdwkpvMJQNSHXUERAFYCGBTZDLWKPV"

	alphabetIdx0   = '1'
	addressPrefix  = 0   //set m as the 1st letter of the account address
	contractPrefix = 59  //c as contract address
	seedPrefix     = 's' //secret prefix 's'
	versionMoac    = 0x00
)

var b58_map = [256]byte{
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 4, 9, 13, 17, 20, 21, 23,
	25, 26, 255, 255, 255, 255, 255, 255,
	255, 44, 49, 47, 52, 42, 45, 48,
	39, 255, 35, 55, 53, 34, 37, 255,
	56, 36, 43, 38, 50, 41, 57, 54,
	40, 46, 51, 255, 255, 255, 255, 255,
	255, 16, 24, 1, 29, 12, 18, 22,
	8, 6, 3, 31, 255, 0, 7, 14,
	32, 5, 15, 2, 27, 11, 33, 30,
	10, 19, 28, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
}

// ErrChecksum indicates that the checksum of a check-encoded string does not verify against
// the checksum.
var ErrChecksum = errors.New("checksum error")
var ErrVersion = errors.New("moac version error")

// ErrInvalidFormat indicates that the check-encoded string has an invalid format.
var ErrInvalidFormat = errors.New("invalid format: version and/or checksum bytes missing")

// checksum: first four bytes of sha256^2
// test with 2 bytes
func checksum(input []byte) (cksum [2]byte) {
	h := sha256.Sum256(input)
	h2 := sha256.Sum256(h[:])
	//fmt.Printf("sha256 len, %d, h2 len %d\n", len(h), len(h2))
	copy(cksum[:], h2[:2]) //change from 4 to 2

	return
}

/*
 * Encode the account address, contract address
 * with checksum
 * This function should be called after check the
 * input HEX address.
 * data - byte array with
 * prefix - 00 for account address,
 *       or 01 for contract address.
 */
func MoacEncode(data []byte, prefix byte) string {

	//remove the header 0x
	if len(data) > 2 {
		if data[0] == '0' && data[1] == 'x' {
			data = data[2:]
		}
	}

	//Add checksum bytes for the input version data
	//The version bytes is used for future processing
	b := make([]byte, 0, 1+len(data)+2) //+4
	b = append(b, prefix)
	b = append(b, data[:]...)

	// fmt.Printf("Encode:%X\n", b)
	cksum := checksum(b)
	b = append(b, cksum[:]...)

	// fmt.Println("Actual Encoded len:", len(b))
	// fmt.Printf("%X\n", b)
	x := new(big.Int)
	x.SetBytes(b)

	answer := make([]byte, 0, len(b)*136/100)
	for x.Cmp(bigZero) > 0 {
		mod := new(big.Int)
		x.DivMod(x, bigRadix, mod)
		answer = append(answer, base58alphabet[mod.Int64()])
	}

	// fmt.Println("Before add Zeros:", len(answer), string(answer))
	// put the zero bytes at the end of string
	for _, i := range b {
		if i != 0 {
			break
		}
		//answer = append(answer, byte(myprefix))
		answer = append(answer, base58alphabet[0])
		// fmt.Println("Got 00 bytes in input")
	}

	//reverse the string
	alen := len(answer)
	for i := 0; i < alen/2; i++ {
		answer[i], answer[alen-1-i] = answer[alen-1-i], answer[i]
	}

	// fmt.Println("Encoded str:", len(answer), string(answer))
	return string(answer)
}

/*
 * Decode the input string into byte array
 * with checksum
 * REturn is a byte array holding the data
 * can be converted to string with hex.EncodeToString
 */
func MoacDecode(input string) (result []byte, version byte, err error) {

	res := big.NewInt(0)
	j := big.NewInt(1)

	scratch := new(big.Int)

	// fmt.Println("Decoded string len:", len(input))

	for i := len(input) - 1; i >= 0; i-- {
		tmp := b58_map[input[i]]

		scratch.SetInt64(int64(tmp))
		scratch.Mul(j, scratch)
		res.Add(res, scratch)
		j.Mul(j, bigRadix)
	}

	tmpval := res.Bytes()

	// fmt.Printf("Decode:%d %X\n", len(tmpval), tmpval)

	//Add prefix 00 bytes if the decode len is shorter
	var numZeros int
	for numZeros = 0; numZeros < len(input); numZeros++ {
		if input[numZeros] != byte(base58alphabet[0]) {
			break
		}
	}

	flen := numZeros + len(tmpval)
	// fmt.Println("Num of zeros", numZeros, len(tmpval), flen, len(tmpval))

	val := make([]byte, flen)
	copy(val[numZeros:], tmpval)

	//remove the prefix byte
	//and check the checksum
	//use 2 bytes
	var cksum [2]byte                //[4]byte
	copy(cksum[:], val[len(val)-2:]) //4
	// fmt.Printf("In Checksum:%X\n", val[:len(val)-2])

	if checksum(val[:len(val)-2]) != cksum {

		fmt.Printf("Checksum error: expect %x, get %x\n", cksum, checksum(val[:len(val)-2]))
		return nil, 0, ErrChecksum
	} else {
		decode_str := make([]byte, len(val)-2-1)
		copy(decode_str[:], val[1:len(val)])
		return decode_str, val[0], nil
	}
}

// CheckDecode decodes a string that was encoded with CheckEncode and verifies the checksum.
// func CheckMoacVersion(input byte) bool {
// 	return input == versionMoac
// }
