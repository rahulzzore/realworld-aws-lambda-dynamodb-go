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
