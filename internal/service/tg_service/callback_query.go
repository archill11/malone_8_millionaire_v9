package tg_service

import (
	"fmt"
	"myapp/internal/models"
	"myapp/pkg/files"
	my_regex "myapp/pkg/regex"
	"strconv"
	"strings"
	"time"
)

func (srv *TgService) HandleCallbackQuery(m models.Update) error {
	cq := m.CallbackQuery
	fromId := cq.From.Id
	fromUsername := cq.From.UserName
	srv.l.Info(fmt.Sprintf("HandleCallbackQuery: fromId: %d, fromUsername: %s, cq.Data: %s", fromId, fromUsername, cq.Data))

	srv.Db.EditStep(fromId, fmt.Sprintf("кнопка: %s", cq.Data))
	srv.SendMsgToServer(fromId, "user", fmt.Sprintf("кнопка: %s", cq.Data))

	go func() {
		if cq.Data != "subscribe" && cq.Data != "zabrat_nagradu" && !strings.HasPrefix(cq.Data, "_win_q") && !strings.HasPrefix(cq.Data, "_lose_q") && !strings.HasPrefix(cq.Data, "show_q") {
			srv.l.Warn("syka")
			time.Sleep(time.Second)
			srv.EditMessageReplyMarkup(fromId, cq.Message.MessageId)
			for i:=cq.Message.MessageId; i >= cq.Message.MessageId-20; i-- {
				user, _ := srv.Db.GetUserById(fromId)
				if i == user.NotDelMessId {
					break
				}
				srv.DeleteMessage(fromId, i)
				time.Sleep(time.Millisecond*300)
			}
			// srv.Db.UpdateLatsActiontime(fromId)
		}
	}()

	// user, err := srv.Db.GetUserById(fromId)
	// if err != nil {
	// 	return fmt.Errorf("HandleCallbackQuery GetUserById err: %v", err)
	// }
	// if user.Id != 0 && user.Lives == 0 {
	// 	return nil
	// }

	if cq.Data == "delete_user_by_username_btn" {
		err := srv.CQ_delete_user_by_username_btn(m)
		if err != nil {
			srv.SendMessage(fromId, ERR_MSG)
			srv.SendMessage(fromId, err.Error())
		}
		return err
	}

	if cq.Data == "delete_user_by_id_btn" {
		err := srv.CQ_delete_user_by_id_btn(m)
		if err != nil {
			srv.SendMessage(fromId, ERR_MSG)
			srv.SendMessage(fromId, err.Error())
		}
		return err
	}

	if cq.Data == "start_game" {
		err := srv.CQ_start_game(m)
		if err != nil {
			srv.SendMessageAndDb(fromId, ERR_MSG)
			srv.SendMessageAndDb(fromId, err.Error())
		}
		srv.Db.UpdateLatsActiontime(fromId)
		return err
	}

	if cq.Data == "restart_game" {
		err := srv.CQ_restart_game(m)
		if err != nil {
			srv.SendMessageAndDb(fromId, ERR_MSG)
			srv.SendMessageAndDb(fromId, err.Error())
		}
		srv.Db.UpdateLatsActiontime(fromId)
		return err
	}

	if cq.Data == "subscribe" {
		err := srv.CQ_subscribe(m)
		if err != nil {
			srv.SendMessageAndDb(fromId, ERR_MSG)
			srv.SendMessageAndDb(fromId, err.Error())
		}
		srv.Db.UpdateLatsActiontime(fromId)
		return err
	}

	if cq.Data == "show_reviews_btn" {
		err := srv.CQ_show_reviews_btn(m)
		if err != nil {
			srv.SendMessageAndDb(fromId, ERR_MSG)
			srv.SendMessageAndDb(fromId, err.Error())
		}
		srv.Db.UpdateLatsActiontime(fromId)
		return err
	}

	if strings.HasPrefix(cq.Data, "show_q_") { // показать mil вопрос
		if strings.HasPrefix(cq.Message.Text, "Ответ неверный") {
			time.Sleep(time.Second)
			srv.DeleteMessage(fromId, cq.Message.MessageId)
			srv.DeleteMessage(fromId, cq.Message.MessageId-1)
		}

		qId := my_regex.GetStringInBetween(cq.Data, "show_q_", "_")
		qIdInt, err := strconv.Atoi(qId)
		if err != nil {
			srv.SendMessageAndDb(fromId, ERR_MSG)
			srv.SendMessageAndDb(fromId, err.Error())
		}
		err = srv.ShowMilQ(fromId, qIdInt)
		
		if err != nil {
			srv.SendMessageAndDb(fromId, ERR_MSG)
			srv.SendMessageAndDb(fromId, err.Error())
		}
		srv.Db.UpdateLatsActiontime(fromId)
		return err
	}

	if strings.HasPrefix(cq.Data, "_lose_q_") { // показать "Попробовать еще раз" на вопрос
		qId := my_regex.GetStringInBetween(cq.Data, "_lose_q_", "_")
		err := srv.ShowQLose(fromId, qId)
		if err != nil {
			srv.SendMessageAndDb(fromId, ERR_MSG)
			srv.SendMessageAndDb(fromId, err.Error())
		}
		srv.Db.UpdateLatsActiontime(fromId)
		return err
	}

	if strings.HasPrefix(cq.Data, "_win_q_") {
		qId := my_regex.GetStringInBetween(cq.Data, "_win_q_", "_")
		err := srv.ShowQWin(fromId, qId)
		if err != nil {
			srv.SendMessageAndDb(fromId, ERR_MSG)
			srv.SendMessageAndDb(fromId, err.Error())
		}
		srv.Db.UpdateLatsActiontime(fromId)
		return err
	}

	if cq.Data == "mailing_copy_btn" {
		err := srv.CQ_mailing_copy_btn(m)
		if err != nil {
			srv.SendMessageAndDb(fromId, ERR_MSG)
			srv.SendMessageAndDb(fromId, err.Error())
		}
		srv.Db.UpdateLatsActiontime(fromId)
		return err
	}

	if cq.Data == "zabrat_instr" {
		err := srv.CQ_zabrat_instr(m)
		if err != nil {
			srv.SendMessageAndDb(fromId, ERR_MSG)
			srv.SendMessageAndDb(fromId, err.Error())
		}
		srv.Db.UpdateLatsActiontime(fromId)
		return err
	}

	if cq.Data == "zabrat_instr_500" {
		err := srv.CQ_zabrat_instr_500(m)
		if err != nil {
			srv.SendMessageAndDb(fromId, ERR_MSG)
			srv.SendMessageAndDb(fromId, err.Error())
		}
		srv.Db.UpdateLatsActiontime(fromId)
		return err
	}

	if cq.Data == "zabrat_nagradu" {
		err := srv.CQ_zabrat_nagradu(m)
		if err != nil {
			srv.SendMessageAndDb(fromId, ERR_MSG)
			srv.SendMessageAndDb(fromId, err.Error())
		}
		srv.Db.UpdateLatsActiontime(fromId)
		return err
	}

	if cq.Data == "get_scheme" {
		err := srv.CQ_get_scheme(m)
		if err != nil {
			srv.SendMessageAndDb(fromId, ERR_MSG)
			srv.SendMessageAndDb(fromId, err.Error())
		}
		srv.Db.UpdateLatsActiontime(fromId)
		return err
	}

	if cq.Data == "obo_nme_btn" {
		err := srv.CQ_obo_nme_btn(m)
		if err != nil {
			srv.SendMessageAndDb(fromId, ERR_MSG)
			srv.SendMessageAndDb(fromId, err.Error())
		}
		srv.Db.UpdateLatsActiontime(fromId)
		return err
	}

	if cq.Data == "info_o_zarabotke_btn" {
		err := srv.CQ_info_o_zarabotke_btn(m)
		if err != nil {
			srv.SendMessageAndDb(fromId, ERR_MSG)
			srv.SendMessageAndDb(fromId, err.Error())
		}
		srv.Db.UpdateLatsActiontime(fromId)
		return err
	}

	if cq.Data == "frequently_questions_btn" {
		err := srv.CQ_frequently_questions_btn(m)
		if err != nil {
			srv.SendMessageAndDb(fromId, ERR_MSG)
			srv.SendMessageAndDb(fromId, err.Error())
		}
		srv.Db.UpdateLatsActiontime(fromId)
		return err
	}

	if cq.Data == "davay_sigraem_btn" {
		err := srv.CQ_davay_sigraem_btn(m)
		if err != nil {
			srv.SendMessageAndDb(fromId, ERR_MSG)
			srv.SendMessageAndDb(fromId, err.Error())
		}
		srv.Db.UpdateLatsActiontime(fromId)
		return err
	}

	if cq.Data == "pognaly_btn" {
		err := srv.CQ_pognaly_btn(m)
		if err != nil {
			srv.SendMessageAndDb(fromId, ERR_MSG)
			srv.SendMessageAndDb(fromId, err.Error())
		}
		srv.Db.UpdateLatsActiontime(fromId)
		return err
	}

	srv.Db.UpdateLatsActiontime(fromId)
	return nil
}

