package repository

import (
	"context"
	"log"
	"strings"

	"strangerbot/repository/model"
	"strangerbot/vars"

	"github.com/jinzhu/gorm"
)

func (r *Repository) UserQuestionDataAdd(ctx context.Context, po *model.UserQuestionData) error {

	if err := r.db.Create(&po).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return nil
		}
		return err
	}

	return nil
}

func (p *Repository) GetUserQuestionDataByOptionAndChat(ctx context.Context, optionId int64, chatId int64) (*model.UserQuestionData, error) {

	po := &model.UserQuestionData{}

	if err := p.db.Where("option_id = ? AND chat_id = ? AND is_del = 0", optionId, chatId).First(&po).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return po, nil
}

func (p *Repository) GetUserQuestionDataByQuestion(ctx context.Context, questionId int64, chatId int64) ([]*model.UserQuestionData, error) {

	q := p.db.Where("question_id = ? and chat_id = ? AND is_del = 0", questionId, chatId)

	var list []*model.UserQuestionData

	if err := q.Model(&model.UserQuestionData{}).Find(&list).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Println(err.Error())
		return nil, err
	}

	return list, nil
}

func (p *Repository) DeleteUserQuestionData(ctx context.Context, po *model.UserQuestionData) error {

	if err := p.db.Delete(&po).Error; err != nil {
		return err
	}

	return nil
}

func (p *Repository) DeleteUserQuestionDataByQuestion(ctx context.Context, chatId int64, questionId int64) error {

	if err := p.db.Where("chat_id = ? AND question_id = ?", chatId, questionId).Delete(&model.UserQuestionData{}).Error; err != nil {
		return err
	}

	return nil
}

func (p *Repository) GetUserQuestionDataByUserQuestion(ctx context.Context, chatId int64, questionId int64) ([]*model.UserQuestionData, error) {
	q := p.db.Where("chat_id = ? AND question_id = ? AND is_del = 0", chatId, questionId)

	var list []*model.UserQuestionData

	if err := q.Model(&model.UserQuestionData{}).Find(&list).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Println(err.Error())
		return nil, err
	}

	return list, nil
}

func (p *Repository) GetUserQuestionDataByUser(ctx context.Context, chatId int64) ([]*model.UserQuestionData, error) {
	q := p.db.Where("chat_id = ? AND is_del = 0", chatId)

	var list []*model.UserQuestionData

	if err := q.Model(&model.UserQuestionData{}).Find(&list).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Println(err.Error())
		return nil, err
	}

	return list, nil
}

func (p *Repository) GetUserQuestionDataByUsers(ctx context.Context, chatId []int64) ([]*model.UserQuestionData, error) {

	q := p.db.Where("chat_id IN(?) AND is_del = 0", chatId)

	var list []*model.UserQuestionData

	if err := q.Model(&model.UserQuestionData{}).Find(&list).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Println(err.Error())
		return nil, err
	}

	return list, nil
}

func (p *Repository) LoadUserQuestionDataByUsers(ctx context.Context, chatIds []int64) ([]*model.UserQuestionData, error) {

	list := make([]*model.UserQuestionData, 0, 1000)

	// 定义每页的大小
	pageSize := 500

	// 计算总页数
	totalPages := (len(chatIds) + pageSize - 1) / pageSize

	// 遍历每一页
	for page := 1; page <= totalPages; page++ {

		// 计算当前页应该遍历的元素的起始位置和结束位置
		start := (page - 1) * pageSize
		end := start + pageSize
		if end > len(chatIds) {
			end = len(chatIds)
		}

		readChatIds := chatIds[start:end]

		rs, err := p.GetUserQuestionDataByUsers(ctx, readChatIds)
		if err != nil {
			return nil, err
		}
		if len(rs) > 0 {
			list = append(list, rs...)
		}

	}

	return list, nil
}

