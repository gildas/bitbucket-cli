package pipeline_test

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"bitbucket.org/gildas_cherruel/bb/cmd/commit"
	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/pipeline"
	"bitbucket.org/gildas_cherruel/bb/cmd/user"
	"github.com/gildas/go-logger"
	"github.com/google/uuid"
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

func (suite *PipelineSuite) UnmarshalData(filename string, v any) error {
	data := suite.LoadTestData(filename)
	suite.Logger.Infof("Loaded %s: %s", filename, string(data))
	return json.Unmarshal(data, v)
}

// *****************************************************************************

func (suite *PipelineSuite) TestCanUnmarshal() {
	var pl pipeline.Pipeline
	err := suite.UnmarshalData("pipeline.json", &pl)
	suite.Require().NoError(err)
	suite.Assert().Equal("{a1b2c3d4-e5f6-7890-abcd-ef1234567890}", pl.ID.String())
	suite.Assert().Equal(uint64(42), pl.BuildNumber)
	suite.Assert().Equal("COMPLETED", pl.State.Name)
	suite.Require().NotNil(pl.State.Result)
	suite.Assert().Equal("SUCCESSFUL", pl.State.Result.Name)
	suite.Assert().Equal(330*time.Second, pl.Duration)
	suite.Assert().Equal(uuid.MustParse("12345678-1234-1234-1234-123456789012"), uuid.UUID(pl.Creator.ID))
	suite.Assert().Equal("557058:12345678-abcd-efgh-ijkl-123456789012", pl.Creator.AccountID)
	suite.Assert().Equal("John Developer", pl.Creator.Name)
	suite.Assert().Equal("johnd", pl.Creator.Nickname)
	suite.Assert().Equal("myworkspace/my-repo", pl.Repository.FullName)

	suite.Require().NotNil(pl.Target)
	suite.Assert().Equal("main", pl.Target.GetDestination())
	suite.Assert().Equal("abc123def456", pl.Target.GetCommit().Hash)

	target, ok := pl.Target.(*pipeline.ReferenceTarget)
	suite.Require().True(ok)
	suite.Assert().Equal("branch", target.ReferenceType)
	suite.Assert().Equal("main", target.ReferenceName)
	suite.Assert().Equal("abc123def456", target.Commit.Hash)
}

func (suite *PipelineSuite) TestCanUnmarshalWithPullRequest() {
	var pl pipeline.Pipeline
	err := suite.UnmarshalData("pipeline-pullrequest.json", &pl)
	suite.Require().NoError(err)
	suite.Assert().Equal("{a1b2c3d4-e5f6-7890-abcd-ef1234567890}", pl.ID.String())
	suite.Assert().Equal(uint64(42), pl.BuildNumber)
	suite.Assert().Equal("COMPLETED", pl.State.Name)
	suite.Require().NotNil(pl.State.Result)
	suite.Assert().Equal("FAILED", pl.State.Result.Name)
	suite.Assert().Equal(330*time.Second, pl.Duration)
	suite.Assert().Equal(uuid.MustParse("12345678-1234-1234-1234-123456789012"), uuid.UUID(pl.Creator.ID))
	suite.Assert().Equal("557058:12345678-abcd-efgh-ijkl-123456789012", pl.Creator.AccountID)
	suite.Assert().Equal("John Developer", pl.Creator.Name)
	suite.Assert().Equal("johnd", pl.Creator.Nickname)
	suite.Assert().Equal("myworkspace/my-repo", pl.Repository.FullName)

	suite.Require().NotNil(pl.Target)
	suite.Assert().Equal("main", pl.Target.GetDestination())
	suite.Assert().Equal("3c80cde6b371", pl.Target.GetCommit().Hash)

	target, ok := pl.Target.(*pipeline.PullRequestReferenceTarget)
	suite.Require().True(ok)
	suite.Assert().Equal("main", target.Destination)
	suite.Assert().Equal("abc123def456", target.DestinationCommit.Hash)
	suite.Assert().Equal("custom", target.Selector.Type)
	suite.Assert().Equal("run-tests", target.Selector.Pattern)
	suite.Assert().Equal("3c80cde6b371", target.Commit.Hash)
	suite.Assert().Equal(uint64(62), target.PullRequest.ID)
	suite.Assert().Equal("feat: add API key authentication", target.PullRequest.Title)
	suite.Assert().False(target.PullRequest.IsDraft)
}

