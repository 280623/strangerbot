package repository

import (
	"context"
	"log"

	"github.com/jinzhu/gorm"
	"strangerbot/repository/model"
)

func (p *Repository) GetMenuList(ctx context.Context, parentId int64) (model.Menus, error) {

	q := p.db.Where("parent_id = ? AND is_del = 0", parentId)

	var list []*model.Menu

	if err := q.Model(&model.Menu{}).Order("sort asc").Find(&list).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Println(err.Error())
		return nil, err
	}

	return list, nil
}

func (p *Repository) GetMenu(ctx context.Context, menuId int64) (*model.Menu, error) {

	po := &model.Menu{}

	if err := p.db.Where("id = ? AND is_del = 0", menuId).First(&po).Error; err != nil {
		return nil, err
	}

	return po, nil
}
func (p *Repository) GetMenuByQuestionId(ctx context.Context, questionId int64) (*model.Menu, error) {

	po := &model.Menu{}

	if err := p.db.Where("question_id = ? AND is_del = 0", questionId).First(&po).Error; err != nil {
		return nil, err
	}

	return po, nil
}
