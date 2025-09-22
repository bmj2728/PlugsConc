package main

import (
	"os"
	"path/filepath"

	"github.com/bmj2728/PlugsConc/internal/config"
	"github.com/bmj2728/PlugsConc/internal/logger"
	"github.com/bmj2728/PlugsConc/internal/mq"
	"github.com/bmj2728/PlugsConc/shared/pkg/animal"

	"github.com/hashicorp/go-hclog"
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
	// Code defined for initial startup prior to configuration loading.
	tempLogger := logger.DefaultLogger()

	cr, err := os.OpenRoot(".")
	if err != nil {
		tempLogger.Error("Failed to open root", logger.KeyError, err)
	}
	defer func(cr *os.Root) {
		err := cr.Close()
		if err != nil {
			tempLogger.Error("Failed to close root", logger.KeyError, err)
		}
	}(cr)

	tempLogger.Info("Root opened")
	conf := config.LoadConfig(cr, "config.yaml")
	tempLogger.Info("Config loaded", "config", conf.Application.AppName)
	// Synchronous logger setup:
	// multilogger is the primary logger for the application.
	// It is a synchronous intercept logger that writes to console and can be configured to write to
	// other io.Writers using sinks.
	multiLogger := logger.MultiLogger(conf.Application.AppName, conf.LogLevel(), hclog.ForceColor, conf.AddSource(), false)
	// Sets the default logger to the multilogger.
	hclog.SetDefault(multiLogger)
	// Read in the configuration for the file logger.
	logRotator := logger.NewRotator(filepath.Join(conf.LogsDir(), conf.LogFilename()),
		conf.LogMaxSize(),
		conf.LogMaxBackups(),
		conf.LogMaxAge(),
		conf.LogCompress())
	// Asynchronous Logger Setup:
	// asyncI is an Intercept logger, the choice of io.Writer is up to the user.
	// This can take additional sinks, similar to the synchronous logger.
	asyncI := logger.AsyncInterceptLogger("async-app-logs", conf.LogLevel(), logRotator, hclog.ColorOff, false, true)
	// This initializes the queue and worker for writing async logs
	q := mq.LogQueue(conf, asyncI)
	// This creates a specialized sink that gets attached to the synchronous logger and is
	// responsible for shipping logs to the queue.
	aLogs := logger.AsyncSink("async-sink", q, conf.LogLevel(), hclog.ColorOff, conf.AddSource(), true)
	// This registers the sink with the synchronous intercept logger.
	multiLogger.RegisterSink(aLogs)
	// We now have a multi-logger configures to write synchronously to the console and asynchronously to a file.
	multiLogger.Info("File logger initialized")

	//pluginsDir := conf.PluginsDir()
	//multiLogger.Info("Plugins directory", "dir", pluginsDir)
	//
	//workerPool := worker.NewPool(500, true, 1000, multiLogger.Named("worker_pool"))
	//
	//workerPool.Run()
	//
	//for i := 0; i < 5; i++ {
	//	// this is how you attach a contextual logger to a job
	//	jobCtx := hclog.WithContext(context.Background(), multiLogger.Named("job_logger").With("job_id", i))
	//	j := worker.NewJob(jobCtx, func(ctx context.Context) (any, error) {
	//		return "done", nil
	//	})
	//	err := workerPool.Submit(j)
	//	if err != nil {
	//		hclog.FromContext(jobCtx).Error("Failed to submit job", logger.KeyError, err)
	//	}
	//}
	//
	//pRoot := filepath.Join(conf.PluginsDir(), "cat")
	//open, err := os.OpenRoot(pRoot)
	//if err != nil {
	//	multiLogger.Error("Failed to open root", logger.KeyError, err)
	//}
	//defer func(open *os.Root) {
	//	err := open.Close()
	//	if err != nil {
	//		multiLogger.Error("Failed to close root", logger.KeyError, err)
	//	}
	//}(open)
	//
	//entrypoint := "cat"
	//csFilename := entrypoint + checksum.CSFileExt
	//
	//secConf, err := checksum.LoadSHA256(open, csFilename)
	//if err != nil {
	//	multiLogger.Error("Failed to load checksum", logger.KeyError, err)
	//} else {
	//	multiLogger.Info("Checksum loaded successfully")
	//}
	//
	//catClient := plugin.NewClient(&plugin.ClientConfig{
	//	HandshakeConfig:  handshakeConfig,
	//	Logger:           multiLogger.Named("cat"),
	//	Plugins:          pluginMap,
	//	Cmd:              exec.Command("./plugins/cat/cat"),
	//	AllowedProtocols: []plugin.Protocol{plugin.ProtocolNetRPC},
	//	AutoMTLS:         true,
	//	SecureConfig:     secConf,
	//})
	//defer catClient.Kill()
	//
	//rpcCatClient, err := catClient.Client()
	//if err != nil {
	//	multiLogger.Error("Failed to create catClient", logger.KeyError, err)
	//	os.Exit(1)
	//}
	//
	//cat, err := rpcCatClient.Dispense("cat")
	//if err != nil {
	//	multiLogger.Error("Failed to dispense cat", logger.KeyError, err)
	//	os.Exit(1)
	//}
	//meow := cat.(animal.Animal).Speak(true)
	//fmt.Printf("The cat says %s\n", meow)
	//
	//// Initialize plugin filewatcher
	//watcher, err := fsnotify.NewWatcher()
	//if err != nil {
	//	multiLogger.Error("Failed to create watcher", logger.KeyError, err)
	//	multiLogger.Warn("Watching for changes will not work")
	//}
	//defer func(watcher *fsnotify.Watcher) {
	//	err := watcher.Close()
	//	if err != nil {
	//		multiLogger.Error("Failed to close watcher", logger.KeyError, err)
	//	}
	//}(watcher)
	//
	//// Start generic watcher
	//// sig
	//go func(watcher *fsnotify.Watcher) {
	//	for {
	//		select {
	//		case event, ok := <-watcher.Events:
	//			if !ok {
	//				return
	//			}
	//			log.Println("event:", event)
	//			if event.Has(fsnotify.Write) {
	//				multiLogger.Info("file changed:", "file", event.Name)
	//			}
	//			if event.Has(fsnotify.Create) {
	//				multiLogger.Info("file created:", "file", event.Name)
	//			}
	//			if event.Has(fsnotify.Remove) {
	//				multiLogger.Info("file removed:", "file", event.Name)
	//			}
	//			if event.Has(fsnotify.Rename) {
	//				multiLogger.Info("file renamed:", "file", event.Name)
	//			}
	//			if event.Has(fsnotify.Chmod) {
	//				multiLogger.Info("file mode changed:", "file", event.Name)
	//			}
	//		case err, ok := <-watcher.Errors:
	//			if !ok {
	//				return
	//			}
	//			multiLogger.Error("filewatcher error: ", logger.KeyError, err)
	//		}
	//	}
	//}(watcher)
	//
	//loader, err := registry.NewPluginLoader(pluginsDir, multiLogger)
	//if err != nil {
	//	multiLogger.Error("Failed to create plugin loader", logger.KeyError, err)
	//	os.Exit(1)
	//}
	//p, e := loader.Load()
	//if len(e) > 0 {
	//	multiLogger.Error("Failed to load plugins", logger.KeyError, e)
	//}
	//for d, m := range p.GetManifests() {
	//	multiLogger.Info("Plugin loaded", "manifest", m.Manifest(), "dir", d)
	//}
	//
	//var pluginMapImported = make(map[string]plugin.Plugin)
	//
	//for _, m := range p.GetManifests() {
	//
	//	// map
	//	validType := registry.AvailablePluginTypesLookup.IsValidPluginType(m.Manifest().PluginData.Type)
	//	if validType {
	//		pt := registry.AvailablePluginTypes.GetByString(m.Manifest().PluginData.Type)
	//		pluginMapImported[m.Manifest().PluginData.Name] = pt
	//	}
	//
	//	pFolder, err := filepath.Abs(filepath.Join(pluginsDir, m.Manifest().PluginData.Name))
	//	if err != nil {
	//		multiLogger.Error("Failed to get absolute path", logger.KeyError, err)
	//	}
	//	err = watcher.Add(pFolder)
	//	if err != nil {
	//		multiLogger.Error("Failed to add watcher", logger.KeyError, err)
	//	}
	//
	//	ld := m.Manifest().ToLaunchDetails()
	//	multiLogger.Info("Plugin loaded", "launch_details", ld)
	//
	//}
	//
	//dogClient := plugin.NewClient(&plugin.ClientConfig{
	//	HandshakeConfig:  handshakeConfig,
	//	Plugins:          pluginMap,
	//	Logger:           multiLogger.Named("dog"),
	//	Cmd:              exec.Command("./plugins/dog/dog"),
	//	AllowedProtocols: []plugin.Protocol{plugin.ProtocolNetRPC},
	//	AutoMTLS:         true,
	//})
	//defer dogClient.Kill()
	//
	//rpcDogClient, err := dogClient.Client()
	//if err != nil {
	//	multiLogger.Error("Failed to create dogClient", logger.KeyError, err)
	//	os.Exit(1)
	//}
	//
	//dog, err := rpcDogClient.Dispense("dog")
	//if err != nil {
	//	multiLogger.Error("Failed to dispense dog", logger.KeyError, err)
	//	os.Exit(1)
	//}
	//woof := dog.(animal.Animal).Speak(true)
	//
	//gDogClient := plugin.NewClient(&plugin.ClientConfig{
	//	HandshakeConfig:  handshakeConfig,
	//	Plugins:          pluginMap,
	//	Logger:           multiLogger.Named("dog-grpc"),
	//	Cmd:              exec.Command("./plugins/dog-grpc/dog-grpc"),
	//	AllowedProtocols: []plugin.Protocol{plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
	//	AutoMTLS:         true,
	//})
	//defer gDogClient.Kill()
	//
	//rpcGDogClient, err := gDogClient.Client()
	//if err != nil {
	//	multiLogger.Error("Failed to create gDogClient", logger.KeyError, err)
	//	os.Exit(1)
	//}
	//gDog, err := rpcGDogClient.Dispense("dog-grpc")
	//if err != nil {
	//	multiLogger.Error("Failed to dispense dog", logger.KeyError, err)
	//}
	//gWoof := gDog.(animal.Animal).Speak(false)
	//
	//fmt.Printf("The dog says %s\n", woof)
	//fmt.Printf("The dog-grpc says %s\n", gWoof)
	//
	//plugin.CleanupClients()

	<-make(chan struct{})
}
