// main_test.go
package main

import (
	"fmt"
	"testing"
)

func test() {
	//	strHandler := func(str string) string {
	//		var index int = 64
	//		for ; index < len(str); index++ {
	//			if str[index] == 0 {
	//				break
	//			}
	//		}
	//		return str[64:index]
	//	}
	//	tmpStr := `0x0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000d66696c6573746f726d2d32544200000000000000000000000000000000000000`
	//	data, _ := hex.DecodeString(tmpStr[2:])
	//	fmt.Println(strHandler(string(data)))

	//	txData := new(requestData.TransactionParameters)
	//	txData.To = "0xb735a842b061e62ce41dcaba33ae2764241b6258"
	//	txData.Data = "0x80f01981"
	//	var tmpRpcClient *Chain3Go.RpcClient
	//	tmpRpcClient = Chain3Go.NewRpcClient(vnodeIp, netType)
	//	fmt.Println(tmpRpcClient.Mc().MC_call(txData, "latest"))

	//	fmt.Println(tmpRpcClient.Mc().ScsRPCMethod_GetNonce("0x24e911d31d82f3482dd36451077d6f481da5167d", "0xd58592114ebd97525856929d5c662b72d58b767b"))
	//	txData, _ := tmpRpcClient.Mc().MC_getTransactionReceipt("0xdae1aed6e15148af816812555bca939a4e7c33974ddb79c17138c388cdf6cc63")
	//	fmt.Printf("%#v\n", txData)

	//	fmt.Println(utils.Sha3Hash("scsArray()"))

	//	testSubChain()

	//	searchAddrBalance()

	//	subTransaction()

	//	searchSubChainAddrBalance()

	//	approve()

	//	buyMintToken()

	//approve
	//buyMintToken
	//redeem
	//tx

	//	jsonStr := `{"address":"d58592114ebd97525856929d5c662b72d58b767b","crypto":{"cipher":"aes-128-ctr","ciphertext":"db96d030406419a3ca0d6e6901b3b688ad4c6f34376f048c8ed56f39d1b37169","cipherparams":{"iv":"9e6aee2025bd919866da952d3df40566"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"f22f59cac654619c29efba26fef1c1b894691b0e2e6a2bf903104dae96d824c4"},"mac":"af807312730a542925badcb3830853a43324c8b19b27cbefff7e72f150e6576f"},"id":"d4bfe897-6bcb-4f5d-9166-cc9a4b5ba488","version":3}`
	//	nonce, err := getAddressNonce("0xd58592114ebd97525856929d5c662b72d58b767b")
	//	//	nonce, err := getAddressInSubChainNonce("0x24e911d31d82f3482dd36451077d6f481da5167d", "0xd58592114ebd97525856929d5c662b72d58b767b")
	//	if err == nil {
	//		fmt.Println(currencyConversion("approve", jsonStr, "Wlr7523286", toAddress, erc20Addr, "0xd58592114ebd97525856929d5c662b72d58b767b", "", 20, nonce))
	//		fmt.Println(currencyConversion("buyMintToken", jsonStr, "Wlr7523286", toAddress, erc20Addr, "0xd58592114ebd97525856929d5c662b72d58b767b", "", 19, nonce+1))

	//		//		fmt.Println(currencyConversion("redeem", jsonStr, "Wlr7523286", "0x24e911d31d82f3482dd36451077d6f481da5167d", "0xd609C9B69EFed83F9eD00486B06198B3b3FD5208", "0xd58592114ebd97525856929d5c662b72d58b767b", "", 1, nonce))

	//		//		fmt.Println(currencyConversion("tx", jsonStr, "Wlr7523286", "0x24e911d31d82f3482dd36451077d6f481da5167d", "0xd609C9B69EFed83F9eD00486B06198B3b3FD5208", "0xd58592114ebd97525856929d5c662b72d58b767b", "0x50463586C483D205F1f15741234F6CD2833e1A59", 1, nonce))
	//	}

	//	fmt.Println(getSubChainHeight("0x24e911d31d82f3482dd36451077d6f481da5167d"))
}
