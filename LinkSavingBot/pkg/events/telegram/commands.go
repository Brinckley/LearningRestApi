package telegram

import (
	"MusicBot/pkg/clients/telegram"
	"MusicBot/pkg/storage"
	"errors"
	"log"
	"net/url"
	"strings"
	"time"
)

const (
	startCmd = "/start"
	helpCmd  = "/help"
	rndCmd   = "/rnd"
)

func (p *TgProcessor) doCmd(text string, chatId int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s'", text, username)

	if isAddCmd(text) {
		err := p.savePage(chatId, text, username)
		if err != nil {
			return err
		}
	}

	// something like router
	switch text {
	case startCmd:
		return p.sendStart(chatId, username)
	case helpCmd:
		return p.sendHelp(chatId, username)
	case rndCmd:
		return p.sendRandom(chatId, username)
	default:
		return p.tg.SendMessage(chatId, msgUnknownCommand)

	}
}

func (p *TgProcessor) sendStart(chatID int, username string) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func (p *TgProcessor) sendHelp(chatID int, username string) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *TgProcessor) sendRandom(chatId int, username string) error {
	page, err := p.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedFiles) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedFiles) {
		return p.tg.SendMessage(chatId, msgNoSavedPages)
	}

	if err := p.tg.SendMessage(chatId, page.URL); err != nil {
		return err
	}

	return p.storage.Remove(page)
}

func newMessageSender(chatID int, tg *telegram.Client) func(string) error {
	return func(msg string) error {
		return tg.SendMessage(chatID, msg)
	}
}

func (p *TgProcessor) savePage(chatID int, pageURL string, username string) error {
	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
		Created:  time.Now(),
	}

	sendMsg := newMessageSender(chatID, p.tg)

	isExists, err := p.storage.IsExists(page)
	if err != nil {
		return err
	}
	if isExists {
		return sendMsg(msgAlreadyExists)
	}

	if err := p.storage.Save(page); err != nil {
		return err
	}

	if err := sendMsg(msgSaved); err != nil {
		return err
	}

	return nil
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
