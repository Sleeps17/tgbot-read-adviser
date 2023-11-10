package event_consumer

import (
	"log"
	"sync"
	"time"

	"tgbot-read-adviser/internal/events"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) *Consumer {
	return &Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c *Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())

			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		if err := c.handleEvents(gotEvents); err != nil {
			log.Print(err)

			continue
		}
	}
}

func (c *Consumer) handleEvents(events []events.Event) error {
	var wg sync.WaitGroup
	for i := range events {
		wg.Add(1)
		log.Printf("got new event: %s", events[i].Text)
		go func(i int) {
			defer wg.Done()
			if err := c.processor.Process(events[i]); err != nil {

				log.Printf("can't handle event: %s", err.Error())

				return
			}
		}(i)
	}
	wg.Wait()
	return nil
}
