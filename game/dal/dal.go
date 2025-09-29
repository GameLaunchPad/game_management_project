package dal

import (
	"context"

	"github.com/yitter/idgenerator-go/idgen"
)

func InitClient(ctx context.Context) {
	initIDGenerator(ctx)
}

func initIDGenerator(ctx context.Context) {
	var options = idgen.NewIdGeneratorOptions(6)
	idgen.SetIdGenerator(options)
}
