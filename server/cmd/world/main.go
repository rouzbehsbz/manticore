package main

import (
	"flag"
	"fmt"
	"runtime"

	"github.com/rouzbehsbz/manticore/server/internal/common"
	"github.com/rouzbehsbz/manticore/server/internal/gameplay"
	"github.com/rouzbehsbz/manticore/server/internal/gameplay/account"
	"github.com/rouzbehsbz/manticore/server/internal/infra/db"
	"github.com/rouzbehsbz/manticore/server/pkg/network"
	"github.com/rouzbehsbz/manticore/server/pkg/network/protocol"
	"github.com/rouzbehsbz/manticore/server/pkg/pool"
)

func main() {
	isDevMode := flag.Bool("dev", true, "Run program in dev mode")
	flag.Parse()

	config, err := common.NewConfig(*isDevMode)
	if err != nil {
		panic(err)
	}

	db, err := db.NewDb(
		config.DbHost,
		config.DbPort,
		config.DbUsername,
		config.DbPassword,
		config.DbName,
		config.DbMaxConnections,
	)
	if err != nil {
		panic(err)
	}

	blockingPool := pool.NewPool(runtime.NumCPU())

	dispatcher := gameplay.NewDispatcher()
	dispatcher.Register(protocol.RegisterRequestPacketId, account.NewLoginHandler(db))
	dispatcher.Register(protocol.LoginRequestPacketId, account.NewRegisterHandler(db))

	server := network.NewServer()

	go func() {
		for rp := range server.SessionsManager.Blocking {
			blockingPool.Jobs <- func() {
				_ = dispatcher.Dispatch(rp)
			}
		}
	}()

	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	server.Listen(addr)
}