func (srv *TgService) CQ_start_game(m models.Update) error {
	cq := m.CallbackQuery
	fromId := cq.From.Id
	fromUsername := cq.From.UserName
	srv.l.Info(fmt.Sprintf("CQ_start_game: fromId: %d, fromUsername: %s", fromId, fromUsername))

	srv.SendAnimMessage("-1", fromId, animTimeout500)
	srv.SendBalance(fromId, "1000", animTimeout500)
	srv.SendAnimMessageHTML("2", fromId, animTimeoutTest)
	srv.SendAnimMessage("4", fromId, animTimeoutTest)
	srv.Db.EditStep(fromId, "5")
	srv.SendAnimMessage("5", fromId, animTimeoutTest)

	err := srv.ShowMilQ(fromId, 1)
	if err != nil {
		return fmt.Errorf("CQ_start_game ShowMilQ1 err: %v", err)
	}

	return nil
}

func (srv *TgService) CQ_mailing_copy_btn(m models.Update) error {
	cq := m.CallbackQuery
	fromId := cq.From.Id
	fromUsername := cq.From.UserName
	srv.l.Info(fmt.Sprintf("CQ_start_game: fromId: %d, fromUsername: %s", fromId, fromUsername))

	srv.SendForceReply(fromId, MAILING_COPY_STEP)

	return nil
}

