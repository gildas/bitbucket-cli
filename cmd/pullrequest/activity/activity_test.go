package activity_test

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"bitbucket.org/gildas_cherruel/bb/cmd/pullrequest/activity"
	"github.com/gildas/go-logger"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
)

type ActivitySuite struct {
	suite.Suite
	Name   string
	Logger *logger.Logger
	Start  time.Time
}

func TestActivitySuite(t *testing.T) {
	suite.Run(t, new(ActivitySuite))
}

// *****************************************************************************
// Suite Tools

func (suite *ActivitySuite) SetupSuite() {
	_ = godotenv.Load()
	suite.Name = strings.TrimSuffix(reflect.TypeOf(suite).Elem().Name(), "Suite")
	suite.Logger = logger.Create("test",
		&logger.FileStream{
			Path:         fmt.Sprintf("./log/test-%s.log", strings.ToLower(suite.Name)),
			Unbuffered:   true,
			SourceInfo:   true,
			FilterLevels: logger.NewLevelSet(logger.TRACE),
		},
	).Child("test", "test")
	suite.Logger.Infof("Suite Start: %s %s", suite.Name, strings.Repeat("=", 80-14-len(suite.Name)))
}

func (suite *ActivitySuite) TearDownSuite() {
	suite.Logger.Debugf("Tearing down")
	if suite.T().Failed() {
		suite.Logger.Warnf("At least one test failed, we are not cleaning")
		suite.T().Log("At least one test failed, we are not cleaning")
	} else {
		suite.Logger.Infof("All tests succeeded, we are cleaning")
	}
	suite.Logger.Infof("Suite End: %s %s", suite.Name, strings.Repeat("=", 80-12-len(suite.Name)))
}

func (suite *ActivitySuite) BeforeTest(suiteName, testName string) {
	suite.Logger.Infof("Test Start: %s %s", testName, strings.Repeat("-", 80-13-len(testName)))
	suite.Start = time.Now()
}

func (suite *ActivitySuite) AfterTest(suiteName, testName string) {
	duration := time.Since(suite.Start)
	if suite.T().Failed() {
		suite.Logger.Errorf("Test %s failed", testName)
	}
	suite.Logger.Record("duration", duration.String()).Infof("Test End: %s %s", testName, strings.Repeat("-", 80-11-len(testName)))
}

func (suite *ActivitySuite) LoadTestData(filename string) []byte {
	data, err := os.ReadFile(fmt.Sprintf("../../../testdata/%s", filename))
	if err != nil {
		suite.T().Fatal(err)
	}
	return data
}

func (suite *ActivitySuite) UnmarshalData(filename string, v interface{}) error {
	data := suite.LoadTestData(filename)
	suite.Logger.Infof("Loaded %s: %s", filename, string(data))
	return json.Unmarshal(data, v)
}

// *****************************************************************************

func (suite *ActivitySuite) TestCanUnmarshalApproval() {
	payload := suite.LoadTestData("activity.json")
	var ac activity.Activity
	err := json.Unmarshal(payload, &ac)
	suite.Require().NoError(err)
	suite.Require().NotNil(ac)
	_, err = json.Marshal(ac)
	suite.Require().NoError(err)
	suite.Assert().NotEmpty(ac.Approval)
	suite.Assert().Empty(ac.Update)
}
func (suite *ActivitySuite) TestCanUnmarshalUpdate() {
	payload := suite.LoadTestData("activity2.json")
	var ac activity.Activity
	err := json.Unmarshal(payload, &ac)
	suite.Require().NoError(err)
	suite.Require().NotNil(ac)
	_, err = json.Marshal(ac)
	suite.Require().NoError(err)
	suite.Assert().Empty(ac.Approval)
	suite.Assert().NotEmpty(ac.Update)
}
