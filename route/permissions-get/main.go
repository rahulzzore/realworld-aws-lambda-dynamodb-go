package main

import (
	"errors"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rahulzzore/realworld-aws-lambda-dynamodb-go/model"
	"github.com/rahulzzore/realworld-aws-lambda-dynamodb-go/service"
	"github.com/rahulzzore/realworld-aws-lambda-dynamodb-go/util"
)

type Response struct {
	Permission []model.Permission `json:"permissions"`
}

func Handle(input events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	_, _, err := service.GetCurrentUser(input.Headers["Authorization"])
	if err != nil {
		return util.NewUnauthorizedResponse()
	}

	offset, err := strconv.Atoi(input.QueryStringParameters["offset"])
	if err != nil {
		offset = 0
	}

	limit, err := strconv.Atoi(input.QueryStringParameters["limit"])
	if err != nil {
		limit = 20
	}

	articleId, err := strconv.ParseInt(input.PathParameters["articleId"], 10, 64)
	if err != nil {
		return util.NewErrorResponse(errors.New("articleId muust be provided in url query param"))
	}

	permissions, err := service.GetPermissions(offset, limit, articleId)

	response := Response{
		Permission: permissions,
	}

	return util.NewSuccessResponse(200, response)
}

func main() {
	lambda.Start(Handle)
}
