package tg_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"myapp/internal/models"
	"myapp/internal/repository/pg"
	"myapp/pkg/logger"
	"myapp/pkg/my_time_parser"
	"net/http"
	"time"

	"github.com/go-co-op/gocron"
)

var (
	mskLoc, _ = time.LoadLocation("Europe/Moscow")
)


var (
	stepsMap = map[string][]string{
		"-1": {"🏦⁣  💰💰💰👈", "🏦⁣ 💰💰💰👈", "🏦⁣💰💰💰👈", "🏦⁣💰💰👈", "🏦⁣💰👈", "🏦⁣👈", "delete"},

		"1": {
			"Я предлагаю тебе сыграть со мной в игру,",
			"Я предлагаю тебе сыграть со мной в игру, где ты будешь получать деньги за каждое правильно выполненное задание 💸",
			"Я предлагаю тебе сыграть со мной в игру, где ты будешь получать деньги за каждое правильно выполненное задание 💸\n\nЕсли пройдешь игру до конца, <b>то сможешь выиграть 100.000 рублей 😳</b>",
			"Я предлагаю тебе сыграть со мной в игру, где ты будешь получать деньги за каждое правильно выполненное задание 💸\n\nЕсли пройдешь игру до конца, <b>то сможешь выиграть 100.000 рублей 😳</b>\nВся игра не займёт более 15 минут ⌛️",
		},
		"2": {
			"<b>Кстати, у тебя есть 3 жизни ❤️</b>",
			"<b>Кстати, у тебя есть 3 жизни ❤️</b>\n\n🥵Если ты долго бездействуешь, то сгорает одна жизнь.",
			"<b>Кстати, у тебя есть 3 жизни ❤️</b>\n\n🥵Если ты долго бездействуешь, то сгорает одна жизнь.\n😔Когда сгорают все 3 жизни, твой баланс обнуляется и ты проигрываешь эту игру",
			"<b>Кстати, у тебя есть 3 жизни ❤️</b>\n\n🥵Если ты долго бездействуешь, то сгорает одна жизнь.\n😔Когда сгорают все 3 жизни, твой баланс обнуляется и ты проигрываешь эту игру\n\nУ тебя будет немного времени, чтобы начать всё заново,",
			"<b>Кстати, у тебя есть 3 жизни ❤️</b>\n\n🥵Если ты долго бездействуешь, то сгорает одна жизнь.\n😔Когда сгорают все 3 жизни, твой баланс обнуляется и ты проигрываешь эту игру\n\nУ тебя будет немного времени, чтобы начать всё заново, но лучше пройди всё с первой попытки🏆",
			"<b>Кстати, у тебя есть 3 жизни ❤️</b>\n\n🥵Если ты долго бездействуешь, то сгорает одна жизнь.\n😔Когда сгорают все 3 жизни, твой баланс обнуляется и ты проигрываешь эту игру\n\nУ тебя будет немного времени, чтобы начать всё заново, но лучше пройди всё с первой попытки🏆\n\nТы прекрасно знаешь, чем обычно заканчивается <i>“откладывание в долгий ящик”🪫</i>",
			"<b>Кстати, у тебя есть 3 жизни ❤️</b>\n\n🥵Если ты долго бездействуешь, то сгорает одна жизнь.\n😔Когда сгорают все 3 жизни, твой баланс обнуляется и ты проигрываешь эту игру\n\nУ тебя будет немного времени, чтобы начать всё заново, но лучше пройди всё с первой попытки🏆\n\nТы прекрасно знаешь, чем обычно заканчивается <i>“откладывание в долгий ящик”🪫</i>\n\nКогда время выйдет полностью - у тебя не будет возможности сыграть,",
			"<b>Кстати, у тебя есть 3 жизни ❤️</b>\n\n🥵Если ты долго бездействуешь, то сгорает одна жизнь.\n😔Когда сгорают все 3 жизни, твой баланс обнуляется и ты проигрываешь эту игру\n\nУ тебя будет немного времени, чтобы начать всё заново, но лучше пройди всё с первой попытки🏆\n\nТы прекрасно знаешь, чем обычно заканчивается <i>“откладывание в долгий ящик”🪫</i>\n\nКогда время выйдет полностью - у тебя не будет возможности сыграть, даже если ты перезапустишь бота ☠️",
		},
		"3": {
			"Поздравляю!",
			"Поздравляю!\nВидишь, как всё просто 🔥",
		},
		"4": {
			"А теперь предлагаю перейти сразу к мясу 🥩",
		},
		"5": {
			"Ответь на следующие 3 вопроса и получишь +10.000₽ к банку! 💸",
		},
		"6": {
			"Воу-воу-воу, палехче 😏",
			"Воу-воу-воу, палехче 😏\n\n+10.000₽ уходят в твой банк за правильные ответы!💸",
		},
		"7": {
			"Поздравляю!🥳",
			"Поздравляю!🥳\nДо этого этапа доходят всего лишь 30% из всех,",
			"Поздравляю!🥳\nДо этого этапа доходят всего лишь 30% из всех, кто запустил бота!",
			"Поздравляю!🥳\nДо этого этапа доходят всего лишь 30% из всех, кто запустил бота!\nСможешь пройти дальше?🤔",
			"Поздравляю!🥳\nДо этого этапа доходят всего лишь 30% из всех, кто запустил бота!\nСможешь пройти дальше?🤔\n\n<b>Ответь еще на 3 вопроса</b>, но учти, они могут быть сложнее предыдущих 😈",
			"Поздравляю!🥳\nДо этого этапа доходят всего лишь 30% из всех, кто запустил бота!\nСможешь пройти дальше?🤔\n\n<b>Ответь еще на 3 вопроса</b>, но учти, они могут быть сложнее предыдущих 😈\n\nПобедитель получит +19.000₽ к банку! 💸",
		},
		"8": {
			"+19.000₽ уходят в твой банк за правильные ответы на вопросы💸",
			"+19.000₽ уходят в твой банк за правильные ответы на вопросы💸\n\n🔐Чтобы разблокировать и забрать награду пришли мне кодовое слово из текста ниже:\n\n<a href=\"https://telegra.ph/KTO-HOCHET-STAT-MILLIONEROM-11-27\">👉🏼《ЧИТАТЬ》👈🏼</a>",
			"+19.000₽ уходят в твой банк за правильные ответы на вопросы💸\n\n🔐Чтобы разблокировать и забрать награду пришли мне кодовое слово из текста ниже:\n\n<a href=\"https://telegra.ph/KTO-HOCHET-STAT-MILLIONEROM-11-27\">👉🏼《ЧИТАТЬ》👈🏼</a>\n\n*Прочтение займёт не более 5 минут.",
			"+19.000₽ уходят в твой банк за правильные ответы на вопросы💸\n\n🔐Чтобы разблокировать и забрать награду пришли мне кодовое слово из текста ниже:\n\n<a href=\"https://telegra.ph/KTO-HOCHET-STAT-MILLIONEROM-11-27\">👉🏼《ЧИТАТЬ》👈🏼</a>\n\n*Прочтение займёт не более 5 минут. После пиши кодовое слово сюда.\nБуду ждать 👇🏻",
		},
		"9": {
			"Я смотрю ты серьезный игрок!",
			"Я смотрю ты серьезный игрок!\nПоэтому повышаю ставки 🔝",
		},
		"10": {
			"+25.000₽ уходят в твой банк за правильные ответы! 💸",
			"+25.000₽ уходят в твой банк за правильные ответы! 💸\n\n🔐Чтобы разблокировать и забрать награду пришли мне кодовое слово из текста ниже:\n\n<a href=\"https://telegra.ph/Rezultaty-i-dokazatelstva-11-27\">👉🏼《ЧИТАТЬ》👈🏼</a>",
			"+25.000₽ уходят в твой банк за правильные ответы! 💸\n\n🔐Чтобы разблокировать и забрать награду пришли мне кодовое слово из текста ниже:\n\n<a href=\"https://telegra.ph/Rezultaty-i-dokazatelstva-11-27\">👉🏼《ЧИТАТЬ》👈🏼</a>\n\n*Прочтение займёт не более 5 минут.",
			"+25.000₽ уходят в твой банк за правильные ответы! 💸\n\n🔐Чтобы разблокировать и забрать награду пришли мне кодовое слово из текста ниже:\n\n<a href=\"https://telegra.ph/Rezultaty-i-dokazatelstva-11-27\">👉🏼《ЧИТАТЬ》👈🏼</a>\n\n*Прочтение займёт не более 5 минут.\nПосле пиши кодовое слово сюда.\nБуду ждать 👇🏻",
		},
		"11": {
			"<b>Ты на завершающем этапе игры🏁</b>",
			"<b>Ты на завершающем этапе игры🏁</b>\n\nДо этого этапа могут пройти только 10% из всех пользователей, которые запустили бота. 😱",
			"<b>Ты на завершающем этапе игры🏁</b>\n\nДо этого этапа могут пройти только 10% из всех пользователей, которые запустили бота. 😱\n\nПоэтому и ставки будут как никогда большими😏",
			"<b>Ты на завершающем этапе игры🏁</b>\n\nДо этого этапа могут пройти только 10% из всех пользователей, которые запустили бота. 😱\n\nПоэтому и ставки будут как никогда большими😏\n\nПобедитель получит +45.000₽ в свой банк💸",
			"<b>Ты на завершающем этапе игры🏁</b>\n\nДо этого этапа могут пройти только 10% из всех пользователей, которые запустили бота. 😱\n\nПоэтому и ставки будут как никогда большими😏\n\nПобедитель получит +45.000₽ в свой банк💸\n\nВсё просто,",
			"<b>Ты на завершающем этапе игры🏁</b>\n\nДо этого этапа могут пройти только 10% из всех пользователей, которые запустили бота. 😱\n\nПоэтому и ставки будут как никогда большими😏\n\nПобедитель получит +45.000₽ в свой банк💸\n\nВсё просто, тебя ждут решающие 3 вопроса, ответив на которые ты наконец-то пройдешь игру и сможешь забрать свою награду 😉",
		},
		"12": {
			"+45.000₽ уходят в твой банк за правильные ответы!💸",
			"+45.000₽ уходят в твой банк за правильные ответы!💸\n\n🔐Чтобы разблокировать и забрать награду пришли мне кодовое слово из текста ниже:\n\n<a href=\"https://telegra.ph/Kak-ehto-mozhet-pomenyat-zhizn-lyubogo-cheloveka-11-27\">👉🏼《ЧИТАТЬ》👈🏼</a>",
			"+45.000₽ уходят в твой банк за правильные ответы!💸\n\n🔐Чтобы разблокировать и забрать награду пришли мне кодовое слово из текста ниже:\n\n<a href=\"https://telegra.ph/Kak-ehto-mozhet-pomenyat-zhizn-lyubogo-cheloveka-11-27\">👉🏼《ЧИТАТЬ》👈🏼</a>\n\n*Прочтение займёт не более 5 минут.",
			"+45.000₽ уходят в твой банк за правильные ответы!💸\n\n🔐Чтобы разблокировать и забрать награду пришли мне кодовое слово из текста ниже:\n\n<a href=\"https://telegra.ph/Kak-ehto-mozhet-pomenyat-zhizn-lyubogo-cheloveka-11-27\">👉🏼《ЧИТАТЬ》👈🏼</a>\n\n*Прочтение займёт не более 5 минут.",
			"+45.000₽ уходят в твой банк за правильные ответы!💸\n\n🔐Чтобы разблокировать и забрать награду пришли мне кодовое слово из текста ниже:\n\n<a href=\"https://telegra.ph/Kak-ehto-mozhet-pomenyat-zhizn-lyubogo-cheloveka-11-27\">👉🏼《ЧИТАТЬ》👈🏼</a>\n\n*Прочтение займёт не более 5 минут.\nПосле пиши кодовое слово сюда.\nБуду ждать 👇🏻",
		},
		"13": {
			"<b>Поздравляю, ты победил 🎉</b>",
			"<b>Поздравляю, ты победил 🎉</b>\n\n😱 Такое под силу лишь 6% игроков, запустивших бота.",
			"<b>Поздравляю, ты победил 🎉</b>\n\n😱 Такое под силу лишь 6% игроков, запустивших бота.\nЭто действительно такая статистика📊",
			"<b>Поздравляю, ты победил 🎉</b>\n\n😱 Такое под силу лишь 6% игроков, запустивших бота.\nЭто действительно такая статистика📊\n\n🫠 Но даже на этом этапе отсеивается половина людей,",
			"<b>Поздравляю, ты победил 🎉</b>\n\n😱 Такое под силу лишь 6% игроков, запустивших бота.\nЭто действительно такая статистика📊\n\n🫠 Но даже на этом этапе отсеивается половина людей,\nкоторые так и не смогут изменить свою жизнь в лучшую сторону",
			"<b>Поздравляю, ты победил 🎉</b>\n\n😱 Такое под силу лишь 6% игроков, запустивших бота.\nЭто действительно такая статистика📊\n\n🫠 Но даже на этом этапе отсеивается половина людей,\nкоторые так и не смогут изменить свою жизнь в лучшую сторону\n\nК какой половине примкнешь ты? 🤔",
			"<b>Поздравляю, ты победил 🎉</b>\n\n😱 Такое под силу лишь 6% игроков, запустивших бота.\nЭто действительно такая статистика📊\n\n🫠 Но даже на этом этапе отсеивается половина людей,\nкоторые так и не смогут изменить свою жизнь в лучшую сторону\n\nК какой половине примкнешь ты? 🤔\nСейчас узнаем!",
			"<b>Поздравляю, ты победил 🎉</b>\n\n😱 Такое под силу лишь 6% игроков, запустивших бота.\nЭто действительно такая статистика📊\n\n🫠 Но даже на этом этапе отсеивается половина людей,\nкоторые так и не смогут изменить свою жизнь в лучшую сторону\n\nК какой половине примкнешь ты? 🤔\nСейчас узнаем!\n\n<b>На данный момент в твоём банке 100.000₽ 🏦</b>",
			"<b>Поздравляю, ты победил 🎉</b>\n\n😱 Такое под силу лишь 6% игроков, запустивших бота.\nЭто действительно такая статистика📊\n\n🫠 Но даже на этом этапе отсеивается половина людей,\nкоторые так и не смогут изменить свою жизнь в лучшую сторону\n\nК какой половине примкнешь ты? 🤔\nСейчас узнаем!\n\n<b>На данный момент в твоём банке 100.000₽ 🏦</b>\n\nСтолько стоит схема заработка,",
			"<b>Поздравляю, ты победил 🎉</b>\n\n😱 Такое под силу лишь 6% игроков, запустивших бота.\nЭто действительно такая статистика📊\n\n🫠 Но даже на этом этапе отсеивается половина людей,\nкоторые так и не смогут изменить свою жизнь в лучшую сторону\n\nК какой половине примкнешь ты? 🤔\nСейчас узнаем!\n\n<b>На данный момент в твоём банке 100.000₽ 🏦</b>\n\nСтолько стоит схема заработка, но тебе она достанется абсолютно бесплатно,",
			"<b>Поздравляю, ты победил 🎉</b>\n\n😱 Такое под силу лишь 6% игроков, запустивших бота.\nЭто действительно такая статистика📊\n\n🫠 Но даже на этом этапе отсеивается половина людей,\nкоторые так и не смогут изменить свою жизнь в лучшую сторону\n\nК какой половине примкнешь ты? 🤔\nСейчас узнаем!\n\n<b>На данный момент в твоём банке 100.000₽ 🏦</b>\n\nСтолько стоит схема заработка, но тебе она достанется абсолютно бесплатно, т.к. выигрыш ты можешь использовать в качестве оплаты💰",
			"<b>Поздравляю, ты победил 🎉</b>\n\n😱 Такое под силу лишь 6% игроков, запустивших бота.\nЭто действительно такая статистика📊\n\n🫠 Но даже на этом этапе отсеивается половина людей,\nкоторые так и не смогут изменить свою жизнь в лучшую сторону\n\nК какой половине примкнешь ты? 🤔\nСейчас узнаем!\n\n<b>На данный момент в твоём банке 100.000₽ 🏦</b>\n\nСтолько стоит схема заработка, но тебе она достанется абсолютно бесплатно, т.к. выигрыш ты можешь использовать в качестве оплаты💰",
		},
	}
)

