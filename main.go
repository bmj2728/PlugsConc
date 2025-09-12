package main

import (
	"log/slog"
	"os"
	"os/exec"

	"PlugsConc/internal/logger"
	"PlugsConc/pkg/exten"

	"github.com/hashicorp/go-plugin"
)

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "ANIMAL_PLUGIN",
	MagicCookieValue: "hello",
}

var pluginMap = map[string]plugin.Plugin{
	"animal": &exten.AnimalPlugin{},
}

func main() {

	logHandler := logger.New(os.Stdout,
		&logger.Options{
			Level:     slog.LevelDebug,
			AddSource: true,
			ColorMap:  logger.DefaultColorMap,
			FullLine:  false},
	)

	slog.SetDefault(slog.New(logHandler))
	slog.Info("Logger initialized")

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		Cmd:             exec.Command("./plugins/dog/dog"),
	})
	defer client.Kill()

	rpcClient, err := client.Client()
	if err != nil {
		slog.Error("Failed to create client", err)
		os.Exit(1)
	}

	raw, err := rpcClient.Dispense("animal")
	if err != nil {
		slog.Error("Failed to dispense dog", err)
		os.Exit(1)
	}

	dog := raw.(exten.Animal)
	woof := dog.Speak()
	slog.Info("Animal says", slog.String("dog", woof))

}
