// Copyright 2016 The MOAC-core Authors
// This file is part of the MOAC-core library.
//
// The MOAC-core library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The MOAC-core library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the MOAC-core library. If not, see <http://www.gnu.org/licenses/>.

package params

import (
	"fmt"
	"math/big"

	"Chain3Go/lib/common"
	"Chain3Go/lib/log"
)

//Updated 2018/04/10
//default value for the SCS service
var (
	SCSService             = false
	ShowToPublic           = false
	VnodeServiceCfg        = "localhost:50062"
	VnodeBeneficialAddress = ""
	VnodeIp                = ""
)

const (
	MainNetworkId = 99
	TestNetworkId = 101
	DevNetworkId  = 100
	NetworkId188  = 188
)

var (
	MainnetGenesisHash = common.HexToHash("0x6b9661646439fab926ffc9bccdf3abb572d5209ae59e3390abf76aee4e2d49cd") // Mainnet genesis hash to enforce below configs on
	TestnetGenesisHash = common.HexToHash("0xb44f499ad420fbba3d4a7e05e9485cddc0e70e3e3622919a5d945e9ed4f7699c") // Testnet genesis hash to enforce below configs on
)

var (
	// MainnetChainConfig is the chain parameters to run a node on the main network.
	MainnetChainConfig = &ChainConfig{
		ChainId:            big.NewInt(MainNetworkId),
		PanguBlock:         big.NewInt(0),
		RemoveEmptyAccount: true,
		// DAOForkBlock:   big.NewInt(1920000),
		// DAOForkSupport: true,
		// EIP150Block:    big.NewInt(2463000),
		// EIP150Hash:     common.HexToHash("0x2086799aeebeae135c246c65021c82b4e15a2c451340993aacfd2751886514f0"),
		// EIP155Block:    big.NewInt(2675000),
		// EIP158Block:    big.NewInt(2675000),
		// ByzantiumBlock: big.NewInt(math.MaxInt64), // Don't enable yet
		NuwaBlock: big.NewInt(647200), // 2018/08/08 UTC 12:00

		Ethash: new(EthashConfig),
	}

	// TestnetChainConfig contains the chain parameters to run a node on the test network.
	TestnetChainConfig = &ChainConfig{
		ChainId:            big.NewInt(TestNetworkId),
		PanguBlock:         big.NewInt(0),
		RemoveEmptyAccount: true,
		// DAOForkBlock:   nil,
		// DAOForkSupport: true,
		// // EIP150Block:    big.NewInt(0),
		// // EIP150Hash:     common.HexToHash("0x41941023680923e0fe4d74a34bdac8141f2540e3ae90623718e47d66d1ca4a2d"),
		// EIP155Block:    big.NewInt(0),
		// EIP158Block:    big.NewInt(0),
		// ByzantiumBlock: big.NewInt(0),
		NuwaBlock: big.NewInt(616700), // 2018/07/30

		Ethash: new(EthashConfig),
	}

	// RinkebyChainConfig contains the chain parameters to run a node on the Rinkeby test network.
	// RinkebyChainConfig = &ChainConfig{
	// 	ChainId:        big.NewInt(4),
	// 	PanguBlock: big.NewInt(1),
	// 	DAOForkBlock:   nil,
	// 	DAOForkSupport: true,
	// 	EIP150Block:    big.NewInt(2),
	// 	EIP150Hash:     common.HexToHash("0x9b095b36c15eaf13044373aef8ee0bd3a382a5abb92e402afa44b8249c3a90e9"),
	// 	EIP155Block:    big.NewInt(3),
	// 	EIP158Block:    big.NewInt(3),
	// 	ByzantiumBlock: big.NewInt(math.MaxInt64), // Don't enable yet

	// 	Clique: &CliqueConfig{
	// 		Period: 15,
	// 		Epoch:  30000,
	// 	},
	// }

	// AllProtocolChanges contains every protocol change (EIPs)
	// introduced and accepted by the MoacNode core developers.
	//
	// This configuration is intentionally not using keyed fields.
	// This configuration must *always* have all forks enabled, which
	// means that all fields must be set at all times. This forces
	// anyone adding flags to the config to also have to set these
	// fields.
	// AllProtocolChanges = &ChainConfig{big.NewInt(1337), big.NewInt(0), nil, false, big.NewInt(0), common.Hash{}, big.NewInt(0), big.NewInt(0), big.NewInt(0), new(EthashConfig), nil}
	// TestChainConfig    = &ChainConfig{big.NewInt(1), big.NewInt(0), nil, false, big.NewInt(0), common.Hash{}, big.NewInt(0), big.NewInt(0), big.NewInt(0), new(EthashConfig), nil}

	AllProtocolChanges = &ChainConfig{big.NewInt(DevNetworkId), big.NewInt(0), true, big.NewInt(0), new(EthashConfig)} //, nil, false, big.NewInt(0), common.Hash{}, big.NewInt(0), big.NewInt(0), big.NewInt(0), new(EthashConfig), nil}
	TestChainConfig    = &ChainConfig{big.NewInt(DevNetworkId), big.NewInt(0), true, big.NewInt(0), new(EthashConfig)} //, nil, false, big.NewInt(0), common.Hash{}, big.NewInt(0), big.NewInt(0), big.NewInt(0), new(EthashConfig), nil}

	TestRules = TestChainConfig.Rules(new(big.Int))
)

