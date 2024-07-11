package tg_service

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func (srv *TgService) ShowAdminPanel(chatId int) error {
	json_data, err := json.Marshal(map[string]any{
		"chat_id": strconv.Itoa(chatId),
		"text":    "Админ панель",
		"reply_markup": `{"inline_keyboard" : [
			[{ "text": "Инфа по Юзеру", "callback_data": "user_info_btn" }]
		]}`,
	})
	if err != nil {
		return fmt.Errorf("ShowAdminPanel Marshal err: %v", err)
	}
	err = srv.SendData(json_data, "sendMessage")
	if err != nil {
		return fmt.Errorf("ShowAdminPanel sendData err: %v", err)
	}

	return nil
}
