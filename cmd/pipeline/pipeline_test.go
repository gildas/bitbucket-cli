package pipeline_test

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"bitbucket.org/gildas_cherruel/bb/cmd/pipeline"
	"github.com/gildas/go-logger"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
)

type PipelineSuite struct {
	suite.Suite
	Name   string
	Logger *logger.Logger
	Start  time.Time
}

func TestPipelineSuite(t *testing.T) {
	suite.Run(t, new(PipelineSuite))
}

// *****************************************************************************
// Suite Tools

func (suite *PipelineSuite) SetupSuite() {
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

func (suite *PipelineSuite) TearDownSuite() {
	suite.Logger.Debugf("Tearing down")
	if suite.T().Failed() {
		suite.Logger.Warnf("At least one test failed, we are not cleaning")
		suite.T().Log("At least one test failed, we are not cleaning")
	}
	suite.Logger.Infof("Suite End: %s %s", suite.Name, strings.Repeat("=", 80-12-len(suite.Name)))
}

func (suite *PipelineSuite) BeforeTest(suiteName, testName string) {
	suite.Logger.Infof("Test Start: %s %s", testName, strings.Repeat("-", 80-13-len(testName)))
	suite.Start = time.Now()
}

func (suite *PipelineSuite) AfterTest(suiteName, testName string) {
	duration := time.Since(suite.Start)
	if suite.T().Failed() {
		suite.Logger.Errorf("Test %s failed", testName)
	}
	suite.Logger.Record("duration", duration.String()).Infof("Test End: %s %s", testName, strings.Repeat("-", 80-11-len(testName)))
}

func (suite *PipelineSuite) LoadTestData(filename string) []byte {
	data, err := os.ReadFile(fmt.Sprintf("../../testdata/%s", filename))
	if err != nil {
		suite.T().Fatal(err)
	}
	return data
}

func (suite *PipelineSuite) UnmarshalData(filename string, v interface{}) error {
	data := suite.LoadTestData(filename)
	suite.Logger.Infof("Loaded %s: %s", filename, string(data))
	return json.Unmarshal(data, v)
}

// *****************************************************************************

func (suite *PipelineSuite) TestCanUnmarshal() {
	payload := suite.LoadTestData("pipeline.json")
	var p pipeline.Pipeline
	err := json.Unmarshal(payload, &p)
	suite.Require().NoError(err)
	suite.Require().NotNil(p)
	suite.Assert().Equal("{a1b2c3d4-e5f6-7890-abcd-ef1234567890}", p.UUID)
	suite.Assert().Equal(42, p.BuildNumber)
	suite.Assert().Equal("COMPLETED", p.State.Name)
	suite.Assert().NotNil(p.State.Result)
	suite.Assert().Equal("SUCCESSFUL", p.State.Result.Name)
	suite.Assert().Equal("branch", p.Target.RefType)
	suite.Assert().Equal("main", p.Target.RefName)
	suite.Assert().Equal("abc123def456", p.Target.Commit.Hash)
	suite.Assert().Equal(330, p.DurationInSeconds)
	suite.Assert().Equal("John Developer", p.Creator.Name)
	suite.Assert().Equal("myworkspace/my-repo", p.Repository.FullName)
}

func (suite *PipelineSuite) TestCanMarshal() {
	payload := suite.LoadTestData("pipeline.json")
	var p pipeline.Pipeline
	err := json.Unmarshal(payload, &p)
	suite.Require().NoError(err)

	data, err := json.Marshal(p)
	suite.Require().NoError(err)
	suite.Assert().NotEmpty(data)

	// Verify we can unmarshal the marshaled data back
	var p2 pipeline.Pipeline
	err = json.Unmarshal(data, &p2)
	suite.Require().NoError(err)
	suite.Assert().Equal(p.UUID, p2.UUID)
	suite.Assert().Equal(p.BuildNumber, p2.BuildNumber)
	suite.Assert().Equal(p.State.Name, p2.State.Name)
}

func (suite *PipelineSuite) TestPipelineString() {
	p := pipeline.Pipeline{BuildNumber: 123}
	suite.Assert().Equal("#123", p.String())
}

func (suite *PipelineSuite) TestPipelineStateWithoutResult() {
	payload := []byte(`{
		"type": "pipeline",
		"uuid": "{test-uuid}",
		"build_number": 1,
		"state": {
			"type": "pipeline_state_in_progress",
			"name": "IN_PROGRESS"
		},
		"target": {
			"type": "pipeline_ref_target",
			"ref_type": "branch",
			"ref_name": "develop"
		},
		"created_on": "2024-01-15T10:30:00.000000+00:00",
		"duration_in_seconds": 0,
		"creator": {
			"type": "user",
			"display_name": "Test User"
		},
		"repository": {
			"type": "repository",
			"name": "test-repo",
			"full_name": "workspace/test-repo"
		},
		"links": {}
	}`)
	var p pipeline.Pipeline
	err := json.Unmarshal(payload, &p)
	suite.Require().NoError(err)
	suite.Assert().Equal("IN_PROGRESS", p.State.Name)
	suite.Assert().Nil(p.State.Result)
}
