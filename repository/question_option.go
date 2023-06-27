package repository

import (
	"context"
	"log"

	"github.com/jinzhu/gorm"
	"strangerbot/repository/model"
)

func (p *Repository) GetOptionByQuestionId(ctx context.Context, questionId int64) ([]*model.QuestionOption, error) {

	q := p.db.Where("question_id = ? AND is_del = 0", questionId)

	var list []*model.QuestionOption

	if err := q.Model(&model.QuestionOption{}).Order("sort asc").Find(&list).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Println(err.Error())
		return nil, err
	}

	return list, nil
}

func (p *Repository) GetOptionById(ctx context.Context, questionId int64) (*model.QuestionOption, error) {
	po := &model.QuestionOption{}

	if err := p.db.Where("id = ? AND is_del = 0", questionId).First(&po).Error; err != nil {
		return nil, err
	}

	return po, nil
}

func (p *Repository) GetAllOption(ctx context.Context) ([]*model.QuestionOption, error) {
	q := p.db.Where("is_del = 0")

	var list []*model.QuestionOption

	if err := q.Model(&model.QuestionOption{}).Find(&list).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Println(err.Error())
		return nil, err
	}

	return list, nil
}
