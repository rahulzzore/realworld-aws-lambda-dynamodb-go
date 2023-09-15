package service

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/rahulzzore/realworld-aws-lambda-dynamodb-go/model"
)

func GetMembershipsbyGroupId(groupId string) ([]model.Membership, error) {
	if groupId == "" {
		return nil, model.NewInputError("groupId", "not specified")
	}

	const limit = 100 // arbitrary limit for now

	// warn: not familiar with nosql, might not be the right way to query this bshrug
	queryMemberships := dynamodb.QueryInput{
		TableName:                 aws.String(MembershipsTableName),
		KeyConditionExpression:    aws.String("GroupId=:groupId"),
		ExpressionAttributeValues: StringKey(":groupId", groupId),
		Limit:                     aws.Int64(int64(limit)),
		//ScanIndexForward:          aws.Bool(false),
	}

	items, err := QueryItems(&queryMemberships, 0, limit)
	if err != nil {
		return nil, err
	}

	memberships := make([]model.Membership, len(items))
	err = dynamodbattribute.UnmarshalListOfMaps(items, &memberships)
	if err != nil {
		return nil, err
	}

	return memberships, nil
}