func (srv *TgService) CQ_zabrat_instr(m models.Update) error {
	cq := m.CallbackQuery
	fromId := cq.From.Id
	fromUsername := cq.From.UserName
	fromFirstName := cq.From.FirstName
	srv.l.Info(fmt.Sprintf("CQ_zabrat_instr: fromId: %d, fromUsername: %s", fromId, fromUsername))

	user, _ := srv.Db.GetUserById(fromId)
	lichka := user.Lichka
	if lichka == "" {
		lichka = "https://t.me/markodinncov"
	}
	// lichkaUrl := fmt.Sprintf("https://t.me/%s", srv.DelAt(lichka))
	scheme, _ := srv.Db.GetsSchemeByLichka(lichka)

	base64Str := srv.CreateBase64UserData(fromId, fromUsername, fromFirstName)
	siteUrl := fmt.Sprintf("%s&data=%s", scheme.Link, base64Str)

	mesgText := srv.GetActualSchema(fromId, siteUrl)

	_, err := srv.SendMessageHTML(fromId, mesgText)
	if err != nil {
		srv.l.Error(fmt.Errorf("CQ_zabrat_instr SendMessageWRM err: %v", err))
	}
	srv.SendMsgToServer(fromId, "bot", mesgText)

	return nil
}

func (srv *TgService) CQ_zabrat_nagradu(m models.Update) error {
	cq := m.CallbackQuery
	fromId := cq.From.Id
	fromUsername := cq.From.UserName
	// fromFirstName := cq.From.FirstName
	srv.l.Info(fmt.Sprintf("CQ_zabrat_nagradu: fromId: %d, fromUsername: %s", fromId, fromUsername))

	user, _ := srv.Db.GetUserById(fromId)
	lichka := user.Lichka
	if lichka == "" {
		lichka = "https://t.me/markodinncov"
	}

	srv.Send3Kruga(fromId)

	return nil
}

func (srv *TgService) CQ_get_scheme(m models.Update) error {
	cq := m.CallbackQuery
	fromId := cq.From.Id
	fromUsername := cq.From.UserName
	fromFirstName := cq.From.FirstName
	srv.l.Info(fmt.Sprintf("CQ_zabrat_nagradu: fromId: %d, fromUsername: %s", fromId, fromUsername))

	user, _ := srv.Db.GetUserById(fromId)
	lichka := user.Lichka
	if lichka == "" {
		lichka = "https://t.me/markodinncov"
	}
	scheme, _ := srv.Db.GetsSchemeByLichka(lichka)

	siteUrl := fmt.Sprintf("%s&data=%s", scheme.Link, srv.CreateBase64UserData(fromId, fromUsername, fromFirstName))
	mesgText := srv.GetActualSchema(fromId, siteUrl)
	_, err := srv.SendMessageHTML(fromId, mesgText)
	if err != nil {
		srv.l.Error(fmt.Errorf("CQ_zabrat_nagradu SendMessageWRM err: %v", err))
	}
	srv.SendMsgToServer(fromId, "bot", mesgText)

	return nil
}

