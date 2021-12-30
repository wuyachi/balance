package watch

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Eth ...
type Eth struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
	accounts   []common.Address
	threshold  uint
	chain      string
	onAlarm    func(string)
	clients    []*ethclient.Client
	i          int
}

// NewEth ...
func NewEth(ctx context.Context, cancelFunc context.CancelFunc, chain string, rpcAddresses, accountsHex []string, threshold uint) *Eth {

	var clients []*ethclient.Client
	for _, rpcAddr := range rpcAddresses {
		client, err := ethclient.Dial(rpcAddr)
		if err != nil {
			panic(fmt.Sprintf("[%s] invalid rpcAddr:%s", chain, rpcAddr))
		}
		clients = append(clients, client)
	}

	var accounts []common.Address
	for _, accountHex := range accountsHex {
		accounts = append(accounts, common.HexToAddress(accountHex))
	}
	return &Eth{ctx: ctx, cancelFunc: cancelFunc, threshold: threshold, chain: chain, accounts: accounts, clients: clients}
}

// SetAlarm ...
func (w *Eth) SetAlarm(onAlarm func(string)) {
	w.onAlarm = onAlarm
}

// Start ...
func (w *Eth) Start(tokens ...string) (err error) {
	if w.onAlarm == nil {
		err = fmt.Errorf("[%s] alarm not set", w.chain)
		return
	}

	defer w.cancelFunc()

	for {

		client := w.clients[w.i%len(w.clients)]

		for _, account := range w.accounts {
			var balance *big.Int
			if len(tokens) > 0 {
				msg := getBalanceMsg(common.HexToAddress(tokens[0]), account)
				balanceByte, err := client.CallContract(context.Background(), msg, nil)
				if err != nil {
					log.Printf("[%s] BalanceOf err:%v", w.chain, err)
					continue
				}
				balance = new(big.Int).SetBytes(balanceByte)
			} else {
				balance, err = client.BalanceAt(w.ctx, account, nil)
				if err != nil {
					log.Printf("[%s] BalanceAt err:%v", w.chain, err)
					continue
				}
			}
			if balance.Uint64() <= uint64(w.threshold) {
				w.onAlarm(fmt.Sprintf("[%s] account %s is out of balance, balance:%d, threshold:%d", w.chain, account.Hex(), balance.Uint64(), w.threshold))
				continue
			}
			log.Printf("[%s] account %s: balance:%d threshold:%d", w.chain, account.Hex(), balance.Uint64(), w.threshold)
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

func getBalanceMsg(token, account common.Address) ethereum.CallMsg {
	balanceOfSig := "70a08231" //balanceOf

	method, _ := hex.DecodeString(balanceOfSig)
	zero, _ := hex.DecodeString(strings.Repeat("0", 24))
	method = append(method, zero...)
	method = append(method, account[:]...)

	msg := ethereum.CallMsg{}
	msg.To = &token
	msg.Data = method

	return msg
}
