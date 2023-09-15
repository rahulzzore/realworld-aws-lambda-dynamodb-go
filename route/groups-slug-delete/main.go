package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rahulzzore/realworld-aws-lambda-dynamodb-go/service"
	"github.com/rahulzzore/realworld-aws-lambda-dynamodb-go/util"
)

func Handle(input events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	user, _, err := service.GetCurrentUser(input.Headers["Authorization"])
	// TODO: use AVP for authorization (when {context.authenticated == true})
	if err != nil {
		return util.NewUnauthorizedResponse()
	}

	group, err := service.GetGroupByGroupId(input.PathParameters["slug"])
	if err != nil {
		return util.NewErrorResponse(err)
	}

	// TODO: use AVP for authorization
	if user.Username != group.Owner {
		util.NewUnauthorizedResponse()
	}

	err = service.DeleteGroup(group)
	if err != nil {
		return util.NewErrorResponse(err)
	}

	return util.NewSuccessResponse(204, nil)
}

func main() {
	lambda.Start(Handle)
}
