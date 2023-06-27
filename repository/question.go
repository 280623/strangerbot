package repository

import (
	"context"
	"log"

	"github.com/jinzhu/gorm"
	"strangerbot/repository/model"
)

func (p *Repository) GetQuestionById(ctx context.Context, questionId int64) (*model.Question, error) {
	po := &model.Question{}

	if err := p.db.Where("id = ? AND is_del = 0", questionId).First(&po).Error; err != nil {
		return nil, err
	}

	return po, nil
}

func (p *Repository) GetAllQuestion(ctx context.Context) ([]*model.Question, error) {
	q := p.db.Where("is_del = 0")

	var list []*model.Question

	if err := q.Model(&model.Question{}).Order("sort asc").Find(&list).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Println(err.Error())
		return nil, err
	}

	return list, nil
}
