package accounts
import (
	"fmt"
	"math/big"
	ethereum "github.com/Cryptochain-VON"
	"github.com/Cryptochain-VON/common"
	"github.com/Cryptochain-VON/core/types"
	"github.com/Cryptochain-VON/event"
	"golang.org/x/crypto/sha3"
)
type Account struct {
	Address common.Address `json:"address"` 
	URL     URL            `json:"url"`     
}
const (
	MimetypeDataWithValidator = "data/validator"
	MimetypeTypedData         = "data/typed"
	MimetypeClique            = "application/x-clique-header"
	MimetypeTextPlain         = "text/plain"
)
type Wallet interface {
	URL() URL
	Status() (string, error)
	Open(passphrase string) error
	Close() error
	Accounts() []Account
	Contains(account Account) bool
	Derive(path DerivationPath, pin bool) (Account, error)
	SelfDerive(bases []DerivationPath, chain ethereum.ChainStateReader)
	SignData(account Account, mimeType string, data []byte) ([]byte, error)
	SignDataWithPassphrase(account Account, passphrase, mimeType string, data []byte) ([]byte, error)
	SignText(account Account, text []byte) ([]byte, error)
	SignTextWithPassphrase(account Account, passphrase string, hash []byte) ([]byte, error)
	SignTx(account Account, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error)
	SignTxWithPassphrase(account Account, passphrase string, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error)
}
type Backend interface {
	Wallets() []Wallet
	Subscribe(sink chan<- WalletEvent) event.Subscription
}
func TextHash(data []byte) []byte {
	hash, _ := TextAndHash(data)
	return hash
}
func TextAndHash(data []byte) ([]byte, string) {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), string(data))
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write([]byte(msg))
	return hasher.Sum(nil), msg
}
type WalletEventType int
const (
	WalletArrived WalletEventType = iota
	WalletOpened
	WalletDropped
)
type WalletEvent struct {
	Wallet Wallet          
	Kind   WalletEventType 
}
