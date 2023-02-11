package event_consumer

import (
	"MusicBot/pkg/events"
	"log"
	"sync"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c *Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("[ERR] consumer : %s", err.Error())
			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		if err := c.handleEvents(gotEvents); err != nil {
			log.Println(err)
		}
	}
}

func (c *Consumer) handleEvents(eventList []events.Event) error {
	var wg sync.WaitGroup
	wg.Add(len(eventList))
	work := func(e events.Event) {
		defer wg.Done()
		log.Printf("new event cought : %s", e.Text)

		if err := c.processor.Process(e); err != nil {
			log.Printf("[ERR] can't handle event : %s", err.Error())
		}
	}

	for _, event := range eventList {
		go work(event)
	}
	wg.Wait()

	return nil
}
