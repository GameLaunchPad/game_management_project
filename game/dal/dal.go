package dal

import (
	"context"

	"github.com/GameLaunchPad/game_management_project/game/constdef"
	"github.com/yitter/idgenerator-go/idgen"
)

func InitClient(ctx context.Context) {
	initIDGenerator(ctx)
	initDB()
}

func initIDGenerator(ctx context.Context) {
	var options = idgen.NewIdGeneratorOptions(constdef.IDWorkers)
	idgen.SetIdGenerator(options)
}
