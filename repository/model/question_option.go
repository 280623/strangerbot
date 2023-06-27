package model

import (
	"context"
	"fmt"
	"sort"

	"strangerbot/keyboard"
	"strangerbot/repository/gorm_global"
	"strangerbot/vars"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type QuestionOption struct {
	gorm_global.ColumnCreateModifyDeleteTime
	QuestionId       int64
	OptionType       int64
	Label            string
	Value            string
	IsMatchingAny    bool
	MatchingOptionId int64
	Sort             int64
	RowIndex         int64
}

func (q *QuestionOption) TableName() string {
	return "bot_option"
}

func (q *QuestionOption) GetOptionLabel(userData *UserQuestionData) string {

	if userData != nil {
		return fmt.Sprintf("%s %s", vars.ChooseMark, q.Label)
	}

	return q.Label
}

func (q *QuestionOption) GetKeyboardButton(userData *UserQuestionData) tgbotapi.InlineKeyboardButton {

	return tgbotapi.InlineKeyboardButton{
		Text: q.GetOptionLabel(userData),
		CallbackData: keyboard.KeyboardCallbackDataPlus{
			ButtonType:  keyboard.BUTTON_TYPE_OPTION,
			ButtonRelId: q.ID,
		}.CallbackData(),
	}

}

type Options []*QuestionOption

func (m Options) GetKeyboardButton(userDatas map[int64]*UserQuestionData) [][]tgbotapi.InlineKeyboardButton {

	rowMap := make(map[int64][]tgbotapi.InlineKeyboardButton)
	for _, item := range m {

		var userData *UserQuestionData
		if v, ok := userDatas[item.ID]; ok {
			userData = v
		}

		if _, ok := rowMap[item.RowIndex]; ok {
			rowMap[item.RowIndex] = append(rowMap[item.RowIndex], item.GetKeyboardButton(userData))
		} else {
			rowMap[item.RowIndex] = []tgbotapi.InlineKeyboardButton{
				item.GetKeyboardButton(userData),
			}
		}

	}

	rrr := make([]struct {
		Key   int64
		Value []tgbotapi.InlineKeyboardButton
	}, 0, len(rowMap))

	for k, v := range rowMap {
		rrr = append(rrr, struct {
			Key   int64
			Value []tgbotapi.InlineKeyboardButton
		}{Key: k, Value: v})
	}

	sort.Slice(rrr, func(i, j int) bool {
		return rrr[i].Key < rrr[j].Key
	})

	rs := make([][]tgbotapi.InlineKeyboardButton, 0, len(rowMap))
	for _, v := range rrr {
		rs = append(rs, v.Value)
	}

	return rs
}

func (m Options) GetQuestionOptions(ctx context.Context, questions Questions) Options {

	rs := make([]*QuestionOption, 0, len(m))
	for i, item := range m {
		for _, question := range questions {
			if item.QuestionId == question.ID {
				rs = append(rs, m[i])
			}
		}
	}

	return rs
}

func (m Options) GetOptionsByIds(ids []int64) []string {

	rs := make([]string, 0, len(ids))

	for _, id := range ids {
		for i, item := range m {
			if item.ID == id {
				rs = append(rs, m[i].Label)
			}
		}
	}

	return rs
}

func (m Options) GetOption(optionId int64) *QuestionOption {
	for i, item := range m {
		if item.ID == optionId {
			return m[i]
		}
	}
	return nil
}

func (m Options) GetOptionByMapping(mappingOptionId int64) *QuestionOption {
	for i, item := range m {
		if item.MatchingOptionId == mappingOptionId {
			return m[i]
		}
	}
	return nil
}

func (m Options) IsHasAnythingOption(questionId int64) *QuestionOption {
	for i, item := range m {
		if item.QuestionId == questionId {
			if item.IsMatchingAny {
				return m[i]
			}
		}
	}
	return nil
}

func (m Options) GetOptionsByQuestionId(questionId int64) Options {

	rs := make([]*QuestionOption, 0, len(m))
	for i, item := range m {
		if item.QuestionId == questionId {
			rs = append(rs, m[i])
		}
	}

	return rs
}
