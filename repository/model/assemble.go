package model

type MatchUserData struct {
	ChatId                 int64
	User                   *User
	PersonalInfoQuestions  map[int64]*ProfileQuestion  `json:"-"`
	MatchCriteriaQuestions map[int64]*MatchingQuestion `json:"-"`
	MatchChatId            int64
	MatchMatchUserData     *MatchUserData `json:"-"`
	VerifyOptionId         int64
}

type ProfileQuestion struct {
	ProfileQuestion *Question
	ProfileOptions  []*QuestionOption
}

type MatchingQuestion struct {
	MatchingQuestion *Question
	MatchingOptions  []*QuestionOption
}
