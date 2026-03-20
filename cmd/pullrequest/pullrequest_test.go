package pullrequest_test

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"bitbucket.org/gildas_cherruel/bb/cmd/pullrequest"
	"github.com/gildas/go-logger"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
)

type PullRequestSuite struct {
	suite.Suite
	Name   string
	Logger *logger.Logger
	Start  time.Time
}

func TestPullRequestSuite(t *testing.T) {
	suite.Run(t, new(PullRequestSuite))
}

// *****************************************************************************
// Suite Tools

func (suite *PullRequestSuite) SetupSuite() {
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

func (suite *PullRequestSuite) TearDownSuite() {
	suite.Logger.Debugf("Tearing down")
	if suite.T().Failed() {
		suite.Logger.Warnf("At least one test failed, we are not cleaning")
		suite.T().Log("At least one test failed, we are not cleaning")
	} else {
		suite.Logger.Infof("All tests succeeded, we are cleaning")
	}
	suite.Logger.Infof("Suite End: %s %s", suite.Name, strings.Repeat("=", 80-12-len(suite.Name)))
}

func (suite *PullRequestSuite) BeforeTest(suiteName, testName string) {
	suite.Logger.Infof("Test Start: %s %s", testName, strings.Repeat("-", 80-13-len(testName)))
	suite.Start = time.Now()
}

func (suite *PullRequestSuite) AfterTest(suiteName, testName string) {
	duration := time.Since(suite.Start)
	if suite.T().Failed() {
		suite.Logger.Errorf("Test %s failed", testName)
	}
	suite.Logger.Record("duration", duration.String()).Infof("Test End: %s %s", testName, strings.Repeat("-", 80-11-len(testName)))
}

func (suite *PullRequestSuite) LoadTestData(filename string) []byte {
	data, err := os.ReadFile(fmt.Sprintf("../../testdata/%s", filename))
	if err != nil {
		suite.T().Fatal(err)
	}
	return data
}

func (suite *PullRequestSuite) UnmarshalData(filename string, v interface{}) error {
	data := suite.LoadTestData(filename)
	suite.Logger.Infof("Loaded %s: %s", filename, string(data))
	return json.Unmarshal(data, v)
}

// *****************************************************************************

func (suite *PullRequestSuite) TestCanUnmarshal() {
	payload := suite.LoadTestData("pullrequest.json")
	var pr pullrequest.PullRequest
	err := json.Unmarshal(payload, &pr)
	suite.Require().NoError(err)
	suite.Require().NotNil(pr)
	data, err := json.Marshal(pr)
	suite.Require().NoError(err)
	suite.Assert().JSONEq(string(payload), string(data))
}

func (suite *PullRequestSuite) TestCanUnmarshalWithNilDestinationRepository() {
	payload := suite.LoadTestData("pullrequest-no-dest-repo.json")
	var pr pullrequest.PullRequest
	err := json.Unmarshal(payload, &pr)
	suite.Require().NoError(err)
	suite.Require().NotNil(pr)
	suite.Assert().Nil(pr.Destination.Repository)
	suite.Assert().NotEmpty(pr.Destination.Branch.Name)
}

func (suite *PullRequestSuite) TestDestinationRepositoryIsNilAfterSettingNewDestination() {
	payload := suite.LoadTestData("pullrequest.json")
	var pr pullrequest.PullRequest
	err := json.Unmarshal(payload, &pr)
	suite.Require().NoError(err)
	suite.Require().NotNil(pr.Destination.Repository)

	pr.Destination = pullrequest.Endpoint{Branch: pullrequest.Branch{Name: "new-branch"}}

	suite.Assert().Nil(pr.Destination.Repository)
	suite.Assert().Equal("new-branch", pr.Destination.Branch.Name)
}
