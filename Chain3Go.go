// Chain3Go project Chain3Go.go
package Chain3Go

import (
	"fmt"
	//	"math/big"

	"Chain3Go/netServe"
	"Chain3Go/requestData"
	"Chain3Go/types"
)

var (
	rpcClient netServe.ProviderInterface

	ipAddr string = "127.0.0.1:8545" //ip地址
	net    int    = 101              //网络状态 99:测试网 101:正式网
)

type RpcClient struct {
	tmpRpcClient netServe.ProviderInterface
}

// 创建Lchain3.RpcClient结构对象
func NewRpcClient(address string, netNumber int) *RpcClient {

	net = netNumber
	ipAddr = address
	rpcClient = netServe.NewHttpProvider(address, 10)
	return &RpcClient{}
}

func DefaultRpcClient() *RpcClient {

	rpcClient = netServe.NewHttpProvider(ipAddr, 10)
	return &RpcClient{tmpRpcClient: rpcClient}
}

//copy 进行多线程网络请求
func (rpcCli *RpcClient) Mc() *RpcClient {

	tmpArray1 := make([]netServe.ProviderInterface, 1)
	tmpArray2 := make([]netServe.ProviderInterface, 1)
	tmpArray2[0] = rpcClient
	copy(tmpArray1, tmpArray2)
	rpcCli.tmpRpcClient = tmpArray1[0]
	return rpcCli
}

//网络请求统一处理
func (rpcCli *RpcClient) netServeHandler(method string, params interface{}) (*requestData.RequestResult, error) {

	pointer := new(requestData.RequestResult)
	err := rpcCli.tmpRpcClient.SendRequest(pointer, method, params, net)
	return pointer, err
}

/*
 *chain3_clientVersion
 *Returns the current client version.
 *Parameters
 *none
 *Returns
 *String - The current client version
 */
func (rpcCli *RpcClient) CHAIN3_clientVersion() (string, error) {

	pointer, err := rpcCli.netServeHandler(CHAIN3_clientVersion, nil)
	if err != nil {
		return "", err
	}
	return pointer.ToString()
}

/*
 *Returns Keccak-256 (not the standardized SHA3-256) of the given data.
 *Parameters
 *DATA - the data to convert into a SHA3 hash
 *params: [
 *  "0x68656c6c6f20776f726c64"
 *]
 *Returns
 *DATA - The SHA3 result of the given string.
 */
func (rpcCli *RpcClient) CHAIN3_sha3(sourceStr types.ComplexString) (string, error) {

	pointer, err := rpcCli.netServeHandler(CHAIN3_sha3, []string{sourceStr.ToHex()})
	if err != nil {
		return "", err
	}
	return pointer.ToString()
}

/*
 *Returns the current network id.
 *Parameters
 *none
 *Returns
 *String - The current network id.
 *"99": LBR Mainnet
 *"101": LBR Testnet
 *"100": Devnet
 */
func (rpcCli *RpcClient) NET_version() (string, error) {

	pointer, err := rpcCli.netServeHandler(NET_version, nil)
	if err != nil {
		return "", err
	}
	return pointer.ToString()
}

/*
 *net_listening
 *Returns true if client is actively listening for network connections.
 *Parameters
 *none
 *Returns
 *Boolean - true when listening, otherwise false.
 */
func (rpcCli *RpcClient) NET_listening() (bool, error) {

	pointer, err := rpcCli.netServeHandler(NET_listening, nil)
	if err != nil {
		return false, err
	}
	return pointer.ToBoolean()
}

/*
 *net_peerCount
 *Returns number of peers currently connected to the client.
 *Parameters
 *none
 *Returns
 *QUANTITY - integer of the number of connected peers.
 */
func (rpcCli *RpcClient) NET_peerCount() (types.ComplexIntResponse, error) {

	pointer, err := rpcCli.netServeHandler(NET_peerCount, nil)
	if err != nil {
		return "", err
	}
	return pointer.ToComplexIntResponse()
}

/*
 *mc_protocolVersion
 *Returns the current LBR protocol version.
 *Parameters
 *none
 *Returns
 *String - The current LBR protocol version
 */
func (rpcCli *RpcClient) MC_protocolVersion() (int64, error) {

	pointer, err := rpcCli.netServeHandler(MC_protocolVersion, nil)
	if err != nil {
		return 0, err
	}
	return pointer.ToInt()
}

/*
 *mc_syncing
 *Returns an object with data about the sync status or false.
 *Parameters
 *none
 *Returns
 *Object|Boolean, An object with sync status data or FALSE, when not syncing:
 *startingBlock: QUANTITY - The block at which the import started (will only be reset, after the sync reached his head)
 *currentBlock: QUANTITY - The current block, same as mc_blockNumber
 *highestBlock: QUANTITY - The estimated highest block
 */
func (rpcCli *RpcClient) MC_syncing() (bool, *requestData.SyncingResponse, error) {

	pointer, err := rpcCli.netServeHandler(MC_syncing, nil)
	if err != nil {
		return false, nil, err
	}
	reData, err := pointer.ToSyncingResponse()
	if reData.CurrentBlock == "" && reData.HighestBlock == "" && reData.StartingBlock == "" {
		return false, nil, nil
	}
	return true, reData, err
}

/*
 *mc_coinbase
 *Returns the client coinbase address.
 *Parameters
 *none
 *Returns
 *DATA, 20 bytes - the current coinbase address.
 */
func (rpcCli *RpcClient) MC_coinbase() (string, error) {

	pointer, err := rpcCli.netServeHandler(MC_coinbase, nil)
	if err != nil {
		return "", err
	}
	return pointer.ToString()
}

/*
 *mc_mining
 *Returns true if client is actively mining new blocks.
 *Parameters
 *none
 *Returns
 *Boolean - returns true of the client is mining, otherwise false.
 */