func (srv *TgService) CQ_zabrat_instr_500(m models.Update) error {
	cq := m.CallbackQuery
	fromId := cq.From.Id
	fromUsername := cq.From.UserName
	fromFirstName := cq.From.FirstName
	srv.l.Info(fmt.Sprintf("CQ_zabrat_instr_500: fromId: %d, fromUsername: %s", fromId, fromUsername))

	user, _ := srv.Db.GetUserById(fromId)
	lichka := user.Lichka
	if lichka == "" {
		lichka = "https://t.me/markodinncov"
	}
	scheme, _ := srv.Db.GetsSchemeByLichka(lichka)

	base64Str := srv.CreateBase64UserData(fromId, fromUsername, fromFirstName)
	siteUrl := fmt.Sprintf("%s&data=%s", scheme.Link, base64Str)

	mesgText := srv.GetActualSchema(fromId, siteUrl)

	_, err := srv.SendMessageHTML(fromId, mesgText)
	if err != nil {
		srv.l.Error(fmt.Errorf("CQ_zabrat_instr SendMessageHTML err: %v", err))
	}
	srv.SendMsgToServer(fromId, "bot", mesgText)

	return nil
}

func (srv *TgService) CQ_restart_game(m models.Update) error {
	cq := m.CallbackQuery
	fromId := cq.From.Id
	fromUsername := cq.From.UserName
	fromFirstName := cq.From.FirstName
	srv.l.Info(fmt.Sprintf("CQ_restart_game: fromId: %d, fromUsername: %s", fromId, fromUsername))

	user, err := srv.Db.GetUserById(fromId)
	if err != nil {
		return fmt.Errorf("CQ_restart_game GetUserById err: %v", err)
	}
	if user.CreatedAt != "" && srv.IsIgnoreUser(fromId) {
		return nil
	}

	err = srv.Db.AddNewUser(fromId, fromUsername, fromFirstName)
	if err != nil {
		return fmt.Errorf("CQ_restart_game AddNewUser err: %v", err)
	}
	srv.Db.EditBotState(fromId, "")
	srv.Db.EditLives(fromId, 3)
	srv.SendMessageAndDb(fromId, fmt.Sprintf("Привет, %s 👋", fromFirstName))

	srv.Db.EditStep(fromId, "1")
	srv.SendAnimMessageHTML("1", fromId, animTimeout3000)

	time.Sleep(time.Millisecond * time.Duration(animTimeoutTest))
	
	text := "Прямо сейчас начинай игру и забирай бонус 1000₽ за уверенный старт! 🚀"
	replyMarkup := `{"inline_keyboard" : [
		[{ "text": "Начать игру", "callback_data": "start_game" }]
	]}`
	srv.SendMessageWRM(fromId, text, replyMarkup)
	
	srv.SendMsgToServer(fromId, "bot", text)
	srv.Db.UpdateLatsActiontime(fromId)

	return nil
}

