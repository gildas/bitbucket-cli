package profile_test

import "bitbucket.org/gildas_cherruel/bb/cmd/profile"

func (suite *ProfileSuite) TestCanUnmarshalErrorAboutPrivileges() {
	var bberr profile.BitBucketError

	err := suite.UnmarshalData("error-privileges.json", &bberr)
	suite.Require().NoError(err)
	suite.Assert().Equal("error", bberr.Type)
	suite.Assert().Equal("Your credentials lack one or more required privilege scopes.", bberr.Message)
	suite.Require().Len(bberr.Fields, 2)
	suite.Require().Contains(bberr.Fields, "required")
	suite.Assert().Contains(bberr.Fields["required"], "project")
	suite.Require().Contains(bberr.Fields, "granted")
	suite.Assert().Contains(bberr.Fields["granted"], "account")
}

func (suite *ProfileSuite) TestCanUnmarshalErrorAboutNoAPI() {
	var bberr profile.BitBucketError

	err := suite.UnmarshalData("error-noapi.json", &bberr)
	suite.Require().NoError(err)
	suite.Assert().Equal("error", bberr.Type)
	suite.Assert().Equal("Resource not found", bberr.Message)
	suite.Assert().Equal("There is no API hosted at this URL", bberr.Detail)
}

func (suite *ProfileSuite) TestCanUnmarshalErrorAboutBadRequest() {
	var bberr profile.BitBucketError

	err := suite.UnmarshalData("error-badrequest.json", &bberr)
	suite.Require().NoError(err)
	suite.Assert().Equal("error", bberr.Type)
	suite.Assert().Equal("Bad request", bberr.Message)
	suite.Require().Len(bberr.Fields, 1)
	suite.Require().Contains(bberr.Fields, "links.avatar")
	suite.Assert().Contains(bberr.Fields["links.avatar"], "required key not provided")
}
