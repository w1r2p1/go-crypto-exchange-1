package providers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/MinterTeam/minter-go-sdk/api"
	mntwallet "github.com/MinterTeam/minter-go-sdk/wallet"
	"github.com/alexander-molina/go-crypto-exchange/config"
	"github.com/alexander-molina/go-crypto-exchange/wallet"
)

var mnt *Provider

func init() {
	NodeURL, err := config.GetVar("MINTER_NODE_URL")
	status := createStatus(Ok)
	if err != nil {
		status = createStatus(WrongNodeUrl)
		log.Printf("WARNING: %s. Minter provider will start with status: %s", err.Error(), status.String())
	}
	client := api.NewApi(NodeURL)
	mnt = &Provider{"", status, client, make((map[string]string)), &wallet.Wallet{}}
}

const (
	Ok           = iota // 0
	WalletError         // 1
	UnknownError        // 2
	WrongNodeUrl        // 3
	WrongChainID        // 4
)

var statusMessages = []string{"Ok", "Wallet error", "Unknown error", "Wrong node URL", "Wrong chain id"}

func createStatus(code int) Status {
	message := statusMessages[code]
	return Status{code, message}
}

// Status of provider
type Status struct {
	code    int
	message string
}

func (s *Status) String() string {
	return "{code: " + strconv.Itoa(s.code) + ",  message: " + s.message + "}"
}

// Provider represents realization of minter provider
type Provider struct {
	chainID       string
	status        Status
	client        *api.Api
	currencies    map[string]string
	reserveWallet *wallet.Wallet
}

// GetInstance returns rpovider instance
func GetInstance() *Provider {
	return mnt
}

// Status returns status of current provider
func (p *Provider) Status() Status {
	return p.status
}

// Currencies returns symbols for currencies in reserve
func (p *Provider) Currencies() map[string]string {
	return p.currencies
}

// ReserveAddr returns reserve address
func (p *Provider) ReserveAddr() string {
	return p.reserveWallet.Addr()
}

// Balance returns balance of required address
func (p *Provider) Balance(address string) (map[string]string, error) {
	if p.status.code >= UnknownError {
		return make(map[string]string), fmt.Errorf("Cannot get balance: error: {code: %d, message: %s}", p.status.code, p.status.message)
	}
	return p.client.Balance(address, api.LatestBlockHeight)
}

// Nonce returns nonce for current address
func (p *Provider) Nonce(address string) (uint64, error) {
	if p.status.code >= UnknownError {
		return 0, fmt.Errorf("Cannot get nonce: error: {code: %d, message: %s}", p.status.code, p.status.message)
	}
	return p.client.Nonce(address)
}

// NodeStatus returns node status
func (p *Provider) NodeStatus() (*api.StatusResult, error) {
	if p.status.code >= UnknownError {
		return &api.StatusResult{}, fmt.Errorf("Cannot get node status: error: {code: %d, message: %s}", p.status.code, p.status.message)
	}
	return p.client.Status()
}

// MinGasPrice returns min gas price
func (p *Provider) MinGasPrice() (string, error) {
	if p.status.code >= UnknownError {
		return "", fmt.Errorf("Cannot get node status: error: {code: %d, message: %s}", p.status.code, p.status.message)
	}
	return p.client.MinGasPrice()
}

func (p *Provider) SignTx() {}

// func (p *Provider) SendTransaction(transaction transaction.SignedTransaction) (*SendResult, error) {}

// GenerateWallet returns minter generated wallet
func (p *Provider) GenerateWallet() (*wallet.Wallet, error) {
	mnemonic, err := mntwallet.NewMnemonic()
	if err != nil {
		return &wallet.Wallet{}, err
	}

	seed, err := mntwallet.Seed(mnemonic)
	if err != nil {
		return &wallet.Wallet{}, err
	}

	prKey, err := mntwallet.PrivateKeyBySeed(seed)
	if err != nil {
		return &wallet.Wallet{}, err
	}

	pubKey, err := mntwallet.PublicKeyByPrivateKey(prKey)
	if err != nil {
		return &wallet.Wallet{}, err
	}

	addr, err := mntwallet.AddressByPublicKey(pubKey)
	if err != nil {
		return &wallet.Wallet{}, err
	}

	return wallet.Create(addr, pubKey, prKey, mnemonic), nil
}

// AddCurrency add currency symbol to reserve
func (p *Provider) AddCurrency(currency string) {
	p.currencies[currency] = currency
}

// LoadWallet loads reserve wallet for minter chain from database
// TODO!
func (p *Provider) LoadWallet() error {
	p.reserveWallet = wallet.Create("Mxc41f4fd29af8f02055821961239bfd1b8cd4c77a",
		"Mpa5895da17ac89db40362510a86a23ad274560a9cf169b8bff2409d42da2bd37eb0780c2ae879f06f33ab0302eb31c69a6c65c1a09b51731a63f37a9e52602dd0",
		"b206665beebeb6cb270618e93c57cbb66ca14a6b15612e797bfc3017025ef09c",
		"opera spare oblige eight boring lady survey photo ugly unfold rib economy")
	return nil
}
