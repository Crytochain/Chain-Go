// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Adapted from: https://golang.org/src/crypto/cipher/xor_test.go

package bitutil

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
)

/*
 * Used to handle Control Flag,
 * SystemFlag 1000 0000 = 0x80 = 128
 * QuesryFlag 0100 0000 = 0x40 = 64
 * ShardingFlag 0010 0000 = 0x20 = 32
 */
func TestControlFlag(t *testing.T) {
	// SystemFlag := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80}   //, 0x00, 0x00, 0x00}
	// QueryFlag := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40}    //, 0x00, 0x00, 0x00}
	// ShardingFlag := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x20} // 0x00, 0x00, 0x00}

	SystemFlag := []byte{0x80}
	QueryFlag := []byte{0x40}
	ShardingFlag := []byte{0x20}

	//Set get flag
	// // b2 := []byte{0xe8, 0x03, 0xd0, 0x07}
	// fmt.Printf("systemflag %x\n", binary.BigEndian.Uint64(SystemFlag))
	// fmt.Printf("systemflag %x\n", binary.BigEndian.Uint64(QueryFlag))
	// fmt.Printf("systemflag %x\n", binary.BigEndian.Uint64(ShardingFlag))

	// var testflag []byte = 224
	var testflag uint64 = 224
	buf := make([]byte, binary.MaxVarintLen64)
	// n := binary.PutUvarint(buf, testflag) //convert the int flag into byte array
	binary.BigEndian.PutUint64(buf, testflag)
	fmt.Printf("Testflag: %x, %x\n", buf, testflag)
	// f1 := binary.LittleEndian.Uint64(b2[0:])
	fmt.Printf("Want %x, Get %x\n", binary.BigEndian.Uint64(buf), testflag)
	// var inval uint64 = 224
	// b1 := []byte{0x00} //test result
	b2 := make([]byte, binary.MaxVarintLen64)

	//Get Flags
	ANDBytes(b2, SystemFlag, buf)
	fmt.Printf("After and Want %x, Get %x\n", b2, SystemFlag)

	if !bytes.Equal(b2, SystemFlag) {

		t.Error("System flag not set")
	}

	//test set flag
	//reset flag to 0
	b2 = make([]byte, binary.MaxVarintLen64)

	binary.LittleEndian.PutUint64(buf, testflag)
	ORBytes(b2, SystemFlag, buf)
	fmt.Printf("Set  : %x, BYTES %x, UINT64: %x\n", SystemFlag, b2, binary.LittleEndian.Uint64(b2))
	XORBytes(b2, SystemFlag, b2)
	fmt.Printf("Clear  : %x, BYTES %x, UINT64: %x\n", SystemFlag, b2, binary.LittleEndian.Uint64(b2))
	return
	binary.BigEndian.PutUint64(buf, testflag)
	ORBytes(b2, SystemFlag, buf)
	fmt.Printf("Set  : %x, BYTES %x, UINT64: %x\n", SystemFlag, b2, binary.BigEndian.Uint64(b2))
	XORBytes(b2, SystemFlag, b2)
	fmt.Printf("Clear  : %x, BYTES %x, UINT64: %x\n", SystemFlag, b2, binary.BigEndian.Uint64(b2))
	return
	ORBytes(b2, QueryFlag, b2)
	fmt.Printf("Set has %x, want %x\n", b2, QueryFlag)
	ORBytes(b2, ShardingFlag, b2)
	fmt.Printf("Set has %x, want %x\n", b2, ShardingFlag)

	// &    bitwise AND            integers
	// |    bitwise OR             integers
	// ^    bitwise XOR            integers
	// &^   bit clear (AND NOT)    integers

	f1 := binary.BigEndian.Uint64(b2[0:])
	// f1 := binary.LittleEndian.Uint64(b2[0:])
	fmt.Printf("Want %x, Get %x\n", f1, testflag)

	if f1 != testflag {
		fmt.Printf("Input %x, Output %x flags not equal\n", f1, testflag)
		t.Error("Input Output flags not equal")
	}
}

