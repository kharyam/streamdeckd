package examples

import (
	"fmt"
	"image"
	"image/draw"
	"log"
	"time"

	"github.com/andreykaipov/goobs"
	"github.com/unix-streamdeck/api"
	"github.com/unix-streamdeck/streamdeckd/handlers"
)

type ObsSceneIconHandler struct {
	Running   bool
	Quit      chan bool
	ObsClient *goobs.Client
}

func (t *ObsSceneIconHandler) Start(k api.Key, info api.StreamDeckInfo, callback func(image image.Image)) {
	t.Running = true
	if t.Quit == nil {
		t.Quit = make(chan bool)
	}

	var err error
	t.ObsClient, err = goobs.New(k.IconHandlerFields["obs_host"]+":"+k.IconHandlerFields["obs_port"], goobs.WithPassword(k.IconHandlerFields["obs_password"]))

	if err != nil {
		log.Println(err)
		t.ObsClient = nil
	}

	go t.obsQueryLoop(k, info, callback)
}

func (t *ObsSceneIconHandler) IsRunning() bool {
	return t.Running
}

func (t *ObsSceneIconHandler) SetRunning(running bool) {
	t.Running = running
}

func (t *ObsSceneIconHandler) Stop() {
	t.Running = false
	t.Quit <- true
}

func (t *ObsSceneIconHandler) obsQueryLoop(k api.Key, info api.StreamDeckInfo, callback func(image image.Image)) {
	for {
		select {
		case <-t.Quit:
			defer t.ObsClient.Disconnect()
			return
		default:
			time.Sleep(time.Second)
			if t.ObsClient == nil {
				var err error
				t.ObsClient, err = goobs.New(k.IconHandlerFields["obs_host"]+":"+k.IconHandlerFields["obs_port"], goobs.WithPassword(k.IconHandlerFields["obs_password"]))

				if err != nil {
					log.Println(err)
					continue
				}
			}

			img := image.NewRGBA(image.Rect(0, 0, info.IconSize, info.IconSize))
			draw.Draw(img, img.Bounds(), image.Black, image.Point{}, draw.Src)

			defer func() {
				if r := recover(); r != nil {
					fmt.Println(r.(error))
					t.ObsClient = nil
					go t.obsQueryLoop(k, info, callback)
				}
			}()

			resp, err := t.ObsClient.Scenes.GetCurrentProgramScene()
			if err != nil {
				log.Println(err)
				continue
			}

			tString := "Scene:\n" + resp.CurrentProgramSceneName
			imgParsed, err := api.DrawText(img, tString, k.TextSize, k.TextAlignment)
			if err != nil {
				log.Println(err)
			} else {
				callback(imgParsed)
			}
		}
	}
}

func RegisterObsScene() handlers.Module {
	return handlers.Module{NewIcon: func() api.IconHandler {
		return &ObsSceneIconHandler{Running: true}
	}, Name: "OBS Scene", IconFields: []api.Field{{Title: "OBS Host", Name: "obs_host", Type: "Text"}, {Title: "OBS Port", Name: "obs_port", Type: "Number"}, {Title: "OBS Password", Name: "obs_password", Type: "Password"}}}
}