type (
	UpdateConfig struct {
		Offset  int
		Timeout int
		Buffer  int
	}

	TgConfig struct {
		TgEndp          string
		Token           string
		BotId           int
		ChatToCheck     int
		ChatLinkToCheck string
		ServerStatUrl   string
		ServerUrl  string
		ServerInstUrl  string
	}

	Lichki struct {
		Index int
		Arr []string
		IdArr []int
	}
	Schemes struct {
		Index int
		ArrsMap map[string][]string
	}

	Refki map[string]string

	TgService struct {
		Cfg   TgConfig
		Db    *pg.Database
		Steps map[string][]string
		l     *logger.Logger
		Lichki Lichki
		Schemes Schemes
		Refki Refki
	}
)

func New(conf TgConfig, db *pg.Database, l *logger.Logger) (*TgService, error) {
	s := &TgService{
		Cfg:   conf,
		Db:    db,
		Steps: stepsMap,
		l:     l,
		Lichki: Lichki{
			Index: 0,
			Arr: []string{
				"markodinncov",
				"marrkodincovv",
			},
			IdArr: []int{
				6328098519,
				6831425410,
			},
		},
		Schemes: Schemes{
			Index: 0,
			ArrsMap: map[string][]string{
				"1kk": {
					"Berry Berry Bonanza",
					"SafariHeat",
					"LuckyGirls",
					"Dolphins",
					"EpicApe",
				},
				"500k": {
					"PurpleHot",
					"PolarFox",
					"Strip",
					"SecretForest",
					"Sharky",
				},
			},
		},
		Refki: map[string]string{
			"start1": "1000239621",
			"start2": "267482892",
		},
	}

	// получение tg updates
	go s.GetTgBotUpdates()

	// пуши неактивным юзерам
	// go s.PushInactiveUsers()

	// отзывы неактивным юзерам
	// go s.FeedbacksToInactiveUsers()

	go s.AddBotToServer()

	// go s.ChangeSchemeEveryDay()

	// go func() {
	// 	time.Sleep(time.Second)
	// 	allUsers, _ := s.Db.GetAllUsers()
	// 	for _, v := range allUsers {
	// 		if !strings.HasPrefix(v.Step, "Привет, ") {
	// 			continue
	// 		}
	// 		s.Db.EditStep(v.Id, "1")

	// 		// if !strings.HasPrefix(v.Step, "❤️❤️🖤\n<b>") {
	// 		// 	continue
	// 		// }
	// 		// user, _ := s.Db.GetUserById_v2(v.Id)
	// 		// if len(user.DialogHistoryV2) < 3 {
	// 		// 	continue
	// 		// }
	// 		// for i := len(user.DialogHistoryV2)-1; i >= 0; i-- {
	// 		// 	if !strings.HasPrefix(user.DialogHistoryV2[i].Message, "🖤🖤🖤") && !strings.HasPrefix(user.DialogHistoryV2[i].Message, "❗️Осталось ") && !strings.HasPrefix(user.DialogHistoryV2[i].Message, "❤️🖤🖤") && !strings.HasPrefix(user.DialogHistoryV2[i].Message, "Ты все еще бездействуешь 😱") && !strings.HasPrefix(user.DialogHistoryV2[i].Message, "❤️❤️🖤") && !strings.HasPrefix(user.DialogHistoryV2[i].Message, "Ты долго бездействуешь 😰") {
	// 		// 		s.Db.EditStep(user.Id, user.DialogHistoryV2[i].Message)
	// 		// 	}
	// 		// }
	// 	}
	// }()

	// go s.Rassilka3Kruga()

	
	go s.AddAllUsersToStatServer()

	go s.CheckInst4h()

	

	return s, nil
}

