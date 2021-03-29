package watch

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/joeqian10/neo-gogogo/tx"
	"github.com/joeqian10/neo-gogogo/wallet"
)

// Neo ...
type Neo struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
	accounts   []string
	threshold  uint
	chain      string
	onAlarm    func(string)
	helpers    []*wallet.WalletHelper
	i          int
}

// NewNeo ...
func NewNeo(ctx context.Context, cancelFunc context.CancelFunc, chain string, rpcAddresses, accountsBase58 []string, threshold uint) *Neo {

	var helpers []*wallet.WalletHelper
	for _, rpcAddr := range rpcAddresses {
		txBuilder := tx.NewTransactionBuilder(rpcAddr)
		walletHelper := wallet.NewWalletHelper(txBuilder, nil)
		helpers = append(helpers, walletHelper)
	}

	return &Neo{ctx: ctx, cancelFunc: cancelFunc, threshold: threshold, chain: chain, accounts: accountsBase58, helpers: helpers}
}

// SetAlarm ...
func (w *Neo) SetAlarm(onAlarm func(string)) {
	w.onAlarm = onAlarm
}

// Start ...
func (w *Neo) Start() (err error) {
	if w.onAlarm == nil {
		err = fmt.Errorf("[%s] alarm not set", w.chain)
		return
	}

	defer w.cancelFunc()

	for {

		helper := w.helpers[w.i%len(w.helpers)]

		for _, account := range w.accounts {
			_, gasBalance, err := helper.GetBalance(account)
			if err != nil {
				log.Printf("[%s] GetBalance err:%v", w.chain, err)
				continue
			}
			if gasBalance <= float64(w.threshold) {
				w.onAlarm(fmt.Sprintf("[%s] account %s is out of balance, balance:%v, threshold:%d", w.chain, account, gasBalance, w.threshold))
				continue
			}
			log.Printf("[%s] account %s: balance:%v threshold:%d", w.chain, account, gasBalance, w.threshold)
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
