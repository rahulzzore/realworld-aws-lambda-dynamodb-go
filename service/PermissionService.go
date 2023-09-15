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

func GetPermissions(offset, limit int, articleId int64) ([]model.Permission, error) {
	queryPermissions := dynamodb.QueryInput{
		TableName:                 aws.String(PermissionTableName),
		IndexName:                 aws.String("ArticleIdIndex"),
		KeyConditionExpression:    aws.String("ArticleId=:articleId"),
		ExpressionAttributeValues: Int64Key(":articleId", articleId),
		Limit:                     aws.Int64(int64(offset + limit)),
		ScanIndexForward:          aws.Bool(false),
	}

	items, err := QueryItems(&queryPermissions, offset, limit)
	if err != nil {
		return nil, err
	}

	permissions := make([]model.Permission, len(items))
	err = dynamodbattribute.UnmarshalListOfMaps(items, &permissions)
	if err != nil {
		return nil, err
	}

	return permissions, nil
}
