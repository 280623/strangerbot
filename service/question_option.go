package service

import (
	"context"
	"fmt"

	"strangerbot/keyboard"
	"strangerbot/repository"
	"strangerbot/repository/model"
	"strangerbot/vars"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func ServiceQuestionOption(ctx context.Context, callbackQuery *tgbotapi.CallbackQuery, chatId int64, data *keyboard.KeyboardCallbackDataPlus) ([]*tgbotapi.MessageConfig, []tgbotapi.CallbackConfig, tgbotapi.Chattable, error) {

	repo := repository.GetRepository()

	// get option
	option, err := repo.GetOptionById(ctx, data.ButtonRelId)
	if err != nil {
		return nil, nil, nil, nil
	}

	// change gender
	if option.ID == vars.FemaleMatchRateLimit.OptionId || option.ID == vars.MaleMatchRateLimit.OptionId {
		if !vars.ChangeGenderEnabled {
			userGender, err := repo.GetUserQuestionDataByUserQuestion(ctx, chatId, option.QuestionId)
			if err != nil {
				return nil, nil, nil, err
			}
			if userGender != nil && len(userGender) > 0 {
				if userGender[0].OptionId != option.ID {
					return nil, []tgbotapi.CallbackConfig{tgbotapi.NewCallback("", vars.ChangeGenderErrorMessage)}, nil, nil
				}
			}
		}
	}

	// get question
	question, err := repo.GetQuestionById(ctx, option.QuestionId)
	if err != nil {
		return nil, nil, nil, nil
	}

	switch question.FrontendType {
	case model.FRONTEND_TYPE_SELECT:

		// delete question old option
		err := repo.DeleteUserQuestionDataByQuestion(ctx, chatId, question.ID)
		if err != nil {
			return nil, nil, nil, nil
		}

	case model.FRONTEND_TYPE_MULTI_SELECT:

		userQuestionData, err := repo.GetUserQuestionDataByOptionAndChat(ctx, option.ID, chatId)
		if err != nil {
			return nil, nil, nil, err
		}

		// delete
		if userQuestionData != nil {

			if err := repo.DeleteUserQuestionData(ctx, userQuestionData); err != nil {
				return nil, nil, nil, nil
			}

			// back to menu
			menu, err := repo.GetMenuByQuestionId(ctx, option.QuestionId)
			if err != nil {
				return nil, nil, nil, nil
			}

			opts, err := repo.GetOptionByQuestionId(ctx, menu.QuestionId)
			if err != nil {
				return nil, nil, nil, nil
			}

			if question == nil || len(opts) == 0 {
				return nil, nil, nil, nil
			}

			userQuestionData, err := repo.GetUserQuestionDataByQuestion(ctx, question.ID, chatId)
			if err != nil {
				return nil, nil, nil, nil
			}

			editMsg := tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, question.GetKeyboardMarkupFrom(menu, opts, userQuestionData))

			return nil, nil, editMsg, nil

		}

		// add len
		if question.MaxMultiLen > 0 {
			uq, err := repo.GetUserQuestionDataByUserQuestion(ctx, chatId, question.ID)
			if err != nil {
				return nil, nil, nil, err
			}

			if len(uq) >= question.MaxMultiLen {
				return nil, []tgbotapi.CallbackConfig{tgbotapi.NewCallback("", fmt.Sprintf(vars.QuestionCallbackTest, question.MaxMultiLen))}, nil, nil
			}
		}
	}

	// add
	po := &model.UserQuestionData{
		ChatId:     chatId,
		QuestionId: option.QuestionId,
		OptionId:   option.ID,
		Value:      option.Value,
	}

	if err := repo.UserQuestionDataAdd(ctx, po); err != nil {
		return nil, nil, nil, err
	}

	if question.FrontendType == model.FRONTEND_TYPE_MULTI_SELECT {

		menu, err := repo.GetMenuByQuestionId(ctx, option.QuestionId)
		if err != nil {
			return nil, nil, nil, nil
		}

		opts, err := repo.GetOptionByQuestionId(ctx, menu.QuestionId)
		if err != nil {
			return nil, nil, nil, nil
		}

		if question == nil || len(opts) == 0 {
			return nil, nil, nil, nil
		}

		userQuestionData, err := repo.GetUserQuestionDataByQuestion(ctx, question.ID, chatId)
		if err != nil {
			return nil, nil, nil, nil
		}

		editMsg := tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, question.GetKeyboardMarkupFrom(menu, opts, userQuestionData))

		return nil, nil, editMsg, nil
	}

	// back to menu
	menu, err := repo.GetMenuByQuestionId(ctx, option.QuestionId)
	if err != nil {
		return nil, nil, nil, nil
	}

	msgs, err := serviceBackMenu(ctx, chatId, repo, menu)

	return msgs, nil, nil, err

}
