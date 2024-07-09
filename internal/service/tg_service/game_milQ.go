package tg_service

import (
	"fmt"
	"time"
)

func (srv *TgService) ShowMilQ(chatId, qNum int) error {
	time.Sleep(time.Millisecond * time.Duration(animTimeoutTest))

	textMap := map[int]string{
		1: "И так, первый вопрос 👆\n\nВыбери правильный ответ 👇",
		2: "Второй вопрос 👆\n\nВыбери правильный ответ 👇",
		3: "Третий вопрос 👆\n\nВыбери правильный ответ 👇",
	}
	fileNameMap := map[int]string{
		1:  "./files/mil_q1.jpg",
		2:  "./files/mil_q2.jpg",
		3:  "./files/mil_q3.jpg",
	}
	replyMarkupMap := map[int]string{
		1: `{"inline_keyboard" : [
			[ { "text": "A", "callback_data": "_lose_q_1_" }, { "text": "B", "callback_data": "_lose_q_1_" }, { "text": "C", "callback_data": "_win_q_1_" }, { "text": "D", "callback_data": "_lose_q_1_" }]
		]}`,
		2: `{"inline_keyboard" : [
			[ { "text": "A", "callback_data": "_win_q_2_" }, { "text": "B", "callback_data": "_lose_q_2_" }, { "text": "C", "callback_data": "_lose_q_2_" }, { "text": "D", "callback_data": "_lose_q_2_" }]
		]}`,
		3: `{"inline_keyboard" : [
			[ { "text": "A", "callback_data": "_lose_q_3_" }, { "text": "B", "callback_data": "_lose_q_3_" }, { "text": "C", "callback_data": "_win_q_3_" }, { "text": "D", "callback_data": "_lose_q_3_" }]
		]}`,
	}

	text := textMap[qNum]
	replyMarkup := replyMarkupMap[qNum]
	fileNameInServer := fileNameMap[qNum]
	_, err := srv.SendPhotoWCaptionWRM(chatId, text, fileNameInServer, replyMarkup)
	if err != nil {
		return fmt.Errorf("ShowMilQ SendPhotoWCaptionWRM err: %v", err)
	}
	srv.Db.EditStep(chatId, text)
	srv.SendMsgToServer(chatId, "bot", text)
	return nil
}

func (srv *TgService) ShowQLose(chatId int, q_num string) error {
	time.Sleep(time.Millisecond * time.Duration(animTimeoutTest))

	text := "Ответ неверный ❌\nК сожалению, ты ошибся, но шанс еще есть!\n\nЖми на кнопку 👇"
	reply_markup := fmt.Sprintf(`{"inline_keyboard" : [
		[{ "text": "Попробовать еще раз", "callback_data": "show_q_%s_" }]
	]}`, q_num)
	srv.SendMessageWRM(chatId, text, reply_markup)

	// srv.SendMsgToServer(chatId, "bot", text)
	return nil
}

