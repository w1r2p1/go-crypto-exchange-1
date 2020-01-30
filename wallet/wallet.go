package wallet

// Create a wallet
func Create(addr, pubKey, prKey string) (w *Wallet) {
	return &Wallet{addr, pubKey, prKey}
}

// Wallet data
type Wallet struct {
	addr   string
	pubKey string
	prKey  string
}

// Addr returns wallet address
func (w *Wallet) Addr() string {
	return w.addr
}

// Pub returns wallet public key
func (w *Wallet) Pub() string {
	return w.pubKey
}

// Priv returns wallet private key
func (w *Wallet) Priv() string {
	return w.prKey
}