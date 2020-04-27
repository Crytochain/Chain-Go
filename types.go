// types project types.go
package types

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

type ComplexIntParameter int64

func (v ComplexIntParameter) ToHex() string {

	return fmt.Sprintf("0x%x", v)
}

type ComplexIntResponse string

func (v ComplexIntResponse) ToUint64() (uint64, error) {

	strValue := string(v)
	cleaned := strings.Replace(strValue, "0x", "", -1)
	return strconv.ParseUint(cleaned, 16, 64)
}

func (v ComplexIntResponse) ToInt64() (int64, error) {

	strValue := string(v)
	cleaned := strings.Replace(strValue, "0x", "", -1)
	return strconv.ParseInt(cleaned, 16, 64)
}

type ComplexString string

func (v ComplexString) ToHex() string {

	if strings.Contains(string(v), "0x") {
		return string(v)
	}
	return fmt.Sprintf("0x%x", v)
}

func (v ComplexString) ToString() (string, error) {

	strValue := string(v)
	cleaned := strings.Replace(strValue, "0x", "", -1)
	sResult, err := hex.DecodeString(cleaned)
	return string(sResult), err
}
