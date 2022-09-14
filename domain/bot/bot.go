package bot

type Bot interface {
	// Add handlers to bot
	AddHandler(...Handler)

	// Start handling message with handlers
	Start() error
}

type Handler func(msg string) (string, error)
