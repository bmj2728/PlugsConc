package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/bmj2728/PlugsConc/internal/checksum"
	"github.com/bmj2728/PlugsConc/internal/config"
	"github.com/bmj2728/PlugsConc/internal/logger"
	"github.com/bmj2728/PlugsConc/internal/mq"
	"github.com/bmj2728/PlugsConc/internal/registry"
	"github.com/bmj2728/PlugsConc/internal/worker"
	"github.com/fsnotify/fsnotify"

	"github.com/bmj2728/PlugsConc/shared/pkg/animal"

	"github.com/hashicorp/go-plugin"
)

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

	cr, err := os.OpenRoot(".")
	conf := config.LoadConfig(cr, "config.yaml")
	slog.Info("Config loaded") //kinda
	fmt.Println(conf)

	//Initialize logger
	colorHandler := logger.New(os.Stdout,
		&logger.Options{
			Level:     conf.LogLevel(),
			AddSource: conf.AddSource(),
			ColorMap:  conf.LoggerColorMap(),
			FullLine:  conf.FullLine()},
	)

	logRotator := logger.NewRotator(filepath.Join(conf.LogsDir(), conf.LogFilename()),
		conf.LogMaxSize(),
		conf.LogMaxBackups(),
		conf.LogMaxAge(),
		conf.LogCompress())

	rotatorOpts := &slog.HandlerOptions{
		AddSource: conf.AddSource(),
		Level:     conf.LogLevel(),
	}

	fileHandler := logger.NewFileLogHandler(logRotator, rotatorOpts)

	handlers := []slog.Handler{
		fileHandler,
		colorHandler,
	}

	multiHandler := logger.NewMultiHandler(handlers)

	slog.SetDefault(slog.New(multiHandler))
	slog.Info("Logger initialized")
	slog.Error("Logger initialized")
	slog.Warn("Logger initialized")
	slog.Debug("Logger initialized")

	pluginsDir := conf.PluginsDir()
	slog.Info("Plugins directory", slog.String("dir", pluginsDir))

	workerPool := worker.NewPool(500, true, 1000)

	workerPool.Run()

	for i := 0; i < 5; i++ {
		j := worker.NewJob(context.Background(), func(ctx context.Context) (any, error) {
			slog.Info("Job started", slog.Int("id", i))
			slog.Info("Job finished", slog.Int("id", i))
			return "done", nil
		})
		err := workerPool.Submit(j)
		if err != nil {
			slog.Error("Failed to submit job", slog.Any(logger.KeyError, err))
		}
	}

	pRoot := filepath.Join(conf.PluginsDir(), "cat")
	open, err := os.OpenRoot(pRoot)
	if err != nil {
		slog.Error("Failed to open root", slog.Any(logger.KeyError, err))
	}
	defer func(open *os.Root) {
		err := open.Close()
		if err != nil {
			slog.Error("Failed to close root", slog.Any(logger.KeyError, err))
		}
	}(open)

	entrypoint := "cat"
	csFilename := entrypoint + checksum.ChecksumFileExt

	secConf, err := checksum.LoadSHA256(open, csFilename)
	if err != nil {
		slog.Error("Failed to load checksum", slog.Any(logger.KeyError, err))
	}
	fmt.Println(secConf.Checksum)

	catClient := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  handshakeConfig,
		Plugins:          pluginMap,
		Cmd:              exec.Command("./plugins/cat/cat"),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolNetRPC},
		AutoMTLS:         true,
		SecureConfig:     secConf,
	})
	defer catClient.Kill()

	rpcCatClient, err := catClient.Client()
	if err != nil {
		slog.Error("Failed to create catClient", slog.Any(logger.KeyError, err))
		os.Exit(1)
	}

	cat, err := rpcCatClient.Dispense("cat")
	if err != nil {
		slog.Error("Failed to dispense cat", slog.Any(logger.KeyError, err))
		os.Exit(1)
	}
	woof := cat.(animal.Animal).Speak(true)
	fmt.Printf("The cat says %s\n", woof)

	// Initialize plugin filewatcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		slog.Error("Failed to create watcher", slog.Any(logger.KeyError, err))
		slog.Warn("Watching for changes will not work")
	}
	defer func(watcher *fsnotify.Watcher) {
		err := watcher.Close()
		if err != nil {
			slog.Error("Failed to close watcher", slog.Any(logger.KeyError, err))
		}
	}(watcher)

	// Start generic watcher
	// sig
	go func(watcher *fsnotify.Watcher) {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Has(fsnotify.Write) {
					slog.Info("file changed:", slog.String("file", event.Name))
				}
				if event.Has(fsnotify.Create) {
					slog.Info("file created:", slog.String("file", event.Name))
				}
				if event.Has(fsnotify.Remove) {
					slog.Info("file removed:", slog.String("file", event.Name))
				}
				if event.Has(fsnotify.Rename) {
					slog.Info("file renamed:", slog.String("file", event.Name))
				}
				if event.Has(fsnotify.Chmod) {
					slog.Info("file mode changed:", slog.String("file", event.Name))
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				slog.Error("filewatcher error: ", slog.Any(logger.KeyError, err))
			}
		}
	}(watcher)

	loader, err := registry.NewPluginLoader(pluginsDir)
	if err != nil {
		slog.Error("Failed to create plugin loader", slog.Any(logger.KeyError, err))
		os.Exit(1)
	}
	p, e := loader.Load()
	if len(e) > 0 {
		slog.Error("Failed to load plugins", slog.Any(logger.KeyError, e))
	}
	slog.Info("Plugins loaded", slog.Any("plugins", p.LogValue()))

	var pluginMapImported = make(map[string]plugin.Plugin)

	for _, m := range p.GetManifests() {

		// map
		validType := registry.AvailablePluginTypesLookup.IsValidPluginType(m.Manifest().PluginData.Type)
		if validType {
			pt := registry.AvailablePluginTypes.GetByString(m.Manifest().PluginData.Type)
			pluginMapImported[m.Manifest().PluginData.Name] = pt
		}

		pFolder, err := filepath.Abs(filepath.Join(pluginsDir, m.Manifest().PluginData.Name))
		if err != nil {
			slog.Error("Failed to get absolute path", slog.Any(logger.KeyError, err))
		}
		err = watcher.Add(pFolder)
		if err != nil {
			slog.Error("Failed to add watcher", slog.Any(logger.KeyError, err))
		}

		ld := m.Manifest().ToLaunchDetails()
		fmt.Println(m.Entrypoint())
		fmt.Println(ld.Cmd())

	}

	fmt.Println(pluginMapImported)
	fmt.Println(pluginMap)

	dogClient := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  handshakeConfig,
		Plugins:          pluginMap,
		Cmd:              exec.Command("./plugins/dog/dog"),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolNetRPC},
		AutoMTLS:         true,
	})
	defer dogClient.Kill()

	rpcDogClient, err := dogClient.Client()
	if err != nil {
		slog.Error("Failed to create dogClient", slog.Any(logger.KeyError, err))
		os.Exit(1)
	}

	dog, err := rpcDogClient.Dispense("dog")
	if err != nil {
		slog.Error("Failed to dispense dog", slog.Any(logger.KeyError, err))
		os.Exit(1)
	}
	woof = dog.(animal.Animal).Speak(true)

	gDogClient := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  handshakeConfig,
		Plugins:          pluginMap,
		Cmd:              exec.Command("./plugins/dog-grpc/dog-grpc"),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
		AutoMTLS:         true,
	})
	defer gDogClient.Kill()

	rpcGDogClient, err := gDogClient.Client()
	if err != nil {
		slog.Error("Failed to create gDogClient", slog.Any(logger.KeyError, err))
		os.Exit(1)
	}
	gDog, err := rpcGDogClient.Dispense("dog-grpc")
	if err != nil {
		slog.Error("Failed to dispense dog", slog.Any(logger.KeyError, err))
	}
	gWoof := gDog.(animal.Animal).Speak(false)

	fmt.Printf("The dog says %s\n", woof)
	fmt.Printf("The dog-grpc says %s\n", gWoof)

	plugin.CleanupClients()

	logQueue := mq.LogQueue()

	logQueue.Add(*mq.NewLoggerJob(slog.LevelInfo, "Logger queue started", []any{"queue_pending", logQueue.NumPending()}))

	<-make(chan struct{})
}
