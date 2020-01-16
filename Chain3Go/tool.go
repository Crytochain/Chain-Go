// Chain3Go project tool.go
package Chain3Go

import (
	//	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"strings"

	"Chain3Go/accounts/abi"
	"Chain3Go/accounts/keystore"

	"Chain3Go/lib/common"
	//	LBRMath "Chain3Go/lib/common/math"

	"Chain3Go/lib/crypto"
	"Chain3Go/lib/rlp"
	"Chain3Go/lib/types"

	"github.com/pborman/uuid"
)

//根据keystore字符串获取私钥
func GetPrivateKey(jsonStr string, password string) (*keystore.Key, error) {

	var jsonBytes []byte
	jsonBytes = []byte(jsonStr)

	var err error
	var storeKey *keystore.Key
	storeKey, err = keystore.DecryptKey(jsonBytes, password)

	//	if err == nil {
	//		privateKey := hex.EncodeToString(LBRMath.PaddedBigBytes(storeKey.PrivateKey.D, storeKey.PrivateKey.Params().BitSize/8))
	//		fmt.Println(privateKey)
	//	}

	return storeKey, err
}

//根据私钥获取keystore字符串
func GetKeystoreStr(privateKey, password string) (string, string, error) {

	var err error
	var testKey keystore.Key
	testKey.PrivateKey, err = crypto.HexToECDSA(privateKey)
	if err != nil {
		return "", "", err
	}
	testKey.Address = crypto.PubkeyToAddress(testKey.PrivateKey.PublicKey)
	testKey.Id = uuid.NewRandom()

	scryptN := keystore.StandardScryptN
	scryptP := keystore.StandardScryptP

	var keyJson []byte
	keyJson, err = keystore.EncryptKey(&testKey, password, scryptN, scryptP)

	return string(keyJson), testKey.Address.Hex(), err
}

//创建地址返回keystore字符串，地址
func CreateKeystoreAddress(tradePassword string) (string, string, error) {
	//创建地址

	scryptN := keystore.StandardScryptN
	scryptP := keystore.StandardScryptP
	key, err := keystore.CreateKey()
	if err == nil {
		keyjson, err := keystore.EncryptKey(key, tradePassword, scryptN, scryptP)
		return string(keyjson), key.Address.Hex(), err
	}
	return "", "", err
}

//交易签名返回签名字符串
func TxSign(keystoreStr, keystorePassword, from, to string, amount, gas, gasPrice *big.Int, shardingFlag uint64, data []byte) (string, error) {

	nonce, nonErr := DefaultRpcClient().MC_getTransactionCount(from, "latest") //获取发送地址交易数

	if nonErr == nil {
		viaAddr, viaErr := DefaultRpcClient().VNODE_address() //获取节点via
		if viaErr == nil {
			via := common.HexToAddress(viaAddr)

			tx := types.NewTransaction(uint64(nonce), common.HexToAddress(to), amount, gas, gasPrice, shardingFlag, &via, data)
			priKey, _ := GetPrivateKey(keystoreStr, keystorePassword) //获取私钥

			signTx, signErr := types.SignTx(tx, types.NewPanguSigner(big.NewInt(101)), priKey.PrivateKey) //根据私钥进行签名
			if signErr != nil {
				return "", signErr
			}

			enc, err := rlp.EncodeToBytes(&signTx.TxData) //解析获取交易Hex
			//			fmt.Println(string(enc))
			return fmt.Sprintf("0x%x", enc), err
		}
		return "", viaErr
	}
	return "", nonErr
}

//交易签名返回签名字符串
func SubChainTxSign(viaAddr string, netType int, keystoreStr, keystorePassword, from, to string, amount, gas, gasPrice *big.Int, shardingFlag uint64, data []byte, nonce uint64) (string, error) {

	//	viaAddr, viaErr := DefaultRpcClient().VNODE_address() //获取节点via
	//	if viaErr == nil {
	via := common.HexToAddress(viaAddr)

	tx := types.NewTransaction(nonce, common.HexToAddress(to), amount, gas, gasPrice, shardingFlag, &via, data)
	priKey, priErr := GetPrivateKey(keystoreStr, keystorePassword) //获取私钥
	if priErr != nil {
		return "", priErr
	}

	signTx, signErr := types.SignTx(tx, types.NewPanguSigner(big.NewInt(int64(netType))), priKey.PrivateKey) //根据私钥进行签名
	if signErr != nil {
		return "", signErr
	}

	enc, err := rlp.EncodeToBytes(&signTx.TxData) //解析获取交易Hex
	return fmt.Sprintf("0x%x", enc), err
	//	}
	//	return "", viaErr
}

func AssembleInput(abiStr, name string, args ...interface{}) ([]byte, error) {

	parsed, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		return nil, errors.New("abiStr error")
	}
	return parsed.Pack(name, args...)
}
