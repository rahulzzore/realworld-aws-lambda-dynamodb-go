package service

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/rahulzzore/realworld-aws-lambda-dynamodb-go/model"
)

func PutPermission(permission model.Permission) error {
	err := permission.ValidatePermission()
	if err != nil {
		return err
	}

	permissionItem, err := dynamodbattribute.MarshalMap(permission)
	if err != nil {
		return err
	}

	transaction := dynamodb.TransactWriteItemsInput{
		TransactItems: []*dynamodb.TransactWriteItem{
			{
				Put: &dynamodb.Put{
					TableName:           aws.String(PermissionTableName),
					Item:                permissionItem,
					ConditionExpression: aws.String("attribute_not_exists(PrincipalId)"),
				},
			},
		},
	}

	_, err = DynamoDB().TransactWriteItems(&transaction)
	if err != nil {
		return err
	}

	return nil
}
