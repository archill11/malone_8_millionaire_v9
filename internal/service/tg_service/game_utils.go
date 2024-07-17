package tg_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"myapp/internal/models"
	"myapp/pkg/files"
	"myapp/pkg/my_time_parser"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	animTimeout500  = 500
	animTimeout1000 = 1000
	animTimeout2000 = 2000
	animTimeout3000 = 3000
	animTimeout4000 = 4000
	animTimeoutTest = 3000
)

func (srv *TgService) Send3Kruga(fromId int) {
	user, _ := srv.Db.GetUserById(fromId)
	lichka := user.Lichka
	if lichka == "" {
		lichka = "https://t.me/markodinncov"
	}
	lichkaUrl := fmt.Sprintf("https://t.me/%s", srv.DelAt(lichka))

	scheme, _ := srv.Db.GetsSchemeByLichka(lichka)

	srv.SendVideoNoteCurrFile(fromId, fmt.Sprintf("./files/krug_3_%s_day_%d.mp4",scheme.Id, scheme.ScIdx))


	siteUrl := fmt.Sprintf("%s&data=%s", scheme.Link, srv.CreateBase64UserData(user.Id, user.Username, user.Firstname))
	mesgText := srv.GetActualSchema(fromId, siteUrl)
	_, err := srv.SendMessageHTML(fromId, mesgText)
	if err != nil {
		srv.l.Error(fmt.Errorf("Send3Kruga SendMessageWRM err: %v", err))
	}

	reply_markup := fmt.Sprintf(`{"inline_keyboard" : [
		[{ "text": "Написать Марку в ЛС", "url": "%s" }],
		[{ "text": "Обо мне", "callback_data": "obo_nme_btn" }],
		[{ "text": "Информация о заработке", "callback_data": "info_o_zarabotke_btn" }],
		[{ "text": "Часто задаваемые вопросы", "callback_data": "frequently_questions_btn" }],
		[{ "text": "Отзывы", "callback_data": "show_reviews_btn" }]
	]}`, lichkaUrl)

	futureJson := map[string]string{
		"video_note":   fmt.Sprintf("@%s", fmt.Sprintf("./files/krug_4_%s_day_%d.mp4",scheme.Id, scheme.ScIdx)),
		"chat_id": strconv.Itoa(fromId),
		"reply_markup": reply_markup,
	}
	cf, body, err := files.CreateForm(futureJson)
	if err != nil {
		err := fmt.Errorf("Send3Kruga CreateForm err: %v", err)
		srv.l.Error(err)
	}
	srv.SendVideoNote(body, cf)
	srv.SendMsgToServer(fromId, "bot", "send_3_kruga")
}

func (srv *TgService) SendMessageAndDb(chat_id int, text string) (models.SendMessageResp, error) {
	resp, err := srv.SendMessageHTML(chat_id, text)
	if err != nil {
		return models.SendMessageResp{}, fmt.Errorf("SendMessageAndDb SendMessage err: %v", err)
	}

	srv.SendMsgToServer(chat_id, "bot", text)
	return resp, nil
}

func (srv *TgService) SendMsgToServer(user_id int, msg_author, msg_txt string) error {
	user, err := srv.Db.GetUserById(user_id)
	if err != nil {
		err := fmt.Errorf("SendMsgToServer GetUserById err: %v", err)
		srv.l.Error(err)
		return err
	}

	botId := srv.Cfg.BotId
	step_id := user.Step
	step_txt := user.Step
	stepTexts := stepsMap[step_id]
	if len(stepTexts) >= 2 {
		step_txt = stepTexts[len(stepTexts)-1]
	}
	if strings.HasPrefix(msg_txt, "/start") {
		refArr := strings.Split(msg_txt, " ")
		ref := ""
		if len(refArr) > 1 {
			ref = refArr[1]
		}
		json_data, _ := json.Marshal(map[string]any{
			"user_id":    strconv.Itoa(user_id),
			"bot_id":     strconv.Itoa(botId),
			"username":    user.Username,
			"fullname":    user.Firstname,
			"step_id":    step_id,
			"step_text":   step_txt,
			"ref":         ref,
		})
		_, err = http.Post(
			fmt.Sprintf("%s/%s", srv.Cfg.ServerStatUrl, "add_user"),
			"application/json",
			bytes.NewBuffer(json_data),
		)
		if err != nil {
			err := fmt.Errorf("SendMsgToServer Post err: %v", err)
			srv.l.Error(err)
			return err
		}
		return nil
	}

	json_data, _ := json.Marshal(map[string]any{
		"user_id":    strconv.Itoa(user_id),
		"bot_id":     strconv.Itoa(botId),
		"username":    user.Username,
		"fullname":    user.Firstname,
		"new_step_id":    step_id,
		"new_step_text":   step_txt,
	})
	_, err = http.Post(
		fmt.Sprintf("%s/%s", srv.Cfg.ServerStatUrl, "update_user_step"),
		"application/json",
		bytes.NewBuffer(json_data),
	)
	if err != nil {
		err := fmt.Errorf("SendMsgToServer Post err: %v", err)
		srv.l.Error(err)
		return err
	}

	json_data, _ = json.Marshal(map[string]any{
		"user_id":    strconv.Itoa(user_id),
		"bot_id":     strconv.Itoa(botId),
		"message":    msg_txt,
		"sender_type": msg_author, // user | bot | admin
	})
	_, err = http.Post(
		fmt.Sprintf("%s/%s", srv.Cfg.ServerStatUrl, "add_message"),
		"application/json",
		bytes.NewBuffer(json_data),
	)
	if err != nil {
		err := fmt.Errorf("SendMsgToServer Post err: %v", err)
		srv.l.Error(err)
		return err
	}
	return nil
}

