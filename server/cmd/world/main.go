package main

import (
	"flag"
	"fmt"
	"runtime"
	"time"

	"github.com/rouzbehsbz/manticore/server/internal/common"
	"github.com/rouzbehsbz/manticore/server/internal/gameplay/account"
	"github.com/rouzbehsbz/manticore/server/internal/gameplay/core"
	"github.com/rouzbehsbz/manticore/server/internal/infra/db"
	"github.com/rouzbehsbz/manticore/server/pkg/network"
	"github.com/rouzbehsbz/manticore/server/pkg/network/protocol"
	"github.com/rouzbehsbz/manticore/server/pkg/pool"
	"github.com/rouzbehsbz/zurvan"
)

const TickRate = 100 * time.Millisecond

func main() {
	isDevMode := flag.Bool("dev", true, "Run program in dev mode")
	flag.Parse()

	config, err := common.NewConfig(*isDevMode)
	if err != nil {
		panic(err)
	}

	numCpus := runtime.NumCPU()

	db, err := db.NewDb(
		config.DbHost,
		config.DbPort,
		config.DbUsername,
		config.DbPassword,
		config.DbName,
		numCpus,
	)
	if err != nil {
		panic(err)
	}

	blockingPool := pool.NewPool(numCpus)

	dispatcher := network.NewDispatcher()
	dispatcher.Register(protocol.RegisterRequestPacketId, account.NewRegisterHandler(db))
	dispatcher.Register(protocol.LoginRequestPacketId, account.NewLoginHandler(db))

	server := network.NewServer()

	go func() {
		for rp := range server.SessionsManager.Blocking {
			blockingPool.Jobs <- func() {
				_ = dispatcher.Dispatch(rp)
			}
		}
	}()

	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	go server.Listen(addr)

	world := zurvan.NewWorld(TickRate)

	world.PushCommands(
		zurvan.NewAddResourceCommand(server.SessionsManager),
		zurvan.NewAddResourceCommand(server.SessionsManager.NonBlocking),
		zurvan.NewAddResourceCommand(dispatcher),
	)

	world.AddSystems(
		zurvan.BuildStageSystems(zurvan.PreUpdateStage,
			&core.NetworkReceiveSystem{},
		),
	)

	world.AddSystems(
		zurvan.BuildStageSystems(zurvan.PostUpdateStage,
			&core.NetworkFlushSystem{},
		),
	)

	world.Run()
}
