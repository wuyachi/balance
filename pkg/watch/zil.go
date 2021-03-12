package watch

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	"github.com/Zilliqa/gozilliqa-sdk/provider"
)

// Zil ...
type Zil struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
	accounts   []string
	threshold  uint
	chain      string
	onAlarm    func(string)
	providers  []*provider.Provider
	i          int
}

// NewZil ...
func NewZil(ctx context.Context, cancelFunc context.CancelFunc, chain string, rpcAddresses, accountsBech32 []string, threshold uint) *Zil {

	var providers []*provider.Provider
	for _, rpcAddr := range rpcAddresses {
		zilSdk := provider.NewProvider(rpcAddr)
		providers = append(providers, zilSdk)
	}

	var accounts []string
	for _, accountBech32 := range accountsBech32 {
		account, err := bech32.FromBech32Addr(accountBech32)
		if err != nil {
			panic(fmt.Sprintf("[%s] invalid accountBech32:%s", chain, accountBech32))
		}
		accounts = append(accounts, account)
	}

	return &Zil{ctx: ctx, cancelFunc: cancelFunc, threshold: threshold, chain: chain, accounts: accounts, providers: providers}
}

// SetAlarm ...
func (w *Zil) SetAlarm(onAlarm func(string)) {
	w.onAlarm = onAlarm
}

// Start ...
func (w *Zil) Start() (err error) {
	if w.onAlarm == nil {
		err = fmt.Errorf("[%s] alarm not set", w.chain)
		return
	}

	defer w.cancelFunc()

	for {

		zilSdk := w.providers[w.i%len(w.providers)]

		for _, account := range w.accounts {
			bn, err := zilSdk.GetBalance(account)
			if err != nil {
				log.Printf("[%s] GetBalance err:%v", w.chain, err)
				continue
			}
			// if bn.Balance <= uint64(w.threshold) {
			// 	w.onAlarm(fmt.Sprintf("[%s] account %s is out of balance, balance:%s, threshold:%d", w.chain, account, balance, w.threshold))
			// 	continue
			// }
			log.Printf("[%s] account %s: balance:%s threshold:%d", w.chain, account, bn.Balance, w.threshold)
		}
		time.Sleep(time.Second)

		select {
		case <-w.ctx.Done():
			err = w.ctx.Err()
			return
		default:
		}

		w.i++
	}
}
