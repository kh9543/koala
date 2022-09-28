package bot

type Bot interface {
	// Add handlers to bot
	AddHandlerFuncs(usePrefix bool, hs ...Handler)

	// Start handling message with handlers
	Start() error

	// Sends direct message
	Send(channelID, message string) error
}

type Handler func(msg string) (string, error)
