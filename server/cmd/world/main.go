package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/rouzbehsbz/manticore/server/internal/gameplay"
	"github.com/rouzbehsbz/manticore/server/internal/gameplay/account"
	"github.com/rouzbehsbz/manticore/server/pkg/network"
	"github.com/rouzbehsbz/manticore/server/pkg/network/protocol"
	"github.com/rouzbehsbz/manticore/server/pkg/pool"
)

func main() {
	// isDevMode := flag.Bool("dev", true, "Run program in dev mode")
	// flag.Parse()

	// config, err := common.NewConfig(*isDevMode)
	// if err != nil {
	// 	panic(err)
	// }

	// db, err := db.NewDb(
	// 	config.DbHost,
	// 	config.DbPort,
	// 	config.DbUsername,
	// 	config.DbPassword,
	// 	config.DbName,
	// 	config.DbMaxConnections,
	// )
	// if err != nil {
	// 	panic(err)
	// }

	blockingPool := pool.NewPool(runtime.NumCPU())

	dispatcher := gameplay.NewDispatcher()
	dispatcher.Register(protocol.RegisterRequestPacketId, account.NewRegisterHandler(nil))
	dispatcher.Register(protocol.LoginRequestPacketId, account.NewLoginHandler(nil))

	server := network.NewServer()

	go func() {
		for rp := range server.SessionsManager.Blocking {
			blockingPool.Jobs <- func() {
				_ = dispatcher.Dispatch(rp)
			}
		}
	}()

	go func() {
		for {
			server.SessionsManager.Flush()

			time.Sleep(100 * time.Millisecond)
		}
	}()

	addr := fmt.Sprintf("%s:%d", "0.0.0.0", 3000)
	server.Listen(addr)
}