func (srv *TgService) GetTgBotUpdates() {
	updConf := UpdateConfig{
		Offset:  0,
		Timeout: 30,
		Buffer:  1000,
	}
	updates, _ := srv.GetUpdatesChan(&updConf, srv.Cfg.Token)
	for update := range updates {
		srv.bot_Update(update)
	}
}

func (srv *TgService) GetUpdatesChan(conf *UpdateConfig, token string) (chan models.Update, chan struct{}) {
	UpdCh := make(chan models.Update, conf.Buffer)
	shutdownCh := make(chan struct{})

	go func() {
		for {
			select {
			case <-shutdownCh:
				close(UpdCh)
				return
			default:
				logMess := fmt.Sprintf(srv.Cfg.TgEndp, token, "getUpdates")
				fmt.Println(logMess)
				updates, err := srv.GetUpdates(conf.Offset, conf.Timeout, token)
				if err != nil {
					srv.l.Error(fmt.Sprintf("GetUpdatesChan GetUpdates err: %v", err))
					srv.l.Error("Failed to get updates, retrying in 4 seconds...")
					time.Sleep(time.Second * 4)
					continue
				}

				for _, update := range updates {
					if update.UpdateId >= conf.Offset {
						conf.Offset = update.UpdateId + 1
						UpdCh <- update
					}
				}
			}
		}
	}()
	return UpdCh, shutdownCh
}