func (srv *TgService) CQ_subscribe(m models.Update) error {
	cq := m.CallbackQuery
	fromId := cq.From.Id
	fromUsername := cq.From.UserName
	// fromFirstName := cq.From.FirstName
	srv.l.Info(fmt.Sprintf("CQ_subscribe: fromId: %d, fromUsername: %s", fromId, fromUsername))

	user, _ := srv.Db.GetUserById(fromId)
	ChatToCheck := -1001654676322
	// if user.Ref == "ref15" {
	// 	ChatToCheck = -1001771020146
	// }

	GetChatMemberResp, err := srv.GetChatMember(fromId, ChatToCheck)
	if err != nil {
		return fmt.Errorf("CQ_subscribe GetChatMember fromId: %d, ChatToCheck: %d, err: %v", fromId, ChatToCheck, err)
	}
	if GetChatMemberResp.Result.Status != "member" && GetChatMemberResp.Result.Status != "creator" {
		logMess := fmt.Sprintf("CQ_subscribe GetChatMember bad resp: %+v", GetChatMemberResp)
		srv.l.Error(logMess)
		mess := "❌ вы не подписаны на канал!"
		srv.SendMessageAndDb(fromId, mess)
		srv.Db.EditStep(fromId, mess)
		return nil
	}

	go func() {
		time.Sleep(time.Second)
		srv.EditMessageReplyMarkup(fromId, cq.Message.MessageId)
	}()

	// srv.SendMessage(fromId, "Отлично! Осталось последнее условие 😎\nСмотри кружочек 👇🏻")
	time.Sleep(time.Second)
	
	// base64Str := srv.CreateBase64UserData(fromId, fromUsername, fromFirstName)
	// reglink := fmt.Sprintf("https://goopro.store/api/v1/redirect/1000153272?register=1&data=%s", base64Str)
	// reply_markup := fmt.Sprintf(`{"inline_keyboard" : [
	// 	[{ "text": "Зарегистрироваться", "url": "%s" }]
	// ]}`, reglink)
	
	// futureJson := map[string]string{
	// 	"video_note":   fmt.Sprintf("@%s", "./files/krug_2.mp4"),
	// 	"chat_id": strconv.Itoa(fromId),
	// 	"reply_markup": reply_markup,
	// }
	// cf, body, err := files.CreateForm(futureJson)
	// if err != nil {
	// 	return fmt.Errorf("CQ_subscribe CreateFormV2 err: %v", err)
	// }
	// srv.SendVideoNote(body, cf)

	// textMess := fmt.Sprintf(
	// 	"Переходи и регистрируйся по ссылке:\n\n%s\n\nДалее присылай сюда почту, на которую регистрировался 👇🏻",
	// 	srv.ChInfoToLinkHTML("https://goopro.store/api/v1/redirect/1000153272?register=1", "ССЫЛКА"),
	// )
	// srv.SendMessageHTML(fromId, textMess)

	// srv.Db.EditBotState(fromId, "wait_email")

	srv.Db.EditBotState(fromId, "")
	// srv.Db.EditEmail(fromId, msgTextEmail)
	// lichka, _,  _ := srv.GetLichka()
	// srv.Db.EditLichka(fromId, lichka)
	// lichkaUrl := fmt.Sprintf("https://t.me/%s", srv.DelAt(lichka))
	// mess := fmt.Sprintf("Ваша личка %s", srv.AddAt(lichka))
	// srv.SendMessage(fromId, mess)

	// url := fmt.Sprintf("%s/api/v1/lichka", srv.Cfg.ServerUrl)
	// jsonBody := []byte(fmt.Sprintf(`{"lichka":"%s", "tg_id":"%d", "tg_username":"%s", "tg_name":"%s", "email":"%s"}`, lichka, tgId, fromUsername, fromFirstName, msgTextEmail))
	// bodyReader := bytes.NewReader(jsonBody)
	// _, err := http.Post(url, "application/json", bodyReader)
	// if err != nil {
	// 	return fmt.Errorf("M_state api/v1/lichka Post err: %v", err)
	// }
	// url = fmt.Sprintf("%s/api/v1/link_ref", srv.Cfg.ServerUrl)
	// ref_id := srv.Refki[user.Ref]
	// if ref_id != "хуй" {
	// 	ref_id = "1000153272"
	// }
	// jsonBody = []byte(fmt.Sprintf(`{"user_email":"%s", "ref_id":"%s"}`, msgTextEmail, ref_id))
	// bodyReader = bytes.NewReader(jsonBody)
	// _, err = http.Post(url, "application/json", bodyReader)
	// if err != nil {
	// 	return fmt.Errorf("CQ_subscribe api/v1/link_ref Post err: %v", err)
	// }

	// gifResp, _ := srv.CopyMessage(fromId, -1002074025173, 86) // https://t.me/c/2074025173/86
	// gifResp, _ := srv.SendVideoWCaption(fromId, "", "./files/gif_1.MOV")
	// time.Sleep(time.Second*6)
	// srv.DeleteMessage(fromId, gifResp.Result.MessageId)

	// otSum := "800.000₽"
	// if user.Ref == "ref15" {
	// 	otSum = "500.000₽"
	// }
	mess := fmt.Sprintf("Поздравляю, ты успешно выполнил все условия и выйграл 5.000 рублей 🎉💰\n\nНаш менеджер свяжется с тобой через этого бота в течение 12 часов ☑️")
	reply_markup := `{"inline_keyboard" : [[{ "text": "Забрать награду", "callback_data": "zabrat_nagradu" }]]}`
	srv.SendMessageWRM(fromId, mess, reply_markup)
	// messId := messResp.Result.MessageId
	// srv.Db.EditNotDelMessId(fromId, messId)
	srv.SendMsgToServer(fromId, "bot", mess)

	// instrLink := "https://telegra.ph/Algoritm-dejstvij-05-04"
	// reply_markup := fmt.Sprintf(`{"inline_keyboard" : [
	// 	[{ "text": "Забрать инструкцию", "url": "%s" }]
	// ]}`, instrLink)
	// reply_markup = `{"inline_keyboard" : [ [{ "text": "Забрать инструкцию", "callback_data": "zabrat_instr" }]]}`

	// futureJson := map[string]string{
	// 	"video_note":   fmt.Sprintf("@%s", "./files/krug_3.mp4"),
	// 	"chat_id": strconv.Itoa(fromId),
	// 	"reply_markup": reply_markup,
	// }
	// cf, body, err := files.CreateForm(futureJson)
	// if err != nil {
	// 	return fmt.Errorf("CQ_subscribe CreateFormV2 err: %v", err)
	// }
	// srv.SendVideoNote(body, cf)

	// go func() {
	// 	time.Sleep(time.Minute)
	// 	// instrLink := "https://telegra.ph/Algoritm-dejstvij-05-04"
	// 	// reply_markup := fmt.Sprintf(`{"inline_keyboard" : [
	// 	// 	[{ "text": "Заработать 500.000₽", "url": "%s" }]
	// 	// ]}`, instrLink)
		
	// 	reply_markup := fmt.Sprintf(`{"inline_keyboard" : [
	// 		[{ "text": "Написать Марку в ЛС", "url": "%s" }],
	// 		[{ "text": "Обо мне", "callback_data": "obo_nme_btn" }],
	// 		[{ "text": "Информация о заработке", "callback_data": "info_o_zarabotke_btn" }],
	// 		[{ "text": "Часто задаваемые вопросы", "callback_data": "frequently_questions_btn" }],
	// 		[{ "text": "Отзывы", "callback_data": "show_reviews_btn" }]
	// 	]}`, lichkaUrl)

	// 	futureJson := map[string]string{
	// 		"video_note":   fmt.Sprintf("@%s", "./files/krug_4.mp4"),
	// 		"chat_id": strconv.Itoa(fromId),
	// 		"reply_markup": reply_markup,
	// 	}
	// 	cf, body, err := files.CreateForm(futureJson)
	// 	if err != nil {
	// 		err := fmt.Errorf("CQ_subscribe CreateFormV2 err: %v", err)
	// 		srv.l.Error(err)
	// 	}
	// 	srv.SendVideoNote(body, cf)
	// }()

	return nil
}