func (rpcCli *RpcClient) MC_mining() (bool, error) {

	pointer, err := rpcCli.netServeHandler(MC_mining, nil)
	if err != nil {
		return false, err
	}
	return pointer.ToBoolean()
}

/*
 *mc_hashrate
 *Returns the number of hashes per second that the node is mining with.
 *Parameters
 *none
 *Returns
 *QUANTITY - number of hashes per second.
 */
func (rpcCli *RpcClient) MC_hashrate() (int64, error) {

	pointer, err := rpcCli.netServeHandler(MC_hashrate, nil)
	if err != nil {
		return 0, err
	}
	return pointer.ToInt()
}

/*
 *mc_gasPrice
 *Returns the current price per gas in sha.
 *Parameters
 *none
 *Returns
 *QUANTITY - integer of the current gas price in sha.
 */
func (rpcCli *RpcClient) MC_gasPrice() (int64, error) {

	pointer, err := rpcCli.netServeHandler(MC_gasPrice, nil)
	if err != nil {
		return 0, err
	}
	return pointer.ToInt()
}

/*
 *mc_accounts
 *Returns a list of addresses owned by client.
 *Parameters
 *none
 *Returns
 *Array of DATA, 20 Bytes - addresses owned by the client.
 */
func (rpcCli *RpcClient) MC_accounts() ([]string, error) {

	pointer, err := rpcCli.netServeHandler(MC_accounts, nil)
	if err != nil {
		return []string{}, err
	}
	return pointer.ToStringArray()
}

/*
 *mc_blockNumber
 *Returns the number of most recent block.
 *Parameters
 *none
 *Returns
 *QUANTITY - integer of the current block number the client is on.
 */
func (rpcCli *RpcClient) MC_blockNumber() (int64, error) {

	pointer, err := rpcCli.netServeHandler(MC_blockNumber, nil)
	if err != nil {
		return 0, err
	}
	return pointer.ToInt()
}

/*
 *mc_getBalance
 *Returns the balance of the account of given address.
 *Parameters
 *DATA, 20 Bytes - address to check for balance.
 *QUANTITY|TAG - integer block number, or the string "latest", "earliest" or "pending"
 *Returns
 *QUANTITY - integer of the current balance in sha.
 */
func (rpcCli *RpcClient) MC_getBalance(address, number string) (string, error) {

	pointer, err := rpcCli.netServeHandler(MC_getBalance, []string{address, number})
	if err != nil {
		return "", err
	}
	return pointer.ToString()
}

/*
 *mc_getStorageAt
 *Returns the value from a storage position at a given address.
 *Parameters
 *DATA, 20 Bytes - address of the storage.
 *QUANTITY - integer of the position in the storage.
 *QUANTITY|TAG - integer block number, or the string "latest", "earliest" or "pending"
 *Returns
 *DATA - the value at this storage position.
 */
func (rpcCli *RpcClient) MC_getStorageAt(address, dataStr, number string) (string, error) {

	pointer, err := rpcCli.netServeHandler(MC_getStorageAt, []string{address, dataStr, number})
	if err != nil {
		return "", err
	}
	return pointer.ToString()
}

/*
 *mc_getTransactionCount
 *Returns the number of transactions sent from an address.
 *Parameters
 *DATA, 20 Bytes - address.
 *QUANTITY|TAG - integer block number, or the string "latest", "earliest" or "pending"
 *params: [
 *   '0x407d73d8a49eeb85d32cf465507dd71d507100c1',
 *   'latest' // state at the latest block
 *]
 *Returns
 *QUANTITY - integer of the number of transactions send from this address.
 */
func (rpcCli *RpcClient) MC_getTransactionCount(address, number string) (int64, error) {

	pointer, err := rpcCli.netServeHandler(MC_getTransactionCount, []string{address, number})
	if err != nil {
		return 0, err
	}
	return pointer.ToInt()
}

/*
 *mc_getBlockTransactionCountByHash
 *Returns the number of transactions in a block from a block matching the given block hash.
 *Parameters
 *DATA, 32 Bytes - hash of a block
 *params: [
 *   '0xb903239f8543d04b5dc1ba6579132b143087c68db1b2168786408fcbce568238'
 *]
 *Returns
 *QUANTITY - integer of the number of transactions in this block.
 */
func (rpcCli *RpcClient) MC_getBlockTransactionCountByHash(blockHash string) (int64, error) {

	pointer, err := rpcCli.netServeHandler(MC_getBlockTransactionCountByHash, []string{blockHash})
	if err != nil {
		return 0, err
	}
	return pointer.ToInt()
}

/*
 *mc_getBlockTransactionCountByNumber
 *Returns the number of transactions in a block matching the given block number.
 *Parameters
 *QUANTITY|TAG - integer of a block number, or the string "earliest", "latest" or "pending"
 *params: [
 *   '0xe8', // 232
 *]
 *Returns
 *QUANTITY - integer of the number of transactions in this block.
 */
func (rpcCli *RpcClient) MC_getBlockTransactionCountByNumber(number string) (int64, error) {

	pointer, err := rpcCli.netServeHandler(MC_getBlockTransactionCountByNumber, []string{number})
	if err != nil {
		return 0, err
	}
	return pointer.ToInt()
}

/*
 *mc_getUncleCountByBlockHash
 *Returns the number of uncles in a block from a block matching the given block hash.
 *Parameters
 *DATA, 32 Bytes - hash of a block
 *params: [
 *   '0xb903239f8543d04b5dc1ba6579132b143087c68db1b2168786408fcbce568238'
 *]
 *Returns
 *QUANTITY - integer of the number of uncles in this block.
 */