func (srv *TgService) SendMessageAndDbAdmin(chat_id int, text string) (models.SendMessageResp, error) {
	resp, err := srv.SendMessage(chat_id, text)
	if err != nil {
		return models.SendMessageResp{}, fmt.Errorf("SendMessageAndDb SendMessage err: %v", err)
	}

	// srv.SendMsgToServer(chat_id, "admin", text)
	return resp, nil
}

func (srv *TgService) SendAnimMessage(steps_id string, chat_id, milisecond int) error {
	steps, ok := srv.Steps[steps_id]
	if !ok {
		return nil
	}

	var messId int
	for i, v := range steps {
		time.Sleep(time.Millisecond * time.Duration(milisecond))
		if i == 0 {
			SendMessageResp, err := srv.SendMessage(chat_id, v)
			if err != nil {
				return fmt.Errorf("SendAnimMessage SendMessage err: %v", err)
			}
			messId = SendMessageResp.Result.MessageId
			continue
		}
		if v == "delete" {
			srv.DeleteMessage(chat_id, messId)
			continue
		}
		srv.EditMessageText(chat_id, messId, v)
	}
	fullMess := steps[len(steps)-1]
	if fullMess == "delete" {
		fullMess = steps[0]
	}

	srv.SendMsgToServer(chat_id, "bot", fullMess)
	return nil
}

func (srv *TgService) SendAnimMessageHTML(steps_id string, chat_id, milisecond int) error {
	steps, ok := srv.Steps[steps_id]
	if !ok {
		return nil
	}

	var messId int
	for i, v := range steps {
		time.Sleep(time.Millisecond * time.Duration(milisecond))
		if i == 0 {
			SendMessageResp, err := srv.SendMessageHTML(chat_id, v)
			if err != nil {
				return fmt.Errorf("SendAnimMessage SendMessageHTML err: %v", err)
			}
			messId = SendMessageResp.Result.MessageId
			continue
		}
		if v == "delete" {
			srv.DeleteMessage(chat_id, messId)
			continue
		}
		srv.EditMessageTextHTML(chat_id, messId, v)
	}
	fullMess := steps[len(steps)-1]
	if fullMess == "delete" {
		fullMess = steps[0]
	}

	srv.SendMsgToServer(chat_id, "bot", fullMess)
	return nil
}

func (srv *TgService) SendBalance(chat_id int, balance string, milisecond int) error {
	time.Sleep(time.Millisecond * time.Duration(milisecond))
	user, err := srv.Db.GetUserById(chat_id)
	if err != nil {
		return fmt.Errorf("SendBalance GetUserById err: %v", err)
	}

	red := strings.Repeat("❤️", user.Lives)
	black := strings.Repeat("🖤", 3-user.Lives)

	logMess := fmt.Sprintf("🏦 Сейчас в банке: %s₽\n🫀Жизни: %s", balance, red+black)
	srv.SendMessageAndDb(chat_id, logMess)

	return nil
}

func (srv *TgService) SendPrePush(chatId int, prePushId int) error {
	PrePushTextMap := map[int]string{
		1: "Ты долго бездействуешь 😰\n❗️ Через 15 минут у тебя сгорит первая жизнь",
		2: "Ты все еще бездействуешь 😱\n❗️ Через 15 минут у тебя сгорит вторая жизнь",
		3: "❗️Осталось 15 минут и у тебя сгорит последняя жизнь\n\nТак у тебя сгорят все 3 жизни и игра закончится ❌",
	}
	text := PrePushTextMap[prePushId]
	srv.SendMessage(chatId, text)

	srv.SendMsgToServer(chatId, "bot", text)
	// srv.Db.EditStep(chatId, "Ты долго бездействуешь")
	srv.Db.EditIsSendPush(chatId, 1)
	srv.Db.UpdateLatsActiontime(chatId)
	srv.Db.UpdateFeedbackTime(chatId)
	return nil
}

