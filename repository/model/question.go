package model

import (
	"fmt"
	"log"

	"strangerbot/repository/gorm_global"
	"strangerbot/vars"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	FRONTEND_TYPE_SELECT       int64 = 1
	FRONTEND_TYPE_MULTI_SELECT int64 = 2
)

const (
	SCENE_TYPE_PROFILE  int64 = 1
	SCENE_TYPE_MATCHING int64 = 2
)

type Question struct {
	gorm_global.ColumnCreateModifyDeleteTime
	SceneType          int64
	HelperTitle        string
	Title              string
	HelperText         string
	FrontendType       int64
	Sort               int64
	MatchingMode       int64
	MatchingQuestionId int64
	MaxMultiLen        int
}

func (q *Question) TableName() string {
	return "bot_question"
}

func (m *Question) GetHelperMessage() string {
	if len(m.HelperText) == 0 {
		return fmt.Sprintf("*%s*\n\n%s", m.HelperTitle, m.Title)
	}
	return fmt.Sprintf("*%s*\n\n%s\n\n%s", m.HelperTitle, m.Title, m.HelperText)
}

func (q *Question) GetKeyboardMarkupFrom(menu *Menu, opts Options, userData []*UserQuestionData) tgbotapi.InlineKeyboardMarkup {

	userDataMap := make(map[int64]*UserQuestionData)
	for i, item := range userData {
		userDataMap[item.OptionId] = userData[i]
	}

	btns := opts.GetKeyboardButton(userDataMap)

	if menu.IsBackEnabled {
		btns = append(btns, menu.GetBackButton())
	}

	return tgbotapi.NewInlineKeyboardMarkup(btns...)
}

type Questions []*Question

func (q Questions) GetQuestionIds() []int64 {

	rs := make([]int64, 0, len(q))
	for _, item := range q {
		rs = append(rs, item.ID)
	}

	return rs
}
func (q Questions) GetProfileQuestion() []*Question {

	rs := make([]*Question, 0, len(q))
	for i, item := range q {
		if item.SceneType == SCENE_TYPE_PROFILE {
			rs = append(rs, q[i])
		}
	}

	return rs
}

func (q Questions) GetMatchingQuestion() []*Question {

	rs := make([]*Question, 0, len(q))
	for i, item := range q {
		if item.SceneType == SCENE_TYPE_MATCHING {
			rs = append(rs, q[i])
		}
	}

	return rs
}

func (q Questions) GetMappingQuestion(matchingQuestion []*Question) []*Question {

	rs := make([]*Question, 0, len(q))
	for _, mq := range matchingQuestion {

		if mq.MatchingQuestionId == 0 {
			continue
		}

		for i, item := range q {
			if item.ID == mq.MatchingQuestionId {
				rs = append(rs, q[i])
			}
		}

	}

	return rs
}

func (q Questions) CheckUserFillFull(userData []*UserQuestionData, isVerify bool) (bool, []*Question) {

	rs := make([]*Question, 0, len(q))
	for i, item := range q {

		if item.ID == vars.VerifyProfileQuestionId {
			if isVerify {
				continue
			}
		}

		found := false
		for _, fill := range userData {
			if fill.QuestionId == item.ID {
				found = true
				break
			}
		}

		if !found {
			rs = append(rs, q[i])
		}
	}

	if len(rs) == 0 {
		return true, nil
	}

	return false, rs
}

func (q Questions) GenMappingQuestion() map[int64]*Question {

	rs := make(map[int64]*Question)
	for i, item := range q {

		rs[item.ID] = q[i]
	}

	return rs
}

func (q Questions) GetMatchingMappingQuestion() map[int64]*Question {

	matchingQuestion := q.GetMatchingQuestion()
	rs := make(map[int64]*Question)
	for _, match := range matchingQuestion {
		for i, item := range q {

			if match.MatchingQuestionId == 0 {
				log.Println("MatchingQuestionId is zero")
				continue
			}

			if item.ID == match.MatchingQuestionId {
				rs[match.ID] = q[i]
			}
		}
	}

	return rs
}

func (q Questions) GetQuestion(questionId int64) *Question {
	for i, item := range q {
		if item.ID == questionId {
			return q[i]
		}
	}
	return nil
}
