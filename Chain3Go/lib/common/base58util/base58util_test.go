package base58util

import (
	// "bytes"
	"encoding/hex"
	"fmt"
	"testing"
)

var checkMoacEncodingStringTests = []struct {
	version byte
	in      string
	out     string
}{
	{0, "", "mcwg"},
	// {1, " ", "mjNMhjvQ"},
	// {0, "-", "mqhcjhqx"},
	// {0, "0", "mq8D9XBV"},
	{0, "0x8605cdbbdb6d264aa742e77020dcbc58fcdce182", "mQ6ArdtrA2dPSLMNo2DX6jkQxUHqoeL"},
	{0, "0xb409bfca27b764315c868034ee406432a7b66168", "mGTnpUPUG52XiTWTCuZ6Qq2f45Nt2f4"},
	{0, "0x000000000000000000000000000001", "mmmmmmmmmmmmmmm7sR"},
	// {1, "0x000000000000000000000000000001", "mmmmmmmmmmmmmmm7sR"},
	{0, "0x100000000000000000000000000000", "mhUFDR6Pd79YZK8Z4LTwx2xG"},
	{0, "0xffffffffffffffffffffffffffffff", "ms7wSkPjMitgWy4hhWBWfUDbx"},
	{0, "0xfffffffffffffffffffffffffffffe", "ms7wSkPjMitgWy4hhWBWfUUAC"},
	// {0, "1", "mqwVGTTQ"},//Should not pass due to the new check
	// {0, "-1", "mgd3z5HGm"},
	{0, "11", "mbzBt5r3H"},
	// {'m', "abc", "m14h3c6cfU92"},
	// {'m', "1234598760", "m1K5zqBMZZTzUbAaWeMzf"},
	// // {'m', "abcdefghijklmnopqrstuvwxyz", "K2RYDcKfupxwXdWhSAxQPCeiULntKm63UXyx5MvEH2"},
	// {0, "0xdbf03b407c01e7cd3cbea99509d93f8dddc8c6fb", "m5sLAvLyhSUfqKemZZ3R9MZdDyAqdCJSeZ"},
	//mY9C12xCrVRurETWTxzm2iAsg6PVaM2
	// {0, "0xab30175830d9912fa768922f9bfaa6e610d14b07", "mY9C12xCrVRurETWTxzm2iAsg6PVaM2"},
	//{59, "0x47d33b27bb249a2dbab4c0612bf9caf4c1950855", "c3mZgz7Q4k8X8rqmWNqAx6vR4ctsoB4J"},
	// {0, "0x47d33b27bb249a2dbab4c0612bf9caf4c1950855", "myzSob7jKBLkNeEEx1yC2B8vYVdjh3t"},
	// {0, "dbf03b407c01e7cd3cbea99509d93f8dddc8c6fb", "mccX3JP9nb1TEvrdP9kzQjpwXHhhEf4g"},
}

func TestMoacAddress(t *testing.T) {
	fmt.Println(" Moac address Encoding/Decoding test......")

	for x, test := range checkMoacEncodingStringTests {
		// test encoding
		var data string

		if len(test.in) > 2 {
			//remove the header 0x
			if test.in[0] == '0' && test.in[1] == 'x' {
				data = test.in[2:]
			} else {
				data = test.in
			}
		} else {
			data = test.in
		}

		in_byte_array, err := hex.DecodeString(data)

		if err != nil {
			fmt.Printf("Input is not hex array:%v\n", err)
		}

		if res := MoacEncode(in_byte_array, test.version); res != test.out {
			fmt.Printf("Base58 str len: %d\n", len(res))
			t.Errorf("MoacEncodeAddress test #%d failed: got %s, want: %s", x, res, test.out)
		}

		// test decoding
		res1, version, err := MoacDecode(test.out)

		res := hex.EncodeToString(res1)

		if err != nil {
			t.Errorf("MoacDecode test #%d failed with err: %v", x, err)
		} else if version != addressPrefix && version != contractPrefix {
			t.Errorf("MoacDecode test #%d failed: got version: %d want: %d", x, version, test.version)
		} else if string(res) != data {
			t.Errorf("MoacDecode test #%d failed: got: %s want: %s", x, res, test.in)
		}
	}

	// test the two decoding failure cases
	// case 1: checksum error
	_, _, err := MoacDecode("3MNQE1Y")
	if err != ErrChecksum {
		t.Error("MoacDecode test failed, expected ErrChecksum")
	}

}
