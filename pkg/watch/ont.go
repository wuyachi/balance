package watch

import (
	"context"
	"fmt"
	"log"
	"time"

	ontSDK "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common"
)

// Ont ...
type Ont struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
	accounts   []common.Address
	threshold  uint
	chain      string
	onAlarm    func(string)
	sdks       []*ontSDK.OntologySdk
	i          int
}

// NewOnt ...
func NewOnt(ctx context.Context, cancelFunc context.CancelFunc, chain string, rpcAddresses, accountsHex []string, threshold uint) *Ont {

	var sdks []*ontSDK.OntologySdk
	for _, rpcAddr := range rpcAddresses {
		sdk := ontSDK.NewOntologySdk()
		client := sdk.NewRpcClient()
		client.SetAddress(rpcAddr)
		sdks = append(sdks, sdk)
	}

	var (
		accounts []common.Address
	)
	for _, accountHex := range accountsHex {
		account, err := str2addr(accountHex, nil)
		if err != nil {
			panic(fmt.Sprintf("[%s] invalid accountHex:%s", chain, accountHex))
		}
		accounts = append(accounts, account)
	}

	return &Ont{ctx: ctx, cancelFunc: cancelFunc, threshold: threshold, chain: chain, accounts: accounts, sdks: sdks}
}

// SetAlarm ...
func (w *Ont) SetAlarm(onAlarm func(string)) {
	w.onAlarm = onAlarm
}

func str2addr(addr string, is58 *bool) (result common.Address, err error) {
	if addr != "" {
		result, err = common.AddressFromBase58(addr)
		if err != nil {
			result, err = common.AddressFromHexString(addr)
			if is58 != nil {
				*is58 = false
			}
		} else {
			if is58 != nil {
				*is58 = true
			}
		}
	} else {
		err = fmt.Errorf("empty addr")
	}

	return
}

// Start ...
func (w *Ont) Start() (err error) {
	if w.onAlarm == nil {
		err = fmt.Errorf("[%s] alarm not set", w.chain)
		return
	}

	defer w.cancelFunc()

	for {

		sdk := w.sdks[w.i%len(w.sdks)]

		for _, account := range w.accounts {
			balance, err := sdk.Native.Ong.BalanceOf(account)
			if err != nil {
				log.Printf("[%s] BalanceOf err:%v", w.chain, err)
				continue
			}
			if balance <= uint64(w.threshold) {
				w.onAlarm(fmt.Sprintf("[%s] account %s is out of balance, balance:%d, threshold:%d", w.chain, account.ToBase58(), balance, w.threshold))
				continue
			}
			log.Printf("[%s] account %s: balance:%d threshold:%d", w.chain, account.ToBase58(), balance, w.threshold)
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
