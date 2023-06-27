package model

type WhiteEmail struct {
	ID    int64  `db:"id"`
	Email string `db:"email"`
}

func (p *WhiteEmail) TableName() string {
	return "e_email_whitelist"
}
