package model

const MaxGroupId = 0x1000000 // exclusive

type Group struct {
	Id          string
	Name        string
	Description string
	CreatedAt   int64
	UpdatedAt   int64
	Owner       string
}

type Membership struct {
	GroupId     string
	UserId      string
	JoinedAt    int64
}

func (group *Group) Validate() error {
	if group.Name == "" {
		return NewInputError("name", "can't be blank")
	}

	if group.Description == "" {
		return NewInputError("description", "can't be blank")
	}

	return nil
}
