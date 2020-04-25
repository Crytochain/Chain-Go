// utils project utils.go
package utils

import (
	"bytes"
	"encoding/json"
	"strconv"

	"github.com/tonnerre/golang-go.crypto/sha3"
)

type JsonRpc struct {
	Version string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Id      int         `json:"id"`
}

//转换成json字符串
func (jRpc *JsonRpc) AsJsonString() string {

	resultBytes, err := json.Marshal(jRpc)
	if err != nil {
		return ""
	}
	return string(resultBytes)
}

func bytesToHex(data []byte) string {
	buffer := new(bytes.Buffer)
	for _, b := range data {
		s := strconv.FormatInt(int64(b&0xFF), 16)
		if len(s) == 1 {
			buffer.WriteString("0")
		}
		buffer.WriteString(s)
	}
	return "0x" + buffer.String()
}

func Sha3Hash(dataStr string) string {

	d := sha3.NewKeccak256()
	d.Write([]byte(dataStr))

	return bytesToHex(d.Sum(nil))[:10]
}