// Tests that bitwise XOR works for various alignments.
func TestXOR(t *testing.T) {
	for alignP := 0; alignP < 2; alignP++ {
		for alignQ := 0; alignQ < 2; alignQ++ {
			for alignD := 0; alignD < 2; alignD++ {
				p := make([]byte, 1023)[alignP:]
				q := make([]byte, 1023)[alignQ:]

				for i := 0; i < len(p); i++ {
					p[i] = byte(i)
				}
				for i := 0; i < len(q); i++ {
					q[i] = byte(len(q) - i)
				}
				d1 := make([]byte, 1023+alignD)[alignD:]
				d2 := make([]byte, 1023+alignD)[alignD:]

				XORBytes(d1, p, q)
				safeXORBytes(d2, p, q)
				if !bytes.Equal(d1, d2) {
					t.Error("not equal", d1, d2)
				}
			}
		}
	}
}

// Tests that bitwise AND works for various alignments.
// Test
func TestAND(t *testing.T) {
	for alignP := 0; alignP < 2; alignP++ {
		for alignQ := 0; alignQ < 2; alignQ++ {
			for alignD := 0; alignD < 2; alignD++ {
				p := make([]byte, 1023)[alignP:]
				q := make([]byte, 1023)[alignQ:]

				for i := 0; i < len(p); i++ {
					p[i] = byte(i)
				}
				for i := 0; i < len(q); i++ {
					q[i] = byte(len(q) - i)
				}
				d1 := make([]byte, 1023+alignD)[alignD:]
				d2 := make([]byte, 1023+alignD)[alignD:]

				ANDBytes(d1, p, q)
				safeANDBytes(d2, p, q)
				if !bytes.Equal(d1, d2) {
					t.Error("not equal")
				}
			}
		}
	}
}

// Tests that bitwise OR works for various alignments.
func TestOR(t *testing.T) {
	for alignP := 0; alignP < 2; alignP++ {
		for alignQ := 0; alignQ < 2; alignQ++ {
			for alignD := 0; alignD < 2; alignD++ {
				p := make([]byte, 1023)[alignP:]
				q := make([]byte, 1023)[alignQ:]

				for i := 0; i < len(p); i++ {
					p[i] = byte(i)
				}
				for i := 0; i < len(q); i++ {
					q[i] = byte(len(q) - i)
				}
				d1 := make([]byte, 1023+alignD)[alignD:]
				d2 := make([]byte, 1023+alignD)[alignD:]

				ORBytes(d1, p, q)
				safeORBytes(d2, p, q)
				if !bytes.Equal(d1, d2) {
					t.Error("not equal")
				}
			}
		}
	}
}

// Tests that bit testing works for various alignments.
func TestTest(t *testing.T) {
	for align := 0; align < 2; align++ {
		// Test for bits set in the bulk part
		p := make([]byte, 1023)[align:]
		p[100] = 1

		if TestBytes(p) != safeTestBytes(p) {
			t.Error("not equal")
		}
		// Test for bits set in the tail part
		q := make([]byte, 1023)[align:]
		q[len(q)-1] = 1

		if TestBytes(q) != safeTestBytes(q) {
			t.Error("not equal")
		}
	}
}

// Benchmarks the potentially optimized XOR performance.
func BenchmarkFastXOR1KB(b *testing.B) { benchmarkFastXOR(b, 1024) }
func BenchmarkFastXOR2KB(b *testing.B) { benchmarkFastXOR(b, 2048) }
func BenchmarkFastXOR4KB(b *testing.B) { benchmarkFastXOR(b, 4096) }

func benchmarkFastXOR(b *testing.B, size int) {
	p, q := make([]byte, size), make([]byte, size)

	for i := 0; i < b.N; i++ {
		XORBytes(p, p, q)
	}
}

