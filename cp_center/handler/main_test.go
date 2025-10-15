package handler_test

import (
	"os"
	"testing"

	"github.com/yitter/idgenerator-go/idgen"
)

// TestMain 是一个特殊的函数，Go 的测试框架会在执行包中任何测试之前先调用它。
// 我们用它来做一些全局的设置，比如初始化 idgenerator。
func TestMain(m *testing.M) {
	// 1. 设置 IdGenerator 的选项，这里的 workerId 设为 1 即可，对于测试不重要。
	options := idgen.NewIdGeneratorOptions(1)
	idgen.SetIdGenerator(options)

	// 2. idgen 初始化完成后，调用 m.Run() 来执行包里所有的单元测试。
	//    这行代码是必须的，它会返回一个退出码。
	exitCode := m.Run()

	// 3. 退出测试。
	os.Exit(exitCode)
}
