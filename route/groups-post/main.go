package main

import (
	"encoding/json"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rahulzzore/realworld-aws-lambda-dynamodb-go/model"
	"github.com/rahulzzore/realworld-aws-lambda-dynamodb-go/service"
	"github.com/rahulzzore/realworld-aws-lambda-dynamodb-go/util"
)

type Request struct {
	Group GroupRequest `json:"group"`
}

type GroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Response struct {
	Group GroupResponse `json:"group"`
}

type GroupResponse struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Owner       string `json:"owner"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

func Handle(input events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	user, _, err := service.GetCurrentUser(input.Headers["Authorization"])
	if err != nil {
		return util.NewUnauthorizedResponse()
	}

	request := Request{}
	err = json.Unmarshal([]byte(input.Body), &request)
	if err != nil {
		return util.NewErrorResponse(err)
	}

	now := time.Now().UTC()
	nowUnixNano := now.UnixNano()
	nowStr := now.Format(model.TimestampFormat)

	// TODO: separate out the API interfaces from the ddb entry
	group := model.Group{
		Name:        request.Group.Name,
		Description: request.Group.Description,
		CreatedAt:   nowUnixNano,
		UpdatedAt:   nowUnixNano,
		Owner:       user.Username,
	}

	err = service.PutGroup(&group)
	if err != nil {
		return util.NewErrorResponse(err)
	}

	response := Response{
		Group: GroupResponse{
			Id:          group.Id,
			Name:        group.Name,
			Description: group.Description,
			CreatedAt:   nowStr,
			UpdatedAt:   nowStr,
			Owner:       group.Owner,
		},
	}

	return util.NewSuccessResponse(201, response)
}

func main() {
	lambda.Start(Handle)
}
