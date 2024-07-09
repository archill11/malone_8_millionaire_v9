package tg_service

import (
	"bytes"
	"fmt"
)

func (srv *TgService) GetActualSchema(fromId int, siteUrl string) (string) {

	user, _ := srv.Db.GetUserById(fromId)
	lichka := user.Lichka
	if lichka == "" {
		lichka = "https://t.me/markodinncov"
	}
	lichkaUrl := fmt.Sprintf("https://t.me/%s", srv.DelAt(lichka))
	scheme, _ := srv.Db.GetsSchemeByLichka(lichka)

	algo := srv.GetActualSchemaAlgo(scheme.ScName)

	var mess bytes.Buffer
	mess.WriteString(fmt.Sprintf("❗️АЛГОРИТМ❗️\n\n"))
	mess.WriteString(fmt.Sprintf("Переходим по ссылке - %s\n(Регистрируемся на сайте)\n\n", srv.ChInfoToLinkHTML(siteUrl, "ССЫЛКА")))
	mess.WriteString(fmt.Sprintf("%s\n\n", algo))
	mess.WriteString(fmt.Sprintf("👉 %s 👈\n\n", srv.ChInfoToLinkHTML(siteUrl, "КРУТИМ ТУТ")))
	mess.WriteString(fmt.Sprintf("🔖 %s 🔖\n\n", srv.ChInfoToLinkHTML(lichkaUrl, "Написать мне")))
	mesgText := mess.String()

	return mesgText
}