func (p *Repository) GetChatByMatching(ctx context.Context, chatId int64, questions model.Questions, options model.Options, userQuestionData model.UserQuestionDataList) ([]int64, model.UserQuestionDataList, error) {

	// matching options data
	matchingQuestion := questions.GetMatchingQuestion()
	profileQuestion := questions.GetMappingQuestion(matchingQuestion)
	matchingOptions := options.GetQuestionOptions(ctx, matchingQuestion)
	profileOptions := options.GetQuestionOptions(ctx, profileQuestion)
	userMatchingData := userQuestionData.GetUserQuestionDataByOptions(ctx, matchingOptions)
	userProfileData := userQuestionData.GetUserQuestionDataByOptions(ctx, profileOptions)

	//selectOptions := make([]int64, 0, len(userMatchingData))

	questionMatchingMap := make(map[int64]bool)
	questionProfileMatchingMap := make(map[int64]bool)
	allMatchingOptions := make([]int64, 0, len(userMatchingData))
	allProfileMatchingOptions := make([]int64, 0, len(userProfileData))
	matchingQuestionNum := 0
	profileMatchingQuestionNum := 0

	// group by question
	for _, item := range userMatchingData {

		option := matchingOptions.GetOption(item.OptionId)
		if option == nil {
			continue
		}

		if option.MatchingOptionId == 0 {
			continue
		}

		if option.QuestionId == vars.MatchingQuestionId {
			continue
		}

		allMatchingOptions = append(allMatchingOptions, option.MatchingOptionId)
		if _, ok := questionMatchingMap[item.QuestionId]; !ok {
			questionMatchingMap[item.QuestionId] = true
			matchingQuestionNum++
		}

	}

	for _, item := range userProfileData {

		// get profile option mapping matching option
		option := matchingOptions.GetOptionByMapping(item.OptionId)
		if option == nil {
			continue
		}

		if option.MatchingOptionId == 0 {
			continue
		}

		if option.QuestionId == vars.VerifyProfileQuestionId {
			continue
		}

		allProfileMatchingOptions = append(allProfileMatchingOptions, option.ID)
		if _, ok := questionProfileMatchingMap[option.QuestionId]; !ok {
			questionProfileMatchingMap[option.QuestionId] = true
			profileMatchingQuestionNum++
		}

	}

	for questionId, _ := range questionProfileMatchingMap {
		anyOption := options.IsHasAnythingOption(questionId)
		if anyOption != nil {
			allProfileMatchingOptions = append(allProfileMatchingOptions, anyOption.ID)
		}
	}

	// build sql
	var sub *gorm.DB

	if vars.RUN_MODE == "debug" {
		sub = p.db.Raw("SELECT chat_id FROM (SELECT chat_id,COUNT(*) AS cnt FROM (SELECT chat_id,question_id FROM bot_user_question_data WHERE option_id IN(?) AND chat_id IN((SELECT chat_id FROM users WHERE available = 1 AND match_chat_id IS NULL)) GROUP BY chat_id,question_id) AS bot_user_question_data GROUP BY chat_id) AS bot_user_question_data WHERE cnt = ?", allMatchingOptions, matchingQuestionNum)
	} else {
		sub = p.db.Raw("SELECT chat_id FROM (SELECT chat_id,COUNT(*) AS cnt FROM (SELECT chat_id,question_id FROM bot_user_question_data WHERE option_id IN(?) AND chat_id IN((SELECT chat_id FROM users WHERE chat_id != ? AND available = 1 AND match_chat_id IS NULL)) GROUP BY chat_id,question_id) AS bot_user_question_data GROUP BY chat_id) AS bot_user_question_data WHERE cnt = ?", allMatchingOptions, chatId, matchingQuestionNum)
	}

	sub = p.db.Raw("SELECT chat_id FROM (SELECT chat_id,COUNT(DISTINCT question_id) AS cnt FROM (SELECT chat_id,question_id FROM bot_user_question_data WHERE chat_id IN(?) AND option_id IN(?)) AS bot_user_question_data GROUP BY chat_id) AS bot_user_question_data WHERE cnt = ?", sub.QueryExpr(), allProfileMatchingOptions, profileMatchingQuestionNum)

	var data []struct {
		ChatId int64
	}

	if err := sub.Scan(&data).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	rs := make([]int64, 0, len(data))
	for _, item := range data {
		rs = append(rs, item.ChatId)
	}

	return rs, userMatchingData, nil
}

func (p *Repository) CheckHasOptionBy(ctx context.Context, chatIds []int64, optionIds []int64) ([]int64, error) {

	q := p.db.Where("chat_id IN(?) AND option_id IN(?) AND is_del = 0", chatIds, optionIds)

	var list []*model.UserQuestionData

	if err := q.Model(&model.UserQuestionData{}).Find(&list).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Println(err.Error())
		return nil, err
	}

	rs := make([]int64, 0, len(list))
	rsMap := make(map[int64]bool, len(list))
	for _, item := range list {

		if _, ok := rsMap[item.ChatId]; ok {
			continue
		}

		rs = append(rs, item.ChatId)
		rsMap[item.ChatId] = true
	}

	return rs, nil
}
