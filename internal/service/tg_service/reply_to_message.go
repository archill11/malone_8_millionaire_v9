package tg_service

import (
	"bytes"
	"fmt"
	"myapp/internal/entity"
	"myapp/internal/models"
	my_regex "myapp/pkg/regex"
	"strconv"
	"strings"
)

func (srv *TgService) HandleReplyToMessage(m models.Update) error {
	rm := m.Message.ReplyToMessage
	replyMes := m.Message.Text
	fromId := m.Message.From.Id
	fromUsername := m.Message.From.UserName
	srv.l.Info(fmt.Sprintf("HandleCallbackQuery: fromId: %d, fromUsername: %s, replyMes: %s, rm.Tex: %s", fromId, fromUsername, replyMes, rm.Text))

	if rm.Text == MAILING_COPY_STEP {
		err := srv.RM__MAILING_COPY_STEP(m)
		if err != nil {
			srv.SendMessage(fromId, ERR_MSG)
			srv.SendMessage(fromId, err.Error())
		}
		return err
	}

	if strings.HasPrefix(rm.Text, "Укажите теперь отправьте сообщение кторое разослать для шага[") { // <ДОБАВЛЕНИЕ КАНАЛОВ В ЗАКУПКУ>
		step := my_regex.GetStringInBetween(rm.Text, "Укажите теперь отправьте сообщение кторое разослать для шага[", "]")
		err := srv.RM__MAILING_COPY(m, step)
		if err != nil {
			srv.SendMessage(fromId, ERR_MSG)
			srv.SendMessage(fromId, err.Error())
		}
		return err
	}

	if rm.Text == DEL_USER_MSG {
		err := srv.RM__DEL_USER_MSG(m)
		if err != nil {
			srv.SendMessage(fromId, ERR_MSG)
			srv.SendMessage(fromId, err.Error())
		}
		return err
	}

	if rm.Text == DEL_USER_ID_MSG {
		err := srv.RM__DEL_USER_ID_MSG(m)
		if err != nil {
			srv.SendMessage(fromId, ERR_MSG)
			srv.SendMessage(fromId, err.Error())
		}
		return err
	}

	if rm.Text == USER_INFO_MSG {
		err := srv.RM__USER_INFO_MSG(m)
		if err != nil {
			srv.SendMessage(fromId, ERR_MSG)
			srv.SendMessage(fromId, err.Error())
		}
		return err
	}

	return nil
}

func (srv *TgService) RM__MAILING_COPY_STEP(m models.Update) error {
	replyMes := m.Message.Text
	fromId := m.Message.From.Id
	fromUsername := m.Message.From.UserName
	srv.l.Info(fmt.Sprintf("RM_add_admin: fromId-%d fromUsername-%s, replyMes-%s", fromId, fromUsername, replyMes))

	srv.SendForceReply(fromId, fmt.Sprintf("Укажите теперь отправьте сообщение кторое разослать для шага[%s]", replyMes))
	return nil
}

func (srv *TgService) RM__MAILING_COPY(m models.Update, step string) error {
	replyMes := m.Message.Text
	fromId := m.Message.From.Id
	fromUsername := m.Message.From.UserName
	srv.l.Info(fmt.Sprintf("RM_add_admin: fromId-%d fromUsername-%s, replyMes-%s", fromId, fromUsername, replyMes))

	users, err := srv.Db.GetUsersByStep(step)
	if err != nil {
		return fmt.Errorf("RM__MAILING_COPY GetUsersByStep err: %v", err)
	}
	for _, v := range users {
		srv.CopyMessage(v.Id, fromId, m.Message.MessageId)
	}

	srv.SendMessage(fromId, fmt.Sprintf("Рассылка закончена, всего: %d", len(users)))

	return nil
}

func (srv *TgService) RM__DEL_USER_MSG(m models.Update) error {
	replyMes := m.Message.Text
	fromId := m.Message.From.Id
	fromUsername := m.Message.From.UserName
	srv.l.Info(fmt.Sprintf("RM__DEL_ADMIN_MSG: fromId: %d, fromUsername: %s, replyMes: %s", fromId, fromUsername, replyMes))

	username := srv.DelAt(replyMes)

	srv.Db.DeleteUserByUsername(username)

	srv.SendMessage(fromId, "юзер удален")
	return nil
}

func (srv *TgService) RM__DEL_USER_ID_MSG(m models.Update) error {
	replyMes := m.Message.Text
	fromId := m.Message.From.Id
	fromUsername := m.Message.From.UserName
	srv.l.Info(fmt.Sprintf("RM__DEL_USER_ID_MSG: fromId: %d, fromUsername: %s, replyMes: %s", fromId, fromUsername, replyMes))

	id, err := strconv.Atoi(replyMes)
	if err != nil {
		return fmt.Errorf("RM__DEL_USER_ID_MSG: Atoi replyMes: %v, err: %v", replyMes, err)
	}

	srv.Db.DeleteUserById(id)

	srv.SendMessage(fromId, "юзер удален")
	return nil
}

func (srv *TgService) RM__USER_INFO_MSG(m models.Update) error {
	replyMes := m.Message.Text
	fromId := m.Message.From.Id
	fromUsername := m.Message.From.UserName
	srv.l.Info(fmt.Sprintf("RM__USER_INFO_MSG: fromId: %d, fromUsername: %s, replyMes: %s", fromId, fromUsername, replyMes))

	userUsername := replyMes
	userId, _ := strconv.Atoi(replyMes)

	var user entity.User
	if userId != 0 {
		user, _ = srv.Db.GetUserById(userId)
	} else {
		user, _ = srv.Db.GetUserByUsername(userUsername)
	}
	usersByRef, _ := srv.Db.GetUsersByRef(strconv.Itoa(user.Id))

	var mess bytes.Buffer
	mess.WriteString(fmt.Sprintf("Юзер: %d | %s", user.Id, srv.AddAt(user.Username)))
	mess.WriteString(fmt.Sprintf("Рефералы: %d шт.", len(usersByRef)))

	srv.SendMessage(fromId, mess.String())
	return nil
}
