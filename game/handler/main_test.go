package handler

import (
	"os"
	"testing"

	"github.com/GameLaunchPad/game_management_project/game/constdef"
	"github.com/yitter/idgenerator-go/idgen"
)

// setupIDGenerator initializes the ID generator with the specified worker ID.
func setupIDGenerator() {
	var options = idgen.NewIdGeneratorOptions(constdef.IDWorkers)
	idgen.SetIdGenerator(options)
}

// TestMain is the entry point for testing in this package.
func TestMain(m *testing.M) {
	setupIDGenerator()

	exitCode := m.Run()

	os.Exit(exitCode)
}
