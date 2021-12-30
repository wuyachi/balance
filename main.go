package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"log"

	"github.com/urfave/cli"
	"github.com/zhiqiangxu/balance/config"
	"github.com/zhiqiangxu/balance/flag"
	watchPkg "github.com/zhiqiangxu/balance/pkg/watch"
	"github.com/zhiqiangxu/util"
)

func setupApp() *cli.App {
	app := cli.NewApp()
	app.Usage = "Account balance watcher"
	app.Action = watch
	app.Flags = []cli.Flag{
		flag.ConfFlag,
	}
	return app
}

func onAlarm(info string) {
	log.Printf("ALARM:%s", info)
}

func watch(ctx *cli.Context) (err error) {
	confBytes, err := ioutil.ReadFile(ctx.String(flag.ConfFlag.Name))
	if err != nil {
		return
	}

	var conf config.Config
	err = json.Unmarshal(confBytes, &conf)
	if err != nil {
		return
	}

	wctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	var wg sync.WaitGroup
	if len(conf.Account.Eth) > 0 {
		eth := watchPkg.NewEth(wctx, cancelFunc, "Eth", conf.Node.Eth, conf.Account.Eth, conf.Threshold.Eth)
		eth.SetAlarm(onAlarm)
		util.GoFunc(&wg, func() {
			err := eth.Start()
			if err != nil {
				log.Printf("Eth watcher quits:%v", err)
			}
		})
	}
	if len(conf.Account.BSC) > 0 {
		eth := watchPkg.NewEth(wctx, cancelFunc, "BSC", conf.Node.BSC, conf.Account.BSC, conf.Threshold.BSC)
		eth.SetAlarm(onAlarm)
		util.GoFunc(&wg, func() {
			err := eth.Start()
			if err != nil {
				log.Printf("BSC watcher quits:%v", err)
			}
		})
	}
	if len(conf.Account.Heco) > 0 {
		eth := watchPkg.NewEth(wctx, cancelFunc, "Heco", conf.Node.Heco, conf.Account.Heco, conf.Threshold.Heco)
		eth.SetAlarm(onAlarm)
		util.GoFunc(&wg, func() {
			err := eth.Start()
			if err != nil {
				log.Printf("Heco watcher quits:%v", err)
			}
		})
	}
	if len(conf.Account.Polygon) > 0 {
		eth := watchPkg.NewEth(wctx, cancelFunc, "Polygon", conf.Node.Polygon, conf.Account.Polygon, conf.Threshold.Polygon)
		eth.SetAlarm(onAlarm)
		util.GoFunc(&wg, func() {
			err := eth.Start()
			if err != nil {
				log.Printf("Polygon watcher quits:%v", err)
			}
		})
	}
	if len(conf.Account.OK) > 0 {
		eth := watchPkg.NewEth(wctx, cancelFunc, "OK", conf.Node.OK, conf.Account.OK, conf.Threshold.OK)
		eth.SetAlarm(onAlarm)
		util.GoFunc(&wg, func() {
			err := eth.Start()
			if err != nil {
				log.Printf("OK watcher quits:%v", err)
			}
		})
	}
	if len(conf.Account.Arb) > 0 {
		eth := watchPkg.NewEth(wctx, cancelFunc, "Arb", conf.Node.Arb, conf.Account.Arb, conf.Threshold.Arb)
		eth.SetAlarm(onAlarm)
		util.GoFunc(&wg, func() {
			err := eth.Start()
			if err != nil {
				log.Printf("Arb watcher quits:%v", err)
			}
		})
	}

	if len(conf.Account.Ont) > 0 {
		ont := watchPkg.NewOnt(wctx, cancelFunc, "Ont", conf.Node.Ont, conf.Account.Ont, conf.Threshold.Ont)
		ont.SetAlarm(onAlarm)
		util.GoFunc(&wg, func() {
			err := ont.Start()
			if err != nil {
				log.Printf("Ont watcher quits:%v", err)
			}
		})
	}

	if len(conf.Account.Zil) > 0 {
		zil := watchPkg.NewZil(wctx, cancelFunc, "Zil", conf.Node.Zil, conf.Account.Zil, conf.Threshold.Zil)
		zil.SetAlarm(onAlarm)
		util.GoFunc(&wg, func() {
			err := zil.Start()
			if err != nil {
				log.Printf("Zil watcher quits:%v", err)
			}
		})
	}

	if len(conf.Account.Neo) > 0 {
		neo := watchPkg.NewNeo(wctx, cancelFunc, "Neo", conf.Node.Neo, conf.Account.Neo, conf.Threshold.Neo)
		neo.SetAlarm(onAlarm)
		util.GoFunc(&wg, func() {
			err := neo.Start()
			if err != nil {
				log.Printf("Neo watcher quits:%v", err)
			}
		})
	}

	if len(conf.Account.Metis) > 0 {
		metis := watchPkg.NewEth(wctx, cancelFunc, "Metis", conf.Node.Metis, conf.Account.Metis, conf.Threshold.Metis)
		metis.SetAlarm(onAlarm)
		util.GoFunc(&wg, func() {
			err := metis.Start("deaddeaddeaddeaddeaddeaddeaddeaddead0000")
			if err != nil {
				log.Printf("Metis watcher quits:%v", err)
			}
		})
	}

	wg.Wait()

	return
}

func main() {
	if err := setupApp().Run(os.Args); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
