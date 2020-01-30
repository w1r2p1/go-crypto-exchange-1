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

var mnt *provider

func init() {
	NodeURL, err := config.GetVar("MINTER_NODE_URL")
	status := createStatus(1)
	if err != nil {
		status = createStatus(2)
		log.Printf("WARNING: %s. Minter provider will start with status: %s", err.Error(), status.String())
	}
	client := api.NewApi(NodeURL)
	mnt = &provider{status, client}
}

const (
	UnknownError = iota // 0
	Ok                  // 1
	WrongNodeUrl        // 2
)

var statusMessages = []string{"Unknown error", "Ok", "Wrong node URL"}

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

type provider struct {
	status Status
	client *api.Api
}

// GetInstance returns rpovider instance
func GetInstance() *provider {
	return mnt
}

func (p *provider) Status() Status {
	return p.status
}

func (p *provider) NodeStatus() (*api.StatusResult, error) {
	if p.status.code != Ok {
		return &api.StatusResult{}, fmt.Errorf("Cannot get node status: error: {code: %d, message: %s}", p.status.code, p.status.message)
	}
	return p.client.Status()
}

func (p *provider) GenerateWallet() (*wallet.Wallet, error) {
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
