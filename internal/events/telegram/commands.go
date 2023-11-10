package telegram

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
	"tgbot-read-adviser/internal/storage"
	e "tgbot-read-adviser/lib/err"
)

const (
	RndCmd    = "/rnd"
	HelpCmd   = "/help"
	StartCmd  = "/start"
	RemoveCmd = "/remove"
	AllCmd    = "/all"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	words := strings.Fields(text)
	for i := range words {
		words[i] = strings.TrimSpace(words[i])
	}

	log.Printf("got new command '%s' from '%s", text, username)

	if isURL(words[0]) {
		return p.savePage(chatID, text, username)
	}

	switch words[0] {
	case RndCmd:
		return p.sendRandom(chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	case RemoveCmd:
		return p.removePage(chatID, words[1], username)
	case AllCmd:
		return p.allPages(chatID, username)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}

	return nil
}

func (p *Processor) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: save page", err) }()

	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	isExists, err := p.storage.IsExists(context.Background(), page)
	if err != nil {
		return err
	}
	if isExists {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}

	if err := p.storage.Save(context.Background(), page); err != nil {
		return err
	}

	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendRandom(chatID int, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send random", err) }()

	page, err := p.storage.PickRandom(context.Background(), username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}

	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return err
	}

	return nil
}

func (p *Processor) removePage(chatID int, pageURL string, username string) error {
	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}
	err := p.storage.Remove(context.Background(), page)
	if err != nil && errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessage(chatID, msgNoThisPage)
	}

	return p.tg.SendMessage(chatID, msgDeleted)
}

func (p *Processor) allPages(chatID int, username string) error {
	pages, err := p.storage.All(context.Background(), username)
	if err != nil && errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessage(chatID, msgEmptyList)
	}
	i := 1
	urls := ""
	for _, page := range pages {
		urls += fmt.Sprintf("%d) %s\n\n", i, page.URL)
		i++
	}

	return p.tg.SendMessage(chatID, urls)
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
