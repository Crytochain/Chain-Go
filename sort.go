package accounts
type AccountsByURL []Account
func (a AccountsByURL) Len() int           { return len(a) }
func (a AccountsByURL) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a AccountsByURL) Less(i, j int) bool { return a[i].URL.Cmp(a[j].URL) < 0 }
type WalletsByURL []Wallet
func (w WalletsByURL) Len() int           { return len(w) }
func (w WalletsByURL) Swap(i, j int)      { w[i], w[j] = w[j], w[i] }
func (w WalletsByURL) Less(i, j int) bool { return w[i].URL().Cmp(w[j].URL()) < 0 }