func (srv *TgService) CQ_show_reviews_btn(m models.Update) error {
	cq := m.CallbackQuery
	fromId := cq.From.Id
	fromUsername := cq.From.UserName
	// fromFirstName := cq.From.FirstName
	srv.l.Info(fmt.Sprintf("CQ_subscribe: fromId: %d, fromUsername: %s", fromId, fromUsername))

	user, _ := srv.Db.GetUserById(fromId)
	lichka := user.Lichka
	if lichka == "" {
		lichka = "https://t.me/markodinncov"
	}
	// lichkaUrl := fmt.Sprintf("https://t.me/%s", srv.DelAt(lichka))

	srv.SendVideoNoteCurrFile(fromId, "./files/krug_reviews.mp4")
	srv.Send3Kruga(fromId)

	return nil
}

func (srv *TgService) CQ_obo_nme_btn(m models.Update) error {
	cq := m.CallbackQuery
	fromId := cq.From.Id
	fromUsername := cq.From.UserName
	srv.l.Info(fmt.Sprintf("CQ_obo_nme_btn: fromId: %d, fromUsername: %s", fromId, fromUsername))

	user, _ := srv.Db.GetUserById(fromId)
	lichka := user.Lichka
	if lichka == "" {
		lichka = "https://t.me/markodinncov"
	}
	// lichkaUrl := fmt.Sprintf("https://t.me/%s", srv.DelAt(lichka))

	futureJson := map[string]string{
		"video_note":   fmt.Sprintf("@%s", "./files/krug_obo_mne.mp4"),
		"chat_id": strconv.Itoa(fromId),
		"reply_markup": `{"inline_keyboard" : [
			[{ "text": "Обо мне (ч1)", "url": "https://telegra.ph/KTO-HOCHET-STAT-MILLIONEROM-05-20" }],
			[{ "text": "Обо мне (ч2)", "url": "https://telegra.ph/Rezultaty-i-dokazatelstva-05-20" }],
			[{ "text": "Обо мне (ч3)", "url": "https://telegra.ph/Kak-ehto-mozhet-pomenyat-zhizn-lyubogo-cheloveka-05-22" }]
		]}`,
	}
	cf, body, err := files.CreateForm(futureJson)
	if err != nil {
		err := fmt.Errorf("CQ_obo_nme_btn CreateForm err: %v", err)
		srv.l.Error(err)
	}
	srv.SendVideoNote(body, cf)

	srv.Send3Kruga(fromId)

	return nil
}

