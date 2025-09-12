package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"

	"github.com/bmj2728/PlugsConc/internal/logger"
	"github.com/bmj2728/PlugsConc/pkg/shared/animal"

	"github.com/hashicorp/go-plugin"
)

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "ANIMAL_PLUGIN",
	MagicCookieValue: "hello",
}

var pluginMap = map[string]plugin.Plugin{
	"dog":   &animal.AnimalPlugin{},
	"pig":   &animal.AnimalPlugin{},
	"cow":   &animal.AnimalPlugin{},
	"cat":   &animal.AnimalPlugin{},
	"horse": &animal.AnimalPlugin{},
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
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolNetRPC}, // add plugin.ProtocolGRPC
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
	woof := dog.(animal.Animal).Speak(true)

	pigClient := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  handshakeConfig,
		Plugins:          pluginMap,
		Cmd:              exec.Command("./plugins/pig/pig"),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolNetRPC}, // add plugin.ProtocolGRPC
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
	oink := pig.(animal.Animal).Speak(false)

	catClient := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  handshakeConfig,
		Plugins:          pluginMap,
		Cmd:              exec.Command("./plugins/cat/cat"),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolNetRPC}, // add plugin.ProtocolGRPC
	})
	defer catClient.Kill()

	rpcCatClient, err := catClient.Client()
	if err != nil {
		slog.Error("Failed to create catClient", slog.Any("err", err))
		os.Exit(1)
	}

	cat, err := rpcCatClient.Dispense("cat")
	if err != nil {
		slog.Error("Failed to dispense cat", slog.Any("err", err))
	}
	meow := cat.(animal.Animal).Speak(true)

	cowClient := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  handshakeConfig,
		Plugins:          pluginMap,
		Cmd:              exec.Command("./plugins/cow/cow"),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolNetRPC}, // add plugin.ProtocolGRPC
	})
	defer cowClient.Kill()

	rpcCowClient, err := cowClient.Client()
	if err != nil {
		slog.Error("Failed to create cowClient", slog.Any("err", err))
		os.Exit(1)
	}
	cow, err := rpcCowClient.Dispense("cow")
	if err != nil {
		slog.Error("Failed to dispense cow", slog.Any("err", err))
	}
	moo := cow.(animal.Animal).Speak(true)

	horseClient := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  handshakeConfig,
		Plugins:          pluginMap,
		Cmd:              exec.Command("./plugins/horse/horse"),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolNetRPC}, // add plugin.ProtocolGRPC
	})
	defer horseClient.Kill()

	rpcHorseClient, err := horseClient.Client()
	if err != nil {
		slog.Error("Failed to create horseClient", slog.Any("err", err))
		os.Exit(1)
	}
	horse, err := rpcHorseClient.Dispense("horse")
	if err != nil {
		slog.Error("Failed to dispense horse", slog.Any("err", err))
	}
	neigh := horse.(animal.Animal).Speak(false)

	fmt.Printf("The dog says %s\n", woof)
	fmt.Printf("The pig says %s\n", oink)
	fmt.Printf("The cat says %s\n", meow)
	fmt.Printf("The cow says %s\n", moo)
	fmt.Printf("The horse says %s\n", neigh)

	plugin.CleanupClients()
}
