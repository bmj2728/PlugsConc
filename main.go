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
	"dog": &exten.AnimalPlugin{},
	"pig": &exten.AnimalPlugin{},
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

	dogClient := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  handshakeConfig,
		Plugins:          pluginMap,
		Cmd:              exec.Command("./plugins/dog/dog"),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolNetRPC},
	})
	defer dogClient.Kill()

	rpcDogClient, err := dogClient.Client()
	if err != nil {
		slog.Error("Failed to create dogClient", slog.Any("err", err))
		os.Exit(1)
	}

	dog, err := rpcDogClient.Dispense("dog")
	if err != nil {
		slog.Error("Failed to dispense dog", slog.Any("err", err))
		os.Exit(1)
	}
	woof := dog.(exten.Animal).Speak()

	pigClient := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  handshakeConfig,
		Plugins:          pluginMap,
		Cmd:              exec.Command("./plugins/pig/pig"),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolNetRPC},
	})
	defer pigClient.Kill()

	rpcPigClient, err := pigClient.Client()
	if err != nil {
		slog.Error("Failed to create pigClient", slog.Any("err", err))
		os.Exit(1)
	}

	pig, err := rpcPigClient.Dispense("pig")
	if err != nil {
		slog.Error("Failed to dispense pig", slog.Any("err", err))
	}
	oink := pig.(exten.Animal).Speak()

	slog.Info("Animal says", slog.String("dog", woof), slog.String("pig", oink))

}
