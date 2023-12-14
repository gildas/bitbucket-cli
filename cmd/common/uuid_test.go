package common_test

import (
	"encoding/json"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
)

func (suite *CommonSuite) TestCanCreateUUID() {
	uuid := common.NewUUID()
	suite.Require().NotNil(uuid)
}

func (suite *CommonSuite) TestCanMarshalUUID() {
	expected := "{c32f719b-6c8a-4c87-93e2-9ba8f5cd90dd}" // a Bitbucket String for UUIDs
	uuid, err := common.ParseUUID(expected)
	suite.Require().NoError(err)
	suite.Require().NotNil(uuid)
	suite.False(uuid.IsNil())
	payload, err := json.Marshal(uuid)
	suite.Require().NoError(err)
	suite.Require().NotNil(payload)
	suite.Equal(`"`+expected+`"`, string(payload))
}

func (suite *CommonSuite) TestCanUnmarshalUUID() {
	expected := "{c32f719b-6c8a-4c87-93e2-9ba8f5cd90dd}" // a Bitbucket String for UUIDs
	var uuid common.UUID
	err := json.Unmarshal([]byte(`"`+expected+`"`), &uuid)
	suite.Require().NoError(err)
	suite.Require().NotNil(uuid)
	suite.False(uuid.IsNil())
	suite.Equal(expected, uuid.String())

	err = json.Unmarshal([]byte(`""`), &uuid)
	suite.Require().NoError(err)
	suite.Require().NotNil(uuid)
	suite.True(uuid.IsNil())
}