// ChainConfig is the core config which determines the blockchain settings.
//
// ChainConfig is stored in the database on a per block basis. This means
// that any network, identified by its genesis block, can have its own
// set of configuration options.
// default is to omitempty Accounts
// EIP158
//
type ChainConfig struct {
	ChainId *big.Int `json:"chainId"` // Chain id identifies the current chain and is used for replay protection

	PanguBlock         *big.Int `json:"panguBlock,omitempty"`         // Pangu switch block (nil = no fork, 0 = already pangu)
	RemoveEmptyAccount bool     `json:"removeEmptyAccount,omitempty"` //Replace EIP158 check and should be set to true

	NuwaBlock *big.Int `json:"nuwaBlock,omitempty"` // nuwa switch block (nil = no fork, 0 = already on nuwa)

	// Various consensus engines
	Ethash *EthashConfig `json:"ethash,omitempty"`
	// Clique *CliqueConfig `json:"clique,omitempty"`
}

// EthashConfig is the consensus engine configs for proof-of-work based sealing.
type EthashConfig struct{}

// String implements the stringer interface, returning the consensus engine details.
func (c *EthashConfig) String() string {
	return "ethash"
}

// Not used
// CliqueConfig is the consensus engine configs for proof-of-authority based sealing.
type CliqueConfig struct {
	Period uint64 `json:"period"` // Number of seconds between blocks to enforce
	Epoch  uint64 `json:"epoch"`  // Epoch length to reset votes and checkpoint
}

// String implements the stringer interface, returning the consensus engine details.
func (c *CliqueConfig) String() string {
	return "clique"
}

// String implements the fmt.Stringer interface.
// Removed the uncessary
func (c *ChainConfig) String() string {
	var engine interface{}
	switch {
	case c.Ethash != nil:
		engine = c.Ethash
	// case c.Clique != nil:
	// 	engine = c.Clique
	default:
		engine = "unknown"
	}
	/*	return fmt.Sprintf("{ChainID: %v Pangu: %v DAO: %v DAOSupport: %v EIP150: %v EIP155: %v EIP158: %v Byzantium: %v Engine: %v}",
		c.ChainId,
		c.PanguBlock,
		c.DAOForkBlock,
		c.DAOForkSupport,
		c.EIP150Block,
		c.EIP155Block,
		c.EIP158Block,
		c.ByzantiumBlock,
		engine,
	)*/
	return fmt.Sprintf("{ChainID: %v Pangu: %v Nuwa: %v Engine: %v}",
		c.ChainId, c.PanguBlock, c.NuwaBlock, engine)
}