func (srv *TgService) CQ_info_o_zarabotke_btn(m models.Update) error {
	cq := m.CallbackQuery
	fromId := cq.From.Id
	fromUsername := cq.From.UserName
	srv.l.Info(fmt.Sprintf("CQ_info_o_zarabotke_btn: fromId: %d, fromUsername: %s", fromId, fromUsername))

	user, _ := srv.Db.GetUserById(fromId)
	lichka := user.Lichka
	if lichka == "" {
		lichka = "https://t.me/markodinncov"
	}
	// lichkaUrl := fmt.Sprintf("https://t.me/%s", srv.DelAt(lichka))

	futureJson := map[string]string{
		"video_note":   fmt.Sprintf("@%s", "./files/krug_info_o_zarabotke.mp4"),
		"chat_id": strconv.Itoa(fromId),
		"reply_markup": `{"inline_keyboard" : [
			[{ "text": "О заработке (ч1)", "url": "https://telegra.ph/KTO-HOCHET-STAT-MILLIONEROM-05-20" }],
			[{ "text": "О заработке (ч2)", "url": "https://telegra.ph/Rezultaty-i-dokazatelstva-05-20" }],
			[{ "text": "О заработке (ч3)", "url": "https://telegra.ph/Kak-ehto-mozhet-pomenyat-zhizn-lyubogo-cheloveka-05-22" }]
		]}`,
	}
	cf, body, err := files.CreateForm(futureJson)
	if err != nil {
		err := fmt.Errorf("CQ_info_o_zarabotke_btn CreateForm err: %v", err)
		srv.l.Error(err)
	}
	srv.SendVideoNote(body, cf)

	srv.Send3Kruga(fromId)

	return nil
}

