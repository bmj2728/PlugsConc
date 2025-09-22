package registry

// PluginState represents various states a plugin can be in during its lifecycle or validation process.
// Values 100+ are reserved for error codes
type PluginState int

const (

	// PluginStateUnknown indicates that the plugin's state is currently unknown or uninitialized during its lifecycle.
	// This is the default state for a plugin and nil value for PluginState.
	PluginStateUnknown PluginState = iota
	// PluginDirectoryDiscovered indicates that a plugin's directory has been located but not yet scanned or validated.
	PluginDirectoryDiscovered
	// PluginDirectoryScanned indicates that the plugin's directory has been successfully scanned for necessary files.
	PluginDirectoryScanned
	// PluginDirectoryValidated indicates that the plugin's directory has been successfully validated.
	PluginDirectoryValidated
	// PluginDataLoaded indicates the state where the plugin's data has been successfully loaded into memory.
	PluginDataLoaded
	// PluginManifestValidated indicates the state where the plugin's manifest has been successfully validated.
	PluginManifestValidated
	// PluginAvailable indicates that the plugin has passed validation to create LaunchDetails.
	// It is ready to be launched by the system.
	PluginAvailable
	// PluginLaunching indicates the state where the plugin is in the process of initializing and starting up.
	PluginLaunching
	// PluginRunning indicates that the plugin is currently active and functioning as intended.
	PluginRunning
	// PluginStopped indicates the state when a plugin has been stopped after running.
	PluginStopped
)
const (
	// PluginMissingManifest is used when a plugin is missing a manifest file
	PluginMissingManifest = PluginState(100)
	// PluginMissingChecksum is used when a plugin is missing a checksum file
	PluginMissingChecksum = PluginState(101)
	// PluginMissingBinary is used when a plugin is missing a binary file
	PluginMissingBinary = PluginState(102)
	// PluginInvalidManifest indicates the plugin's manifest file is present but invalid,
	// such as malformed or containing errors.
	PluginInvalidManifest = PluginState(103)
	// PluginInvalidLaunchDetails indicates the plugin's launch details are invalid or improperly configured.
	PluginInvalidLaunchDetails = PluginState(104)
	// PluginInvalidChecksum indicates the plugin's checksum file is present but invalid,
	// such as being corrupted or tampered with.
	PluginInvalidChecksum = PluginState(105)
	// PluginInvalidBinary indicates the plugin's binary file is present but invalid, for instance, not executable.
	PluginInvalidBinary = PluginState(106)
	// PluginBadChecksum indicates the plugin's checksum file is present but does not match the plugin's binary file.
	PluginBadChecksum = PluginState(107)
	// PluginFailedToLaunch indicates the plugin failed to launch due to an error during its initialization process.
	PluginFailedToLaunch = PluginState(108)
	// PluginFailedToStop indicates that a plugin could not be terminated successfully during its lifecycle.
	PluginFailedToStop = PluginState(109)
	// PluginStoppedUnexpectedly indicates that the plugin ceased running unexpectedly due to an unforeseen issue.
	PluginStoppedUnexpectedly = PluginState(110)
)