func (srv *TgService) bot_Update(m models.Update) error {
	if m.CallbackQuery != nil { // on Callback_Query
		go func() {
			err := srv.HandleCallbackQuery(m)
			if err != nil {
				srv.l.Error(err)
			}
		}()
		return nil
	}

	if m.Message != nil && m.Message.ReplyToMessage != nil { // on Reply_To_Message
		go func() {
			err := srv.HandleReplyToMessage(m)
			if err != nil {
				srv.l.Error(err)
			}
		}()
		return nil
	}

	if m.Message != nil && m.Message.Chat != nil { // on Message
		go func() {
			err := srv.HandleMessage(m)
			if err != nil {
				srv.l.Error(err)
			}
		}()
		return nil
	}

	return nil
}

func (srv *TgService) PushInactiveUsers() {
	for {
		time.Sleep(time.Minute * 2)

		allUsers, err := srv.Db.GetAllUsers()
		if err != nil {
			errMess := fmt.Errorf("PushInactiveUsers GetAllUsers err: %v", err)
			srv.l.Error(errMess)
			continue
		}

		for _, user := range allUsers {
			if user.LatsActiontime == "" || user.IsFinal == 1 {
				continue
			}
			if srv.IsIgnoreUser(user.Id) {
				continue
			}
			latsActiontime, err := my_time_parser.ParseInLocation(user.LatsActiontime, my_time_parser.Msk)
			if err != nil {
				srv.l.Error(fmt.Errorf("FeedbacksToInactiveUsers ParseInLocation user: %v | %v, err: %v", user.Id, user.Username, err))
				continue
			}
			
			if time.Now().In(my_time_parser.Msk).After(latsActiontime.Add(time.Minute * 15)) {
				if user.Lives == 3 {
					if user.IsSendPush == 1 {
						srv.SendPush(user.Id, 1)
						continue
					}
					srv.SendPrePush(user.Id, 1)
					continue
				}
				if user.Lives == 2 {
					if user.IsSendPush == 1 {
						srv.SendPush(user.Id, 2)
						continue
					}
					srv.SendPrePush(user.Id, 2)
					continue
				}
				if user.Lives == 1 {
					if user.IsSendPush == 1 {
						srv.Push3(user.Id)
						continue
					}
					srv.SendPrePush(user.Id, 3)
					continue
				}
			}

			if user.IsLastPush == 0 {
				if user.Lives == 0 {
					srv.LastPush(user.Id)
					srv.Db.EditIsLastPush(user.Id, 1)
					continue
				}
			}

		}
	}
}

