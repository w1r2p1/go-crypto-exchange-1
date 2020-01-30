package wallet

// Create a wallet
func Create(addr, pubKey, prKey, seed string) (w *Wallet) {
	return &Wallet{addr, pubKey, prKey, seed}
}

// Wallet data
type Wallet struct {
	addr   string
	pubKey string
	prKey  string
	seed   string
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

// Seed returns wallet seed phrase
func (w *Wallet) Seed() string {
	return w.seed
}
