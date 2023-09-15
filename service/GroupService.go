package service

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/rahulzzore/realworld-aws-lambda-dynamodb-go/model"
)

func PutGroup(group *model.Group) error {
	// TODO: validate name/description

	const maxAttempt = 5

	// Try to find a unique group id
	for attempt := 0; ; attempt++ {
		err := putGroupWithRandomId(group)

		if err == nil {
			return nil
		}

		if attempt >= maxAttempt {
			return err
		}

		if !IsConditionalCheckFailed(err) {
			return err
		}

		GroupIdRand.RenewSeed()
	}
}

func putGroupWithRandomId(group *model.Group) error {
	group.Id = (string)(1 + GroupIdRand.Get().Int63n(model.MaxGroupId-1)) // range: [1, MaxGroupId)

	groupItem, err := dynamodbattribute.MarshalMap(group)
	if err != nil {
		return err
	}

	transactItems := []*dynamodb.TransactWriteItem{&dynamodb.TransactWriteItem{
		Put: &dynamodb.Put{
			TableName:           aws.String(GroupTableName),
			Item:                groupItem,
			ConditionExpression: aws.String("attribute_not_exists(Id)"),
		},
	}}

	_, err = DynamoDB().TransactWriteItems(&dynamodb.TransactWriteItemsInput{
		TransactItems: transactItems,
	})

	return err
}