func (rpcCli *RpcClient) MC_getUncleCountByBlockHash(blockHash string) (int64, error) {

	pointer, err := rpcCli.netServeHandler(MC_getUncleCountByBlockHash, []string{blockHash})
	if err != nil {
		return 0, err
	}
	return pointer.ToInt()
}

/*mc_getUncleCountByBlockNumber
 *Returns the number of uncles in a block from a block matching the given block number.
 *Parameters
 *QUANTITY|TAG - integer of a block number, or the string "latest", "earliest" or "pending"
 *params: [
 *   '0xe8', // 232
 *]
 *Returns
 *QUANTITY - integer of the number of uncles in this block.
 */
func (rpcCli *RpcClient) MC_getUncleCountByBlockNumber(number string) (int64, error) {

	pointer, err := rpcCli.netServeHandler(MC_getUncleCountByBlockNumber, []string{number})
	if err != nil {
		return 0, err
	}
	return pointer.ToInt()
}

/*
 *mc_getCode
 *Returns code at a given address.
 *Parameters
 *DATA, 20 Bytes - address
 *QUANTITY|TAG - integer block number, or the string "latest", "earliest" or "pending"
 *params: [
 *   '0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b',
 *   '0x2'  // 2
 *]
 *Returns
 *DATA - the code from the given address.
 */
func (rpcCli *RpcClient) MC_getCode(address, number string) (string, error) {

	pointer, err := rpcCli.netServeHandler(MC_getCode, []string{address, number})
	if err != nil {
		return "", err
	}
	return pointer.ToString()
}

/*
 *mc_sign
 *The sign method calculates an LBR specific signature with: sign(keccak256("\x19LBR Signed Message:\n" + len(message) + message))).
 *By adding a prefix to the message makes the calculated signature recognisable as an LBR specific signature. This prevents misuse where a malicious DApp can sign arbitrary data (e.g. transaction) and use the signature to impersonate the victim.
 *Note the address to sign with must be unlocked.
 *Parameters
 *account, message
 *DATA, 20 Bytes - address
 *DATA, N Bytes - message to sign
 *Returns
 *DATA: Signature
 */
func (rpcCli *RpcClient) MC_sign(address, dataStr string) (string, error) {

	pointer, err := rpcCli.netServeHandler(MC_sign, []string{address, dataStr})
	if err != nil {
		return "", err
	}
	return pointer.ToString()
}

/*
 *mc_sendTransaction
 *Creates new message call transaction or a contract creation, if the data field contains code.
 *Note the address to sign with must be unlocked.
 *Parameters
 *Object - The transaction object
 *from: DATA, 20 Bytes - The address the transaction is send from.
 *to: DATA, 20 Bytes - (optional when creating new contract) The address the transaction is directed to.
 *gas: QUANTITY - (optional, default: 90000) Integer of the gas provided for the transaction execution. It will return unused gas.
 *gasPrice: QUANTITY - (optional, default: To-Be-Determined) Integer of the gasPrice used for each paid gas
 *value: QUANTITY - (optional) Integer of the value sent with this transaction
 *data: DATA - The compiled code of a contract OR the hash of the invoked method signature and encoded parameters. For details see Ethereum Contract ABI
 *nonce: QUANTITY - (optional) Integer of a nonce. This allows to overwrite your own pending transactions that use the same nonce.
 *params: [{
 *  "from": "0xb60e8dd61c5d32be8058bb8eb970870f07233155",
 *  "to": "0xd46e8dd67c5d32be8058bb8eb970870f07244567",
 *  "gas": "0x76c0", // 30400
 *  "gasPrice": "0x9184e72a000", // 10000000000000
 *  "value": "0x9184e72a", // 2441406250
 *  "data": "0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675"
 *}]
 *Returns
 *DATA, 32 Bytes - the transaction hash, or the zero hash if the transaction is not yet available.
 *Use mc_getTransactionReceipt to get the contract address, after the transaction was mined, when you created a contract.
 */
func (rpcCli *RpcClient) MC_sendTransaction(txData *requestData.TransactionParameters) (string, error) {

	pointer, err := rpcCli.netServeHandler(MC_sendTransaction, []interface{}{txData.Transform()})
	if err != nil {
		return "", err
	}
	return pointer.ToString()
}

/*
 *mc_sendRawTransaction
 *Creates new message call transaction or a contract creation for signed transactions.
 *Parameters
 *DATA, The signed transaction data.
 *params: ["0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675"]
 *Returns
 *DATA, 32 Bytes - the transaction hash, or the zero hash if the transaction is not yet available.
 */
func (rpcCli *RpcClient) MC_sendRawTransaction(signTxStr string) (string, error) {

	pointer, err := rpcCli.netServeHandler(MC_sendRawTransaction, []string{signTxStr})
	if err != nil {
		return "", err
	}
	return pointer.ToString()
}

/*
 *mc_call
 *Executes a new message call immediately without creating a transaction on the block chain.
 *Parameters
 *Object - The transaction call object
 *from: DATA, 20 Bytes - (optional) The address the transaction is sent from.
 *to: DATA, 20 Bytes - The address the transaction is directed to.
 *gas: QUANTITY - (optional) Integer of the gas provided for the transaction execution. mc_call consumes zero gas, but this parameter may be needed by some executions.
 *gasPrice: QUANTITY - (optional) Integer of the gasPrice used for each paid gas
 *value: QUANTITY - (optional) Integer of the value sent with this transaction
 *data: DATA - (optional) Hash of the method signature and encoded parameters. For details see Ethereum Contract ABI
 *QUANTITY|TAG - integer block number, or the string "latest", "earliest" or "pending"
 *Returns
 *DATA - the return value of executed contract.
 */
func (rpcCli *RpcClient) MC_call(txData *requestData.TransactionParameters, number string) (string, error) {

	pointer, err := rpcCli.netServeHandler(MC_call, []interface{}{requestData.Struct2Map(*(txData.Transform())), number})
	if err != nil {
		return "", err
	}
	return pointer.ToString()
}

