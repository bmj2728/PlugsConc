PlugsConc

A concise, extensible playground that demonstrates a production‑style architecture for running external plugins with a strong developer experience: structured logging with rotation, an optional persistent MQ for async logs, a worker pool with retries and metrics, a simple plugin registry based on manifests, and a file‑watcher to observe plugin folders.

This README walks you through what’s built, how it all fits together, and how to use or extend it.


Key capabilities

- Plugin system
  - HashiCorp go‑plugin transports (net/rpc and gRPC) with a simple shared interface (Animal) used by sample plugins.
  - Per‑plugin manifest (YAML) describing name, type, version, entrypoint, transport format, language, handshake, and security options.
  - Loader scans the plugins directory, validates entrypoints, parses manifests, and exposes loaded metadata.
  - Security features:
    - Optional Auto mTLS supported by go‑plugin (AutoMTLS flag from manifests or direct usage).
    - SHA‑256 checksum loader to provide SecureConfig to plugin clients.
    - Handshake protocol settings (magic cookie key/value + protocol version) from manifest.

- Structured logger
  - Based on hashicorp/go-hclog with multi-sink support.
  - Clear split between realtime sinks (console, files, etc.) and async logging via a persistent queue.
  - Realtime: MultiLogger is the primary intercept logger that fans out to registered sinks immediately.
  - Async: the log queue has its own InterceptLogger (created with AsyncInterceptLogger) that has its own primary output and sinks. A realtime MultiLogger can register an AsyncSink that enqueues logs into the queue.
  - Typical flow: multi -> AsyncSink -> AsyncWriter -> queue -> async intercept logger -> async sinks (e.g., rotated file).
  - File rotation via lumberjack; size/backups/age/compress configurable.
  - JSON or colored console output; include location toggle.

Logging initialization example (see main.go lines 60–77)

```
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
	q := logger.LogQueue(conf, asyncI)
	
	// This creates a specialized sink that gets attached to the synchronous logger and is 
	// responsible for shipping logs to the queue.
	aLogs := logger.AsyncSink("async-sink", q, conf.LogLevel(), hclog.ColorOff, conf.AddSource(), true)
	
	// This registers the sink with the synchronous intercept logger.
	multiLogger.RegisterSink(aLogs)
	
	// We now have a multi-logger configures to write synchronously to the console and asynchronously to a file.
	multiLogger.Info("File logger initialized")
```

- File watcher
  - fsnotify‑based watcher observes plugin directories (create/modify/remove/rename/chmod) and logs changes.

- Worker pool
  - Fixed pool with N workers and buffered/unbuffered channels.
  - Job abstraction supports context propagation, cancellation (with or without cause), deadlines/timeouts, configurable retries with backoff delay, and panic safety.
  - Metrics for pool lifecycle and per‑job timing: submissions, failed submissions, successes, failures, started/stopped/completed timestamps, duration, attempts, etc.
  - Graceful stop/shutdown/terminate modes and safe result delivery.

- MQ (message queue)
  - Optional persistent queue for offloading log events; a worker drains the queue and logs at the requested level.

- Configuration system
  - YAML‑based config with sensible defaults and helpers for reading application/logging/dirs/file‑watcher/worker‑pool/MQ sections.


Repository layout (selected)