func (srv *TgService) CQ_frequently_questions_btn(m models.Update) error {
	cq := m.CallbackQuery
	fromId := cq.From.Id
	fromUsername := cq.From.UserName
	srv.l.Info(fmt.Sprintf("CQ_info_o_zarabotke_btn: fromId: %d, fromUsername: %s", fromId, fromUsername))

	user, _ := srv.Db.GetUserById(fromId)
	lichka := user.Lichka
	if lichka == "" {
		lichka = "https://t.me/markodinncov"
	}
	// lichkaUrl := fmt.Sprintf("https://t.me/%s", srv.DelAt(lichka))


	
	messTxt := `❓Ответы на часто задаваемые вопросы:

	<b>• Как я могу понять, что схема работает?</b>
	
	- Проверить мои схемы вы можете в демо-режиме, открутив их несколько раз и набить руку.
	Так же в своем канале я публикую подробные откруты, на которых видно, что все схемы полностью рабочие
	
	<b>• Зачем тебе это все? В чем твоя выгода?</b>
	
	- Я не строю из себя благодетеля, а прямым текстом говорю, что делаю это, исходя из своей выгоды. Вы откручиваете схему и отправляете мне 20% с выигрыша. Справедливая сделка win-win
	
	<b>• Как я могу быть уверен, что ты не мошенник?</b>
	
	- Я предоставляю реальный заработок и не беру никаких денег до того момента, пока вы не сделаете вывод себе на карту. 
	Для начала можете зайти в демо и прокрутить схему там, алгоритм рабочий и работает всегда, нет разницы демо либо реальный счет, но убедиться в этом вы можете именно на демо счете. Так же я не скрываю ни своего лица, ни своего местонахождения. А на моем канале вы можете найти кучу отзывов от довольных членов моей команды. При необходимости могу созвониться с вами.
	Комментарии в своем канале я не могу открыть по элементарным причинам - казино сразу же начинает обваливать на меня массовый спам ботами, которые пишут гневные комментарии. Если вы хотите получить контакты людей, которые уже крутили схему - напишите мне в лс и я без проблем поделюсь с вами. В канале эти ссылки опубликовать не могу, так как вы начнете заваливать сообщениями моих ребят, а это ни к чему)
	
	<b>• Как часто можно крутить схему?</b>
	
	- С одного устройства и аккаунта можно крутить не более одного раза в неделю, чтобы не вызывать подозрений у тех.поддержки казика
	
	<b>• А как казино до сих пор не спалило твои схемы? Там же столько выводов каждый день, уже бы давно закрыли всё или там какие-то дураки сидят по-твоему?</b>
	
	- Для этого мы с командой каждый день обновляем схемы, алгоритмы, суммы пополнения и т.д. Так же там есть люди, которые просто крутят слоты и даже не догадываются о моем существовании. Лудоманы проигрывают в казиках миллионы долларов каждый день. Поэтому наши выводы для них - как иголка в стоге сена.
	
	<b>• Почему ты сам просто не крутишь своими схемы много раз в день?</b>
	
	- Я выстраиваю структуру своей работы так, чтобы мне не приходилось самому делать фактически ничего, кроме того, как заниматься разработкой схем. Я бы мог и сам спокойно крутить их целыми днями кучу раз, но это сопровождается возней с аккаунтами, картами, банками и т.д. Поэтому мне проще набирать людей в команду, которые будут стабильно работать по моим схемам и скидывать мне процент.
	
	<b>• Почему ты не одалживаешь и не даешь деньги на открут схемы?</b>
	
	- Сам посмотри на абсурд всей ситуации. Ты приходишь ко мне в команду на все готовенькое. Все что от тебя требуется - это найти небольшую сумму, открутить по схеме, вывести бабки и отправить 20%. Но в то же время, люди еще умудряются клянчить у меня денег на депозит для схемы. Это все очень меня злит и огорчает, поэтому даже не советую заниматься подобным в общении со мной.`
	
	_, err := srv.SendMessageHTML(fromId, messTxt)
	if err != nil {
		srv.l.Error(fmt.Sprintf("CQ_frequently_questions_btn SendMessageHTML err: %v", err))
	}

	srv.Send3Kruga(fromId)

	return nil
}

func (srv *TgService) CQ_delete_user_by_username_btn(m models.Update) error {
	cq := m.CallbackQuery
	fromId := cq.From.Id
	fromUsername := cq.From.UserName
	srv.l.Info(fmt.Sprintf("CQ_delete_user_by_username_btn: fromId: %d, fromUsername: %s", fromId, fromUsername))

	srv.SendForceReply(fromId, DEL_USER_MSG)
	return nil
}

func (srv *TgService) CQ_delete_user_by_id_btn(m models.Update) error {
	cq := m.CallbackQuery
	fromId := cq.From.Id
	fromUsername := cq.From.UserName
	srv.l.Info(fmt.Sprintf("CQ_delete_user_by_id_btn: fromId: %d, fromUsername: %s", fromId, fromUsername))

	srv.SendForceReply(fromId, DEL_USER_ID_MSG)
	return nil
}


func (srv *TgService) CQ_davay_sigraem_btn(m models.Update) error {
	cq := m.CallbackQuery
	fromId := cq.From.Id
	fromUsername := cq.From.UserName
	srv.l.Info(fmt.Sprintf("CQ_davay_sigraem_btn: fromId: %d, fromUsername: %s", fromId, fromUsername))

	mess := "В кратце расскажу тебе правила:\n\nДля победы тебе нужно дать правильные ответы на 3 вопроса!\nУ тебя имеется 3 попытки, но постарайся пройти с первой 😎\n\nЖелаю удачи! ✊\nЖми кнопку 👇"
	replyMarkup :=`{"inline_keyboard" : [
		[ { "text": "Погнали!", "callback_data": "pognaly_btn" } ]
	]}`
	srv.SendMessageWRM(fromId, mess, replyMarkup)
	return nil
}

func (srv *TgService) CQ_pognaly_btn(m models.Update) error {
	cq := m.CallbackQuery
	fromId := cq.From.Id
	fromUsername := cq.From.UserName
	srv.l.Info(fmt.Sprintf("CQ_davay_sigraem_btn: fromId: %d, fromUsername: %s", fromId, fromUsername))

	srv.ShowMilQ(fromId, 1)

	return nil
}