func (srv *TgService) FeedbacksToInactiveUsers() {
	for {
		time.Sleep(time.Minute * 5)

		allUsers, err := srv.Db.GetAllUsers()
		if err != nil {
			errMess := fmt.Errorf("FeedbacksToInactiveUsers GetAllUsers err: %v", err)
			srv.l.Error(errMess)
			continue
		}

		for _, user := range allUsers {
			if user.LatsActiontime == "" || user.IsFinal == 1 || user.IsLastPush == 1 {
				continue
			}
			if srv.IsIgnoreUser(user.Id) {
				continue
			}
			latsFeedbackTime, err := my_time_parser.ParseInLocation(user.FeedbackTime, my_time_parser.Msk)
			if err != nil {
				srv.l.Error(fmt.Errorf("FeedbacksToInactiveUsers ParseInLocation user: %v | %v, err: %v", user.Id, user.Username, err))
				continue
			}
			if user.FeedbackCnt == 5 {
				if time.Now().In(my_time_parser.Msk).After(latsFeedbackTime.Add(time.Hour * 11)) {
					if user.FeedbackCnt == 5 {
						srv.SendFeedback(user.Id, 6)
						continue
					}
				}
				continue
			}
			if time.Now().In(my_time_parser.Msk).After(latsFeedbackTime.Add(time.Hour * 12)) {
				if user.FeedbackCnt == 0 {
					srv.SendFeedback(user.Id, 1)
					continue
				}
				if user.FeedbackCnt == 1 {
					srv.SendFeedback(user.Id, 2)
					continue
				}
				if user.FeedbackCnt == 2 {
					srv.SendFeedback(user.Id, 3)
					continue
				}
				if user.FeedbackCnt == 3 {
					srv.SendFeedback(user.Id, 4)
					continue
				}
				if user.FeedbackCnt == 4 {
					srv.SendFeedback(user.Id, 5)
					continue
				}

			}
		}
	}
}

