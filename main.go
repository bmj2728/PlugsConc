package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/bmj2728/PlugsConc/internal/checksum"
	"github.com/bmj2728/PlugsConc/internal/logger"
	"github.com/bmj2728/PlugsConc/internal/registry"
	"github.com/bmj2728/PlugsConc/shared/pkg/animal"
	"github.com/fsnotify/fsnotify"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

const (
	ConfigDir  = "."
	ConfigFile = "config.yaml"
)

var catHandshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "CAT_PLUGIN",
	MagicCookieValue: "lFLmoCE3ckw6erJxYxcRd6keedUodVMctD3XOGj9bLMYsFZi1Qh0vKEJftppo5ek",
}

var dogHandshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "DOG_PLUGIN",
	MagicCookieValue: "2ggRd5S9bhHottawB6eXwghiOAhekGORmOfIczh5b1D3AYlmrRWIXdbqwDHDJmjq",
}

func main() {
	/*
		Logger Setup Example w/ config
	*/

	// Synchronous logger setup:
	// multilogger is the primary logger for the application.
	// It is a synchronous intercept logger that writes to console and can be configured to write to
	// other io.Writers using sinks.
	multiLogger := logger.MultiLogger("app-name", hclog.Info, hclog.ForceColor, true, false)
	// Sets the default logger to the multilogger.
	hclog.SetDefault(multiLogger)
	//// Read in the configuration for the file logger.
	//logRotator := logger.NewRotator(filepath.Join("./logs", "app.log"),
	//	2,
	//	25,
	//	90,
	//	true)
	//// Asynchronous Logger Setup:
	//// asyncI is an Intercept logger, the choice of io.Writer is up to the user.
	//// This can take additional sinks, similar to the synchronous logger.
	//asyncI := logger.AsyncInterceptLogger("async-app-logs", hclog.Info, logRotator, hclog.ColorOff, false, true)
	//// This initializes the queue and worker for writing async logs
	//q := logger.LogQueue(conf, asyncI)
	//// This creates a specialized sink that gets attached to the synchronous logger and is
	//// responsible for shipping logs to the queue.
	//aLogs := logger.AsyncSink("async-sink", q, hclog.Info, hclog.ColorOff, true, true)
	//// This registers the sink with the synchronous intercept logger.
	//multiLogger.RegisterSink(aLogs)
	//// We now have a multi-logger configures to write synchronously to the console and asynchronously to a file.
	//multiLogger.Info("File logger initialized")

	/*
		Example General Worker Pool
	*/

	//workerPool := worker.NewPool(500, true, 1000, multiLogger.Named("worker_pool"))
	//
	//workerPool.Run()
	//
	//for i := 0; i < 5; i++ {
	//	// this is how you attach a contextual logger to a job
	//	jobCtx := hclog.WithContext(context.Background(), multiLogger.Named("job_logger").With("job_id", i))
	//	j := worker.NewJob(jobCtx, func(ctx context.Context) (any, error) {
	//		t := time.Now().Unix()
	//		x := 1.0 / float64(t)
	//		hclog.FromContext(jobCtx).Info("Done", "time", t)
	//		return x, nil
	//	})
	//	err := workerPool.Submit(j)
	//	if err != nil {
	//		hclog.FromContext(jobCtx).Error("Failed to submit job", logger.KeyError, err)
	//	}
	//}

	/*
		Example File Watcher
	*/

	// Initialize plugin filewatcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		multiLogger.Error("Failed to create watcher", logger.KeyError, err)
		multiLogger.Warn("Watching for changes will not work")
	}
	defer func(watcher *fsnotify.Watcher) {
		err := watcher.Close()
		if err != nil {
			multiLogger.Error("Failed to close watcher", logger.KeyError, err)
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
					multiLogger.Info("file changed:", "file", event.Name)
				}
				if event.Has(fsnotify.Create) {
					multiLogger.Info("file created:", "file", event.Name)
				}
				if event.Has(fsnotify.Remove) {
					multiLogger.Info("file removed:", "file", event.Name)
				}
				if event.Has(fsnotify.Rename) {
					multiLogger.Info("file renamed:", "file", event.Name)
				}
				if event.Has(fsnotify.Chmod) {
					multiLogger.Info("file mode changed:", "file", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				multiLogger.Error("filewatcher error: ", logger.KeyError, err)
			}
		}
	}(watcher)

	/*
		Plugin Loading
	*/

	// Get the configured plugins directory
	pluginsDir := "./plugins"
	multiLogger.Info("Plugins directory", "dir", pluginsDir)
	// Add the plugins directory to the file watcher using the absolute path
	pAbs, err := filepath.Abs(pluginsDir)
	if err != nil {
		multiLogger.Error("Failed to get absolute path for plugins", logger.KeyError, err)
	}
	err = watcher.Add(pAbs)
	if err != nil {
		multiLogger.Error("Failed to add plugins directory", logger.KeyError, err)
	}

	// Load plugins
	loader, err := registry.NewPluginLoader(pluginsDir, multiLogger)
	if err != nil {
		multiLogger.Error("Failed to create plugin loader", logger.KeyError, err)
		os.Exit(1)
	}
	p, e := loader.Load()
	if len(e) > 0 {
		multiLogger.Error("Failed to load plugins", logger.KeyError, e)
	}
	for d, m := range p.GetManifests() {
		multiLogger.Info("Plugin loaded", "manifest", m.Manifest(), "dir", d)
	}

	var pluginMapImported = make(map[string]plugin.Plugin)

	for _, m := range p.GetManifests() {

		// Generates the PluginMap
		validType := registry.AvailablePluginTypesLookup.IsValidPluginType(m.Manifest().PluginData.Type)
		if validType {
			pt := registry.AvailablePluginTypes.GetByString(m.Manifest().PluginData.Type)
			pluginMapImported[m.Manifest().PluginData.Name] = pt
		}

		// Establish plugin root
		pFolder, err := filepath.Abs(filepath.Join(pluginsDir, m.Manifest().PluginData.Name))
		if err != nil {
			multiLogger.Error("Failed to get absolute path", logger.KeyError, err)
		}
		// Add this plugin dir to the file watcher
		err = watcher.Add(pFolder)
		if err != nil {
			multiLogger.Error("Failed to add watcher", logger.KeyError, err)
		}

		// Convert Manifest to LaunchDetails
		ld := m.Manifest().ToLaunchDetails()
		if ld != nil {
			multiLogger.Info("Plugin loaded", "launch_details", ld.HandshakeConfig)
		}

		fsCap := m.Manifest().Capabilities.Filesystem
		for _, f := range fsCap {
			multiLogger.Info("Filesystem capability detected", "filesystem", f)
		}

		network := m.Manifest().Capabilities.Network
		if network != nil {
			multiLogger.Info("Network capability detected", "network", *network)
		}

		procCap := m.Manifest().Capabilities.Process
		for _, p := range procCap.Exec {
			multiLogger.Info("Process capability detected", "exec", p)
		}
		for _, p := range procCap.Kill {
			multiLogger.Info("Process capability detected", "kill", p)
		}
		for _, p := range procCap.List {
			multiLogger.Info("Process capability detected", "list", p)
		}
		for _, p := range procCap.Signal {
			multiLogger.Info("Process capability detected", "signal", p)
		}

	}

	cSHA, err := checksum.NewSHA256File("./plugins/cat")
	if err != nil {
		multiLogger.Error("Failed to load checksum", logger.KeyError, err)
		return
	}
	err = cSHA.Parse()
	if err != nil {
		multiLogger.Error("Failed to parse checksum", logger.KeyError, err)
		return
	}

	multiLogger.Info("Checksum parsed successfully", "hex", cSHA.Hash(), "file", cSHA.FileName())

	secConf, err := cSHA.SecConf()
	if err != nil {
		multiLogger.Error("Failed to get secure config", logger.KeyError, err)
		return
	}

	catClient := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  catHandshake,
		Logger:           multiLogger.Named("cat"),
		Plugins:          pluginMapImported,
		Cmd:              exec.Command("./plugins/cat/cat"),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolNetRPC},
		AutoMTLS:         true,
		SecureConfig:     secConf,
	})
	defer catClient.Kill()

	rpcCatClient, err := catClient.Client()
	if err != nil {
		multiLogger.Error("Failed to create catClient", logger.KeyError, err)
		os.Exit(1)
	}

	cat, err := rpcCatClient.Dispense("cat")
	if err != nil {
		multiLogger.Error("Failed to dispense cat", logger.KeyError, err)
		os.Exit(1)
	}
	meow := cat.(animal.Animal).Speak(true)
	fmt.Printf("The cat says %s\n", meow)

	gDogClient := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  dogHandshake,
		Plugins:          pluginMapImported,
		Logger:           multiLogger.Named("dog-grpc"),
		Cmd:              exec.Command("./plugins/dog-grpc/dog"),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
		AutoMTLS:         true,
	})
	defer gDogClient.Kill()

	rpcGDogClient, err := gDogClient.Client()
	if err != nil {
		multiLogger.Error("Failed to create gDogClient", logger.KeyError, err)
		os.Exit(1)
	}
	gDog, err := rpcGDogClient.Dispense("dog-grpc")
	if err != nil {
		multiLogger.Error("Failed to dispense dog", logger.KeyError, err)
	}
	gWoof := gDog.(animal.Animal).Speak(false)

	fmt.Printf("The dog-grpc says %s\n", gWoof)

	plugin.CleanupClients()

	<-make(chan struct{})
}
