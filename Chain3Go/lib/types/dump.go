package types

import (
	"Chain3Go/lib/common"
	pb "Chain3Go/lib/proto"
	"strconv"
)

// type ContractInfoReq struct {
// 	Reqtype      int//查询类型0: 查看合约全部变量 , 1: 查看合约某一个数组变量 , 2: 查看合约某一个mapping变量 , 3: 查看合约某一个结构体变量, 4: 查看合约某一简单类型变量（单倍长度存储的变量）, 5: 查看合约某一变长变量（如string、bytes）
// 	Key          string//64位定长十六进制字符串，查询的变量在合约里面的index ，查询全部变量时可以不填
// 	Position     string//64位定长十六进制字符串，当Reqtype==1时，Position为数组里第几个变量（从0开始）；当Reqtype==2时，Position为mapping下标
// 	Structformat []byte//当出现结构体变量时，此变量描述结构，结构体只允许出现1-single（简单类型变量单倍长度存储的变量）, 2-list（简单类型数组变量）, 3-string变长变量（如string、bytes），若结构变量为ContractInfoReq，Structformat = []byte{‘1’,’3’,’3’,’3’}
// }

type GetContractInfoReq struct {
	SubChainAddr common.Address
	Request      []*pb.StorageRequest
}

func ScreeningStorage(storage map[string]string, request []*pb.StorageRequest) map[string]string {
	resp := make(map[string]string)
	for _, val := range request {
		key := common.Bytes2Hex(val.Storagekey)
		position := common.Bytes2Hex(val.Position)
		structformat := val.Structformat
		switch val.Reqtype {
		case 0:
			for k, value := range storage {
				resp[k] = value
			}
		case 1:
			if len(position) == 0 {
				var num int64
				strlen := storage[key]
				if len(strlen) > 2 {
					num, _ = strconv.ParseInt(strlen[2:], 16, 64)
				} else {
					num, _ = strconv.ParseInt(strlen, 16, 64)
				}
				resp[key] = storage[key]
				keys := common.KeytoKey(key)
				for i := int64(0); i < num; i++ {
					if len(structformat) != 0 {
						key0 := keys
						for j := 0; j < len(structformat); j++ {
							if structformat[j] == '1' {
								resp[key0] = storage[key0]
							} else if structformat[j] == '2' {
								resp[key0] = storage[key0]
								var num0 int64
								strlen0 := storage[key0]
								if len(strlen0) > 2 {
									num0, _ = strconv.ParseInt(strlen0[2:], 16, 64)
								} else {
									num0, _ = strconv.ParseInt(strlen0, 16, 64)
								}
								key1 := common.KeytoKey(key0)
								for k := int64(0); k < num0; k++ {
									resp[key1] = storage[key1]
									key1 = common.IncreaseHexByOne(key1)
								}
							} else if structformat[j] == '3' {
								nlen := len(storage[key0])
								if nlen == 66 {
									resp[key0] = storage[key0]
								} else if nlen == 2 {
									resp[key0] = storage[key0]
									key1 := common.KeytoKey(key0)
									resp[key1] = storage[key1]
									key1 = common.IncreaseHexByOne(key1)
									resp[key1] = storage[key1]
								} else if nlen > 2 && nlen < 66 {
									resp[key0] = storage[key0]
									if nlen < 7 {
										num, _ := strconv.ParseInt(storage[key0][2:], 16, 64)
										key1 := common.KeytoKey(key0)
										for i := num - 1; i > 0; {
											resp[key1] = storage[key1]
											key1 = common.IncreaseHexByOne(key1)
											i = i - 64
										}
									}
								}
							}
							key0 = common.IncreaseHexByOne(key0)
						}
					} else {
						// resp[keys] = storage[keys]
						key0 := keys
						nlen := len(storage[key0])
						if nlen == 66 {
							resp[key0] = storage[key0]
						} else if nlen == 2 {
							resp[key0] = storage[key0]
							key1 := common.KeytoKey(key0)
							resp[key1] = storage[key1]
							key1 = common.IncreaseHexByOne(key1)
							resp[key1] = storage[key1]
						} else if nlen > 2 && nlen < 66 {
							resp[key0] = storage[key0]
							if nlen < 7 {
								num, _ := strconv.ParseInt(storage[key0][2:], 16, 64)
								key1 := common.KeytoKey(key0)
								for i := num - 1; i > 0; {
									if storage[key1] != "" {
										resp[key1] = storage[key1]
									}
									key1 = common.IncreaseHexByOne(key1)
									i = i - 64
								}
							}
						}
					}
					keys = common.IncreaseHexByOne(keys)
				}

			} else {
				num, _ := strconv.ParseInt(position, 16, 64)
				keys := common.KeytoKey(key)
				keys = common.IncreaseHexByNum(num, keys)
				if len(structformat) != 0 {
					key0 := keys
					for j := 0; j < len(structformat); j++ {
						if structformat[j] == '1' {
							resp[key0] = storage[key0]
						} else if structformat[j] == '2' {
							resp[key0] = storage[key0]
							var num0 int64
							strlen0 := storage[key0]
							if len(strlen0) > 2 {
								num0, _ = strconv.ParseInt(strlen0[2:], 16, 64)
							} else {
								num0, _ = strconv.ParseInt(strlen0, 16, 64)
							}
							key1 := common.KeytoKey(key0)
							for k := int64(0); k < num0; k++ {
								resp[key1] = storage[key1]
								key1 = common.IncreaseHexByOne(key1)
							}
						} else if structformat[j] == '3' {
							nlen := len(storage[key0])
							if nlen == 66 {
								resp[key0] = storage[key0]
							} else if nlen == 2 {
								resp[key0] = storage[key0]
								key1 := common.KeytoKey(key0)
								resp[key1] = storage[key1]
								key1 = common.IncreaseHexByOne(key1)
								resp[key1] = storage[key1]
							} else if nlen > 2 && nlen < 66 {
								resp[key0] = storage[key0]
								if nlen < 7 {
									num, _ := strconv.ParseInt(storage[key0][2:], 16, 64)
									key1 := common.KeytoKey(key0)
									for i := num - 1; i > 0; {
										resp[key1] = storage[key1]
										key1 = common.IncreaseHexByOne(key1)
										i = i - 64
									}
								}
							}
						}
						key0 = common.IncreaseHexByOne(key0)
					}
				} else {
					// resp[keys] = storage[keys]
					key0 := keys
					nlen := len(storage[key0])
					if nlen == 66 {
						resp[key0] = storage[key0]
					} else if nlen == 2 {
						resp[key0] = storage[key0]
						key1 := common.KeytoKey(key0)
						resp[key1] = storage[key1]
						key1 = common.IncreaseHexByOne(key1)
						resp[key1] = storage[key1]
					} else if nlen > 2 && nlen < 66 {
						resp[key0] = storage[key0]
						if nlen < 7 {
							num, _ := strconv.ParseInt(storage[key0][2:], 16, 64)
							key1 := common.KeytoKey(key0)
							for i := num - 1; i > 0; {
								if storage[key1] != "" {
									resp[key1] = storage[key1]
								}
								key1 = common.IncreaseHexByOne(key1)
								i = i - 64
							}
						}
					}
				}
			}
		case 2:
			keys := common.KeytoKey(position + key)
			if len(structformat) != 0 {
				key0 := keys
				for j := 0; j < len(structformat); j++ {
					if structformat[j] == '1' {
						resp[key0] = storage[key0]
					} else if structformat[j] == '2' {
						resp[key0] = storage[key0]
						var num0 int64
						strlen0 := storage[key0]
						if len(strlen0) > 2 {
							num0, _ = strconv.ParseInt(strlen0[2:], 16, 64)
						} else {
							num0, _ = strconv.ParseInt(strlen0, 16, 64)
						}

						key1 := common.KeytoKey(key0)
						for k := int64(0); k < num0; k++ {
							resp[key1] = storage[key1]
							key1 = common.IncreaseHexByOne(key1)
						}
					} else if structformat[j] == '3' {
						nlen := len(storage[key0])
						if nlen == 66 {
							resp[key0] = storage[key0]
						} else if nlen == 2 {
							resp[key0] = storage[key0]
							key1 := common.KeytoKey(key0)
							resp[key1] = storage[key1]
							key1 = common.IncreaseHexByOne(key1)
							resp[key1] = storage[key1]
						} else if nlen > 2 && nlen < 66 {
							resp[key0] = storage[key0]
							if nlen < 7 {
								num, _ := strconv.ParseInt(storage[key0][2:], 16, 64)
								key1 := common.KeytoKey(key0)
								for i := num - 1; i > 0; {
									resp[key1] = storage[key1]
									key1 = common.IncreaseHexByOne(key1)
									i = i - 64
								}
							}
						}
					}
					key0 = common.IncreaseHexByOne(key0)
				}
			} else {
				// resp[keys] = storage[keys]
				key0 := keys
				nlen := len(storage[key0])
				if nlen == 66 {
					resp[key0] = storage[key0]
				} else if nlen == 2 {
					resp[key0] = storage[key0]
					key1 := common.KeytoKey(key0)
					resp[key1] = storage[key1]
					key1 = common.IncreaseHexByOne(key1)
					resp[key1] = storage[key1]
				} else if nlen > 2 && nlen < 66 {
					resp[key0] = storage[key0]
					if nlen < 7 {
						num, _ := strconv.ParseInt(storage[key0][2:], 16, 64)
						key1 := common.KeytoKey(key0)
						for i := num - 1; i > 0; {
							if storage[key1] != "" {
								resp[key1] = storage[key1]
							}
							key1 = common.IncreaseHexByOne(key1)
							i = i - 64
						}
					}
				}
			}
		case 3:
			key0 := key
			for j := 0; j < len(structformat); j++ {
				if structformat[j] == '1' {
					resp[key0] = storage[key0]
				} else if structformat[j] == '2' {
					resp[key0] = storage[key0]
					var num0 int64
					strlen0 := storage[key0]
					if len(strlen0) > 2 {
						num0, _ = strconv.ParseInt(strlen0[2:], 16, 64)
					} else {
						num0, _ = strconv.ParseInt(strlen0, 16, 64)
					}
					key1 := common.KeytoKey(key0)
					for k := int64(0); k < num0; k++ {
						resp[key1] = storage[key1]
						key1 = common.IncreaseHexByOne(key1)
					}
				} else if structformat[j] == '3' {
					nlen := len(storage[key0])
					if nlen == 66 {
						resp[key0] = storage[key0]
					} else if nlen == 2 {
						resp[key0] = storage[key0]
						key1 := common.KeytoKey(key0)
						resp[key1] = storage[key1]
						key1 = common.IncreaseHexByOne(key1)
						resp[key1] = storage[key1]
					} else if nlen > 2 && nlen < 66 {
						resp[key0] = storage[key0]
						if nlen < 7 {
							num, _ := strconv.ParseInt(storage[key0][2:], 16, 64)
							key1 := common.KeytoKey(key0)
							for i := num - 1; i > 0; {
								resp[key1] = storage[key1]
								key1 = common.IncreaseHexByOne(key1)
								i = i - 64
							}
						}
					}
				}
				key0 = common.IncreaseHexByOne(key0)
			}
		case 4:
			resp[key] = storage[key]
		case 5:
			nlen := len(storage[key])
			if nlen == 66 {
				resp[key] = storage[key]
			} else if nlen == 2 {
				resp[key] = storage[key]
				key0 := common.KeytoKey(key)
				resp[key0] = storage[key0]
				key0 = common.IncreaseHexByOne(key0)
				resp[key0] = storage[key0]
			} else if nlen > 2 && nlen < 66 {
				resp[key] = storage[key]
				if nlen < 7 {
					num, _ := strconv.ParseInt(storage[key][2:], 16, 64)
					key0 := common.KeytoKey(key)
					for i := num - 1; i > 0; {
						resp[key0] = storage[key0]
						key0 = common.IncreaseHexByOne(key0)
						i = i - 64
					}
				}
			}
		default:

		}
	}
	return resp
}