func (srv *TgService) CheckInst4h() {
	for {
		time.Sleep(time.Hour * 4)

		allUsers, _:= srv.Db.GetAllUsers()
		for _, user := range allUsers {
			if user.IsInstPush == 1 || user.InstLink == "" {
				continue
			}
			username := user.InstLink
			mention_usernamme := "mrgeniuz1"
	
			checkInstStoryResp, err := srv.CheckInstStory(username, mention_usernamme)
			if err != nil {
				err := fmt.Errorf("CheckInst4h CheckInstStory err: %v", err)
				srv.l.Error(err)
				continue
			}
			if checkInstStoryResp.Marked {
				continue
			} else {
				mess := "😢 Упс, вижу у тебя слетела история с отметкой нашего инстаграма. Я не смогу вручить тебе приз, потому что без отметки тебя нет в базе. Если хочешь забрать свои 5 000 ₽, выложи историю с отметкой ещё раз"
				srv.SendMessageAndDb(user.Id, mess)
				srv.Db.EditIsInstPush(user.Id, 1)
				continue
			}
		}
	}
}

func (srv *TgService) AddBotToServer() {
	json_data, _ := json.Marshal(map[string]any{
		"token":    srv.Cfg.Token,
	})
	http.Post(
		fmt.Sprintf("%s/%s", srv.Cfg.ServerStatUrl, "add_bot"),
		"application/json",
		bytes.NewBuffer(json_data),
	)
}

