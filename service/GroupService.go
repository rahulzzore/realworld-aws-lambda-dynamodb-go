package service

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/rahulzzore/realworld-aws-lambda-dynamodb-go/model"
)

func PutGroup(group *model.Group) error {
	err := group.Validate()
	if err != nil {
		return err
	}

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

func GetGroupByGroupId(groupId string) (model.Group, error) {
	group := model.Group{}
	found, err := GetItemByKey(GroupTableName, StringKey("Id", groupId), &group)

	if err != nil {
		return model.Group{}, err
	}

	if !found {
		return model.Group{}, model.NewInputError("id", "not found")
	}

	return group, nil
}

// probs want to redo this without using transactItems but meh
func UpdateGroup(oldGroup model.Group, newGroup *model.Group) error {
	err := newGroup.Validate()
	if err != nil {
		return err
	}

	transactItems := make([]*dynamodb.TransactWriteItem, 0, 1)

	expr, err := buildGroupUpdateExpression(oldGroup, *newGroup)
	if err != nil {
		return err
	}

	// No field changed
	if expr.Update() == nil {
		return nil
	}

	transactItems = append(transactItems, &dynamodb.TransactWriteItem{
		Update: &dynamodb.Update{
			TableName:                 aws.String(GroupTableName),
			Key:                       StringKey("Id", oldGroup.Id),
			ConditionExpression:       aws.String("attribute_exists(Id)"),
			UpdateExpression:          expr.Update(),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
		},
	})

	_, err = DynamoDB().TransactWriteItems(&dynamodb.TransactWriteItemsInput{
		TransactItems: transactItems,
	})
	if err != nil {
		return err
	}

	return nil
}

func buildGroupUpdateExpression(oldGroup model.Group, newGroup model.Group) (expression.Expression, error) {
	update := expression.UpdateBuilder{}

	if oldGroup.Name != newGroup.Name {
		update = update.Set(expression.Name("Name"), expression.Value(newGroup.Name))
	}

	if oldGroup.Description != newGroup.Description {
		update = update.Set(expression.Name("Description"), expression.Value(newGroup.Description))
	}

	if oldGroup.UpdatedAt != newGroup.UpdatedAt {
		update = update.Set(expression.Name("UpdatedAt"), expression.Value(newGroup.UpdatedAt))
	}

	if IsUpdateBuilderEmpty(update) {
		return expression.Expression{}, nil
	}

	builder := expression.NewBuilder().WithUpdate(update)
	return builder.Build()
}

// TODO: delete all associated memberships & permissions in the transaction, delete all associated avp policies
func DeleteGroup(group model.Group) error {
	transactItems := make([]*dynamodb.TransactWriteItem, 0, 1)

	transactItems = append(transactItems, &dynamodb.TransactWriteItem{
		Delete: &dynamodb.Delete{
			TableName:                 aws.String(GroupTableName),
			Key:                       StringKey("GroupId", group.Id),
		},
	})

	_, err := DynamoDB().TransactWriteItems(&dynamodb.TransactWriteItemsInput{
		TransactItems: transactItems,
	})
	if err != nil {
		return err
	}

	return nil
}
