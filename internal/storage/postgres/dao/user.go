package dao

import "review-assigner/internal/model"

// User maps to 'users' table.
type User struct {
	ID       string `db:"id"`
	Username string `db:"username"`
	TeamName string `db:"team_name"`
	IsActive bool   `db:"is_active"`
}

func (u User) ToModel() model.User {
	return model.User{
		UserID:   u.ID,
		Username: u.Username,
		TeamName: u.TeamName,
		IsActive: u.IsActive,
	}
}
