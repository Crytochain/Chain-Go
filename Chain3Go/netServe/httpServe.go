// netServe project httpServe.go
package netServe

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"fmt"

	"serverLayer/Chain3Go/utils"
)

//网络请求结构体
type HttpProvider struct {
	httpResponse *http.Response
	address      string //http(https) + ip + 端口
	timeout      int8   //请求超时时间
}

//创建网络请求结构体
func NewHttpProvider(address string, timeout int8) *HttpProvider {

	prvider := new(HttpProvider)
	prvider.address = address
	prvider.timeout = timeout
	return prvider
}

//发送网络请求
func (hProvider *HttpProvider) SendRequest(relust interface{}, method string, params interface{}, netNumber int) error {

	bodyStr := utils.JsonRpc{
		Version: "2.0",
		Method:  method,
		Params:  params,
		Id:      netNumber,
	}

	body := strings.NewReader(bodyStr.AsJsonString())

	//log打印发送的json字符串
	fmt.Printf("【发送】%+v\n", bodyStr)

	req, err := http.NewRequest("POST", hProvider.address, body)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	netCli := &http.Client{
		Timeout: time.Second * time.Duration(hProvider.timeout),
	}

	hProvider.httpResponse, err = netCli.Do(req)
	if err != nil {
		return err
	}

	defer hProvider.httpResponse.Body.Close()
	defer func() {
		hProvider.httpResponse = nil
	}()

	var bodyBytes []byte
	if hProvider.httpResponse.StatusCode == 200 {
		bodyBytes, err = ioutil.ReadAll(hProvider.httpResponse.Body)
		//log打印接收的json字符串
		fmt.Println("【接收】", string(bodyBytes))
		if err != nil {
			return err
		}
	}
	return json.Unmarshal(bodyBytes, relust)
}

//关闭网络请求
func (hProvider *HttpProvider) Close() error {

	if hProvider.httpResponse == nil {
		return errors.New("http Response is nil")
	}
	return hProvider.httpResponse.Body.Close()
}
