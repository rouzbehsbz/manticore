package main

import (
	"flag"
	"fmt"
	"runtime"
	"time"

	"github.com/rouzbehsbz/manticore/server/internal/common"
	"github.com/rouzbehsbz/manticore/server/internal/gameplay/account"
	"github.com/rouzbehsbz/manticore/server/internal/gameplay/character"
	"github.com/rouzbehsbz/manticore/server/internal/gameplay/combat"
	"github.com/rouzbehsbz/manticore/server/internal/gameplay/core"
	"github.com/rouzbehsbz/manticore/server/internal/infra/db"
	"github.com/rouzbehsbz/manticore/server/internal/models"
	"github.com/rouzbehsbz/manticore/server/pkg/network"
	"github.com/rouzbehsbz/manticore/server/pkg/network/protocol"
	"github.com/rouzbehsbz/manticore/server/pkg/pool"
	"github.com/rouzbehsbz/manticore/server/pkg/util"
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

	world := zurvan.NewWorld(TickRate)

	blockingPool := pool.NewPool(numCpus)

	dispatcher := network.NewDispatcher()
	dispatcher.Register(protocol.RegisterReqPacketId, account.NewRegisterHandler(db))
	dispatcher.Register(protocol.LoginReqPacketId, account.NewLoginHandler(db))
	dispatcher.Register(protocol.MyCharactersListReqPacketId, account.NewMyCharactersListHandler(db))
	dispatcher.Register(protocol.CharacterJoinReqPacketId, account.NewCharacterJoinHandler(db, world))
	dispatcher.Register(protocol.CastSpellReqPacketId, combat.NewCastSpellHandler(world))

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

	world.PushCommands(
		zurvan.NewAddResourceCommand(server.SessionsManager),
		zurvan.NewAddResourceCommand(server.SessionsManager.NonBlocking),
		zurvan.NewAddResourceCommand(dispatcher),
		zurvan.NewAddResourceCommand(util.NewSyncMap[uint32, zurvan.Entity]()),
		zurvan.NewAddResourceCommand(util.NewSyncMap[uint32, models.Spell]()),
	)

	world.AddSystems(
		zurvan.BuildStageSystems(zurvan.PreUpdateStage,
			&core.NetworkReceiveSystem{},
			&character.JoinWorldSystem{},
			&character.ExperienceSystem{},
			&character.LevelUpSystem{},
		),
	)

	world.AddSystems(
		zurvan.BuildStageSystems(zurvan.FixedUpdateStage,
			&character.StatCalculationSystem{},
			&combat.CastSpellSystem{},
			&combat.FireSpellSystem{},
			&combat.TakeDamageSystem{},
			&combat.TakeHealSystem{},
		),
	)

	world.AddSystems(
		zurvan.BuildStageSystems(zurvan.UpdateStage,
			&combat.CastingSpellSystem{},
			&combat.TakeOverTimeSystem{},
			&combat.RegenerationSystem{},
		),
	)

	world.AddSystems(
		zurvan.BuildStageSystems(zurvan.PostUpdateStage,
			&core.NetworkFlushSystem{},
		),
	)

	world.Run()
}
