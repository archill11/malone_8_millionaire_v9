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
			[{ "text": "Рассылка по шагу [copy]", "callback_data": "mailing_copy_btn" }],
			[{ "text": "Удалить юзера по username", "callback_data": "delete_user_by_username_btn" }],
			[{ "text": "Удалить юзера по id", "callback_data": "delete_user_by_id_btn" }]
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
