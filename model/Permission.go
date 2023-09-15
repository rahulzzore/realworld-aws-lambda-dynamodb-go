package model

const (
	READ    = "READ"
	COMMENT = "COMMENT"
	EDIT    = "EDIT"
)

type Permission struct {
	PrincipalId string
	ArticleId   int64
	AccessLevel string
	AVPPolicyId string
}

func (p *Permission) ValidatePermission() error {
	if !(p.AccessLevel == READ || p.AccessLevel == COMMENT || p.AccessLevel == EDIT) {
		return NewInputError("accesslevel", "can only be READ, COMMENT or EDIT")
	}
	return nil
}
