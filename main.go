package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	//"os/exec"

	"github.com/bmj2728/PlugsConc/internal/logger"
	"github.com/bmj2728/PlugsConc/internal/registry"
	"github.com/bmj2728/PlugsConc/shared/pkg/animal"

	"github.com/hashicorp/go-plugin"
)

// pluginDir is the path to the directory containing the plugins.
// This will be configurable in the future.
var pluginDir = "./plugins"

// generic handshake configuration for testing.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "ANIMAL_PLUGIN",
	MagicCookieValue: "hello",
}

// manually generated list to be replaced by dynamic discovery.
var pluginMap = map[string]plugin.Plugin{
	"dog":      &animal.AnimalPlugin{},
	"dog-grpc": &animal.AnimalGRPCPlugin{},
	"pig":      &animal.AnimalPlugin{},
	"cow":      &animal.AnimalPlugin{},
	"cat":      &animal.AnimalPlugin{},
	"horse":    &animal.AnimalPlugin{},
}

func main() {

	logHandler := logger.New(os.Stdout,
		&logger.Options{
			Level:     slog.LevelInfo,
			AddSource: true,
			ColorMap:  logger.DefaultColorMap,
			FullLine:  false},
	)

	slog.SetDefault(slog.New(logHandler))
	slog.Info("Logger initialized")

	loader, err := registry.NewPluginLoader("./plugins")
	if err != nil {
		slog.Error("Failed to create plugin loader", slog.Any("err", err))
		os.Exit(1)
	}
	p, e := loader.Load()
	if len(e) > 0 {
		slog.Error("Failed to load plugins", slog.Any("err", e))
	}
	slog.Info("Plugins loaded", slog.Any("plugins", p.LogValue()))

	var pluginMapImported = make(map[string]plugin.Plugin)

	for _, m := range p.GetManifests() {
		config, err := m.Manifest().Handshake.ToConfig()
		if err != nil {
			slog.Error("Failed to convert manifest to config", slog.Any("err", err))
		}
		validType := registry.AvailablePluginTypesLookup.IsValidPluginType(m.Manifest().PluginType)
		if validType {
			pt := registry.AvailablePluginTypes.GetByString(m.Manifest().PluginType)
			pluginMapImported[m.Manifest().PluginName] = pt
		}
		entrypoint, err := filepath.Abs(filepath.Join(pluginDir, m.Manifest().PluginEntrypoint))
		if err != nil {
			slog.Error("Failed to get absolute path", slog.Any("err", err))
		}

		validFormat := registry.AvailablePluginFormatLookup.IsValidFormat(m.Manifest().PluginFormat)
		if validFormat {
			pf := registry.AvailablePluginFormats.GetByString(m.Manifest().PluginFormat)
			fmt.Println(pf)
		}
		fmt.Println(entrypoint)
		fmt.Println(config)
	}

	fmt.Println(pluginMapImported)
	fmt.Println(pluginMap)

	//dogClient := plugin.NewClient(&plugin.ClientConfig{
	//	HandshakeConfig:  handshakeConfig,
	//	Plugins:          pluginMap,
	//	Cmd:              exec.Command("./plugins/dog/dog"),
	//	AllowedProtocols: []plugin.Protocol{plugin.ProtocolNetRPC},
	//})
	//defer dogClient.Kill()
	//
	//rpcDogClient, err := dogClient.Client()
	//if err != nil {
	//	slog.Error("Failed to create dogClient", slog.Any("err", err))
	//	os.Exit(1)
	//}
	//
	//dog, err := rpcDogClient.Dispense("dog")
	//if err != nil {
	//	slog.Error("Failed to dispense dog", slog.Any("err", err))
	//	os.Exit(1)
	//}
	//woof := dog.(animal.Animal).Speak(true)

	//pigClient := plugin.NewClient(&plugin.ClientConfig{
	//	HandshakeConfig:  handshakeConfig,
	//	Plugins:          pluginMap,
	//	Cmd:              exec.Command("/home/brian/GolandProjects/PlugsConc/plugins/pig/pig"),
	//	AllowedProtocols: []plugin.Protocol{plugin.ProtocolNetRPC}, // add plugin.ProtocolGRPC
	//})
	//defer pigClient.Kill()
	//
	//rpcPigClient, err := pigClient.Client()
	//if err != nil {
	//	slog.Error("Failed to create pigClient", slog.Any("err", err))
	//	os.Exit(1)
	//}
	//
	//pig, err := rpcPigClient.Dispense("pig")
	//if err != nil {
	//	slog.Error("Failed to dispense pig", slog.Any("err", err))
	//}
	//oink := pig.(animal.Animal).Speak(false)

	//catClient := plugin.NewClient(&plugin.ClientConfig{
	//	HandshakeConfig:  handshakeConfig,
	//	Plugins:          pluginMap,
	//	Cmd:              exec.Command("./plugins/cat/cat"),
	//	AllowedProtocols: []plugin.Protocol{plugin.ProtocolNetRPC}, // add plugin.ProtocolGRPC
	//})
	//defer catClient.Kill()
	//
	//rpcCatClient, err := catClient.Client()
	//if err != nil {
	//	slog.Error("Failed to create catClient", slog.Any("err", err))
	//	os.Exit(1)
	//}
	//
	//cat, err := rpcCatClient.Dispense("cat")
	//if err != nil {
	//	slog.Error("Failed to dispense cat", slog.Any("err", err))
	//}
	//meow := cat.(animal.Animal).Speak(true)
	//
	//cowClient := plugin.NewClient(&plugin.ClientConfig{
	//	HandshakeConfig:  handshakeConfig,
	//	Plugins:          pluginMap,
	//	Cmd:              exec.Command("./plugins/cow/cow"),
	//	AllowedProtocols: []plugin.Protocol{plugin.ProtocolNetRPC}, // add plugin.ProtocolGRPC
	//})
	//defer cowClient.Kill()
	//
	//rpcCowClient, err := cowClient.Client()
	//if err != nil {
	//	slog.Error("Failed to create cowClient", slog.Any("err", err))
	//	os.Exit(1)
	//}
	//cow, err := rpcCowClient.Dispense("cow")
	//if err != nil {
	//	slog.Error("Failed to dispense cow", slog.Any("err", err))
	//}
	//moo := cow.(animal.Animal).Speak(true)
	//
	//horseClient := plugin.NewClient(&plugin.ClientConfig{
	//	HandshakeConfig:  handshakeConfig,
	//	Plugins:          pluginMap,
	//	Cmd:              exec.Command("./plugins/horse/horse"),
	//	AllowedProtocols: []plugin.Protocol{plugin.ProtocolNetRPC}, // add plugin.ProtocolGRPC
	//})
	//defer horseClient.Kill()
	//
	//rpcHorseClient, err := horseClient.Client()
	//if err != nil {
	//	slog.Error("Failed to create horseClient", slog.Any("err", err))
	//	os.Exit(1)
	//}
	//horse, err := rpcHorseClient.Dispense("horse")
	//if err != nil {
	//	slog.Error("Failed to dispense horse", slog.Any("err", err))
	//}
	//neigh := horse.(animal.Animal).Speak(false)
	//
	//gDogClient := plugin.NewClient(&plugin.ClientConfig{
	//	HandshakeConfig:  handshakeConfig,
	//	Plugins:          pluginMap,
	//	Cmd:              exec.Command("./plugins/dog-grpc/dog-grpc"),
	//	AllowedProtocols: []plugin.Protocol{plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
	//})
	//defer gDogClient.Kill()
	//
	//rpcGDogClient, err := gDogClient.Client()
	//if err != nil {
	//	slog.Error("Failed to create gDogClient", slog.Any("err", err))
	//	os.Exit(1)
	//}
	//gDog, err := rpcGDogClient.Dispense("dog-grpc")
	//if err != nil {
	//	slog.Error("Failed to dispense dog", slog.Any("err", err))
	//}
	//gWoof := gDog.(animal.Animal).Speak(false)
	//
	//fmt.Printf("The dog says %s\n", woof)
	//fmt.Printf("The pig says %s\n", oink)
	//fmt.Printf("The cat says %s\n", meow)
	//fmt.Printf("The cow says %s\n", moo)
	//fmt.Printf("The horse says %s\n", neigh)
	//fmt.Printf("The dog-grpc says %s\n", gWoof)
	//
	//plugin.CleanupClients()
}
