package examples

import "github.com/unix-streamdeck/streamdeckd/handlers"

func RegisterBaseModules() {
	handlers.RegisterModule(RegisterGif())
	handlers.RegisterModule(RegisterTime())
	handlers.RegisterModule(RegisterTime12())
	handlers.RegisterModule(RegisterCounter())
	handlers.RegisterModule(RegisterSpotify())
}