/*
 *mc_estimateGas
 *Generates and returns an estimate of how much gas is necessary to allow the transaction to complete. The transaction will not be added to the blockchain. Note that the estimate may be significantly more than the amount of gas actually used by the transaction, for a variety of reasons including EVM mechanics and node performance.
 *Parameters
 *See mc_call parameters, expect that all properties are optional. If no gas limit is specified LBR uses the block gas limit from the pending block as an upper bound. As a result the returned estimate might not be enough to executed the call/transaction when the amount of gas is higher than the pending block gas limit.
 *Returns
 *QUANTITY - the amount of gas used.
 */
func (rpcCli *RpcClient) MC_estimateGas(txData *requestData.TransactionParameters) (int64, error) {

	pointer, err := rpcCli.netServeHandler(MC_estimateGas, []interface{}{txData.Transform()})
	if err != nil {
		return 0, err
	}
	return pointer.ToInt()
}

/*
 *mc_getBlockByHash
 *Returns information about a block by hash.
 *Parameters
 *DATA, 32 Bytes - Hash of a block.
 *Boolean - If true it returns the full transaction objects, if false only the hashes of the transactions.
 *params: [
 *   '0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331',
 *   true
 *]
 *Returns
 *Object - A block object, or null when no block was found:
 *number: QUANTITY - the block number. null when its pending block.
 *hash: DATA, 32 Bytes - hash of the block. null when its pending block.
 *parentHash: DATA, 32 Bytes - hash of the parent block.
 *nonce: DATA, 8 Bytes - hash of the generated proof-of-work. null when its pending block.
 *sha3Uncles: DATA, 32 Bytes - SHA3 of the uncles data in the block.
 *logsBloom: DATA, 256 Bytes - the bloom filter for the logs of the block. null when its pending block.
 *transactionsRoot: DATA, 32 Bytes - the root of the transaction trie of the block.
 *stateRoot: DATA, 32 Bytes - the root of the final state trie of the block.
 *receiptsRoot: DATA, 32 Bytes - the root of the receipts trie of the block.
 *miner: DATA, 20 Bytes - the address of the beneficiary to whom the mining rewards were given.
 *difficulty: QUANTITY - integer of the difficulty for this block.
 *totalDifficulty: QUANTITY - integer of the total difficulty of the chain until this block.
 *extraData: DATA - the "extra data" field of this block.
 *size: QUANTITY - integer the size of this block in bytes.
 *gasLimit: QUANTITY - the maximum gas allowed in this block.
 *gasUsed: QUANTITY - the total used gas by all transactions in this block.
 *timestamp: QUANTITY - the unix timestamp for when the block was collated.
 *transactions: Array - Array of transaction objects, or 32 Bytes transaction hashes depending on the last given parameter.
 *uncles: Array - Array of uncle hashes.
 */
func (rpcCli *RpcClient) MC_getBlockByHash(blockHash string, objFlag bool) (*requestData.Block, error) {

	pointer, err := rpcCli.netServeHandler(MC_getBlockByHash, []interface{}{blockHash, objFlag})
	if err != nil {
		return nil, err
	}

	return pointer.ToBlock(objFlag)
}

/*
 *mc_getBlockByNumber
 *Returns information about a block by block number.
 *Parameters
 *QUANTITY|TAG - integer of a block number, or the string "earliest", "latest" or "pending", as in the default block parameter.
 *Boolean - If true it returns the full transaction objects, if false only the hashes of the transactions.
 *params: [
 *   '0x1b4', // 436
 *   true
 *]
 *Returns
 *See mc_getBlockByHash
 */
func (rpcCli *RpcClient) MC_getBlockByNumber(number types.ComplexIntParameter, objFlag bool) (*requestData.Block, error) {

	pointer, err := rpcCli.netServeHandler(MC_getBlockByNumber, []interface{}{number.ToHex(), objFlag})
	if err != nil {
		return nil, err
	}

	return pointer.ToBlock(objFlag)
}

/*
 *mc_getTransactionByHash
 *Returns the information about a transaction requested by transaction hash.
 *Parameters
 *DATA, 32 Bytes - hash of a transaction
 *params: [
 *   "0xb903239f8543d04b5dc1ba6579132b143087c68db1b2168786408fcbce568238"
 *]
 *Returns
 *Object - A transaction object, or null when no transaction was found:
 *hash: DATA, 32 Bytes - hash of the transaction.
 *nonce: QUANTITY - the number of transactions made by the sender prior to this one.
 *blockHash: DATA, 32 Bytes - hash of the block where this transaction was in. null when its pending.
 *blockNumber: QUANTITY - block number where this transaction was in. null when its pending.
 *transactionIndex: QUANTITY - integer of the transactions index position in the block. null when its pending.
 *from: DATA, 20 Bytes - address of the sender.
 *to: DATA, 20 Bytes - address of the receiver. null when its a contract creation transaction.
 *value: QUANTITY - value transferred in Wei.
 *gasPrice: QUANTITY - gas price provided by the sender in Wei.
 *gas: QUANTITY - gas provided by the sender.
 *input: DATA - the data send along with the transaction.
 */
func (rpcCli *RpcClient) MC_getTransactionByHash(hash string) (*requestData.TransactionResponse, error) {

	pointer, err := rpcCli.netServeHandler(MC_getTransactionByHash, []string{hash})
	if err != nil {
		return nil, err
	}

	return pointer.ToTransactionResponse()
}

/*
 *mc_getTransactionByBlockHashAndIndex
 *Returns information about a transaction by block hash and transaction index position.
 *Parameters
 *DATA, 32 Bytes - hash of a block.
 *QUANTITY - integer of the transaction index position.
 *params: [
 *   '0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331',
 *   '0x0' // 0
 *]
 *Returns
 *See mc_getTransactionByHash
 */
