package configuration_test

import (
	"context"

	"os"
	"testing"

	"main/src/books/infrastructure/configuration"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stretchr/testify/suite"
)

type DynamoDBBookConfigSuite struct {
	suite.Suite
	client    *dynamodb.Client
	ctx       context.Context
	tableName string
}

func (suite *DynamoDBBookConfigSuite) SetupSuite() {
	suite.ctx = context.TODO()
	var err error
	suite.client, err = configuration.GetDynamoDBClient(suite.ctx)
	suite.NoError(err)
	suite.tableName = configuration.GetDynamoDBBookTable()
}

func (suite *DynamoDBBookConfigSuite) TestGetDynamoDBCostTableWithEnvSet() {
	testValue := "MyUserTable"
	os.Setenv("BOOK_TABLE", testValue)
	defer os.Unsetenv("BOOK_TABLE")
	result := configuration.GetDynamoDBBookTable()
	suite.Equal(testValue, result)
}

func (suite *DynamoDBBookConfigSuite) TestGetDynamoDBCostTableWithoutEnvSet() {
	result := configuration.GetDynamoDBBookTable()
	suite.Equal("Test_Book_Table", result)
}

func (suite *DynamoDBBookConfigSuite) TestCreateDescribeDeleteLocalDynamoDBCostTable() {
	tableName := "New_Book_Table_Testing"
	err := configuration.CreateLocalDynamoDBBookTable(suite.ctx, suite.client, tableName)
	suite.NoError(err)

	exists, err := configuration.DescribeBookTable(suite.ctx, suite.client, tableName)
	suite.NoError(err)
	suite.True(exists, "DescribeRoleTable exists")

	err = configuration.DeleteLocalDynamoDBBookTable(suite.ctx, suite.client, tableName)
	suite.NoError(err)
}

func (suite *DynamoDBBookConfigSuite) TestDeleteTableNotExistsDynamoTable() {
	tableName := "NonExistentTable"
	err := configuration.DeleteLocalDynamoDBBookTable(suite.ctx, suite.client, tableName)
	suite.Error(err)
	suite.Contains(err.Error(), "does not exist")
}

func TestDynamoDBBookConfigSuite(t *testing.T) {
	suite.Run(t, new(DynamoDBBookConfigSuite))
}