func (srv *TgService) ShowQWin(chatId int, q_num string) error {
	time.Sleep(time.Millisecond * time.Duration(animTimeoutTest))
	
	textMap := map[string]string{
		"1":  "Правильный ответ! ✅\n\nА ты неплох, осталось 2 вопроса и ты сможешь забрать заветные 5.000 рублей 👏\n\nДавай дальше 👇",
		"2":  "Правильный ответ! ✅\n\nКрасавчик! Остался последний вопрос, отделяющий тебя от приза 💰\n\nДавай дальше 👇",
		"3": "ТЫ ПОБЕДИЛ ✅\n\nПОЗДРАВЛЯЮ, ТЫ ЭТО СДЕЛАЛ, ТВОЯ ИНТУИЦИЯ И МОЗГ ТЕБЯ НЕ ПОДВЕЛИ! 🎉🎉🎉\n\nПрими мои поздравления, ты выйграл 5.000₽! 😇🎉🎁",
	}

	if q_num == "1" {
		srv.SendMessageAndDb(chatId, textMap[q_num])
		// time.Sleep(time.Millisecond * 2000)
		srv.ShowMilQ(chatId, 2)
		return nil
	}
	if q_num == "2" {
		srv.SendMessageAndDb(chatId, textMap[q_num])
		// time.Sleep(time.Millisecond * 2000)
		srv.ShowMilQ(chatId, 3)
		return nil
	}
	if q_num == "3" {
		srv.SendMessageAndDb(chatId, textMap[q_num])
		// time.Sleep(time.Millisecond * 2000)
		// srv.Db.EditStep(chatId, "6")
		// srv.SendAnimMessage("6", chatId, animTimeoutTest)
		// time.Sleep(time.Second)

		// user, _ := srv.Db.GetUserById(chatId)
		chLink := "https://t.me/+rLIklQb0ALNhZjEx"
		chLink2 := "https://t.me/geniusgiveaway"

		messText := fmt.Sprintf("Чтобы забрать свой приз , тебе необходимо выполнить 2 простых условия 😎\n\nПервое условие:\nТебе нужно подписаться на этот канал👇\n\n %s\n %s", chLink, chLink2)
		reply_markup := `{"inline_keyboard" : [
			[{ "text": "Подписался☑️", "callback_data": "subscribe" }]
		]}`
		srv.SendMessageWRM(chatId, messText, reply_markup)

		srv.SendMsgToServer(chatId, "bot", "Проверка подписки")
		return nil
	}
	if q_num == "4" {
		srv.SendMessageAndDb(chatId, textMap[q_num])
		time.Sleep(time.Second * 2)
		srv.ShowMilQ(chatId, 5)
		return nil
	}
	if q_num == "5" {
		srv.SendMessageAndDb(chatId, textMap[q_num])
		time.Sleep(time.Second * 2)
		srv.ShowMilQ(chatId, 6)
		return nil
	}
	if q_num == "6" {
		srv.SendMessageAndDb(chatId, textMap[q_num])
		time.Sleep(time.Second * 2)

		srv.Db.EditStep(chatId, "8")
		srv.SendAnimMessageHTML("8", chatId, 2000)
		srv.Db.EditBotState(chatId, "read_article_after_KNB_win")
		return nil
	}
	if q_num == "7" {
		srv.SendMessageAndDb(chatId, textMap[q_num])
		time.Sleep(time.Second * 2)
		srv.ShowMilQ(chatId, 8)
		return nil
	}
	if q_num == "8" {
		srv.SendMessageAndDb(chatId, textMap[q_num])
		time.Sleep(time.Second * 2)
		srv.ShowMilQ(chatId, 9)
		return nil
	}
	if q_num == "9" {
		srv.SendMessageAndDb(chatId, textMap[q_num])
		time.Sleep(time.Second * 2)

		srv.Db.EditStep(chatId, "10")
		srv.SendAnimMessageHTML("10", chatId, 2000)
		srv.Db.EditBotState(chatId, "read_article_after_OIR_win")
		return nil
	}
	if q_num == "10" {
		srv.SendMessageAndDb(chatId, textMap[q_num])
		time.Sleep(time.Second * 2)
		srv.ShowMilQ(chatId, 11)
		return nil
	}
	if q_num == "11" {
		srv.SendMessageAndDb(chatId, textMap[q_num])
		time.Sleep(time.Second * 2)
		srv.ShowMilQ(chatId, 12)
		return nil
	}
	if q_num == "12" {
		srv.SendMessageAndDb(chatId, textMap[q_num])
		time.Sleep(time.Second * 2)

		srv.Db.EditStep(chatId, "12")
		srv.SendAnimMessageHTML("12", chatId, 2000)
		srv.Db.EditBotState(chatId, "read_article_after_TrurOrFalse_win")
		return nil
	}
	return nil
}
