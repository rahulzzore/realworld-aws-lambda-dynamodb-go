package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rahulzzore/realworld-aws-lambda-dynamodb-go/model"
	"github.com/rahulzzore/realworld-aws-lambda-dynamodb-go/service"
	"github.com/rahulzzore/realworld-aws-lambda-dynamodb-go/util"
)

type Request struct {
	Permission PermissionRequest `json:"permission"`
}

type PermissionRequest struct {
	PrincipalId string `json:"principalId"`
	ArticleId   int64  `json:"articleId"`
	AccessLevel string `json:"accessLevel"`
	AVPPolicyId string `json:"avpPolicyId"`
}

func Handle(input events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	_, _, err := service.GetCurrentUser(input.Headers["Authorization"])
	if err != nil {
		return util.NewUnauthorizedResponse()
	}

	request := Request{}
	err = json.Unmarshal([]byte(input.Body), &request)
	if err != nil {
		return util.NewErrorResponse(err)
	}

	permission := model.Permission{
		PrincipalId: request.Permission.PrincipalId,
		ArticleId:   request.Permission.ArticleId,
		AccessLevel: request.Permission.AccessLevel,
		AVPPolicyId: request.Permission.AVPPolicyId,
	}

	err = service.PutPermission(permission)
	if err != nil {
		return util.NewErrorResponse(err)
	}

	return util.NewSuccessResponse(200, request)
}

func main() {
	lambda.Start(Handle)
}
