package model

import (
	"fmt"
	"sort"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strangerbot/keyboard"
	"strangerbot/repository/gorm_global"
)

const (
	TARGET_TYPE_MENU     int64 = 1 // menu to menu
	TARGET_TYPE_QUESTION int64 = 2 // menu to question
	TARGET_TYPE_COMMAND  int64 = 3 // menu to command
)

type Menu struct {
	gorm_global.ColumnCreateModifyDeleteTime
	Name           string
	ParentId       int64
	QuestionId     int64
	Sort           int64
	RowIndex       int64
	RunCommand     string
	TargetType     int64
	HelperTitle    string
	HelperText     string
	IsBackEnabled  bool
	BackButtonText string
}

func (m *Menu) TableName() string {
	return "bot_menu"
}

func (m *Menu) GetKeyboardButton() tgbotapi.InlineKeyboardButton {

	return tgbotapi.InlineKeyboardButton{
		Text: m.Name,
		CallbackData: keyboard.KeyboardCallbackDataPlus{
			ButtonType:  keyboard.BUTTON_TYPE_MENU,
			ButtonRelId: m.ID,
		}.CallbackData(),
	}

}

func (m *Menu) GetHelperMessage() string {
	return fmt.Sprintf("*%s*\n\n%s", m.HelperTitle, m.HelperText)
}

func (m *Menu) GetBackButton() []tgbotapi.InlineKeyboardButton {
	return []tgbotapi.InlineKeyboardButton{
		{
			Text: m.BackButtonText,
			CallbackData: keyboard.KeyboardCallbackDataPlus{
				ButtonType:   keyboard.BUTTON_TYPE_MENU,
				ButtonRelId:  m.ID,
				IsBackButton: true,
			}.CallbackData(),
		},
	}
}

func (m *Menu) GetSubMenusKeyboardMarkup(subMenus Menus) tgbotapi.InlineKeyboardMarkup {

	btns := subMenus.GetKeyboardButton()

	if m.IsBackEnabled {
		btns = append(btns, m.GetBackButton())
	}

	return tgbotapi.NewInlineKeyboardMarkup(btns...)
}

type Menus []*Menu

func (m Menus) GetKeyboardButton() [][]tgbotapi.InlineKeyboardButton {

	rowMap := make(map[int64][]tgbotapi.InlineKeyboardButton)
	rowIndex := make([]int64, 0, len(m))
	for _, item := range m {
		if _, ok := rowMap[item.RowIndex]; ok {
			rowMap[item.RowIndex] = append(rowMap[item.RowIndex], item.GetKeyboardButton())
		} else {
			rowMap[item.RowIndex] = []tgbotapi.InlineKeyboardButton{
				item.GetKeyboardButton(),
			}
			rowIndex = append(rowIndex, item.RowIndex)
		}
	}

	// sort row index
	sort.Slice(rowIndex, func(i, j int) bool {
		return i < j
	})

	rs := make([][]tgbotapi.InlineKeyboardButton, 0, len(rowMap))
	for _, v := range rowIndex {
		rs = append(rs, rowMap[v])
	}

	return rs
}

func (m Menus) GetKeyboardMarkup() tgbotapi.InlineKeyboardMarkup {

	btn := m.GetKeyboardButton()

	return tgbotapi.NewInlineKeyboardMarkup(btn...)
}
