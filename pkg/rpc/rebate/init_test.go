package rebate

import (
	"os"
	"testing"

	"dev-gitlab.wanxingrowth.com/fanli/rebate/pkg/utils/test"
)

func TestMain(m *testing.M) {

	test.Init()

	code := m.Run()

	test.Release()

	os.Exit(code)
}
