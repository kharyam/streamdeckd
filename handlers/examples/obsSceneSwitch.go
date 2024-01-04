package examples

import (
	"fmt"
	"os/exec"
	"syscall"

	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/scenes"
	"github.com/unix-streamdeck/api"
	"github.com/unix-streamdeck/streamdeckd/handlers"
)

type ObsSceneSwitchKeyHandler struct {
	ObsClient *goobs.Client
}

func (t *ObsSceneSwitchKeyHandler) Key(key api.Key, info api.StreamDeckInfo) {
	if key.KeyHandler != "ObsSceneSwitch" {
		fmt.Println(key.KeyHandler)
		return
	}
	var err error
	t.ObsClient, err = goobs.New(key.KeyHandlerFields["obs_host"]+":"+key.KeyHandlerFields["obs_port"], goobs.WithPassword(key.KeyHandlerFields["obs_password"]))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer t.ObsClient.Disconnect()

	t.ObsClient.Scenes.SetCurrentProgramScene(scenes.NewSetCurrentProgramSceneParams().WithSceneName(key.KeyHandlerFields["scene_name"]))
	if key.KeyHandlerFields["scene_command"] != "" {
		runCommand(key.KeyHandlerFields["scene_command"])
	}

}

func RegisterObsSceneSwitch() handlers.Module {
	return handlers.Module{NewKey: func() api.KeyHandler {
		return &ObsSceneSwitchKeyHandler{}
	}, Name: "ObsSceneSwitch",
		KeyFields: []api.Field{{Title: "OBS Host", Name: "obs_host", Type: "Text"}, {Title: "OBS Port", Name: "obs_port", Type: "Number"}, {Title: "OBS Password", Name: "obs_password", Type: "Text"}, {Title: "Scene Name", Name: "scene_name", Type: "Text"}, {Title: "Command", Name: "scene_command", Type: "Text"}}}
}

// TODO This was copied from main.  Put it in a common package and import it into main and this package.
func runCommand(command string) {
	go func() {
		cmd := exec.Command("/bin/sh", "-c", "/usr/bin/nohup "+command)

		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setpgid:   true,
			Pgid:      0,
			Pdeathsig: syscall.SIGHUP,
		}
		if err := cmd.Start(); err != nil {
			fmt.Println("There was a problem running ", command, ":", err)
		} else {
			pid := cmd.Process.Pid
			cmd.Process.Release()
			fmt.Println(command, " has been started with pid", pid)
		}
	}()
}