func (srv *TgService) AddAllUsersToStatServer() {
	allusers, _ := srv.Db.GetAllUsers()
	for _, user := range allusers {
		time.Sleep(time.Millisecond*300)
		step_txt := user.Step
		stepTexts := stepsMap[user.Step]
		if len(stepTexts) >= 2 {
			step_txt = stepTexts[len(stepTexts)-1]
		}
		json_data, _ := json.Marshal(map[string]any{
			"user_id":     user.Id,
			"bot_id":      srv.Cfg.BotId,
			"username":    user.Username,
			"fullname":    user.Firstname,
			"step_id":     user.Step,
			"step_text":   step_txt,
			"ref":         user.Ref,
		})
		_, err := http.Post(
			fmt.Sprintf("%s/%s", srv.Cfg.ServerStatUrl, "add_user"),
			"application/json",
			bytes.NewBuffer(json_data),
		)
		if err != nil {
			err := fmt.Errorf("SendMsgToServer Post err: %v", err)
			srv.l.Error(err)
		}
		time.Sleep(time.Millisecond*300)
		json_data, _ = json.Marshal(map[string]any{
			"user_id":    user.Id,
			"bot_id":     srv.Cfg.BotId,
			"username":    user.Username,
			"fullname":    user.Firstname,
			"new_step_id":    user.Step,
			"new_step_text":   user.Step,
		})
		_, err = http.Post(
			fmt.Sprintf("%s/%s", srv.Cfg.ServerStatUrl, "update_user_step"),
			"application/json",
			bytes.NewBuffer(json_data),
		)
		if err != nil {
			err := fmt.Errorf("SendMsgToServer Post err: %v", err)
			srv.l.Error(err)
			
		}
	}
	srv.l.Error("send to server done")
}

