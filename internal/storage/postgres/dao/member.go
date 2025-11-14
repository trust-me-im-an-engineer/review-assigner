package dao

import "review-assigner/internal/model"

type Member struct {
	UserID   string `db:"id"`
	Username string `db:"username"`
	IsActive bool   `db:"is_active"`
}

func (m Member) ToModel() model.TeamMember {
	return model.TeamMember{
		UserID:   m.UserID,
		Username: m.Username,
		IsActive: m.IsActive,
	}
}
