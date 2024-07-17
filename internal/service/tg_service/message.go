package tg_service

import (
	"encoding/json"
	"fmt"
	"io"
	"myapp/internal/models"
	"myapp/pkg/files"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (srv *TgService) HandleMessage(m models.Update) error {
	msgText := m.Message.Text
	fromUsername := m.Message.From.UserName
	fromId := m.Message.From.Id
	srv.l.Info(fmt.Sprintf("HandleMessage: fromId-%d fromUsername-%s, msgText-%s", fromId, fromUsername, msgText))

	srv.SendMsgToServer(fromId, "user", msgText)

	if msgText == "/admin" {
		err := srv.M_admin(m)
		if err != nil {
			srv.SendMessage(fromId, ERR_MSG)
			srv.SendMessage(fromId, err.Error())
		}
		return err
	}

	// user, err := srv.Db.GetUserById(fromId)
	// if err != nil {
	// 	return fmt.Errorf("HandleMessage GetUserById err: %v", err)
	// }
	// if user.Id != 0 && user.Lives == 0 {
	// 	return nil
	// }

	if msgText == "/help" {
		srv.SendMessageAndDb(fromId, "@millioner_support\nвот контакт для связи")
		srv.Db.UpdateLatsActiontime(fromId)
		srv.Db.UpdateFeedbackTime(fromId)
		return nil
	}

	if msgText == "Условия розыгрыша" {
		userPersonalRef := srv.GetUserPersonalRef(fromId)
		chLink := "https://t.me/+rLIklQb0ALNhZjEx"
		chLink2 := "https://t.me/geniusgiveaway"
		messText := fmt.Sprintf("Условие первое:\nПодпишись на эти каналы 👇\n\n %s\n %s", chLink, chLink2)
		mess := fmt.Sprintf("Условие второе:\nВыложи себе в инстаграм stories нашу картинку и в соответствующем поле отметь инстаграм-аккаунт раздачи: %s %s.\n\nИли:\nПригласи двух друзей по своей уникальной ссылке: %s. Отправь ссылку друзьям.", "@mrgeniuz1", srv.ChInfoToLinkHTML("https://www.instagram.com/mrgeniuz1", "(прямая ссылка на профиль)"), userPersonalRef)
		
		fullMess := fmt.Sprintf("%s\n\n%s", messText, mess)

		reply_markup := `{"inline_keyboard" : [
			[ { "text": "☑️ Отметил", "callback_data": "otmetil_btn" }, { "text": "☑️ Пригласил", "callback_data": "priglasil_btn" } ]
		]}`
		fileNameInServer := "./files/inst_story_draft.jpeg"
		_, err := srv.SendDocumentWCaptionWRM(fromId, fullMess, fileNameInServer, reply_markup)
		if err != nil {
			return fmt.Errorf("HandleMessage SendDocumentWCaptionWRM err: %v", err)
		}

		// srv.SendMessageAndDb(fromId, fullMess)

		srv.Db.UpdateLatsActiontime(fromId)
		srv.Db.UpdateFeedbackTime(fromId)
		return nil
	}

	if msgText == "/ref" || msgText == "Мои рефералы" {
		usersByRef, _ := srv.Db.GetUsersByRef(strconv.Itoa(fromId))
		userPersonalRef := srv.GetUserPersonalRef(fromId)
		srv.SendMessageAndDb(fromId, fmt.Sprintf("Ваша рефералка: %s\nВаши рефералы: %d шт.", userPersonalRef, len(usersByRef)))
		srv.Db.UpdateLatsActiontime(fromId)
		srv.Db.UpdateFeedbackTime(fromId)
		return nil
	}

	// if user.IsLastPush == 1 {
	// 	srv.SendMessageAndDb(fromId, "бот вам больше не доступен")
	// 	return nil
	// }

	if strings.HasPrefix(msgText, "/start") { // https://t.me/tgbotusername?start=ref01 -> /start ref01
		err := srv.M_start(m)
		if err != nil {
			srv.SendMessageAndDb(fromId, ERR_MSG)
			srv.SendMessageAndDb(fromId, err.Error())
		}
		srv.Db.UpdateLatsActiontime(fromId)
		srv.Db.UpdateFeedbackTime(fromId)
		return err
	}

	err := srv.M_state(m)
	if err != nil {
		srv.SendMessageAndDb(fromId, ERR_MSG)
		srv.SendMessageAndDb(fromId, err.Error())
	}
	srv.Db.UpdateLatsActiontime(fromId)
	srv.Db.UpdateFeedbackTime(fromId)
	return err
}

func (srv *TgService) M_start(m models.Update) error {
	fromId := m.Message.Chat.Id
	msgText := m.Message.Text
	fromFirstName := m.Message.From.FirstName
	fromUsername := m.Message.From.UserName
	srv.l.Info(fmt.Sprintf("M_start: fromId: %d, fromUsername: %s, msgText: %s", fromId, fromUsername, msgText))

	refArr := strings.Split(msgText, " ")
	ref := ""
	if len(refArr) > 1 {
		ref = refArr[1]
	}

	// user, err := srv.Db.GetUserById(fromId)
	// if err != nil {
	// 	return fmt.Errorf("M_start GetUserById err: %v", err)
	// }
	// if user.CreatedAt != "" && srv.IsIgnoreUser(fromId) {
	// 	return nil
	// }

	err := srv.Db.AddNewUser(fromId, fromUsername, fromFirstName)
	if err != nil {
		return fmt.Errorf("M_start AddNewUser err: %v", err)
	}
	srv.Db.EditRef(fromId, ref)
	// lichka := "odincovmarkk"
	// if ref == "ref15" {
	// 	lichka = "markodinncov"
	// }
	// srv.Db.EditLichka(fromId, lichka)
	if fromId == 1394096901 {
		srv.Db.EditAdmin(fromId, 1)
	}
	srv.Db.EditBotState(fromId, "")
	// srv.Db.EditLives(fromId, 3)
	// srv.Db.EditStep(fromId, "1")
	mess := fmt.Sprintf("Привет, %s \n\nПредлагаю сыграть тебе в 'Кто хочет стать миллионером?' 🌀\n\nЕсли сможешь выиграть - отправлю тебе 5000 рублей 💸\n\nПопробуем? 😏", fromFirstName)
	replyMarkup :=`{"inline_keyboard" : [
		[ { "text": "Давай сыграем!", "callback_data": "davay_sigraem_btn" } ]
	]}`
	srv.SendMessageWRM(fromId, mess, replyMarkup)


	return nil
}

func (srv *TgService) PushKrugToUsers(m models.Update) {
	time.Sleep(time.Hour * 24)

	srv.CQ_obo_nme_btn(m)

	time.Sleep(time.Hour * 3)

	srv.CQ_info_o_zarabotke_btn(m)

	time.Sleep(time.Hour * 3)

	srv.CQ_show_reviews_btn(m)

	time.Sleep(time.Hour * 3)

	srv.CQ_frequently_questions_btn(m)
}

func (srv *TgService) M_state(m models.Update) error {
	fromId := m.Message.Chat.Id
	msgText := m.Message.Text
	fromUsername := m.Message.From.UserName
	// fromFirstName := m.Message.From.FirstName
	srv.l.Info(fmt.Sprintf("M_state: fromId: %d, fromUsername: %s, msgText: %s", fromId, fromUsername, msgText))

	user, err := srv.Db.GetUserById(fromId)
	if err != nil {
		srv.l.Warn(fmt.Errorf("M_state GetUserById err: %v", err))
	}
	// srv.Db.UpdateLatsActiontime(fromId)
	// if user.BotState == "" {
	// 	return nil
	// }

	if user.BotState == "read_article_after_KNB_win" { // Го, ко, коу, гоу, гэу
		if !strings.HasPrefix(strings.ToLower(msgText), "гоу") && !strings.HasPrefix(strings.ToLower(msgText), "го") && !strings.HasPrefix(strings.ToLower(msgText), "ко") && !strings.HasPrefix(strings.ToLower(msgText), "коу") && !strings.HasPrefix(strings.ToLower(msgText), "гэу") && !strings.HasPrefix(strings.ToLower(msgText), "go") {
			srv.SendMessageAndDb(fromId, "❌ Вы неверно ввели кодовое слово, сверьтесь с кружочком и попробуйте еще раз")
			return nil
		}
		srv.Db.EditBotState(fromId, "")
		// srv.SendAnimMessage("-1", fromId, animTimeout500)
		// srv.SendBalance(fromId, "30.000", animTimeoutTest)
		// srv.SendAnimMessage("9", fromId, animTimeoutTest)
		// srv.Db.EditStep(fromId, "9")

		otSum := "800.000₽"
		if user.Ref == "ref15" {
			otSum = "500.000₽"
		}
		text := fmt.Sprintf("Ну что, поехали, ответь правильно на 3 вопроса и уже сегодня сможешь заработать %s 😏", otSum)
		// replyMarkup :=`{"inline_keyboard" : [
		// 	[ { "text": "Давай попробуем", "callback_data": "show_q_7_" } ]
		// ]}`
		srv.SendMessageAndDb(fromId, text)

		err = srv.ShowMilQ(fromId, 1)
		if err != nil {
			return fmt.Errorf("M_state ShowMilQ err: %v", err)
		}

		srv.Db.EditStep(fromId, text)
		return nil
	}

	if user.BotState == "wait_email" {
		msgTextEmail := msgText
		url := fmt.Sprintf("%s/api/v1/user?email=%s", srv.Cfg.ServerUrl, msgTextEmail)
		srv.l.Info("M_state wait_email иду к API", url)
		response, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("M_state wait_email Post err: %v", err)
		}
		srv.l.Info("M_state wait_email сходил к API")
		defer response.Body.Close()
	
		if response.StatusCode != http.StatusOK {
			bodyBytes, err := io.ReadAll(response.Body)
			if err != nil {
				return fmt.Errorf("M_state wait_email ReadAll err: %v", err)
			}
			return fmt.Errorf("M_state wait_email post %s bad response: [%d] %v", url, response.StatusCode, string(bodyBytes))
		}
	
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("M_state wait_email ReadAll err: %v", err)
		}
	
		resp := struct{
			Status string `json:"status"`
			Data   string `json:"data"`
		}{}
		json.Unmarshal(bodyBytes, &resp)
	
		if resp.Status == "success" {

			srv.Db.EditBotState(fromId, "")
			srv.Db.EditEmail(fromId, msgTextEmail)
			// lichka, tgId,  _ := srv.GetLichka()
			// srv.Db.EditLichka(fromId, lichka)
			// mess := fmt.Sprintf("Ваша личка %s", srv.AddAt(lichka))
			// srv.SendMessage(fromId, mess)

			// url := fmt.Sprintf("%s/api/v1/lichka", srv.Cfg.ServerUrl)
			// jsonBody := []byte(fmt.Sprintf(`{"lichka":"%s", "tg_id":"%d", "tg_username":"%s", "tg_name":"%s", "email":"%s"}`, lichka, tgId, fromUsername, fromFirstName, msgTextEmail))
			// bodyReader := bytes.NewReader(jsonBody)
			// _, err := http.Post(url, "application/json", bodyReader)
			// if err != nil {
			// 	return fmt.Errorf("M_state api/v1/lichka Post err: %v", err)
			// }
			url = fmt.Sprintf("%s/api/v1/link_ref", srv.Cfg.ServerUrl)
			ref_id := srv.Refki[user.Ref]
			if ref_id != "хуй" {
				ref_id = "1000153272"
			}
			// jsonBody = []byte(fmt.Sprintf(`{"user_email":"%s", "ref_id":"%s"}`, msgTextEmail, ref_id))
			// bodyReader = bytes.NewReader(jsonBody)
			// _, err = http.Post(url, "application/json", bodyReader)
			// if err != nil {
			// 	return fmt.Errorf("M_state api/v1/link_ref Post err: %v", err)
			// }

			gifResp, _ := srv.CopyMessage(fromId, -1002074025173, 86) // https://t.me/c/2074025173/86
			// gifResp, _ := srv.SendVideoWCaption(fromId, "", "./files/gif_1.MOV")
			time.Sleep(time.Second*6)
			srv.DeleteMessage(fromId, gifResp.Result.MessageId)

			// mess = "Поздравляю тебя! 🎉\n\nЯ уже проверил сегодняшний алгоритм и прописал необходимые настройки, благодаря которым ты уже сегодня сможешь вытащить солидную прибыль.\n\nНиже отправляю тебе инструкцию, повторив которую ты уже сегодня заработаешь 500.000₽👇\n\nВсё работает на 1.000%! Попробуй и убедись🤝"
			// srv.SendMessageAndDb(fromId, mess)

			// instrLink := "https://telegra.ph/Algoritm-dejstvij-05-04"
			// reply_markup := fmt.Sprintf(`{"inline_keyboard" : [
			// 	[{ "text": "Забрать инструкцию", "url": "%s" }]
			// ]}`, instrLink)
			reply_markup := `{"inline_keyboard" : [ [{ "text": "Забрать инструкцию", "callback_data": "zabrat_instr" }]]}`

			futureJson := map[string]string{
				"video_note":   fmt.Sprintf("@%s", "./files/krug_3.mp4"),
				"chat_id": strconv.Itoa(fromId),
				"reply_markup": reply_markup,
			}
			cf, body, err := files.CreateForm(futureJson)
			if err != nil {
				return fmt.Errorf("M_state CreateFormV2 err: %v", err)
			}
			srv.SendVideoNote(body, cf)

			go func() {
				time.Sleep(time.Minute)
				// instrLink := "https://telegra.ph/Algoritm-dejstvij-05-04"
				// reply_markup := fmt.Sprintf(`{"inline_keyboard" : [
				// 	[{ "text": "Заработать 500.000₽", "url": "%s" }]
				// ]}`, instrLink)
				reply_markup := `{"inline_keyboard" : [ [{ "text": "Заработать 500.000₽", "callback_data": "zabrat_instr_500" }]]}`

				futureJson = map[string]string{
					"video_note":   fmt.Sprintf("@%s", "./files/krug_4.mp4"),
					"chat_id": strconv.Itoa(fromId),
					"reply_markup": reply_markup,
				}
				cf, body, err = files.CreateForm(futureJson)
				if err != nil {
					err := fmt.Errorf("M_state CreateFormV2 err: %v", err)
					srv.l.Error(err)
				}
				srv.SendVideoNote(body, cf)
			}()
		} else {
			srv.SendMessage(fromId, "❌ Почта не найдена")
		}
	}

	if user.BotState == "read_article_after_OIR_win" {
		if !strings.HasPrefix(strings.ToLower(msgText), "рез") && !strings.HasPrefix(strings.ToLower(msgText), "риз") {
			srv.SendMessageAndDb(fromId, "❌ Вы неверно ввели кодовое слово, сверьтесь со статьей и попробуйте еще раз")
			return nil
		}
		srv.Db.EditBotState(fromId, "")
		srv.SendAnimMessage("-1", fromId, animTimeout500)
		srv.SendBalance(fromId, "55.000", animTimeoutTest)
		srv.SendAnimMessageHTML("11", fromId, animTimeoutTest)
		srv.Db.EditStep(fromId, "11")
		time.Sleep(time.Second)

		text :=  "Готов перейти к первому вопросу? 😏"
		replyMarkup := `{"inline_keyboard" : [
			[{ "text": "Ествественно! Погнали!", "callback_data": "show_q_10_" }]
		]}`
		srv.SendMessageWRM(fromId, text, replyMarkup)
		srv.Db.EditStep(fromId, text)

		return nil
	}

	if user.BotState == "read_article_after_TrurOrFalse_win" {
		if !strings.HasPrefix(strings.ToLower(msgText), "син") {
			srv.SendMessageAndDb(fromId, "❌ Вы неверно ввели кодовое слово, сверьтесь со статьей и попробуйте еще раз")
			return nil
		}
		srv.Db.EditBotState(fromId, "")
		srv.SendAnimMessage("-1", fromId, animTimeout500)
		srv.SendBalance(fromId, "100.000", animTimeoutTest)
		srv.SendAnimMessageHTML("13", fromId, animTimeoutTest)
		srv.Db.EditStep(fromId, "13")
		// srv.CopyMessage(fromId, 1394096901, 925)
		srv.CopyMessage(fromId, -1002074025173, 22)
		time.Sleep(time.Second)
		
		text :=  "Если тобою прочитаны все статьи, то ты прямо сейчас можешь забрать свою награду стоимостью в 100.000₽ 💸\n\nЖми кнопку ниже ⬇️"
		replyMarkup := fmt.Sprintf(`{"inline_keyboard" : [
			[{ "text": "пустая рефка (%s)", "url": "https://t.me/threeprocentsclub_bot" }]
		]}`, user.Ref)
		if user.Ref == "bot1" {
			replyMarkup = fmt.Sprintf(`{"inline_keyboard" : [
				[{ "text": "Забрать награду (%s)", "url": "https://t.me/threeprocentsclub_bot" }]
			]}`, user.Ref)
		}
		if user.Ref == "bot2" {
			replyMarkup = fmt.Sprintf(`{"inline_keyboard" : [
				[{ "text": "Забрать награду (%s)", "url": "https://t.me/threeprocentsclub2_bot" }]
			]}`, user.Ref)
		}
		srv.SendMessageWRM(fromId, text, replyMarkup)

		srv.SendMsgToServer(fromId, "bot", text)
		srv.Db.EditLatsActiontime(fromId, "")
		srv.Db.EditIsFinal(fromId, 1)

		return nil
	}

	if user.BotState == "wait_inst_link" {

		username := srv.DelAt(msgText)
		mention_usernamme := "mrgeniuz1"

		checkInstStoryResp, err := srv.CheckInstStory(username, mention_usernamme)
		if err != nil {
			err := fmt.Errorf("M_state CheckInstStory err: %v", err)
			return err
		}
		srv.l.Info(fmt.Sprintf("M_state checkInstStoryResp: %+v", checkInstStoryResp))
		if checkInstStoryResp.Marked {
			mess := fmt.Sprintf("🎉 Поздравляю, ты участвуешь в розыгрыше 5 000 ₽! Переходи в канал раздачи, там объявим победителей в прямом эфире. Обязательно приходи👇")
			replyMarkup := `{"inline_keyboard" : [
				[{ "text": "Узнать итоги", "url": "https://t.me/geniusgiveaway" }]
			]}`
			_, err := srv.SendMessageWRM(fromId, mess, replyMarkup)
			if err != nil {
				return fmt.Errorf("M_state SendMessageWRM err: %v", err)
			}
			srv.SendMsgToServer(fromId, "bot", mess)
	
			srv.Db.EditBotState(fromId, "")
			return nil
		} else {
			mess := "🤔Что-то не так. Обычно такое случается, если ссылка на историю неверная или не было отметки. Попробуй ещё раз."
			srv.SendMessageAndDb(fromId, mess)
	
			return nil
		}

	}

	return nil
}

func (srv *TgService) M_admin(m models.Update) error {
	fromId := m.Message.Chat.Id
	msgText := m.Message.Text
	fromUsername := m.Message.From.UserName
	srv.l.Info(fmt.Sprintf("M_admin: fromId: %d, fromUsername: %s, msgText: %s", fromId, fromUsername, msgText))

	u, err := srv.Db.GetUserById(fromId)
	if err != nil {
		return fmt.Errorf("M_admin GetUserById err: %v", err)
	}
	if u.Id == 0 {
		srv.SendMessage(fromId, "Нажмите сначала /start")
	}
	if u.IsAdmin != 1 {
		return fmt.Errorf("_!_")
	}
	err = srv.ShowAdminPanel(fromId)

	return err
}
