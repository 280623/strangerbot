package repository

import (
	"context"

	"strangerbot/repository/model"
)

func (p *Repository) GetWhiteEmailAll(ctx context.Context) ([]*model.WhiteEmail, error) {

	var list []*model.WhiteEmail
	err := p.db.Find(&list).Error
	if err != nil {
		return nil, err
	}

	return list, nil
}