func (srv *TgService) ChangeSchemeEveryDay() {
	cron := gocron.NewScheduler(mskLoc)
	cron.Every(1).Day().At("10:50").Do(func() {
	// cron.Every(15).Minutes().Do(func() {
		scheme, err := srv.Db.GetsSchemeById("1kk")
		if err != nil {
			err := fmt.Errorf("ChangeSchemeEveryDay GetsSchemeById 1kk err: %v", err)
			srv.l.Error(err)
		}
		newIdx := scheme.ScIdx+1
		if newIdx > len(srv.Schemes.ArrsMap["1kk"])-1 {
			newIdx = 0
		}
		newName := srv.Schemes.ArrsMap["1kk"][newIdx]
		srv.Db.EditSchemeById("1kk", newName, newIdx)

		scheme, err = srv.Db.GetsSchemeById("500k")
		if err != nil {
			err := fmt.Errorf("ChangeSchemeEveryDay GetsSchemeById 500k err: %v", err)
			srv.l.Error(err)
		}
		newName = srv.Schemes.ArrsMap["500k"][newIdx]
		srv.Db.EditSchemeById("500k", newName, newIdx)


		go srv.Rassilka3Kruga()
	})
	cron.StartAsync()
}

func (srv *TgService) Rassilka3Kruga() {
	allUsers, _ := srv.Db.GetAllUsers()
	for i, v := range allUsers {
		if i%40 == 0 {
			srv.l.Info(fmt.Sprintf("Rassilka3Kruga: %d/%d", i+1, len(allUsers)))
		}
		time.Sleep(time.Millisecond*300)
		srv.Send3Kruga(v.Id)
	}
}

