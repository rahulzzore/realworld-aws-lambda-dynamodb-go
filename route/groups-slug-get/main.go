package main

import (
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rahulzzore/realworld-aws-lambda-dynamodb-go/model"
	"github.com/rahulzzore/realworld-aws-lambda-dynamodb-go/service"
	"github.com/rahulzzore/realworld-aws-lambda-dynamodb-go/util"
)

type Response struct {
	Group GroupResponse `json:"group"`
}

type GroupResponse struct {
	Id          string               `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Owner       string               `json:"owner"`
	CreatedAt   string               `json:"createdAt"`
	UpdatedAt   string               `json:"updatedAt"`
	Memberships []MembershipResponse `json:"memberships"`
}

type MembershipResponse struct {
	UserId    string    `json:"userId"`
	JoinedAt  string    `json:"joinedAt"`
}

func Handle(input events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	group, err := service.GetGroupByGroupId(input.PathParameters["slug"])
	if err != nil {
		return util.NewErrorResponse(err)
	}

	memberships, err := service.GetMembershipsbyGroupId(group.Id)
	if err != nil {
		return util.NewErrorResponse(err)
	}

	membershipResponses := make([]MembershipResponse, 0, len(memberships))

	for _, membership := range memberships {
		membershipResponses = append(membershipResponses, MembershipResponse{
			UserId:    membership.UserId,
			JoinedAt:  time.Unix(0, membership.JoinedAt).Format(model.TimestampFormat),
		})
	}

	response := Response{
		Group: GroupResponse{
			Id:            group.Id,
			Name:          group.Name,
			Description:   group.Description,
			CreatedAt:     time.Unix(0, group.CreatedAt).Format(model.TimestampFormat),
			UpdatedAt:     time.Unix(0, group.UpdatedAt).Format(model.TimestampFormat),
			Owner:         group.Owner,
			Memberships:   membershipResponses,
		},
	}

	return util.NewSuccessResponse(200, response)
}

func main() {
	lambda.Start(Handle)
}
