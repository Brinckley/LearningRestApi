package events

type Fetcher interface { // functional interface to get events
	Fetch(limit int) ([]Event, error)
}

type Processor interface {
	Process(e Event) error
}

type Type int

const ( // classification for messages
	Unknown Type = iota
	Message
)

type Event struct {
	Type Type
	Text string
	Meta interface{} // we can put here whatever we want. It depends on messenger options
}
