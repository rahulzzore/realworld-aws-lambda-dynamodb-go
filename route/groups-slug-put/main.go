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
	// TODO: use AVP for authorization (when {context.authenticated == true})
	if err != nil {
		return util.NewUnauthorizedResponse()
	}

	oldGroup, err := service.GetGroupByGroupId(input.PathParameters["slug"])
	if err != nil {
		return util.NewErrorResponse(err)
	}

	// TODO: use AVP for authorization
	if user.Username != oldGroup.Owner {
		util.NewUnauthorizedResponse()
	}

	request := Request{}
	err = json.Unmarshal([]byte(input.Body), &request)
	if err != nil {
		return util.NewErrorResponse(err)
	}

	newGroup := createNewGroup(request, oldGroup)

	err = service.UpdateGroup(oldGroup, &newGroup)
	if err != nil {
		return util.NewErrorResponse(err)
	}

	response := Response{
		Group: GroupResponse{
			Id:            newGroup.Id,
			Name:          newGroup.Name,
			Description:   newGroup.Description,
			CreatedAt:     time.Unix(0, newGroup.CreatedAt).Format(model.TimestampFormat),
			UpdatedAt:     time.Unix(0, newGroup.UpdatedAt).Format(model.TimestampFormat),
			Owner:         newGroup.Owner,
		},
	}

	return util.NewSuccessResponse(200, response)
}

func createNewGroup(request Request, oldGroup model.Group) model.Group {
	newGroup := model.Group{
		Id:             oldGroup.Id,
		Name:           request.Group.Name,
		Description:    request.Group.Description,
		CreatedAt:      oldGroup.CreatedAt,
		UpdatedAt:      time.Now().UTC().UnixNano(),
		Owner:          oldGroup.Owner,
	}

	if newGroup.Name == "" {
		newGroup.Name = oldGroup.Name
	}

	if newGroup.Description == "" {
		newGroup.Description = oldGroup.Description
	}

	return newGroup
}

func main() {
	lambda.Start(Handle)
}
