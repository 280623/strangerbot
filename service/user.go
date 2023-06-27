package service

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"sort"
	"strangerbot/repository/gorm_global"
	"time"

	"strangerbot/repository"
	"strangerbot/repository/model"
	"strangerbot/vars"
)

var (
	ErrUserNotFillAllQuestion = errors.New("user not fill all question")
)

func ServiceCheckUserFillFull(ctx context.Context, chatId int64, isVerify bool) (bool, error) {

	repo := repository.GetRepository()

	// find all question
	questions, err := repo.GetAllQuestion(ctx)
	if err != nil {
		return false, err
	}

	// find user all user question data
	userQuestionData, err := repo.GetUserQuestionDataByUser(ctx, chatId)
	if err != nil {
		return false, err
	}

	// check user fill full
	fillFull, _ := model.Questions(questions).CheckUserFillFull(userQuestionData, isVerify)
	if fillFull {
		return true, nil
	}

	return false, nil
}

type boolgen struct {
	src       rand.Source
	cache     int64
	remaining int
}

func (b *boolgen) Bool() bool {
	if b.remaining == 0 {
		b.cache, b.remaining = b.src.Int63(), 63
	}

	result := b.cache&0x01 == 1
	b.cache >>= 1
	b.remaining--

	return result
}

func New() *boolgen {
	return &boolgen{src: rand.NewSource(time.Now().UnixNano())}
}

func ServiceGlobalMatch(ctx context.Context) ([]*model.MatchUserData, error) {

	repo := repository.GetRepository()

	// find all User
	users, err := repo.LoadAllAvailableUsers(ctx)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, nil
	}

	// Sort the users slice by member status,
	// with members coming before non-members.
	randBool := New()
	sort.Slice(users, func(i, j int) bool {

		if users[i].MemberLevel == users[j].MemberLevel {
			return randBool.Bool()
		}

		return users[i].MemberLevel > users[j].MemberLevel
	})

	// create all user id slice
	chatIds := make([]int64, 0, len(users))
	userIds := make([]int64, 0, len(users))
	chatUserMap := make(map[int64]*model.User)
	for i, u := range users {
		userIds = append(userIds, u.ID)
		chatIds = append(chatIds, u.ChatID)
		chatUserMap[u.ChatID] = users[i]
	}

	// load users data
	userQuestionData, err := repo.LoadUserQuestionDataByUsers(ctx, chatIds)
	if err != nil {
		return nil, err
	}

	if len(userQuestionData) == 0 {
		return nil, nil
	}

	// find all question
	questions, err := repo.GetAllQuestion(ctx)
	if err != nil {
		return nil, err
	}

	// find all question option
	options, err := repo.GetAllOption(ctx)
	if err != nil {
		return nil, err
	}

	profileQuestion := model.Questions(questions).GetProfileQuestion()
	profileQuestionMap := model.Questions(profileQuestion).GenMappingQuestion()
	profileOptions := model.Options(options).GetQuestionOptions(ctx, profileQuestion)
	matchQuestion := model.Questions(questions).GetMatchingQuestion()
	matchQuestionMap := model.Questions(matchQuestion).GenMappingQuestion()
	matchOptions := model.Options(options).GetQuestionOptions(ctx, matchQuestion)

	// assemble matching MatchUserData
	rs := make(map[int64]*model.MatchUserData)

	for i, u := range users {

		// this user all question data
		userData := model.UserQuestionDataList(userQuestionData).GetByChatId(u.ChatID)
		if len(userData) == 0 {
			continue
		}

		// verified user matching
		tmp := &model.MatchUserData{
			ChatId:         u.ChatID,
			User:           users[i],
			MatchChatId:    0,
			VerifyOptionId: model.UserQuestionDataList(userQuestionData).GetFirstOptionIdByQuestionId(vars.MatchingQuestionId),
		}

		profileQuestions := make(map[int64]*model.ProfileQuestion)
		matchQuestions := make(map[int64]*model.MatchingQuestion)
		for _, item := range userData {

			// PersonalInfoQuestions
			if q, ok := profileQuestionMap[item.QuestionId]; ok {

				qo := profileOptions.GetOption(item.OptionId)
				if qo == nil {
					continue
				}

				if _, ok := profileQuestions[q.ID]; !ok {
					profileQuestions[q.ID] = &model.ProfileQuestion{
						ProfileQuestion: q,
						ProfileOptions:  []*model.QuestionOption{qo},
					}
				} else {
					profileQuestions[q.ID].ProfileOptions = append(profileQuestions[q.ID].ProfileOptions, qo)
				}

			}

			// MatchCriteriaQuestions
			if q, ok := matchQuestionMap[item.QuestionId]; ok {

				qo := matchOptions.GetOption(item.OptionId)
				if qo == nil {
					continue
				}

				if _, ok := matchQuestions[q.ID]; !ok {
					matchQuestions[q.ID] = &model.MatchingQuestion{
						MatchingQuestion: q,
						MatchingOptions:  []*model.QuestionOption{qo},
					}
				} else {
					matchQuestions[q.ID].MatchingOptions = append(matchQuestions[q.ID].MatchingOptions, qo)
				}

			}

		}

		tmp.PersonalInfoQuestions = profileQuestions
		tmp.MatchCriteriaQuestions = matchQuestions

		rs[u.ID] = tmp
	}

	matchedResult := make([]*model.MatchUserData, 0, 100)
	for _, u := range users {

		mud, ok := rs[u.ID]
		if !ok {
			continue
		}

		if mud.MatchChatId != 0 {
			continue
		}

		for _, mu := range users {

			if mud.User.ID == mu.ID {
				continue
			}

			mud2, ok2 := rs[mu.ID]
			if !ok2 {
				continue
			}

			if mud2.MatchChatId != 0 {
				continue
			}

			if CheckMatch(mud, mud2) && CheckMatch(mud2, mud) {
				mud.MatchChatId = mud2.ChatId
				mud2.MatchChatId = mud.ChatId
				mud.MatchMatchUserData = mud2
				mud2.MatchMatchUserData = mud
				matchedResult = append(matchedResult, mud)
			}
		}
	}

	return matchedResult, nil
}