func (srv *TgService) SendPush(chatId int, pushId int) error {
	pushTextMap := map[int]string{
		1: "❤️❤️🖤\n<b>У тебя сгорела одна жизнь, потому что 15 минут не было никаких действий 😔</b>\n\nПомни - ты не кошка, у тебя не 9 жизней 😸\nЕсли сгорят все жизни, то доступ ко всем статьям закроется, деньги на балансе пропадут и ты так и будешь до конца жизни ждать лучшего момента 🤷🏻‍♂️\n\nИли всё-таки изменишь свою жизнь к лучшему?\n\nУ тебя еще есть время, чтобы продолжить и дойти до своей награды в 100.000₽ 💸",
		2: "❤️🖤🖤\n<b>У тебя сгорела вторая жизнь, и это очень печально 😒</b>\n\nЕсли продолжишь в таком же темпе, то ни к чему хорошему это не приведёт.\n\nНеужели ты сейчас занимаешься чем-то более увлекательным, что готов упустить такой шанс изменить свою жизнь к лучшему? 🧐\n\nОсталась одна жизнь и один шанс получить награду\nЖми кнопку выше ☝🏻",
	}
	text := pushTextMap[pushId]
	fileNameInServer := fmt.Sprintf("./files/push_%d.jpg", pushId)
	_, err := srv.SendPhotoWCaption(chatId, text, fileNameInServer)
	if err != nil {
		return fmt.Errorf("Push3 SendPhotoWCaption err: %v", err)
	}
	
	srv.Db.EditLives(chatId, (3 - pushId))
	srv.SendMsgToServer(chatId, "bot", text)
	// srv.Db.EditStep(chatId, text)
	srv.Db.EditIsSendPush(chatId, 0)
	srv.Db.UpdateLatsActiontime(chatId)
	srv.Db.UpdateFeedbackTime(chatId)
	return nil
}

func (srv *TgService) Push3(chatId int) error {
	text := "🖤🖤🖤\n<b>Все жизни сгорели 🥶\nНо у тебя еще есть шанс восстановить их.</b>\n\nГарантирую, что такого ценного предложения, как даю тебе я, ты еще не получал 😉"
	fileNameInServer := "./files/push_3.jpg"
	_, err := srv.SendPhotoWCaption(chatId, text, fileNameInServer)
	if err != nil {
		return fmt.Errorf("Push3 SendPhotoWCaption err: %v", err)
	}

	srv.Db.EditLives(chatId, 0)
	srv.SendMsgToServer(chatId, "bot", text)
	srv.Db.EditIsSendPush(chatId, 0)
	srv.Db.UpdateLatsActiontime(chatId)
	srv.Db.UpdateFeedbackTime(chatId)
	return nil
}

func (srv *TgService) LastPush(chatId int) error {
	text := "<b>К сожалению, ты не успел начать игру заново в отведенное время, впредь бот закрыт для тебя навсегда! 😵</b>\n\nЕсли у тебя остались какие-то вопросы - ты можешь обратиться в тех.поддержку через команду /help ✍️"
	fileNameInServer := "./files/push_4.jpg"
	_, err := srv.SendPhotoWCaption(chatId, text, fileNameInServer)
	if err != nil {
		return fmt.Errorf("Push3 SendPhotoWCaption err: %v", err)
	}

	srv.Db.EditLives(chatId, 0)
	srv.SendMsgToServer(chatId, "bot", text)
	time.Sleep(time.Second * 3)

	huersStr, _ := srv.GetUserLeftTime(chatId)
	text = fmt.Sprintf("❗️У тебя есть %s на то, чтобы начать игру заново♻️\n\nЕсли ты не успеешь запустить игру за это время, то доступ к боту будет закрыт навсегда. Перезапуск бота не поможет, он просто перестанет работать для тебя ⛔️", huersStr)
	replyMarkup := `{"inline_keyboard" : [
		[{ "text": "ЗАБРАТЬ 100.000₽", "callback_data": "restart_game" }]
	]}`
	srv.SendMessageWRM(chatId, text, replyMarkup)

	srv.SendMsgToServer(chatId, "bot", text)
	srv.Db.UpdateLatsActiontime(chatId)
	srv.Db.UpdateFeedbackTime(chatId)
	return nil
}

