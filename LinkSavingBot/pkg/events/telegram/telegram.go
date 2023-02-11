package telegram

import (
	"MusicBot/pkg/clients/telegram"
	"MusicBot/pkg/events"
	"MusicBot/pkg/storage"
	"errors"
	"fmt"
)

var (
	NoUpdatesFound      = errors.New("no updates")
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

type TgProcessor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct { // special field for tg
	ChatID   int
	Username string
}

func NewTgProcessor(client *telegram.Client, storage storage.Storage) *TgProcessor {
	return &TgProcessor{
		tg:      client,
		storage: storage,
	}
}

func (p *TgProcessor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, fmt.Errorf("can't fetch upds : %s", err.Error())
	}

	if len(updates) == 0 {
		return nil, NoUpdatesFound
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].UpdateID + 1

	return res, nil
}

func (p *TgProcessor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return fmt.Errorf("can't process msg : %s", ErrUnknownEventType)
	}
}

func (p *TgProcessor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return fmt.Errorf("can't process msg : %s", err.Error())
	}

	if err := p.doCmd(event.Text, meta.ChatID, meta.Username); err != nil {
		return fmt.Errorf("can't process msg : %s", err.Error())
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta) // checking if there is any meta-data
	if !ok {
		return Meta{}, fmt.Errorf("can't get meta : %S", ErrUnknownMetaType)
	}

	return res, nil
}

func event(update telegram.Update) events.Event { // creating new event, based on update info
	updType := fetchType(update)

	res := events.Event{
		Type: updType,
		Text: fetchText(update),
	}

	if updType == events.Message { // extracting info that is unique for tg
		res.Meta = Meta{
			ChatID:   update.Message.Chat.ID,
			Username: update.Message.From.Username,
		}
	}

	return res
}

func fetchType(update telegram.Update) events.Type {
	if update.Message == nil {
		return events.Unknown
	}
	return events.Message
}

func fetchText(update telegram.Update) string {
	if update.Message == nil {
		return ""
	}
	return update.Message.Text
}
