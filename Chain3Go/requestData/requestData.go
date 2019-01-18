// requestData project requestData.go
package requestData

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"serverLayer/Chain3Go/constants"
	"serverLayer/Chain3Go/types"
)

type RequestResult struct {
	Id      int         `json:"id"`
	Version string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	Error   *Error      `json:"error,omitempty"`
	Data    string      `json:"data,omitempty"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func (pointer *RequestResult) ToStringArray() ([]string, error) {

	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.([]interface{})

	new := make([]string, len(result))
	for i, v := range result {
		new[i] = v.(string)
	}

	return new, nil
}

func (pointer *RequestResult) ToComplexString() (types.ComplexString, error) {

	if err := pointer.checkResponse(); err != nil {
		return "", err
	}

	result := (pointer).Result.(interface{})

	return result.(types.ComplexString), nil
}

func (pointer *RequestResult) ToString() (string, error) {

	if err := pointer.checkResponse(); err != nil {
		return "", err
	}

	result := (pointer).Result.(interface{})

	return result.(string), nil
}

func (pointer *RequestResult) ToInt() (int64, error) {

	if err := pointer.checkResponse(); err != nil {
		return 0, err
	}

	result := (pointer).Result.(interface{})

	hex := result.(string)

	numericResult, err := strconv.ParseInt(strings.Replace(hex, "0x", "", -1), 16, 64)

	return numericResult, err
}

func (pointer *RequestResult) ToComplexIntResponse() (types.ComplexIntResponse, error) {

	if err := pointer.checkResponse(); err != nil {
		return types.ComplexIntResponse(0), err
	}

	result := (pointer).Result.(interface{})

	var hex string

	switch v := result.(type) {
	//Testrpc returns a float64
	case float64:
		hex = strconv.FormatFloat(v, 'E', 16, 64)
		break
	default:
		hex = result.(string)
	}

	cleaned := strings.Replace(strings.Replace(hex, "0x", "", -1), "0x", "", -1)

	return types.ComplexIntResponse(cleaned), nil
}

func (pointer *RequestResult) ToBoolean() (bool, error) {

	if err := pointer.checkResponse(); err != nil {
		return false, err
	}

	result := (pointer).Result.(interface{})

	return result.(bool), nil
}

func (pointer *RequestResult) ToTransactionResponse() (*TransactionResponse, error) {

	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.(map[string]interface{})

	if len(result) == 0 {
		return nil, customerror.EMPTYRESPONSE
	}

	transactionResponse := &TransactionResponse{}

	marshal, err := json.Marshal(result)

	if err != nil {
		return nil, customerror.UNPARSEABLEINTERFACE
	}

	json.Unmarshal([]byte(marshal), transactionResponse)

	return transactionResponse, nil
}

func (pointer *RequestResult) ToTransactionReceipt() (*TransactionReceipt, error) {

	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.(map[string]interface{})

	if len(result) == 0 {
		return nil, customerror.EMPTYRESPONSE
	}

	transactionReceipt := &TransactionReceipt{}

	marshal, err := json.Marshal(result)

	if err != nil {
		return nil, customerror.UNPARSEABLEINTERFACE
	}

	json.Unmarshal([]byte(marshal), transactionReceipt)

	return transactionReceipt, nil
}

func (pointer *RequestResult) ToBlock(objFlag bool) (*Block, error) {

	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.(map[string]interface{})

	if len(result) == 0 {
		return nil, customerror.EMPTYRESPONSE
	}

	block := &Block{}

	marshal, err := json.Marshal(result)

	if err != nil {
		return nil, customerror.UNPARSEABLEINTERFACE
	}

	err = json.Unmarshal([]byte(marshal), block)

	if objFlag {
		//交易信息
		txD := make([]TransactionResponse, 0)
		marshal, err := json.Marshal((pointer.Result.(map[string]interface{}))["transactions"])
		if err != nil {
			return nil, customerror.UNPARSEABLEINTERFACE
		}
		err = json.Unmarshal([]byte(marshal), &txD)
		block.TransactionsTrue = txD
	} else {
		//交易hash
		txHash := ((pointer.Result.(map[string]interface{}))["transactions"]).([]interface{})
		tF := make([]string, len(txHash))
		for index, vlaue := range txHash {
			tF[index] = vlaue.(string)
		}
		block.TransactionsFalse = tF
	}

	return block, err
}

func (pointer *RequestResult) ToScsBlock() (*ScsBlock, error) {

	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.(map[string]interface{})

	if len(result) == 0 {
		return nil, customerror.EMPTYRESPONSE
	}

	scsBlock := &ScsBlock{}

	marshal, err := json.Marshal(result)

	if err != nil {
		return nil, customerror.UNPARSEABLEINTERFACE
	}

	err = json.Unmarshal([]byte(marshal), scsBlock)

	return scsBlock, err
}

func (pointer *RequestResult) ToSyncingResponse() (*SyncingResponse, error) {

	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	var result map[string]interface{}

	switch (pointer).Result.(type) {
	case bool:
		return &SyncingResponse{}, nil
	case map[string]interface{}:
		result = (pointer).Result.(map[string]interface{})
	default:
		return nil, customerror.UNPARSEABLEINTERFACE
	}

	if len(result) == 0 {
		return nil, customerror.EMPTYRESPONSE
	}

	syncingResponse := &SyncingResponse{}

	marshal, err := json.Marshal(result)

	if err != nil {
		return nil, customerror.UNPARSEABLEINTERFACE
	}

	json.Unmarshal([]byte(marshal), syncingResponse)

	return syncingResponse, nil
}

// To avoid a conversion of a nil interface
func (pointer *RequestResult) checkResponse() error {

	if pointer.Error != nil {
		return errors.New(pointer.Error.Message)
	}
	if pointer.Result == nil {
		return customerror.EMPTYRESPONSE
	}
	return nil
}

func (pointer *RequestResult) ToFiltersLogs() ([]FiltersLogs, error) {

	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.([]interface{})

	if len(result) == 0 {
		return nil, customerror.EMPTYRESPONSE
	}

	tmpLogs := make([]FiltersLogs, len(result))
	var reErr error

	for index, value := range result {
		var log FiltersLogs
		marshal, err := json.Marshal(value)
		if err != nil {
			return nil, customerror.UNPARSEABLEINTERFACE
		}
		err = json.Unmarshal([]byte(marshal), &log)
		if err != nil {
			reErr = err
			break
		}
		tmpLogs[index] = log
	}

	return tmpLogs, reErr
}

func (pointer *RequestResult) ToScsTransactionReceipt() (*ScsTransactionReceipt, error) {

	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.(map[string]interface{})

	if len(result) == 0 {
		return nil, customerror.EMPTYRESPONSE
	}

	scsTrReceipt := &ScsTransactionReceipt{}

	marshal, err := json.Marshal(result)

	if err != nil {
		return nil, customerror.UNPARSEABLEINTERFACE
	}

	err = json.Unmarshal([]byte(marshal), scsTrReceipt)

	return scsTrReceipt, err
}
