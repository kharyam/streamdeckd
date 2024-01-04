package examples

import (
	"image"
	"image/draw"
	"log"
	"time"

	"github.com/unix-streamdeck/api"
	"github.com/unix-streamdeck/streamdeckd/handlers"
)

type Time12IconHandler struct {
	Running bool
	Quit    chan bool
}

func (t *Time12IconHandler) Start(k api.Key, info api.StreamDeckInfo, callback func(image image.Image)) {
	t.Running = true
	if t.Quit == nil {
		t.Quit = make(chan bool)
	}
	go t.timeLoop(k, info, callback)
}

func (t *Time12IconHandler) IsRunning() bool {
	return t.Running
}

func (t *Time12IconHandler) SetRunning(running bool) {
	t.Running = running
}

func (t *Time12IconHandler) Stop() {
	t.Running = false
	t.Quit <- true
}

func (t *Time12IconHandler) timeLoop(k api.Key, info api.StreamDeckInfo, callback func(image image.Image)) {
	for {
		select {
		case <-t.Quit:
			return
		default:
			img := image.NewRGBA(image.Rect(0, 0, info.IconSize, info.IconSize))
			draw.Draw(img, img.Bounds(), image.Black, image.Point{}, draw.Src)
			t := time.Now()
			tString := t.Format("3:04:05 pm")
			imgParsed, err := api.DrawText(img, tString, k.TextSize, k.TextAlignment)
			if err != nil {
				log.Println(err)
			} else {
				callback(imgParsed)
			}
			time.Sleep(time.Second)
		}
	}
}

func RegisterTime12() handlers.Module {
	return handlers.Module{NewIcon: func() api.IconHandler {
		return &Time12IconHandler{Running: true}
	}, Name: "Time12"}
}