func CheckMatch(requestUser *model.MatchUserData, targetUser *model.MatchUserData) bool {

	// verify
	switch requestUser.VerifyOptionId {
	case vars.MatchingVerifiedOptionId:
		if !targetUser.User.IsVerify {
			return false
		}
	case vars.MatchingUnverifiedOptionId:
		if targetUser.User.IsVerify {
			return false
		}
	}

	// question match
	for _, mq := range requestUser.MatchCriteriaQuestions {

		if mq.MatchingQuestion.MatchingQuestionId == 0 || mq.MatchingQuestion.MatchingQuestionId == vars.MatchingQuestionId {
			continue
		}

		// find profile question
		profileQuestion := false
		pq, ok := targetUser.PersonalInfoQuestions[mq.MatchingQuestion.MatchingQuestionId]
		if ok {
			profileQuestion = true
		}

		matchQuestion := false
		matchq, ok := targetUser.MatchCriteriaQuestions[mq.MatchingQuestion.MatchingQuestionId]
		if ok {
			matchQuestion = true
		}

		if (!matchQuestion) && (!profileQuestion) {
			continue
		}

		qMatch := false
		for _, mo := range mq.MatchingOptions {

			if mo.IsMatchingAny {
				qMatch = true
				break
			}

			found := false
			if profileQuestion {
				for _, po := range pq.ProfileOptions {
					if mo.MatchingOptionId == po.ID {
						found = true
						break
					}
				}
			} else {
				for _, po := range matchq.MatchingOptions {
					if mo.MatchingOptionId == po.ID {
						found = true
						break
					}
				}
			}

			if found {
				qMatch = true
				break
			}

		}

		if !qMatch {
			return false
		}

	}

	return true
}