// IsPangu returns whether num is either equal to the pangu block or greater.
func (c *ChainConfig) IsPangu(num *big.Int) bool {
	flag := isForked(c.PanguBlock, num)
	log.Debugf("[params/config.go->IsPangu] PanguBlock:%v num:%v flag:%v", c.PanguBlock, num, flag)
	return flag
}

//Remove the Empty account
//To replace the EIP158Block check
//Default is to remove all empty Accounts
func (c *ChainConfig) IsRemoveEmptyAccount(num *big.Int) bool {

	return isForked(c.PanguBlock, num)
	//Compare the current header with internal
	// return isForked(c.EIP158Block, num)
}

// IsDAO returns whether num is either equal to the DAO fork block or greater.
// func (c *ChainConfig) IsDAOFork(num *big.Int) bool {
// 	return isForked(c.DAOForkBlock, num)
// }

// func (c *ChainConfig) IsEIP150(num *big.Int) bool {
// 	return isForked(c.EIP150Block, num)
// }

// func (c *ChainConfig) IsEIP155(num *big.Int) bool {
// 	return isForked(c.EIP155Block, num)
// }

// func (c *ChainConfig) IsEIP158(num *big.Int) bool {
// 	return isForked(c.EIP158Block, num)
// }

// func (c *ChainConfig) IsByzantium(num *big.Int) bool {
// 	return isForked(c.ByzantiumBlock, num)
// }

func (c *ChainConfig) IsNuwa(num *big.Int) bool {
	log.Debugf("[params/config.go->IsNuwa] NuwaBlock:%v num:%v", c.NuwaBlock, num)
	return isForked(c.NuwaBlock, num)
}

// GasTable returns the gas table corresponding to the current phase (pangu or pangu reprice).
//
// The returned GasTable's fields shouldn't, under any circumstances, be changed.
// Used the
func (c *ChainConfig) GasTable(num *big.Int) GasTable {
	if num == nil {
		return GasTablePangu
	}
	// switch {
	// case c.IsEIP158(num):
	// 	return GasTableEIP158
	// case c.IsEIP150(num):
	// 	return GasTablePanguGasRepriceFork
	// default:
	return GasTablePangu
	// }
}

// CheckCompatible checks whether scheduled fork transitions have been imported
// with a mismatching chain configuration.
func (c *ChainConfig) CheckCompatible(newcfg *ChainConfig, height uint64) *ConfigCompatError {
	bhead := new(big.Int).SetUint64(height)

	// Iterate checkCompatible to find the lowest conflict.
	var lasterr *ConfigCompatError
	for {
		err := c.checkCompatible(newcfg, bhead)
		if err == nil || (lasterr != nil && err.RewindTo == lasterr.RewindTo) {
			break
		}
		lasterr = err
		bhead.SetUint64(err.RewindTo)
	}
	return lasterr
}

