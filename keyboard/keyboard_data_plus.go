package keyboard

import "encoding/json"

const (
	BUTTON_TYPE_MENU     int64 = 1
	BUTTON_TYPE_QUESTION int64 = 2
	BUTTON_TYPE_OPTION   int64 = 3
)

type KeyboardCallbackDataPlus struct {
	ButtonType   int64
	ButtonRelId  int64
	IsBackButton bool
}

func (k KeyboardCallbackDataPlus) CallbackData() *string {
	buf, _ := json.Marshal(k)
	str := string(buf)
	return &str
}