func (suite *PipelineSuite) TestCanMarshal() {
	expected := suite.LoadTestData("pipeline.json")
	pipeline := &pipeline.Pipeline{
		ID:          common.UUID(uuid.MustParse("a1b2c3d4-e5f6-7890-abcd-ef1234567890")),
		BuildNumber: 42,
		State: pipeline.PipelineState{
			Type: "pipeline_state_completed",
			Name: "COMPLETED",
			Result: &pipeline.PipelineResult{
				Type: "pipeline_state_completed_successful",
				Name: "SUCCESSFUL",
			},
		},
		Target: pipeline.ReferenceTarget{
			Type:          "pipeline_ref_target",
			ReferenceType: "branch",
			ReferenceName: "main",
			Commit:        commit.CommitReference{Hash: "abc123def456"},
			Selector:      &common.Selector{Type: "default"},
		},
		CreatedOn:   time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		CompletedOn: time.Date(2024, 1, 15, 10, 35, 30, 0, time.UTC),
		Duration:    330 * time.Second,
		Creator: user.User{
			Type:      "user",
			ID:        common.UUID(uuid.MustParse("12345678-1234-1234-1234-123456789012")),
			AccountID: "557058:12345678-abcd-efgh-ijkl-123456789012",
			Name:      "John Developer",
			Nickname:  "johnd",
			Links: common.Links{
				Self:   &common.Link{HREF: url.URL{Scheme: "https", Host: "api.bitbucket.org", Path: "/2.0/users/{12345678-1234-1234-1234-123456789012}"}},
				Avatar: &common.Link{HREF: url.URL{Scheme: "https", Host: "secure.gravatar.com", Path: "/avatar/abc123"}},
				HTML:   &common.Link{HREF: url.URL{Scheme: "https", Host: "bitbucket.org", Path: "/{12345678-1234-1234-1234-123456789012}/"}},
			},
		},
		Repository: pipeline.Repository{
			Type:     "repository",
			UUID:     "{12854e6a-e6f8-44ac-b006-1521931d4c0d}",
			Name:     "my-repo",
			FullName: "myworkspace/my-repo",
			Links: common.Links{
				Self:   &common.Link{HREF: url.URL{Scheme: "https", Host: "api.bitbucket.org", Path: "/2.0/repositories/myworkspace/my-repo"}},
				HTML:   &common.Link{HREF: url.URL{Scheme: "https", Host: "bitbucket.org", Path: "/myworkspace/my-repo"}},
				Avatar: &common.Link{HREF: url.URL{Scheme: "https", Host: "bytebucket.org", Path: "/ravatar/{12854e6a-e6f8-44ac-b006-1521931d4c0d}"}},
			},
		},
		Links: common.Links{
			Self:  &common.Link{HREF: url.URL{Scheme: "https", Host: "api.bitbucket.org", Path: "/2.0/repositories/myworkspace/my-repo/pipelines/{a1b2c3d4-e5f6-7890-abcd-ef1234567890}"}},
			Steps: &common.Link{HREF: url.URL{Scheme: "https", Host: "api.bitbucket.org", Path: "/2.0/repositories/myworkspace/my-repo/pipelines/{a1b2c3d4-e5f6-7890-abcd-ef1234567890}/steps/"}},
		},
	}

	data, err := json.Marshal(pipeline)
	suite.Require().NoError(err)
	suite.Assert().NotEmpty(data)
	suite.Require().JSONEq(string(expected), string(data))
}

func (suite *PipelineSuite) TestPipelineString() {
	p := pipeline.Pipeline{BuildNumber: 123}
	suite.Assert().Equal("#123", p.String())
}

func (suite *PipelineSuite) TestPipelineStateWithoutResult() {
	payload := []byte(`{
		"type": "pipeline",
		"uuid": "{a1b2c3d4-e5f6-7890-abcd-ef1234567890}",
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
