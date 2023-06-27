package service

import (
	"context"

	"strangerbot/keyboard"
	"strangerbot/repository"
	"strangerbot/repository/model"
	"strangerbot/vars"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func ServiceMenu(ctx context.Context, chatId int64, data *keyboard.KeyboardCallbackDataPlus, userVerify bool) ([]*tgbotapi.MessageConfig, error) {

	repo := repository.GetRepository()
	menu, err := repo.GetMenu(ctx, data.ButtonRelId)
	if err != nil {
		return nil, err
	}

	// back pop menu
	if data.IsBackButton {
		return serviceBackMenu(ctx, chatId, repo, menu)
	}

	switch menu.TargetType {
	case model.TARGET_TYPE_MENU:

		// sub menus
		menuList, err := repo.GetMenuList(ctx, menu.ID)
		if err != nil {
			return nil, err
		}

		if len(menuList) == 0 {
			return nil, nil
		}

		msg := tgbotapi.NewMessage(chatId, menu.GetHelperMessage())
		msg.ParseMode = "Markdown"
		msg.ReplyMarkup = menu.GetSubMenusKeyboardMarkup(menuList)

		return []*tgbotapi.MessageConfig{&msg}, nil

	case model.TARGET_TYPE_QUESTION:

		if menu.QuestionId == 0 {
			return nil, nil
		}

		question, err := repo.GetQuestionById(ctx, menu.QuestionId)
		if err != nil {
			return nil, err
		}

		opts, err := repo.GetOptionByQuestionId(ctx, menu.QuestionId)
		if err != nil {
			return nil, err
		}

		if question == nil || len(opts) == 0 {
			return nil, nil
		}

		userQuestionData, err := repo.GetUserQuestionDataByQuestion(ctx, question.ID, chatId)
		if err != nil {
			return nil, err
		}

		if question.ID == vars.VerifyProfileQuestionId {
			if userVerify {
				userQuestionData = []*model.UserQuestionData{
					{
						ChatId:     chatId,
						QuestionId: question.ID,
						OptionId:   vars.VerifyOptionId,
						Value:      "",
					},
				}
			}
		}

		msgs := make([]*tgbotapi.MessageConfig, 0, 2)
		//if len(question.HelperText) > 0 {
		//	msg := tgbotapi.NewMessage(chatId, question.HelperText)
		//	msg.ParseMode = "Markdown"
		//	msgs = append(msgs, &msg)
		//}

		msg := tgbotapi.NewMessage(chatId, question.GetHelperMessage())
		msg.ParseMode = "Markdown"
		msg.ReplyMarkup = question.GetKeyboardMarkupFrom(menu, opts, userQuestionData)
		msgs = append(msgs, &msg)

		return msgs, nil

	}

	return nil, nil
}

func serviceBackMenu(ctx context.Context, chatId int64, repo *repository.Repository, menu *model.Menu) ([]*tgbotapi.MessageConfig, error) {

	menuList, err := repo.GetMenuList(ctx, menu.ParentId)
	if err != nil {
		return nil, err
	}

	if len(menuList) == 0 {
		return nil, nil
	}

	if menu.ParentId > 0 {

		parentMenu, err := repo.GetMenu(ctx, menu.ParentId)
		if err != nil {
			return nil, err
		}

		msg := tgbotapi.NewMessage(chatId, parentMenu.GetHelperMessage())
		msg.ParseMode = "Markdown"
		msg.ReplyMarkup = parentMenu.GetSubMenusKeyboardMarkup(menuList)

		return []*tgbotapi.MessageConfig{&msg}, nil
	}

	msg := tgbotapi.NewMessage(chatId, vars.TopMenuMessage)
	msg.ReplyMarkup = menuList.GetKeyboardMarkup()

	return []*tgbotapi.MessageConfig{&msg}, nil
}
