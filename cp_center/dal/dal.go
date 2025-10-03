package dal

import (
	"context"

	"github.com/GameLaunchPad/game_management_project/constdef"
	"github.com/yitter/idgenerator-go/idgen"
)

func InitClient(ctx context.Context) {
	initIDGenerator(ctx)
}

func initIDGenerator(ctx context.Context) {
	var options = idgen.NewIdGeneratorOptions(constdef.IDWorkers)
	idgen.SetIdGenerator(options)
}