- main.go — application bootstrap: config, logging sinks, worker pool demo jobs, plugin loader, specific plugin clients (dog, cat, dog‑grpc), fsnotify watcher, and MQ log example.
- internal/logger — multi‑sink logger, console/file helpers, async writer abstraction, constants for structured fields.
- internal/worker — pool, worker, job and metrics; context helpers for job/pool metadata; retry/cancellation logic.
- internal/registry — manifest types/loader; plugin formats/types/languages lookups; validation helpers; launch config derivation.
- internal/mq — persistent logging queue integration (sqliteq + varmq) and job types.
- internal/checksum — SHA‑256 checksum file loader for plugin binaries.
- internal/watcher — placeholder for general watcher interface (fsnotify used directly in main.go for now).
- internal/config — config models/defaults/loader and accessor helpers.
- shared/pkg/animal — shared plugin interfaces, and RPC/gRPC shims used by the example plugins.
- plugins/* — example plugin folders (cat, dog, dog‑grpc, pig, cow, horse), each with a manifest and an entrypoint binary.


How things work together

1) Startup
- Load configuration from config.yaml (using an fs rooted at repository root via os.OpenRoot).
- Initialize a realtime MultiLogger (console).
- Create the async side: Rotator -> AsyncInterceptLogger -> LogQueue; register an AsyncSink on the realtime MultiLogger to enqueue logs (see main.go lines 60–77).
- Start the persistent log MQ worker, which drains the queue and logs using the async intercept logger.
- Start a worker pool and submit some sample jobs to demonstrate logging from within job contexts.

2) Plugin discovery and use
- The registry scans the configured plugins directory; for each subdirectory it attempts to read manifest.yaml, compute a manifest hash, validate the entrypoint, and add it to the in‑memory Manifests registry.
- The app also demonstrates creating explicit go‑plugin clients for cat, dog (RPC) and dog‑grpc (gRPC) and calling Animal.Speak(true/false).
- For cat, a SecureConfig is provided by reading plugins/cat/cat.sha256 to enable checksum verification.

3) File watching
- fsnotify watcher is started and each valid plugin folder is added. File change events are logged so you can observe hot‑changes to plugin folders.

4) Worker pool
- A pool is created with configurable size and logging; jobs illustrate context‑based structured logs and result handling. The pool records metrics for lifecycle and per‑job durations and supports graceful stop/shutdown/terminate semantics.


Configuration

File: config.yaml (see internal/config/models.go for schema, internal/config/default.go for defaults)

- application
  - app_name: string
  - app_mode: string (dev/prod/etc.)
  - app_version:
    - major/minor/patch/full/codename

- directories
  - plugins_dir: path to the plugins directory (default ./plugins)
  - plugin_configs_dir: path to per‑plugin configs (unused by core but available)
  - logs_dir: path to logs directory (default ./logs)

- logging
  - log_level: trace|debug|info|warn|error
  - log_filename: log file name (default app.log)
  - log_max_size: max MB before rotation (capped by file helper to <= 2MB if unset or too large)
  - log_max_backups: number of rotated files to keep
  - log_max_age: days to keep rotated files
  - log_compress: whether to gzip rotated files
  - log_include_location: include source locations in logs
  - mq:
    - log_enable_persistent_queue: bool (enable async logging MQ)
    - log_db_file: SQLite file name (e.g., logs.db)
    - log_queue: queue name
    - log_remove_on_complete: remove items after processing

- file_watcher
  - fw_enabled: bool (global toggle, main.go currently starts regardless and logs if creation fails)
  - fw_watch_plugins: bool (when true, plugin folders are added to the watcher)

- worker_pool
  - wp_max_workers: maximum pool workers (code also considers CPU limits)


Logging

- Console logger is created via MultiLogger and set as default for hclog. You can register additional sinks (e.g., file sink with rotation) using logger.FileSink.
- File rotation is handled by lumberjack with configurable size/backups/age/compression.
- Async/Persistent logging (optional):
  - internal/mq.LogQueue(conf, log) initializes a sqlite‑backed varmq persistent queue and returns a queue handle.
  - You can enqueue messages as mq.NewLoggerJob(level, "message", key, value, ...); a worker consumes entries and logs with the provided level.


Worker pool

Types
- Pool: manages workers, a jobs channel, a results channel, a metrics channel, and lifecycle state.
- Job: unit of work with context, metrics, retry settings, and cancel/timeouts. Create with worker.NewJob(ctx, func(ctx) (any, error) {...}).
- JobResult: result, error, and copied metrics for a completed job.
- PoolMetrics: lifecycle times and counters; thread‑safe updates.

Highlights
- Retry support: Job.WithRetry(maxRetries, retryDelayMs); worker loops until success or attempts exhausted, honoring cancel.
- Cancellation/timeouts: Job.WithCancel(), WithCancelCause(), WithTimeout(d), WithTimeoutCause(d, cause), WithDeadline(t), WithDeadlineCause(t, cause).
- Panic safety: job execution protected; panics converted to errors with stack trace.
- Graceful lifecycle: Stop (waits, keeps result chan open), Shutdown (waits + closes channels), Terminate (fast cancel/close). Metrics record started/stopped/completed/duration.
- Metrics fan‑in: workers send success/failure to a pool metrics channel, aggregated under lock.

Observability via context
- internal/worker/ctx.go stores and retrieves keys such as job_id, retry counts, submitted/started/finished times, duration, worker_id, pool metrics snapshots, etc., mirroring constants in internal/logger/constants.go.


Plugin system

Manifests
- Each plugin folder should provide manifest.yaml with at least:

  plugin:
    name: "cat"
    type: "animal"            # see internal/registry/plugin_types.go
    format: "rpc"              # "rpc" or "grpc" (internal/registry/plugin_formats.go)
    entrypoint: "./plugins/cat/cat"
    language: "go"
    version: "1.0.0"
  about:
    description: "A cat that speaks"
    maintainer: "you@example.com"
    url: "https://example.com/cat"
  handshake:
    protocol_version: 1
    magic_cookie_key: "PLUGIN_MAGIC"
    magic_cookie_value: "secret"
  security:
    auto_mtls: true

- The loader computes an MD5 of the manifest content (for quick change detection) and validates the entrypoint is present in PATH/relative.
- Launch details are derived from the manifest, including handshake config and allowed protocols.

Types and formats
- internal/registry/plugin_types.go maps logical plugin "types" to go‑plugin Plugin implementations. The sample exposes:
  - type: "animal" -> net/rpc (AnimalPlugin)
  - type: "animal‑grpc" -> gRPC (AnimalGRPCPlugin)
- internal/registry/plugin_formats.go maps "rpc" or "grpc" to allowed go‑plugin protocols.

Security: checksums + handshake
- internal/checksum.LoadSHA256 reads a .sha256 file and returns a go‑plugin SecureConfig with SHA‑256 checksum. main.go shows providing this for the cat plugin.
- HandshakeConfig is built from manifest values; missing required fields produce errors.
- Optional AutoMTLS can be enabled.

Capabilities and sandboxing
- Purpose: Sensitive operations (filesystem, network, process) are mediated by host services. Plugins do not touch the OS directly; instead, they request operations via host‑provided services. Requested capabilities in the plugin manifest inform what the host may allow at runtime, enabling robust sandboxing and least‑privilege defaults.
- Types: See internal/capability/capability.go for the schema mirrored in manifests.
  - Filesystem: a list of path‑scoped grants with permissions.
    - path: file or directory path
    - permissions: any of [read, write, list, create, delete]
    - recursive: whether the grant applies to subdirectories (dirs only)
  - Network: split into egress and ingress rule sets.
    - egress: [{ protocol: tcp|udp, hosts: [hostname|IP], ports: [int] }]
    - ingress: [{ protocol: tcp|udp, ports: [int], allowed_origins: [IP/CIDR/host] }]
  - Process: controls executing and managing processes.
    - exec: { command: string, args: [string] } — whitelist of allowed commands/args the plugin may ask the host to run
    - kill: [scopes] — e.g., ["children"] to restrict to processes spawned for this plugin
    - list: [scopes] — e.g., ["children"]
    - signal: [scopes] — if used, same scoping semantics
- Relation to manifests: The Manifest struct includes `capabilities` and is parsed by the registry loader. Enforcement is performed by host services (future/ongoing work). If a capability is not requested (or the section is omitted), the default is deny.

Declaring capabilities in a plugin manifest
- Add a top‑level `capabilities` section alongside `plugin`, `about`, `handshake`, and `security`.
- Example (abridged):

  capabilities:
    filesystem:
      - path: "/home/user/data/"
        permissions: [read, write, list, create, delete]
        recursive: true
      - path: "/etc/config.json"
        permissions: [read]
    network:
      egress:
        - protocol: tcp
          hosts: [ "api.example.com", "192.168.1.100" ]
          ports: [ 443, 8080 ]
      ingress:
        - protocol: tcp
          ports: [ 9000 ]
          allowed_origins: [ "192.168.1.100", "192.168.1.101" ]
    process:
      - exec:
          command: "/usr/bin/rsync"
          args: [ "-a", "--delete", "*" ]
      - kill: [ children ]
      - list: [ children ]

Host services expectations
- FilesystemService: validates path + permission + recursion before reads/writes/creates/deletes.
- NetworkService: validates egress by (protocol, host, port) and ingress by (protocol, port, origin) before connecting/listening.
- ProcessService: validates exec requests against the whitelist; kill/list/signal requests are scoped (e.g., only children). Host may log and deny non‑granted operations.

Best practices for plugin authors
- Request the minimum needed privileges; manifests are reviewed and can be rejected by the host.
- Use explicit paths and ports; avoid wildcards unless absolutely necessary.
- Keep capability requests stable; changing them will trigger re‑review in managed environments.

Examples in this repo
- See manifest.example.yaml for a complete sample.
- The sample cat and dog‑grpc manifests under plugins/ include a capabilities section demonstrating filesystem/network/process requests.

Using plugins from main
- A manual pluginMap with known keys (cat, dog, pig, cow, horse, dog‑grpc) is provided for demo.
- registry.NewPluginLoader(conf.PluginsDir(), log).Load() discovers manifests and logs launch details.
- Examples:
  - RPC: client := plugin.NewClient(&plugin.ClientConfig{HandshakeConfig: handshake, Plugins: pluginMap, Cmd: exec.Command("./plugins/dog/dog"), AllowedProtocols: {ProtocolNetRPC}, AutoMTLS: true})
  - gRPC: same but AllowedProtocols includes ProtocolGRPC and the plugin key is "dog‑grpc".


File watching

- fsnotify.NewWatcher() is created; for each valid manifest, the plugin folder path is added. A goroutine logs Events (Write, Create, Remove, Rename, Chmod) and Errors.
- The internal/watcher package is a placeholder for a fuller abstraction.


MQ usage example

- Enable in config:
  logging:
    log_enable_persistent_queue: true
    log_db_file: logs.db
    log_queue: app_logs
    log_remove_on_complete: true

- In code (already in main.go):
  q := mq.LogQueue(conf, log)
  q.Add(mq.NewLoggerJob(hclog.Info, "Hello, world!", "key", "value"))


Building and running

Prerequisites
- Go 1.22+
- Plugins compiled for your platform and located under the configured plugins_dir.

Run
- go run ./
- Observe console logs and logs/app.log. If MQ logging is enabled, you’ll also see logs/logs.db created.
- The sample prints the outputs of the cat, dog, and dog‑grpc plugins Speak() calls.


Developing a plugin (minimal guide)

1) Implement the shared interface
- See shared/pkg/animal. Your implementation must satisfy the Animal interface: Speak(isLoud bool) string.
- For RPC: implement Server()/Client() using AnimalPlugin; for gRPC: register with AnimalGRPCPlugin.

2) Build an entrypoint binary
- Place the compiled binary under plugins/<name>/<name>.

3) Create manifest.yaml in the same folder
- Follow the schema shown above; ensure entrypoint points to your binary.

4) (Optional) Create a SHA‑256 file
- Name: <entrypoint>.sha256 with the standard `<hex>  <filename>` format. When provided, main can supply SecureConfig to go‑plugin for checksum verification.

5) Start the app
- The loader will pick up your plugin’s folder, watch it for changes, and you can create a client to Dispense by its key.


Troubleshooting

- Plugin fails to start
  - Check manifest entrypoint path and permissions; verify the binary runs standalone.
- Manifest errors
  - The loader logs YAML/unmarshal errors and records problematic entries; check logs/app.log.
- Checksum errors
  - Verify the .sha256 file format; ensure it matches the current binary.
- MQ not writing
  - Confirm logging.mq.* settings and that SQLite file path is writable.
- Logs not rotating
  - Ensure log_max_size/backups/age are set; note internal/logger/file.go caps max size to 2MB by default.
- Worker pool not processing
  - Ensure you called Run() and are submitting jobs before closing the pool.


Reference: important types and helpers

- Logger: logger.MultiLogger, logger.FileSink, logger.NewRotator, logger.AsyncWriter; constants in internal/logger/constants.go.
- Worker: worker.NewPool, pool.Submit/SubmitBatch, pool.Results(), pool.Stop/Shutdown/Terminate; worker.NewJob and WithRetry/WithCancel*/WithTimeout*/WithDeadline* helpers; worker.JobMetrics and PoolMetrics accessors.
- Registry: registry.NewPluginLoader, loader.Load() -> Manifests; Manifest.ToLaunchDetails(); Plugin types and format lookups.
- MQ: mq.LogQueue(conf, log) and mq.NewLoggerJob.
- Config: config.LoadConfig(), getters like LogLevel(), LogsDir(), PluginsDir(), WorkerPoolMaxWorkers(), LogMQEnabled().
- Security: checksum.LoadSHA256() -> plugin.SecureConfig.


License

For educational/demo purposes. Integrate and adapt freely within your own projects.
