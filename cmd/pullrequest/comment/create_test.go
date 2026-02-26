package comment_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"bitbucket.org/gildas_cherruel/bb/cmd/pullrequest/comment"
	"github.com/gildas/go-logger"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
)

type CommentCreateSuite struct {
	suite.Suite
	Name   string
	Logger *logger.Logger
	Start  time.Time
}

func TestCommentCreateSuite(t *testing.T) {
	suite.Run(t, new(CommentCreateSuite))
}

// *****************************************************************************
// Suite Tools

func (suite *CommentCreateSuite) SetupSuite() {
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

func (suite *CommentCreateSuite) TearDownSuite() {
	suite.Logger.Debugf("Tearing down")
	if suite.T().Failed() {
		suite.Logger.Warnf("At least one test failed, we are not cleaning")
		suite.T().Log("At least one test failed, we are not cleaning")
	} else {
		suite.Logger.Infof("All tests succeeded, we are cleaning")
	}
	suite.Logger.Infof("Suite End: %s %s", suite.Name, strings.Repeat("=", 80-12-len(suite.Name)))
}

func (suite *CommentCreateSuite) BeforeTest(suiteName, testName string) {
	suite.Logger.Infof("Test Start: %s %s", testName, strings.Repeat("-", 80-13-len(testName)))
	suite.Start = time.Now()
}

func (suite *CommentCreateSuite) AfterTest(suiteName, testName string) {
	duration := time.Since(suite.Start)
	if suite.T().Failed() {
		suite.Logger.Errorf("Test %s failed", testName)
	}
	suite.Logger.Record("duration", duration.String()).Infof("Test End: %s %s", testName, strings.Repeat("-", 80-11-len(testName)))
}

// *****************************************************************************

func (suite *CommentCreateSuite) TestCanMarshalCommentCreatorWithParent() {
	creator := comment.CommentCreator{
		Content: comment.ContentCreator{Raw: "This is a reply"},
		Parent:  &comment.ParentRef{ID: 759578390},
	}

	data, err := json.Marshal(creator)
	suite.Require().NoError(err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	suite.Require().NoError(err)

	content, ok := result["content"].(map[string]interface{})
	suite.Require().True(ok, "content should be present")
	suite.Assert().Equal("This is a reply", content["raw"])

	parent, ok := result["parent"].(map[string]interface{})
	suite.Require().True(ok, "parent should be present")
	suite.Assert().Equal(float64(759578390), parent["id"])
}

func (suite *CommentCreateSuite) TestCanMarshalCommentCreatorWithoutParent() {
	creator := comment.CommentCreator{
		Content: comment.ContentCreator{Raw: "This is a top-level comment"},
	}

	data, err := json.Marshal(creator)
	suite.Require().NoError(err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	suite.Require().NoError(err)

	content, ok := result["content"].(map[string]interface{})
	suite.Require().True(ok, "content should be present")
	suite.Assert().Equal("This is a top-level comment", content["raw"])

	_, ok = result["parent"]
	suite.Assert().False(ok, "parent should not be present when nil")
}

func (suite *CommentCreateSuite) TestCommentCreatorJSONMatchesBitbucketAPIFormat() {
	creator := comment.CommentCreator{
		Content: comment.ContentCreator{Raw: "Done!"},
		Parent:  &comment.ParentRef{ID: 759578390},
	}

	data, err := json.Marshal(creator)
	suite.Require().NoError(err)

	expected := `{"content":{"raw":"Done!"},"parent":{"id":759578390}}`
	suite.Assert().JSONEq(expected, string(data))
}
