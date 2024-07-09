package tg_service

import (
	"fmt"
	"strings"
)

func (srv *TgService) GetLichka() (string, int, error) {
	index := srv.Lichki.Index
	var newIndex int
	if index >= len(srv.Lichki.Arr)-1 {
		newIndex = 0
	} else {
		newIndex = index + 1
	}
	srv.Lichki.Index = newIndex
	lichka := srv.Lichki.Arr[index]
	id := srv.Lichki.IdArr[index]
	return lichka, id, nil
}

func (srv *TgService) ChInfoToLinkHTML(link, title string) string {
	if strings.HasPrefix(link, "@") {
		link = fmt.Sprintf("https://t.me/%s", link)
	}
	return fmt.Sprintf("<a href=\"%s\">%s</a>", link, title)
}

func (srv *TgService) DelAt(username string) string {
	usernameRunes := []rune(username)
	if len(usernameRunes) == 0 {
		return ""
	}
	if usernameRunes[0] == '@' {
		usernameRunes = usernameRunes[1:]
	}
	return string(usernameRunes)
}

func (srv *TgService) AddAt(username string) string {
	usernameRunes := []rune(username)
	if len(usernameRunes) == 0 {
		return string('@')
	}
	if usernameRunes[0] != '@' {
		usernameRunes = append([]rune{'@'}, usernameRunes...)
	}
	return string(usernameRunes)
}