func (c *ChainConfig) checkCompatible(newcfg *ChainConfig, head *big.Int) *ConfigCompatError {
	if isForkIncompatible(c.PanguBlock, newcfg.PanguBlock, head) {
		return newCompatError("Pangu fork block", c.PanguBlock, newcfg.PanguBlock)
	}
	// if isForkIncompatible(c.DAOForkBlock, newcfg.DAOForkBlock, head) {
	// 	return newCompatError("DAO fork block", c.DAOForkBlock, newcfg.DAOForkBlock)
	// }
	// if c.IsDAOFork(head) && c.DAOForkSupport != newcfg.DAOForkSupport {
	// 	return newCompatError("DAO fork support flag", c.DAOForkBlock, newcfg.DAOForkBlock)
	// }
	// if isForkIncompatible(c.EIP150Block, newcfg.EIP150Block, head) {
	// 	return newCompatError("EIP150 fork block", c.EIP150Block, newcfg.EIP150Block)
	// }
	// if isForkIncompatible(c.EIP155Block, newcfg.EIP155Block, head) {
	// 	return newCompatError("EIP155 fork block", c.EIP155Block, newcfg.EIP155Block)
	// }
	// if isForkIncompatible(c.EIP158Block, newcfg.EIP158Block, head) {
	// 	return newCompatError("EIP158 fork block", c.EIP158Block, newcfg.EIP158Block)
	// }
	// if c.IsEIP158(head) && !configNumEqual(c.ChainId, newcfg.ChainId) {
	// 	return newCompatError("EIP158 chain ID", c.EIP158Block, newcfg.EIP158Block)
	// }
	// if isForkIncompatible(c.ByzantiumBlock, newcfg.ByzantiumBlock, head) {
	// 	return newCompatError("Byzantium fork block", c.ByzantiumBlock, newcfg.ByzantiumBlock)
	// }
	return nil
}

// isForkIncompatible returns true if a fork scheduled at s1 cannot be rescheduled to
// block s2 because head is already past the fork.
func isForkIncompatible(s1, s2, head *big.Int) bool {
	return (isForked(s1, head) || isForked(s2, head)) && !configNumEqual(s1, s2)
}

// isForked returns whether a fork scheduled at block s is active at the given head block.
func isForked(s, head *big.Int) bool {
	if s == nil || head == nil {
		return false
	}
	return s.Cmp(head) <= 0
}

func configNumEqual(x, y *big.Int) bool {
	if x == nil {
		return y == nil
	}
	if y == nil {
		return x == nil
	}
	return x.Cmp(y) == 0
}

// ConfigCompatError is raised if the locally-stored blockchain is initialised with a
// ChainConfig that would alter the past.
type ConfigCompatError struct {
	What string
	// block numbers of the stored and new configurations
	StoredConfig, NewConfig *big.Int
	// the block number to which the local chain must be rewound to correct the error
	RewindTo uint64
}

func newCompatError(what string, storedblock, newblock *big.Int) *ConfigCompatError {
	var rew *big.Int
	switch {
	case storedblock == nil:
		rew = newblock
	case newblock == nil || storedblock.Cmp(newblock) < 0:
		rew = storedblock
	default:
		rew = newblock
	}
	err := &ConfigCompatError{what, storedblock, newblock, 0}
	if rew != nil && rew.Sign() > 0 {
		err.RewindTo = rew.Uint64() - 1
	}
	return err
}

func (err *ConfigCompatError) Error() string {
	return fmt.Sprintf("mismatching %s in database (have %d, want %d, rewindto %d)", err.What, err.StoredConfig, err.NewConfig, err.RewindTo)
}

// Rules wraps ChainConfig and is merely syntatic sugar or can be used for functions
// that do not have or require information about the block.
//
// Rules is a one time interface meaning that it shouldn't be used in between transition
// phases.
type Rules struct {
	ChainId *big.Int
	IsPangu bool
	IsNuwa  bool
	// IsPangu, IsEIP150, IsEIP155, IsEIP158 bool
	// IsByzantium                           bool
}

// Pangu 0.8 version
func (c *ChainConfig) Rules(num *big.Int) Rules {
	chainId := c.ChainId
	if chainId == nil {
		chainId = new(big.Int)
	}
	return Rules{ChainId: new(big.Int).Set(chainId), IsPangu: c.IsPangu(num), IsNuwa: c.IsNuwa(num)}
	// return Rules{ChainId: new(big.Int).Set(chainId), IsPangu: c.IsPangu(num), IsEIP150: c.IsEIP150(num), IsEIP155: c.IsEIP155(num), IsEIP158: c.IsEIP158(num), IsByzantium: c.IsByzantium(num)}
}
