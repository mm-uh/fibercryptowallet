package skycoin

import (
	"github.com/fibercrypto/fibercryptowallet/src/util/logging"
	"os"
	"testing"

	"github.com/fibercrypto/fibercryptowallet/src/core"
	util "github.com/fibercrypto/fibercryptowallet/src/util"
	"github.com/stretchr/testify/mock"
)

var global_mock *SkycoinApiMock

var logModelTest = logging.MustGetLogger("Skycoin Model Test")
// CleanGlobalMock util when is needed to change the values of an

// API method used in other test with different values.
func CleanGlobalMock() {
	global_mock.ExpectedCalls = []*mock.Call{}
}

//Prepare the mock API for all test
func TestMain(m *testing.M) {
	if global_mock == nil {
		global_mock = new(SkycoinApiMock)
	}
	err := core.GetMultiPool().CreateSection(PoolSection, global_mock)
	if err != nil {
		logModelTest.WithError(err).Warn("Error creating pool section")
		return
	}
	util.RegisterAltcoin(NewSkyFiberPlugin(SkycoinMainNetParams))
	os.Exit(m.Run())
}
