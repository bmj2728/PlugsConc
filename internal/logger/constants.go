package logger

const (
	// KeyJobID is the constant key used to represent a unique identifier for a job in operations and logging.
	KeyJobID = "job_id"
	// KeyMaxRetries represents the maximum number of retry attempts for a job in case of failure.
	KeyMaxRetries = "max_retries"
	// KeyRetryDelay represents the key used to identify the retry delay duration for a job in context
	// or logging operations.
	KeyRetryDelay = "retry_delay"
	// KeyRetryCount represents the constant key for tracking the number of retries a job has undergone.
	KeyRetryCount = "retry_count"
	// KeyJobSubmittedAt is a constant key representing the timestamp when a job was submitted.
	KeyJobSubmittedAt = "submitted_at"
	// KeyJobStartedAt is a constant key used to store or retrieve the timestamp of when a job started from a context.
	KeyJobStartedAt = "started_at"
	// KeyJobFinishedAt is a constant key used to store or retrieve the job's completion time in context or logs.
	KeyJobFinishedAt = "finished_at"
	// KeyJobDuration is a constant key representing the duration of a job in seconds, used for context
	// or logging operations.
	KeyJobDuration = "job_duration_seconds"
	// KeyWorkerCount denotes the number of workers in the pool.
	KeyWorkerCount = "worker_count"
	// KeySubmittedJobs represents the total number of jobs submitted to the pool.
	KeySubmittedJobs = "jobs_submitted"
	// KeyFailedSubmissions indicates the count of job submissions that failed.
	KeyFailedSubmissions = "failed_submissions"
	// KeyPoolStartedAt records the timestamp when the pool was started.
	KeyPoolStartedAt = "pool_started_at"
	// KeyPoolStoppedAt holds the timestamp when the pool was stopped.
	KeyPoolStoppedAt = "pool_stopped_at"
	// KeyPoolCompletedAt captures the timestamp when the pool completed processing.
	KeyPoolCompletedAt = "pool_completed_at"
	// KeyPoolDuration refers to the total duration of the pool's operation in seconds.
	KeyPoolDuration = "pool_duration_seconds"
	// KeyPoolClosed signifies whether the pool has been closed.
	KeyPoolClosed = "pool_closed"
	// KeySuccessfulJobs represents the number of successfully processed jobs.
	KeySuccessfulJobs = "successful_jobs"
	// KeyFailedJobs indicates the count of jobs that failed during processing.
	KeyFailedJobs = "failed_jobs"
	// KeyPoolMetrics provides the metrics collected for the pool.
	KeyPoolMetrics = "pool_metrics"
	// KeyWorkerID is a constant key used to associate a worker's unique ID with context or logging operations.
	KeyWorkerID = "worker_id"
	// KeyBatchErrors represents the logging key for storing or referencing batch error information.
	KeyBatchErrors = "batch_errors"
	// KeyJobMetrics represents the identifier key for job-related metrics in the system.
	KeyJobMetrics = "job_metrics"
	// KeyJobValue represents the value associated with a specific job in the job processing system.
	KeyJobValue = "job_value"
	// KeyJobError represents the key used to record or identify errors associated
	// with a specific job during processing.
	KeyJobError = "job_error"
	// KeyPluginLoadErrors is a constant key used to log or identify errors encountered during plugin loading processes.
	KeyPluginLoadErrors = "plugin_load_errors"
	// KeyPluginMap is a constant key used to store or retrieve the plugin map from context or logging operations.
	KeyPluginMap = "plugin_map"
	// KeyError represents the constant string key used for storing or identifying errors within the system.
	KeyError = "err"
	// KeyGroupPlugin represents the grouping key for plugin-specific metadata and attributes in logging or
	// configuration.
	KeyGroupPlugin = "plugin"
	// KeyPluginType represents the type of a plugin within the system, used to classify plugins by their functional
	// category.
	KeyPluginType = "type"
	// KeyPluginName represents the key identifier for a plugin's name field in logging or metadata structures.
	KeyPluginName = "name"
	// KeyPluginVersion represents the version of the plugin as a string constant.
	KeyPluginVersion = "version"
	// KeyPluginFormat represents the format of the plugin (e.g., binary, script, etc.).
	KeyPluginFormat = "format"
	// KeyPluginLanguage represents the programming language used by the plugin, such as "Go", "Python", or "Java".
	KeyPluginLanguage = "language"
	// KeyPluginEntrypoint defines the entrypoint for a plugin, specifying the executable or script to initialize it.
	KeyPluginEntrypoint = "entrypoint"
	// KeyGroupAbout defines the constant value used as a key for organizing plugin metadata about description,
	// maintainer, and URL.
	KeyGroupAbout = "about"
	// KeyPluginDescription represents the key for a plugin's description in metadata or structured logging.
	KeyPluginDescription = "description"
	// KeyPluginMaintainer represents the maintainer of the plugin as a string constant.
	KeyPluginMaintainer = "maintainer"
	// KeyPluginURL is used to reference the URL providing additional information about the plugin.
	KeyPluginURL = "url"
	// KeyGroupHandshakeConfig represents the group key for handshake configuration constants in the plugin metadata.
	KeyGroupHandshakeConfig = "handshake_config"
	// KeyHandshakeProtocolVersion is the key used to define the protocol version in the handshake configuration.
	KeyHandshakeProtocolVersion = "protocol_version"
	// KeyHandshakeMagicCookieKey represents the key used to identify the handshake magic cookie in plugin
	// configuration.
	KeyHandshakeMagicCookieKey = "magic_cookie_key"
	// KeyHandshakeMagicCookieValue represents the constant key for the handshake magic cookie value used during plugin
	// communication.
	KeyHandshakeMagicCookieValue = "magic_cookie_value"
	// KeyGroupSecurity represents a constant key for grouping security-related information in structured logging or data.
	KeyGroupSecurity = "security"
	// KeyPluginAutoMTLS represents the configuration key for enabling or disabling automatic mTLS in plugins.
	KeyPluginAutoMTLS = "auto_mtls"
)