func (rpcCli *RpcClient) MC_getTransactionByBlockHashAndIndex(blockHash string, index types.ComplexIntParameter) (*requestData.TransactionResponse, error) {

	pointer, err := rpcCli.netServeHandler(MC_getTransactionByBlockHashAndIndex, []string{blockHash, index.ToHex()})
	if err != nil {
		return nil, err
	}

	return pointer.ToTransactionResponse()
}

/*
 *mc_getTransactionByBlockNumberAndIndex
 *Returns information about a transaction by block number and transaction index position.
 *Parameters
 *QUANTITY|TAG - a block number, or the string "earliest", "latest" or "pending", as in the default block parameter.
 *QUANTITY - the transaction index position.
 *params: [
 *   '0x29c', // 668
 *   '0x0' // 0
 *]
 *Returns
 *See mc_getTransactionByHash
 */
func (rpcCli *RpcClient) MC_getTransactionByBlockNumberAndIndex(number types.ComplexIntParameter, index types.ComplexIntParameter) (*requestData.TransactionResponse, error) {

	pointer, err := rpcCli.netServeHandler(MC_getTransactionByBlockNumberAndIndex, []string{number.ToHex(), index.ToHex()})
	if err != nil {
		return nil, err
	}

	return pointer.ToTransactionResponse()
}

/*
 *mc_getTransactionReceipt
 *Returns the receipt of a transaction by transaction hash.
 *Note That the receipt is not available for pending transactions.
 *Parameters
 *DATA, 32 Bytes - hash of a transaction
 *params: [
 *   '0x7bb694c3462764cb113e9b742faaf06adc728e70b607f8b7aa95207ee32b1c5e'
 *]
 *Returns
 *Object - A transaction receipt object, or null when no receipt was found:
 *transactionHash : DATA, 32 Bytes - hash of the transaction.
 *transactionIndex: QUANTITY - integer of the transactions index position in the block.
 *blockHash: DATA, 32 Bytes - hash of the block where this transaction was in.
 *blockNumber: QUANTITY - block number where this transaction was in.
 *cumulativeGasUsed : QUANTITY - The total amount of gas used when this transaction was executed in the block.
 *gasUsed : QUANTITY - The amount of gas used by this specific transaction alone.
 *contractAddress : DATA, 20 Bytes - The contract address created, if the transaction was a contract creation, otherwise null.
 *logs: Array - Array of log objects, which this transaction generated.
 *logsBloom: DATA, 256 Bytes - Bloom filter for light clients to quickly retrieve related logs.
 *It also returns either :
 *root : DATA 32 bytes of post-transaction stateroot (pre Byzantium)
 *status: QUANTITY either 1 (success) or 0 (failure)
 */
func (rpcCli *RpcClient) MC_getTransactionReceipt(hash string) (*requestData.TransactionReceipt, error) {

	pointer, err := rpcCli.netServeHandler(MC_getTransactionReceipt, []string{hash})
	if err != nil {
		return nil, err
	}

	return pointer.ToTransactionReceipt()
}

/*
 *mc_getUncleByBlockHashAndIndex
 *Returns information about a uncle of a block by hash and uncle index position.
 *Parameters
 *DATA, 32 Bytes - hash a block.
 *QUANTITY - the uncle's index position.
 *params: [
 *   '0xc6ef2fc5426d6ad6fd9e2a26abeab0aa2411b7ab17f30a99d3cb96aed1d1055b',
 *   '0x0' // 0
 *]
 *Returns
 *See mc_getBlockByHash
 *Note: An uncle doesn't contain individual transactions.
 */
func (rpcCli *RpcClient) MC_getUncleByBlockHashAndIndex(blockHash string, index types.ComplexIntParameter) (*requestData.Block, error) {

	pointer, err := rpcCli.netServeHandler(MC_getUncleByBlockHashAndIndex, []interface{}{blockHash, index.ToHex()})
	if err != nil {
		return nil, err
	}

	return pointer.ToBlock(true)
}

/*
 *mc_getUncleByBlockNumberAndIndex
 *Returns information about a uncle of a block by number and uncle index position.
 *Parameters
 *QUANTITY|TAG - a block number, or the string "earliest", "latest" or "pending", as in the default block parameter.
 *QUANTITY - the uncle's index position.
 *params: [
 *   '0x29c', // 668
 *   '0x0' // 0
 *]
 *Returns
 *See mc_getBlockByHash
 *Note: An uncle doesn't contain individual transactions.
 */
func (rpcCli *RpcClient) MC_getUncleByBlockNumberAndIndex(number types.ComplexString, index types.ComplexIntParameter) (*requestData.TransactionResponse, error) {

	pointer, err := rpcCli.netServeHandler(MC_getUncleByBlockNumberAndIndex, []string{number.ToHex(), index.ToHex()})
	if err != nil {
		return nil, err
	}

	return pointer.ToTransactionResponse()
}

/*
 *mc_newFilter
 *Creates a filter object, based on filter options, to notify when the state changes (logs). To check if the state has changed, call mc_getFilterChanges.
 *A note on specifying topic filters:
 *Topics are order-dependent. A transaction with a log with topics [A, B] will be matched by the following topic filters:
 *[] "anything"
 *[A] "A in first position (and anything after)"
 *[null, B] "anything in first position AND B in second position (and anything after)"
 *[A, B] "A in first position AND B in second position (and anything after)"
 *[[A, B], [A, B]] "(A OR B) in first position AND (A OR B) in second position (and anything after)"
 *Parameters
 *Object - The filter options:
 *fromBlock: QUANTITY|TAG - (optional, default: "latest") Integer block number, or "latest" for the last mined block or "pending", "earliest" for not yet mined transactions.
 *toBlock: QUANTITY|TAG - (optional, default: "latest") Integer block number, or "latest" for the last mined block or "pending", "earliest" for not yet mined transactions.
 *address: DATA|Array, 20 Bytes - (optional) Contract address or a list of addresses from which logs should originate.
 *topics: Array of DATA, - (optional) Array of 32 Bytes DATA topics. Topics are order-dependent. Each topic can also be an array of DATA with "or" options.
 */
