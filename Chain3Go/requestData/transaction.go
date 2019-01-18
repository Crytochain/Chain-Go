// requestData project transaction.go
package requestData

import (
	"reflect"

	"Chain3Go/types"
)

// TransactionParameters GO transaction to make more easy controll the parameters
type TransactionParameters struct {
	From     string
	To       string
	Gas      types.ComplexIntParameter
	GasPrice types.ComplexIntParameter
	Value    types.ComplexIntParameter
	Data     types.ComplexString
}

// RequestTransactionParameters JSON
type RequestTransactionParameters struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Gas      string `json:"gas,omitempty"`
	GasPrice string `json:"gasPrice,omitempty"`
	Value    string `json:"value"`
	Data     string `json:"data,omitempty"`
}

// Transform the GO transactions parameters to json style
func (params *TransactionParameters) Transform() *RequestTransactionParameters {
	request := new(RequestTransactionParameters)
	request.From = params.From
	request.To = params.To
	if params.Gas != 0 {
		request.Gas = params.Gas.ToHex()
	} else {
		request.Gas = "0x0"
	}
	if params.GasPrice != 0 {
		request.GasPrice = params.GasPrice.ToHex()
	} else {
		request.GasPrice = "0x0"
	}
	if params.Value != 0 {
		request.Value = params.Value.ToHex()
	} else {
		request.Value = "0x0"
	}
	if params.Data != "" {
		request.Data = params.Data.ToHex()
	} else {
		request.Data = "0x0"
	}
	return request
}

type TransactionResponse struct {
	//	Hash             string                   `json:"hash"`
	//	Nonce            int                      `json:"nonce"`
	//	BlockHash        string                   `json:"blockHash"`
	//	BlockNumber      int64                    `json:"blockNumber"`
	//	TransactionIndex int64                    `json:"transactionIndex"`
	//	From             string                   `json:"from"`
	//	To               string                   `json:"to"`
	//	Value            types.ComplexIntResponse `json:"value"`
	//	GasPrice         types.ComplexIntResponse `json:"gasPrice,omitempty"`
	//	Gas              types.ComplexIntResponse `json:"gas,omitempty"`
	//	Data             types.ComplexString      `json:"input,omitempty"`
	//	ShardingFlag     uint8                    `json:"shardingFlag"`
	Hash             string `json:"hash"`
	BlockHash        string `json:"blockHash"`
	BlockNumber      string `json:"blockNumber"`
	TransactionIndex string `json:"transactionIndex"`
	AccountNonce     string `json:"nonce"`
	GasPrice         string `json:"gasPrice"`
	GasLimit         string `json:"gas"`
	To               string `json:"to"`
	From             string `json:"from"`
	Amount           string `json:"value"`
	Payload          string `json:"input"`

	// Signature values
	V string `json:"v"`
	R string `json:"r"`
	S string `json:"s"`

	ShardingFlag string `json:"shardingFlag"`
}

type TransactionReceipt struct {
	TransactionHash   string   `json:"transactionHash"`
	TransactionIndex  string   `json:"transactionIndex"`
	BlockHash         string   `json:"blockHash"`
	BlockNumber       string   `json:"blockNumber"`
	CumulativeGasUsed string   `json:"cumulativeGasUsed"`
	GasUsed           string   `json:"gasUsed"`
	ContractAddress   string   `json:"contractAddress"`
	Logs              []string `json:"logs"`
	To                string   `json:"to"`
	From              string   `json:"from"`
	Bloom             string   `json:"logsBloom"`
	Status            string   `json:"status"`
}

type FiltersLogs struct {
	LogIndex         string   `json:"logIndex"`
	BlockNumber      string   `json:"blockNumber"`
	BlockHash        string   `json:"blockHash"`
	TransactionHash  string   `json:"transactionHash"`
	TransactionIndex string   `json:"transactionIndex"`
	Address          string   `json:"address"`
	Data             string   `json:"data"`
	Topics           []string `json:"topics"`
}

type ScsTransactionReceipt struct {
	LogsArray       []Logs `json:"logs"`
	LogsBloom       string `json:"logsBloom"`
	Status          string `json:"status"`
	TransactionHash string `json:"transactionHash"`
}

type Logs struct {
	Address          string   `json:"address"`
	Topics           []string `json:"topics"`
	Data             string   `json:"data"`
	BlockNumber      string   `json:"blockNumber"`
	TransactionHash  string   `json:"transactionHash"`
	TransactionIndex string   `json:"transactionIndex"`
	BlockHash        string   `json:"blockHash"`
	LogIndex         string   `json:"logIndex"`
	Removed          string   `json:"removed"`
}

func Struct2Map(obj interface{}) map[string]interface{} {

	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		vType := t.Field(i).Type.String()
		vlaue := v.Field(i).Interface()
		data[t.Field(i).Name] = vlaue
		switch vType {
		case "string":
			if vlaue.(string) == "" || vlaue.(string) == "0x0" {
				delete(data, t.Field(i).Name)
			}
		case "int":
			if vlaue.(int) == 0 {
				delete(data, t.Field(i).Name)
			}
		default:
		}
	}

	return data
}
