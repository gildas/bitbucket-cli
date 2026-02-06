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
	var pipeline pipeline.Pipeline
	err := suite.UnmarshalData("pipeline.json", &pipeline)
	suite.Require().NoError(err)
	suite.Require().NotNil(pipeline)
	suite.Assert().Equal("{a1b2c3d4-e5f6-7890-abcd-ef1234567890}", pipeline.ID.String())
	suite.Assert().Equal(uint64(42), pipeline.BuildNumber)
	suite.Assert().Equal("COMPLETED", pipeline.State.Name)
	suite.Assert().NotNil(pipeline.State.Result)
	suite.Assert().Equal("SUCCESSFUL", pipeline.State.Result.Name)
	suite.Assert().Equal("branch", pipeline.Target.RefType)
	suite.Assert().Equal("main", pipeline.Target.RefName)
	suite.Assert().Equal("abc123def456", pipeline.Target.Commit.Hash)
	suite.Assert().Equal(330*time.Second, pipeline.Duration)
	suite.Assert().Equal("John Developer", pipeline.Creator.Name)
	suite.Assert().Equal("myworkspace/my-repo", pipeline.Repository.FullName)
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
		Target: pipeline.Target{
			Type:    "pipeline_ref_target",
			RefType: "branch",
			RefName: "main",
			Commit: &commit.Commit{
				Hash: "abc123def456",
			},
			Selector: &pipeline.Selector{
				Type: "default",
			},
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
			UUID:     "{repo-uuid-1234-5678-abcd}",
			Name:     "my-repo",
			FullName: "myworkspace/my-repo",
			Links: common.Links{
				Self:   &common.Link{HREF: url.URL{Scheme: "https", Host: "api.bitbucket.org", Path: "/2.0/repositories/myworkspace/my-repo"}},
				HTML:   &common.Link{HREF: url.URL{Scheme: "https", Host: "bitbucket.org", Path: "/myworkspace/my-repo"}},
				Avatar: &common.Link{HREF: url.URL{Scheme: "https", Host: "bytebucket.org", Path: "/ravatar/{repo-uuid}"}},
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

func (suite *PipelineSuite) TestCanUnmarshalPRPipeline() {
	var p pipeline.Pipeline
	err := suite.UnmarshalData("pipeline-pr.json", &p)
	suite.Require().NoError(err)
	suite.Require().NotNil(p)
	suite.Assert().Equal("{b2c3d4e5-f6a7-8901-bcde-f12345678901}", p.ID.String())
	suite.Assert().Equal(uint64(2808), p.BuildNumber)
	suite.Assert().Equal("COMPLETED", p.State.Name)
	suite.Assert().NotNil(p.State.Result)
	suite.Assert().Equal("FAILED", p.State.Result.Name)
	suite.Assert().Equal("pipeline_pullrequest_target", p.Target.Type)
	// RefName should be populated from source branch for display compatibility
	suite.Assert().Equal("frontend-develop-non-delete-key", p.Target.RefName)
	// Source and Destination are strings (branch names)
	suite.Assert().Equal("frontend-develop-non-delete-key", p.Target.Source)
	suite.Assert().Equal("main", p.Target.Destination)
	// DestinationCommit should be populated
	suite.Require().NotNil(p.Target.DestinationCommit)
	suite.Assert().Equal("8dc910c779d5", p.Target.DestinationCommit.Hash)
	// Commit (source commit) should be populated
	suite.Require().NotNil(p.Target.Commit)
	suite.Assert().Equal("3c80cde6b371", p.Target.Commit.Hash)
	// PullRequest should be populated
	suite.Require().NotNil(p.Target.PullRequest)
	suite.Assert().Equal(62, p.Target.PullRequest.ID)
	suite.Assert().Equal("feat: add API key authentication", p.Target.PullRequest.Title)
	// Selector should be populated (PR pipelines can have custom selectors too)
	suite.Require().NotNil(p.Target.Selector)
	suite.Assert().Equal("custom", p.Target.Selector.Type)
	suite.Assert().Equal("run-tests", p.Target.Selector.Pattern)
	suite.Assert().Equal(630*time.Second, p.Duration)
}

func (suite *PipelineSuite) TestCanUnmarshalPRPipelineGetRowBranch() {
	var p pipeline.Pipeline
	err := suite.UnmarshalData("pipeline-pr.json", &p)
	suite.Require().NoError(err)
	// GetRow with "branch" header should return the PR source branch
	row := p.GetRow([]string{"branch"})
	suite.Require().Len(row, 1)
	suite.Assert().Equal("frontend-develop-non-delete-key", row[0])
}

func (suite *PipelineSuite) TestTargetUnmarshalUnknownType() {
	payload := []byte(`{
		"type": "pipeline_unknown_target",
		"ref_type": "branch",
		"ref_name": "main"
	}`)
	var target pipeline.Target
	err := json.Unmarshal(payload, &target)
	suite.Assert().Error(err)
}

func (suite *PipelineSuite) TestTargetGetTypeDefault() {
	target := pipeline.Target{}
	suite.Assert().Equal("pipeline_ref_target", target.GetType())
}

func (suite *PipelineSuite) TestTargetGetTypePreserved() {
	target := pipeline.Target{Type: "pipeline_pullrequest_target"}
	suite.Assert().Equal("pipeline_pullrequest_target", target.GetType())
}