func (rpcCli *RpcClient) MC_newFilter(fromBlock, toBlock, address string, dataAray []interface{}) (int64, error) {

	pointer, err := rpcCli.netServeHandler(MC_newFilter, []interface{}{fromBlock, toBlock, address, dataAray})
	if err != nil {
		return 0, err
	}

	return pointer.ToInt()
}

/*
 *mc_newBlockFilter
 *Creates a filter in the node, to notify when a new block arrives. To check if the state has changed, call mc_getFilterChanges.
 *Parameters
 *None
 *Returns
 *QUANTITY - A filter id.
 */
func (rpcCli *RpcClient) MC_newBlockFilter() (string, error) {

	pointer, err := rpcCli.netServeHandler(MC_newBlockFilter, nil)
	if err != nil {
		return "", err
	}
	return pointer.ToString()
}

/*
 *mc_newPendingTransactionFilter
 *Creates a filter in the node, to notify when new pending transactions arrive. To check if the state has changed, call mc_getFilterChanges.
 *Parameters
 *None
 *Returns
 *QUANTITY - A filter id.
 */
func (rpcCli *RpcClient) MC_newPendingTransactionFilter() (string, error) {

	pointer, err := rpcCli.netServeHandler(MC_newPendingTransactionFilter, nil)
	if err != nil {
		return "", err
	}
	return pointer.ToString()
}

/*
 *mc_uninstallFilter
 *Uninstalls a filter with given id. Should always be called when watch is no longer needed. Additonally Filters timeout when they aren't requested with mc_getFilterChanges for a period of time.
 *Parameters
 *QUANTITY - The filter id.
 *params: [
 *  "0xb" // 11
 *]
 *Returns
 *Boolean - true if the filter was successfully uninstalled, otherwise false.
 */
func (rpcCli *RpcClient) MC_uninstallFilter(filterId string) (bool, error) {

	pointer, err := rpcCli.netServeHandler(MC_uninstallFilter, []string{filterId})
	if err != nil {
		return false, err
	}
	return pointer.ToBoolean()
}

/*
 *mc_getFilterChanges
 *Polling method for a filter, which returns an array of logs which occurred since last poll.
 *Parameters
 *QUANTITY - the filter id.
 *params: [
 *  "0x16" // 22
 *]
 *Returns
 *Array - Array of log objects, or an empty array if nothing has changed since last poll.
 *For filters created with mc_newBlockFilter the return are block hashes (DATA, 32 Bytes), e.g. ["0x3454645634534..."].
 *For filters created with mc_newPendingTransactionFilter the return are transaction hashes (DATA, 32 Bytes), e.g. ["0x6345343454645..."].
 *For filters created with mc_newFilter logs are objects with following params:
 *removed: TAG - true when the log was removed, due to a chain reorganization. false if its a valid log.
 *logIndex: QUANTITY - integer of the log index position in the block. null when its pending log.
 *transactionIndex: QUANTITY - integer of the transactions index position log was created from. null when its pending log.
 *transactionHash: DATA, 32 Bytes - hash of the transactions this log was created from. null when its pending log.
 *blockHash: DATA, 32 Bytes - hash of the block where this log was in. null when its pending. null when its pending log.
 *blockNumber: QUANTITY - the block number where this log was in. null when its pending. null when its pending log.
 *address: DATA, 20 Bytes - address from which this log originated.
 *data: DATA - contains one or more 32 Bytes non-indexed arguments of the log.
 *topics: Array of DATA - Array of 0 to 4 32 Bytes DATA of indexed log arguments. (In solidity: The first topic is the hash of the signature of the event (e.g. Deposit(address,bytes32,uint256)), except you declared the event with the anonymous specifier.)
 */
func (rpcCli *RpcClient) MC_getFilterChanges(filterId string) ([]requestData.FiltersLogs, error) {

	pointer, err := rpcCli.netServeHandler(MC_getFilterChanges, []string{filterId})
	if err != nil {
		return nil, err
	}
	return pointer.ToFiltersLogs()
}

/*
 *mc_getFilterLogs
 *Returns an array of all logs matching filter with given id.
 *Parameters
 *QUANTITY - The filter id.
 *params: [
 *  "0x16" // 22
 *]
 *Returns
 *See mc_getFilterChanges
 */
func (rpcCli *RpcClient) MC_getFilterLogs(filterId string) ([]requestData.FiltersLogs, error) {

	pointer, err := rpcCli.netServeHandler(MC_getFilterLogs, []string{filterId})
	if err != nil {
		return nil, err
	}
	return pointer.ToFiltersLogs()
}

/*
 *mc_getLogs
 *Returns an array of all logs matching a given filter object.
 *Parameters
 *Object - the filter object, see mc_newFilter parameters.
 *params: [{
 *  "topics": ["0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b"]
 *}]
 *Returns
 *See mc_getFilterChanges
 */
func (rpcCli *RpcClient) MC_getLogs(topics []string) ([]requestData.FiltersLogs, error) {

	pointer, err := rpcCli.netServeHandler(MC_getLogs, []map[string]([]string){map[string]([]string){"topics": topics}})
	if err != nil {
		return nil, err
	}

	return pointer.ToFiltersLogs()
}