func (srv *TgService) SendFeedback(chatId, feedbackNum int) error {
	feedbacksMap := map[int]string{
		1: "+1 счастливчик уже получил главный приз, но к сожалению это не ты😞\n\nНу же, перезапусти игру, я уверен у тебя получится пройти её, я уже подготовил приз и для тебя, но <b>у тебя осталось всего %s</b> 🫣",
		2: "Игру мешает пройти лень человека, а не её сложность. Я уверен тебе по силам пройти её, просто соберись🤜🏼🤛🏼\n\nЯ напомню, что <b>у тебя осталось %s</b> на то, чтобы перезапустить её и попробовать запуск! Дальше бот для тебя отключится😵",
		3: "Люблю я исполнять желания людей, которые этого заслуживают🥰\n\nСтоит просто пройти игру и я исполню и твоё желание, но <b>у тебя на это осталось %s</b> 😯",
		4: "94%% людей не проходят бота только потому, что откладывают до последнего на потом, а потом становится поздно🥶\n\n<b>У тебя осталось %s</b> чтобы не пополнить толпу зевак🥱",
		5: "Голосовые с благодарностью - это мёд для моих ушей🫠\n\nПочему я до сих пор не получил восторженного голосового сообщения от тебя?🧐 Ах да, игра еще тобой не пройдена, но <b>у тебя еще осталось %s</b>. Поторопись!",
		6: "Человек прошел игру в самый последний момент и спас дочь от отчисления из универа🥳\n\nКстати ты тоже еще можешь запрыгнуть в уходящий вагон поезда, <b>у тебя осталось %s</b> до того, как бот отключится! Если начнешь прямо сейчас, то успеешь пройти🥵",
	}

	huersStr, _ := srv.GetUserLeftTime(chatId)
	text := fmt.Sprintf(feedbacksMap[feedbackNum], huersStr)
	fileNameInServer := fmt.Sprintf("./files/feedback_%d.jpg", feedbackNum)
	replyMarkup := `{"inline_keyboard" : [
		[{ "text": "ЗАБРАТЬ 100.000₽", "callback_data": "restart_game" }]
	]}`
	_, err := srv.SendPhotoWCaptionWRM(chatId, text, fileNameInServer, replyMarkup)
	if err != nil {
		return fmt.Errorf("SendFeedback SendPhotoWCaptionWRM err: %v", err)
	}

	srv.Db.EditdFeedbackCnt(chatId, feedbackNum)
	srv.Db.UpdateFeedbackTime(chatId)
	srv.SendMsgToServer(chatId, "bot", text)
	// srv.Db.EditStep(chatId, text)
	return nil
}


func AbsTimeStrToRusStr(time string) string {
	bef, _, _ := strings.Cut(time, "m")
	bef = strings.Replace(bef, "h", "ч ", -1)
	bef = fmt.Sprintf("%s мин", bef)
	bef = strings.Replace(bef, "часа мин", "часа ", -1)
	return bef
}

func (srv *TgService) GetUserLeftTime(chatId int) (string, error) {
	user, err := srv.Db.GetUserById(chatId)
	if err != nil {
		return "", fmt.Errorf("GetUserLeftTime GetUserById err: %v", err)
	}
	createdAt, _ := my_time_parser.ParseInLocation(user.CreatedAt, my_time_parser.Msk)
	maxTimeToUseBot := createdAt.Add(time.Hour * 72)
	tn := time.Now().In(my_time_parser.Msk)
	if tn.After(maxTimeToUseBot) {
		return "", nil
	}
	leffft := maxTimeToUseBot.Sub(tn) // сколько осталось юзать бота
	if leffft < 0 {
		return "", nil
	}
	huersStr := AbsTimeStrToRusStr(leffft.Round(time.Second).String())
	return huersStr, nil
}

func (srv *TgService) IsIgnoreUser(chatId int) bool {
	user, err := srv.Db.GetUserById(chatId)
	if err != nil {
		return true
	}
	if user.Id == 0 {
		return false
	}
	createdAt, _ := my_time_parser.ParseInLocation(user.CreatedAt, my_time_parser.Msk)
	maxTimeToUseBot := createdAt.Add(time.Hour*24)
	tn := time.Now().In(my_time_parser.Msk)
	if tn.After(maxTimeToUseBot) {
		return true
	}
	return false
}

func (srv *TgService) GetUserPersonalRef(fromId int) string {
	getMeResp, _ := srv.GetMe(srv.Cfg.Token)
	return fmt.Sprintf("https://t.me/%s?start=%d", getMeResp.Result.UserName, fromId)
}