func (srv *TgService) GetActualSchemaAlgo(sc_name string) (string) {
	var mess bytes.Buffer

	if sc_name == "Berry Berry Bonanza" {
		mess.WriteString(fmt.Sprintf("❗️Схему можно проверить в демо-режиме\n\n"))
		mess.WriteString(fmt.Sprintf("Если все сделаете правильно - <b>Куш составит 1.000.000 рублей !</b> 💰\n\n"))
		mess.WriteString(fmt.Sprintf("- Сумма пополнения: 1535₽ 🔥\n\n"))
		mess.WriteString(fmt.Sprintf("- Проходимость схемы 100%%✅\n- Схема для ТЕЛЕФОНА и ПК\n"))
		mess.WriteString(fmt.Sprintf("- Игра «Berry Berry Bonanza»\n\n⚙️<b>Алгоритм ставок:</b>\n\n"))
		mess.WriteString(fmt.Sprintf("Переходим на главную страницу 👉🏻 затем нажимаем на «Лупу» 👉🏻 и через Поиск находим игру:\n<b>«Berry Berry Bonanza»</b>\n\n"))
		mess.WriteString(fmt.Sprintf("<b>🕹 LINES BET — 5 | BET/LINE — 150 | Крутим 3 раза\n\n🕹 LINES BET — 8 | BET/LINE — 240 | Крутим 2 раза\n\n🕹 LINES BET — 9 | BET/LINE — 135 | Крутим 3 раза\n\n🕹 LINES BET — 1 | BET/LINE — 100 | Крутим 2 раза</b>\n\n"))
		mess.WriteString(fmt.Sprintf("<b>👉🏻Крутим схему 1 раз! С одного пополнения хватает на весь круг 👌🏻</b>\n\n"))
	}
	if sc_name == "SafariHeat" {
		mess.WriteString(fmt.Sprintf("❗️Схему можно проверить в демо-режиме\n\n"))
		mess.WriteString(fmt.Sprintf("Если все сделаете правильно - <b>Куш составит 912.500 рублей ! </b>💰\n\n"))
		mess.WriteString(fmt.Sprintf("- Сумма пополнения: 1504₽ 🔥\n\n"))
		mess.WriteString(fmt.Sprintf("- Проходимость схемы 100%%✅\n- Схема для ТЕЛЕФОНА и ПК\n"))
		mess.WriteString(fmt.Sprintf("- Игра «SafariHeat»\n\n⚙️<b>Алгоритм ставок:</b>\n\n"))
		mess.WriteString(fmt.Sprintf("Переходим на главную страницу 👉🏻 затем нажимаем на «Лупу» 👉🏻 и через Поиск находим игру:\n<b>«SafariHeat»</b>\n\n"))
		mess.WriteString(fmt.Sprintf("<b>🕹 LINES BET — 5 | BET/LINE — 50 | Крутим 1 раз\n\n🕹 LINES BET — 7 | BET/LINE — 10 | Крутим 5 раз\n\n🕹 LINES BET — 1 | BET/LINE — 70 | Крутим 4 раза\n\n🕹 LINES BET — 9 | BET/LINE — 6 | Крутим 6 раз\n\n🕹 LINES BET — 3 | BET/LINE — 100 | Крутим 1 раз</b>\n\n"))
		mess.WriteString(fmt.Sprintf("<b>👉🏻Крутим схему 1 раз! С одного пополнения хватает на весь круг 👌🏻</b>\n\n"))
	}
	if sc_name == "LuckyGirls" {
		mess.WriteString(fmt.Sprintf("❗️Схему можно проверить в демо-режиме\n\n"))
		mess.WriteString(fmt.Sprintf("Если все сделаете правильно - <b>Куш составит 900.000 рублей ! </b>💰\n\n"))
		mess.WriteString(fmt.Sprintf("- Сумма пополнения: 1540₽ 🔥\n\n"))
		mess.WriteString(fmt.Sprintf("- Проходимость схемы 100%%✅\n- Схема для ТЕЛЕФОНА и ПК\n"))
		mess.WriteString(fmt.Sprintf("- Игра «LuckyGirls»\n\n⚙️<b>Алгоритм ставок:</b>\n\n"))
		mess.WriteString(fmt.Sprintf("Переходим на главную страницу 👉🏻 затем нажимаем на «Лупу» 👉🏻 и через Поиск находим игру:\n<b>«LuckyGirls»</b>\n\n"))
		mess.WriteString(fmt.Sprintf("<b>🕹 LINES BET — 9 | BET/LINE — 40 | Крутим 2 раза\n\n🕹 LINES BET — 7 | BET/LINE — 10 | Крутим 3 раза\n\n🕹 LINES BET — 3 | BET/LINE — 20 | Крутим 4 раза\n\n🕹 LINES BET — 9 | BET/LINE — 15 | Крутим 2 раза\n\n🕹 LINES BET — 1 | BET/LINE — 100 | Крутим 1 раз</b>\n\n"))
		mess.WriteString(fmt.Sprintf("<b>👉🏻Крутим схему 1 раз! С одного пополнения хватает на весь круг 👌🏻</b>\n\n"))
	}
	if sc_name == "Dolphins" {
		mess.WriteString(fmt.Sprintf("❗️Схему можно проверить в демо-режиме\n\n"))
		mess.WriteString(fmt.Sprintf("Если все сделаете правильно - <b>Куш составит 900.000 рублей ! </b>💰\n\n"))
		mess.WriteString(fmt.Sprintf("- Сумма пополнения: 1540₽ 🔥\n\n"))
		mess.WriteString(fmt.Sprintf("- Проходимость схемы 100%%✅\n- Схема для ТЕЛЕФОНА и ПК\n"))
		mess.WriteString(fmt.Sprintf("- Игра «Dolphins»\n\n⚙️<b>Алгоритм ставок:</b>\n\n"))
		mess.WriteString(fmt.Sprintf("Переходим на главную страницу 👉🏻 затем нажимаем на «Лупу» 👉🏻 и через Поиск находим игру:\n<b>«Dolphins»</b>\n\n"))
		mess.WriteString(fmt.Sprintf("<b>🕹 LINES BET — 3 | BET/LINE — 40 | Крутим 2 раза\n\n🕹 LINES BET — 9 | BET/LINE — 5 | Крутим 5 раз\n\n🕹 LINES BET — 3 | BET/LINE — 50 | Крутим 5 раз\n\n🕹 LINES BET — 5 | BET/LINE — 15 | Крутим 3 раза\n\n🕹 LINES BET — 1 | BET/LINE — 100 | Крутим 1 раз</b>\n\n"))
		mess.WriteString(fmt.Sprintf("<b>👉🏻Крутим схему 1 раз! С одного пополнения хватает на весь круг 👌🏻</b>\n\n"))
	}
	if sc_name == "EpicApe" {
		mess.WriteString(fmt.Sprintf("❗️Схему можно проверить в демо-режиме\n\n"))
		mess.WriteString(fmt.Sprintf("Если все сделаете правильно - <b>Куш составит 864.000 рублей ! </b>💰\n\n"))
		mess.WriteString(fmt.Sprintf("- Сумма пополнения: 1520₽ 🔥\n\n"))
		mess.WriteString(fmt.Sprintf("- Проходимость схемы 100%%✅\n- Схема для ТЕЛЕФОНА и ПК\n"))
		mess.WriteString(fmt.Sprintf("- Игра «EpicApe»\n\n⚙️<b>Алгоритм ставок:</b>\n\n"))
		mess.WriteString(fmt.Sprintf("Переходим на главную страницу 👉🏻 затем нажимаем на «Лупу» 👉🏻 и через Поиск находим игру:\n<b>«EpicApe»</b>\n\n"))
		mess.WriteString(fmt.Sprintf("<b>🕹 COIN VALUE - 2 | BET- 80 | Крутим 4 раза\n\n🕹 COIN VALUE - 3 | BET- 120 | Крутим 3 раза\n\n🕹 COIN VALUE - 7 | BET- 280 | Крутим 1 раза\n\n🕹 COIN VALUE -1 | BET- 40 | Крутим 4 раза\n\n🕹 COIN VALUE - 10 | BET- 400 | Крутим 1 раз</b>\n\n"))
		mess.WriteString(fmt.Sprintf("<b>👉🏻Крутим схему 1 раз! С одного пополнения хватает на весь круг 👌🏻</b>\n\n"))
	}
	if sc_name == "PurpleHot" {
		mess.WriteString(fmt.Sprintf("❗️Схему можно проверить в демо-режиме\n\n"))
		mess.WriteString(fmt.Sprintf("Если все сделаете правильно - <b>Куш составит 500.000 рублей ! </b>💰\n\n"))
		mess.WriteString(fmt.Sprintf("- Сумма пополнения: 1530₽ 🔥\n\n"))
		mess.WriteString(fmt.Sprintf("- Проходимость схемы 100%%✅\n- Схема для ТЕЛЕФОНА и ПК\n"))
		mess.WriteString(fmt.Sprintf("- Игра «PurpleHot»\n\n⚙️<b>Алгоритм ставок:</b>\n\n"))
		mess.WriteString(fmt.Sprintf("Переходим на главную страницу 👉🏻 затем нажимаем на «Лупу» 👉🏻 и через Поиск находим игру:\n<b>«PurpleHot»</b>\n\n"))
		mess.WriteString(fmt.Sprintf("<b>🕹 LINES BET -  5  | BET/LINE - 20 | Крутим 4 раза\n\n🕹 LINES BET -  3  | BET/LINE - 30 | Крутим 2 раза\n\n🕹 LINES BET -  5 | BET/LINE - 30 | Крутим 5 раза\n\n🕹 LINES BET -  1  | BET/LINE - 100 | Крутим 2 раза</b>\n\n"))
		mess.WriteString(fmt.Sprintf("<b>👉🏻Крутим схему 1 раз! С одного пополнения хватает на весь круг 👌🏻</b>\n\n"))
	}
	if sc_name == "PolarFox" {
		mess.WriteString(fmt.Sprintf("❗️Схему можно проверить в демо-режиме\n\n"))
		mess.WriteString(fmt.Sprintf("Если все сделаете правильно - <b>Куш составит 541.500 рублей ! </b>💰\n\n"))
		mess.WriteString(fmt.Sprintf("- Сумма пополнения: 1535₽ 🔥\n\n"))
		mess.WriteString(fmt.Sprintf("- Проходимость схемы 100%%✅\n- Схема для ТЕЛЕФОНА и ПК\n"))
		mess.WriteString(fmt.Sprintf("- Игра «PolarFox»\n\n⚙️<b>Алгоритм ставок:</b>\n\n"))
		mess.WriteString(fmt.Sprintf("Переходим на главную страницу 👉🏻 затем нажимаем на «Лупу» 👉🏻 и через Поиск находим игру:\n<b>«PolarFox»</b>\n\n"))
		mess.WriteString(fmt.Sprintf("<b>🕹 LINES BET -  7 | BET/LINE - 10 | Крутим 1 раза\n\n🕹 LINES BET -  1  | BET/LINE - 5 | Крутим 5 раза\n\n🕹 LINES BET -  3 | BET/LINE - 50 | Крутим 3 раза\n\n🕹 LINES BET -  9  | BET/LINE - 30| Крутим 3 раза\n\n🕹 LINES BET -  3  | BET/LINE - 60| Крутим 1 раз</b>\n\n"))
		mess.WriteString(fmt.Sprintf("<b>👉🏻Крутим схему 1 раз! С одного пополнения хватает на весь круг 👌🏻</b>\n\n"))
	}
	if sc_name == "Strip" {
		mess.WriteString(fmt.Sprintf("❗️Схему можно проверить в демо-режиме\n\n"))
		mess.WriteString(fmt.Sprintf("Если все сделаете правильно - <b>Куш составит 510.000 рублей ! </b>💰\n\n"))
		mess.WriteString(fmt.Sprintf("- Сумма пополнения: 1810₽ 🔥\n\n"))
		mess.WriteString(fmt.Sprintf("- Проходимость схемы 100%%✅\n- Схема для ТЕЛЕФОНА и ПК\n"))
		mess.WriteString(fmt.Sprintf("- Игра «Strip»\n\n⚙️<b>Алгоритм ставок:</b>\n\n"))
		mess.WriteString(fmt.Sprintf("Переходим на главную страницу 👉🏻 затем нажимаем на «Лупу» 👉🏻 и через Поиск находим игру:\n<b>«Strip»</b>\n\n"))
		mess.WriteString(fmt.Sprintf("<b>🕹 LINES BET -  3 | BET/LINE - 60 | Крутим 2 раза\n\n🕹 LINES BET -  5  | BET/LINE - 30 | Крутим 3 раза\n\n🕹 LINES BET -  1 | BET/LINE - 80 | Крутим 5 раза\n\n🕹 LINES BET -  5  | BET/LINE - 20 | Крутим 3 раза\n\n🕹 LINES BET -  3  | BET/LINE - 100 | Крутим 1 раз</b>\n\n"))
		mess.WriteString(fmt.Sprintf("<b>👉🏻Крутим схему 1 раз! С одного пополнения хватает на весь круг 👌🏻</b>\n\n"))
	}
	if sc_name == "SecretForest" {
		mess.WriteString(fmt.Sprintf("❗️Схему можно проверить в демо-режиме\n\n"))
		mess.WriteString(fmt.Sprintf("Если все сделаете правильно - <b>Куш составит 549.240 рублей ! </b>💰\n\n"))
		mess.WriteString(fmt.Sprintf("- Сумма пополнения: 1517₽ 🔥\n\n"))
		mess.WriteString(fmt.Sprintf("- Проходимость схемы 100%%✅\n- Схема для ТЕЛЕФОНА и ПК\n"))
		mess.WriteString(fmt.Sprintf("- Игра «SecretForest»\n\n⚙️<b>Алгоритм ставок:</b>\n\n"))
		mess.WriteString(fmt.Sprintf("Переходим на главную страницу 👉🏻 затем нажимаем на «Лупу» 👉🏻 и через Поиск находим игру:\n<b>«SecretForest»</b>\n\n"))
		mess.WriteString(fmt.Sprintf("<b>🕹 LINES BET -  5 | BET/LINE - 40 | Крутим 2 раза\n\n🕹 LINES BET -  9  | BET/LINE - 5 | Крутим 5 раза\n\n🕹 LINES BET -  1 | BET/LINE - 50 | Крутим 5 раза\n\n🕹 LINES BET -  7  | BET/LINE - 2 | Крутим 3 раза\n\n🕹 LINES BET -  5  | BET/LINE - 60 | Крутим 2 раз</b>\n\n"))
		mess.WriteString(fmt.Sprintf("<b>👉🏻Крутим схему 1 раз! С одного пополнения хватает на весь круг 👌🏻</b>\n\n"))
	}
	if sc_name == "Sharky" {
		mess.WriteString(fmt.Sprintf("❗️Схему можно проверить в демо-режиме\n\n"))
		mess.WriteString(fmt.Sprintf("Если все сделаете правильно - <b>Куш составит 500.000 рублей ! </b>💰\n\n"))
		mess.WriteString(fmt.Sprintf("- Сумма пополнения: 1750₽ 🔥\n\n"))
		mess.WriteString(fmt.Sprintf("- Проходимость схемы 100%%✅\n- Схема для ТЕЛЕФОНА и ПК\n"))
		mess.WriteString(fmt.Sprintf("- Игра «Sharky»\n\n⚙️<b>Алгоритм ставок:</b>\n\n"))
		mess.WriteString(fmt.Sprintf("Переходим на главную страницу 👉🏻 затем нажимаем на «Лупу» 👉🏻 и через Поиск находим игру:\n<b>«Sharky»</b>\n\n"))
		mess.WriteString(fmt.Sprintf("<b>🕹 LINES BET — 3 | BET/LINE — 50 | Крутим 3 раза\n\n🕹 LINES BET — 5 | BET/LINE — 15 | Крутим 4 раза\n\n🕹 LINES BET — 3 | BET/LINE — 6 | Крутим 5 раза\n\n🕹 LINES BET — 9 | BET/LINE — 30 | Крутим 3 раза\n\n🕹 LINES BET — 1 | BET/LINE — 100 | Крутим 1 раз</b>\n\n"))
		mess.WriteString(fmt.Sprintf("<b>👉🏻Крутим схему 1 раз! С одного пополнения хватает на весь круг 👌🏻</b>\n\n"))
	}
	mesgText := mess.String()

	return mesgText
}