/*
 *mc_getWork
 *Returns the hash of the current block, the seedHash, and the boundary condition to be met ("target").
 *Parameters
 *none
 *Returns
 *Array - Array with the following properties:
 *DATA, 32 Bytes - current block header pow-hash
 *DATA, 32 Bytes - the seed hash used for the DAG.
 *DATA, 32 Bytes - the boundary condition ("target"), 2^256 / difficulty.
 */
func (rpcCli *RpcClient) MC_getWork() ([]string, error) {

	pointer, err := rpcCli.netServeHandler(MC_getWork, nil)
	if err != nil {
		return []string{}, err
	}

	return pointer.ToStringArray()
}

/*
 *mc_submitWork
 *Used for submitting a proof-of-work solution.
 *Parameters
 *DATA, 8 Bytes - The nonce found (64 bits)
 *DATA, 32 Bytes - The header's pow-hash (256 bits)
 *DATA, 32 Bytes - The mix digest (256 bits)
 *params: [
 *  "0x0000000000000001",
 *  "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
 *  "0xD1FE5700000000000000000000000000D1FE5700000000000000000000000000"
 *]
 *Returns
 *Boolean - returns true if the provided solution is valid, otherwise false.
 */
func (rpcCli *RpcClient) MC_submitWork(paramsData []string) (bool, error) {

	pointer, err := rpcCli.netServeHandler(MC_submitWork, paramsData)
	if err != nil {
		return false, err
	}

	return pointer.ToBoolean()
}

/*
 *vnode_address
 *Returns the VNODE benificial address. This is needed for SCS to use in the communication.
 *Parameters
 *none
 *Returns
 *address: DATA, 20 Bytes - address from which the VNODE settings in the vnodeconfig.json file.
 */
func (rpcCli *RpcClient) VNODE_address() (string, error) {

	pointer, err := rpcCli.netServeHandler(VNODE_address, nil)
	if err != nil {
		return "", err
	}

	return pointer.ToString()
}

/*
 *vnode_scsService
 *Returns if the VNODE enable the service for SCS servers.
 *Parameters
 *none
 *Returns
 *Bool - true, enable the SCS service; false, not open.
 */
func (rpcCli *RpcClient) VNODE_scsService() (bool, error) {

	pointer, err := rpcCli.netServeHandler(VNODE_scsService, nil)
	if err != nil {
		return false, err
	}

	return pointer.ToBoolean()
}

/*
 *vnode_serviceCfg
 *Returns the VNODE SCS service ports for connecting with.
 *Parameters
 *none
 *Returns
 *String - The current service port set in the vnodeconfig.json.
 */
func (rpcCli *RpcClient) VNODE_serviceCfg() (string, error) {

	pointer, err := rpcCli.netServeHandler(VNODE_serviceCfg, nil)
	if err != nil {
		return "", err
	}

	return pointer.ToString()
}

/*
 *vnode_showToPublic
 *Returns if the VNODE enable the public view.
 *Parameters
 *none
 *Returns
 *Bool - true, open to the public; false, not open.
 */
func (rpcCli *RpcClient) VNODE_showToPublic() (bool, error) {

	pointer, err := rpcCli.netServeHandler(VNODE_showToPublic, nil)
	if err != nil {
		return false, err
	}

	return pointer.ToBoolean()
}

/*
 *vnode_vnodeIP
 *Returns VNODE IP for users to access. Note for cloud servers, this could be different from the cloud server IP.
 *Parameters
 *none
 *Returns
 *String - The current IP that can be used to access. This is set in the vnodeconfig.json.
 */
func (rpcCli *RpcClient) VNODE_vnodeIP() (string, error) {

	pointer, err := rpcCli.netServeHandler(VNODE_vnodeIP, nil)
	if err != nil {
		return "", err
	}

	return pointer.ToString()
}

/*
 *scs_directCall
 *Executes a new constant call of the MicroChain Dapp function without creating a transaction on the MicroChain. This RPC call is used by API/lib to call MicroChain Dapp functions.
 *Parameters
 *Object - The transaction call object
 *from: DATA, 20 Bytes - (optional) The address the transaction is sent from.
 *to: DATA, 20 Bytes - The address the transaction is directed to. This parameter is the MicroChain address.
 *data: DATA - (optional) Hash of the method signature and encoded parameters.
 *Returns
 *DATA - the return value of executed Dapp constant function call.
 */
func (rpcCli *RpcClient) SCS_directCall(txData *requestData.TransactionParameters) (string, error) {

	pointer, err := rpcCli.netServeHandler(SCS_directCall, []interface{}{requestData.Struct2Map(*(txData.Transform()))})
	if err != nil {
		return "", err
	}

	return pointer.ToString()
}

/*
 *scs_getBlock
 *Returns information about a MicroChain block by block number.
 *Parameters
 *String - the address of the MicroChain that Dapp is on.
 *QUANTITY|TAG - integer of a block number, or the string "earliest" or "latest", as in the default block parameter. Note, scs_getBlock does not support "pending".
 *Returns
 *OBJ - object of the block on the MicroChain.
 */
func (rpcCli *RpcClient) SCS_getBlock(dappAddress, number string) (*requestData.ScsBlock, error) {

	pointer, err := rpcCli.netServeHandler(SCS_getBlock, []string{dappAddress, number})
	if err != nil {
		return nil, err
	}

	return pointer.ToScsBlock()
}

/*
 *scs_getBlockNumber
 *Returns the number of most recent block .
 *Parameters
 *String - the address of the MicroChain that Dapp is on.
 *Returns
 *NUMBER - integer of the current block number the client is on.
 */
func (rpcCli *RpcClient) SCS_getBlockNumber(dappAddress string) (int64, error) {

	pointer, err := rpcCli.netServeHandler(SCS_getBlockNumber, []string{dappAddress})
	if err != nil {
		return 0, err
	}

	return pointer.ToInt()
}

