package handler_test

import (
	"os"
	"testing"

	"github.com/yitter/idgenerator-go/idgen"
)

func TestMain(m *testing.M) {
	options := idgen.NewIdGeneratorOptions(1)
	idgen.SetIdGenerator(options)
	exitCode := m.Run()
	os.Exit(exitCode)
}