func ServiceSaveMatch(ctx context.Context, mud *model.MatchUserData) error {

	repo := repository.NewTRepository(repository.GetRepository().GetDB())

	if err := repo.Begin(); err != nil {
		return err
	}
	defer func() {
		if err := repo.Rollback(); err != nil {
		}
	}()

	if err := repo.UpdateMatchId(ctx, mud.User.ID, mud.MatchChatId); err != nil {
		return err
	}

	if err := repo.UpdateMatchId(ctx, mud.MatchMatchUserData.User.ID, mud.MatchMatchUserData.MatchChatId); err != nil {
		return err
	}

	// add
	po := &model.MatchedDetail{
		ColumnCreateModifyDeleteTime: gorm_global.ColumnCreateModifyDeleteTime{
			CreateTime: time.Now().Unix(),
			ModifyTime: time.Now().Unix(),
		},
		ChatId:      mud.ChatId,
		MatchChatId: mud.MatchChatId,
	}

	if err := repo.MatchedDetailAdd(ctx, po); err != nil {
		return err
	}

	// add
	po2 := &model.MatchedDetail{
		ColumnCreateModifyDeleteTime: gorm_global.ColumnCreateModifyDeleteTime{
			CreateTime: time.Now().Unix(),
			ModifyTime: time.Now().Unix(),
		},
		ChatId:      mud.MatchMatchUserData.ChatId,
		MatchChatId: mud.MatchMatchUserData.MatchChatId,
	}

	if err := repo.MatchedDetailAdd(ctx, po2); err != nil {
		return err
	}

	if err := repo.Commit(); err != nil {
		return err
	}

	return nil

}

func ServiceMatch(ctx context.Context, chatId int64, isVerify bool) (*model.User, error) {

	repo := repository.GetRepository()

	// find all question
	questions, err := repo.GetAllQuestion(ctx)
	if err != nil {
		return nil, err
	}

	// find all question option
	options, err := repo.GetAllOption(ctx)
	if err != nil {
		return nil, err
	}

	// find user all user question data
	userQuestionData, err := repo.GetUserQuestionDataByUser(ctx, chatId)
	if err != nil {
		return nil, err
	}

	// find matching options and user matching options value
	chatIds, _, err := repo.GetChatByMatching(ctx, chatId, questions, options, userQuestionData)
	if err != nil {
		return nil, err
	}

	if len(chatIds) == 0 {
		return nil, nil
	}

	log.Printf("chat_id: %d matching chat ids: %v \n", chatId, chatIds)

	// verified user matching
	userDataList := model.UserQuestionDataList(userQuestionData)
	vMOI := userDataList.GetFirstOptionIdByQuestionId(vars.MatchingQuestionId)
	switch vMOI {
	case vars.MatchingVerifiedOptionId:

		chatIds, err = repo.GetVerifyUser(ctx, chatIds, true)
		if err != nil {
			return nil, err
		}

	case vars.MatchingUnverifiedOptionId:

		chatIds, err = repo.GetVerifyUser(ctx, chatIds, false)
		if err != nil {
			return nil, err
		}

	}

	log.Printf("chat_id: %d verified matching chat ids: %v \n", chatId, chatIds)

	if isVerify {
		chatIds, err = repo.CheckHasOptionBy(ctx, chatIds, []int64{vars.MatchingVerifiedOptionId, vars.MatchingAnyOptionId})
		if err != nil {
			return nil, err
		}
	} else {
		chatIds, err = repo.CheckHasOptionBy(ctx, chatIds, []int64{vars.MatchingUnverifiedOptionId, vars.MatchingAnyOptionId})
		if err != nil {
			return nil, err
		}
	}

	log.Printf("chat_id: %d CheckHasOptionBy chat ids: %v \n", chatId, chatIds)

	if len(chatIds) == 0 {
		return nil, nil
	}

	// shuffle chat id
	matchChatId := shuffleChatId(chatIds)

	// find user
	user, err := repo.GetUserByChatId(ctx, matchChatId)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func shuffleChatId(a []int64) int64 {
	for i := range a {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
	return a[0]
}

func ServiceCheckEmailUnique(ctx context.Context, email string) error {

	repo := repository.GetRepository()

	cnt, err := repo.GetEmailCnt(ctx, email)
	if err != nil {
		return err
	}

	if cnt > 0 {
		return errors.New("Sorry, This email has already been registered.")
	}

	return nil
}