// Benchmarks the baseline XOR performance.
func BenchmarkBaseXOR1KB(b *testing.B) { benchmarkBaseXOR(b, 1024) }
func BenchmarkBaseXOR2KB(b *testing.B) { benchmarkBaseXOR(b, 2048) }
func BenchmarkBaseXOR4KB(b *testing.B) { benchmarkBaseXOR(b, 4096) }

func benchmarkBaseXOR(b *testing.B, size int) {
	p, q := make([]byte, size), make([]byte, size)

	for i := 0; i < b.N; i++ {
		safeXORBytes(p, p, q)
	}
}

// Benchmarks the potentially optimized AND performance.
func BenchmarkFastAND1KB(b *testing.B) { benchmarkFastAND(b, 1024) }
func BenchmarkFastAND2KB(b *testing.B) { benchmarkFastAND(b, 2048) }
func BenchmarkFastAND4KB(b *testing.B) { benchmarkFastAND(b, 4096) }

func benchmarkFastAND(b *testing.B, size int) {
	p, q := make([]byte, size), make([]byte, size)

	for i := 0; i < b.N; i++ {
		ANDBytes(p, p, q)
	}
}

// Benchmarks the baseline AND performance.
func BenchmarkBaseAND1KB(b *testing.B) { benchmarkBaseAND(b, 1024) }
func BenchmarkBaseAND2KB(b *testing.B) { benchmarkBaseAND(b, 2048) }
func BenchmarkBaseAND4KB(b *testing.B) { benchmarkBaseAND(b, 4096) }

func benchmarkBaseAND(b *testing.B, size int) {
	p, q := make([]byte, size), make([]byte, size)

	for i := 0; i < b.N; i++ {
		safeANDBytes(p, p, q)
	}
}

// Benchmarks the potentially optimized OR performance.
func BenchmarkFastOR1KB(b *testing.B) { benchmarkFastOR(b, 1024) }
func BenchmarkFastOR2KB(b *testing.B) { benchmarkFastOR(b, 2048) }
func BenchmarkFastOR4KB(b *testing.B) { benchmarkFastOR(b, 4096) }

func benchmarkFastOR(b *testing.B, size int) {
	p, q := make([]byte, size), make([]byte, size)

	for i := 0; i < b.N; i++ {
		ORBytes(p, p, q)
	}
}

// Benchmarks the baseline OR performance.
func BenchmarkBaseOR1KB(b *testing.B) { benchmarkBaseOR(b, 1024) }
func BenchmarkBaseOR2KB(b *testing.B) { benchmarkBaseOR(b, 2048) }
func BenchmarkBaseOR4KB(b *testing.B) { benchmarkBaseOR(b, 4096) }

func benchmarkBaseOR(b *testing.B, size int) {
	p, q := make([]byte, size), make([]byte, size)

	for i := 0; i < b.N; i++ {
		safeORBytes(p, p, q)
	}
}

// Benchmarks the potentially optimized bit testing performance.
func BenchmarkFastTest1KB(b *testing.B) { benchmarkFastTest(b, 1024) }
func BenchmarkFastTest2KB(b *testing.B) { benchmarkFastTest(b, 2048) }
func BenchmarkFastTest4KB(b *testing.B) { benchmarkFastTest(b, 4096) }

func benchmarkFastTest(b *testing.B, size int) {
	p := make([]byte, size)
	for i := 0; i < b.N; i++ {
		TestBytes(p)
	}
}

// Benchmarks the baseline bit testing performance.
func BenchmarkBaseTest1KB(b *testing.B) { benchmarkBaseTest(b, 1024) }
func BenchmarkBaseTest2KB(b *testing.B) { benchmarkBaseTest(b, 2048) }
func BenchmarkBaseTest4KB(b *testing.B) { benchmarkBaseTest(b, 4096) }

func benchmarkBaseTest(b *testing.B, size int) {
	p := make([]byte, size)
	for i := 0; i < b.N; i++ {
		safeTestBytes(p)
	}
}