/*
 *scs_getDappState
 *Returns the Dapp state on the MicroChain.
 *Parameters
 *String - the address of the MicroChain that Dapp is on.
 *Returns
 *Number - 0, no DAPP is deployed on the MicroChain; 1, DAPP is deployed on the MicroChain.
 */
func (rpcCli *RpcClient) SCS_getDappState(dappAddress string) (float64, error) {

	pointer, err := rpcCli.netServeHandler(SCS_getDappState, []string{dappAddress})
	if err != nil {
		return 0, err
	}

	return pointer.Result.(float64), nil
}

/*
 *scs_getMicroChainList
 *Returns the Dapp state on the MicroChain.
 *Parameters
 *None
 *Returns
 *Array - A list of Micro Chain addresses on the SCS.
 */
func (rpcCli *RpcClient) SCS_getMicroChainList() ([]string, error) {

	pointer, err := rpcCli.netServeHandler(SCS_getMicroChainList, nil)
	if err != nil {
		return []string{}, err
	}

	return pointer.ToStringArray()
}

/*
 *scs_getNonce
 *Returns the Dapp state on the MicroChain.
 *Parameters
 *String - the address of the MicroChain that Dapp is on.
 *String - the address of the accountn.
 *Returns
 *Array - 0, no DAPP is deployed on the MicroChain; 1, DAPP is deployed on the MicroChain.
 */
func (rpcCli *RpcClient) SCS_getNonce(dappAddress, address string) (float64, error) {

	pointer, err := rpcCli.netServeHandler(SCS_getNonce, []string{dappAddress, address})
	fmt.Println(pointer.Result)
	if err != nil {
		return 0, err
	}
	return pointer.Result.(float64), nil
}

/*
 *scs_getSCSId
 *Returns the SCS id.
 *Parameters
 *None
 *Returns
 *String - SCS id in the scskeystore directory, used for SCS identification to send deposit and receive MicroChain mining rewards.
 */
func (rpcCli *RpcClient) SCS_getSCSId() (string, error) {

	pointer, err := rpcCli.netServeHandler(SCS_getSCSId, nil)
	if err != nil {
		return "", err
	}

	return pointer.ToString()
}

/*
 *scs_getTransactionReceipt
 *Returns the receipt of a transaction by transaction hash. Note That the receipt is not available for pending transactions.
 *Parameters
 *String - The MicroChain address. String - The transaction hash. Function - (optional) If you pass a callback the HTTP request is made asynchronous.
 *Returns
 *Object - A transaction receipt object, or null when no receipt was found:.
 */
func (rpcCli *RpcClient) SCS_getTransactionReceipt(scsAddress, hash string) (*requestData.ScsTransactionReceipt, error) {

	pointer, err := rpcCli.netServeHandler(SCS_getTransactionReceipt, []string{scsAddress, hash})
	if err != nil {
		return nil, err
	}
	//	fmt.Println(pointer.Result)
	return pointer.ToScsTransactionReceipt()
}

func (rpcCli *RpcClient) PERSONAL_unlockAccount(address, password string) (bool, error) {

	pointer, err := rpcCli.netServeHandler(PERSONAL_unlockAccount, []string{address, password})
	if err != nil {
		return false, err
	}
	return pointer.ToBoolean()
}

func (rpcCli *RpcClient) PERSONAL_lockAccount(address string) (bool, error) {

	pointer, err := rpcCli.netServeHandler(PERSONAL_lockAccount, []string{address})
	if err != nil {
		return false, err
	}
	return pointer.ToBoolean()
}

func (rpcCli *RpcClient) PERSONAL_newAccount(password string) (string, error) {

	pointer, err := rpcCli.netServeHandler(PERSONAL_newAccount, []string{password})
	if err != nil {
		return "", err
	}
	return pointer.ToString()
}

func (rpcCli *RpcClient) PERSONAL_sign(data, address, password string) (string, error) {

	pointer, err := rpcCli.netServeHandler(PERSONAL_sign, []string{data, address, password})
	if err != nil {
		return "", err
	}
	return pointer.ToString()
}

func (rpcCli *RpcClient) PERSONAL_ecRecover(signData, signature string) (string, error) {

	pointer, err := rpcCli.netServeHandler(PERSONAL_ecRecover, []string{signData, signature})
	if err != nil {
		return "", err
	}
	return pointer.ToString()
}

func (rpcCli *RpcClient) PERSONAL_sendTransaction(txData *requestData.TransactionParameters, password string) (string, error) {

	pointer, err := rpcCli.netServeHandler(PERSONAL_sendTransaction, []interface{}{requestData.Struct2Map(*(txData.Transform())), password})
	if err != nil {
		return "", err
	}
	return pointer.ToString()
}

func (rpcCli *RpcClient) ScsRPCMethod_GetBalance(subChainAddr, sender string) (string, error) {

	pointer, err := rpcCli.netServeHandler(ScsRPCMethod_GetBalance, map[string]interface{}{"SubChainAddr": subChainAddr, "Sender": sender})
	if err != nil {
		return "", err
	}
	return pointer.ToString()
}

func (rpcCli *RpcClient) ScsRPCMethod_GetNonce(subChainAddr, sender string) (float64, error) {

	pointer, err := rpcCli.netServeHandler(ScsRPCMethod_GetNonce, map[string]interface{}{"SubChainAddr": subChainAddr, "Sender": sender})
	fmt.Println(pointer.Result)
	if err != nil {
		return 0, err
	}
	return pointer.Result.(float64), nil
}

func (rpcCli *RpcClient) ScsRPCMethod_GetBlockNumber(subChainAddr string) (float64, error) {

	pointer, err := rpcCli.netServeHandler(ScsRPCMethod_GetBlockNumber, map[string]interface{}{"subChainAddr": subChainAddr})
	fmt.Println(pointer.Result)
	if err != nil {
		return 0, err
	}
	return pointer.Result.(float64), nil
}
