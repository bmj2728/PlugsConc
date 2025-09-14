package registry

import (
	"sync"

	"github.com/bmj2728/PlugsConc/shared/pkg/animal"
	"github.com/hashicorp/go-plugin"
)

// PluginType represents a custom type used for defining various plugin classifications within the system.
type PluginType int

// AnimalPlugin represents a standard animal-related plugin type.
// AnimalGRPCPlugin represents an animal-related plugin type using gRPC.
const (
	AnimalPlugin PluginType = iota
	AnimalGRPCPlugin
)

// PluginTypes provides thread-safe storage and retrieval of plugin types, mapped from PluginType to their
// implementations.
type PluginTypes struct {
	types map[PluginType]plugin.Plugin
	mu    sync.RWMutex
}

// AvailablePluginTypes is a global instance of PluginTypes containing mappings of PluginType to their respective
// implementations.
var AvailablePluginTypes = PluginTypes{
	types: map[PluginType]plugin.Plugin{
		AnimalPlugin:     &animal.AnimalPlugin{},
		AnimalGRPCPlugin: &animal.AnimalGRPCPlugin{},
	},
	mu: sync.RWMutex{},
}

// Get retrieves the value associated with the given PluginType from the types map in a thread-safe manner.
func (pt *PluginTypes) Get(pluginType PluginType) plugin.Plugin {
	pt.mu.RLock()
	defer pt.mu.RUnlock()
	return pt.types[pluginType]
}

// GetByString retrieves the value associated with a plugin type string from the PluginTypes map if it is valid.
func (pt *PluginTypes) GetByString(pluginType string) plugin.Plugin {
	pt.mu.RLock()
	defer pt.mu.RUnlock()
	if AvailablePluginTypesLookup.IsValidPluginType(pluginType) {
		return pt.types[AvailablePluginTypesLookup.GetPluginType(pluginType)]
	} else {
		return nil
	}
}

// PluginTypeLookup is a thread-safe structure that maps string keys to PluginType objects for plugin type management.
type PluginTypeLookup struct {
	types map[string]PluginType
	mu    sync.RWMutex
}

// AvailablePluginTypesLookup is a mapping of plugin type names to their corresponding PluginType values.
var AvailablePluginTypesLookup = PluginTypeLookup{
	types: map[string]PluginType{
		"animal":      AnimalPlugin,
		"animal-grpc": AnimalGRPCPlugin,
	},
	mu: sync.RWMutex{},
}

// GetPluginType retrieves the PluginType associated with the provided pluginType key from the lookup map.
func (ptl *PluginTypeLookup) GetPluginType(pluginType string) PluginType {
	ptl.mu.RLock()
	defer ptl.mu.RUnlock()
	return ptl.types[pluginType]
}

// IsValidPluginType checks if the given plugin type string exists in the PluginTypeLookup's types map.
func (ptl *PluginTypeLookup) IsValidPluginType(pluginType string) bool {
	ptl.mu.RLock()
	defer ptl.mu.RUnlock()
	_, ok := ptl.types[pluginType]
	return ok
}

/**
 * Plugin Format Types
**/

// PluginFormat represents the type for defining various plugin communication formats.
type PluginFormat int

// GRPC represents a plugin format using gRPC.
// RPC represents a plugin format using RPC.
const (
	GRPC PluginFormat = iota
	RPC
)

// PluginFormats is a struct that manages a thread-safe map of PluginFormat values to their string representations.
type PluginFormats struct {
	formats map[PluginFormat][]plugin.Protocol
	mu      sync.RWMutex
}

// AvailablePluginFormats defines a mapping between PluginFormat constants and their string representations.
var AvailablePluginFormats = PluginFormats{
	formats: map[PluginFormat][]plugin.Protocol{
		GRPC: {plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
		RPC:  {plugin.ProtocolNetRPC},
	},
	mu: sync.RWMutex{},
}

func (pf *PluginFormats) Get(format PluginFormat) []plugin.Protocol {
	pf.mu.RLock()
	defer pf.mu.RUnlock()
	return pf.formats[format]
}

func (pf *PluginFormats) GetByString(format string) []plugin.Protocol {
	pf.mu.RLock()
	defer pf.mu.RUnlock()
	return pf.formats[AvailablePluginFormatLookup.GetPluginFormat(format)]
}

type PluginFormatLookup struct {
	mu      sync.RWMutex
	formats map[string]PluginFormat
}

// AvailablePluginFormatLookup is a pre-initialized PluginFormatLookup containing supported plugin formats
// with thread safety.
var AvailablePluginFormatLookup = PluginFormatLookup{
	formats: map[string]PluginFormat{
		"grpc": GRPC,
		"rpc":  RPC,
	},
	mu: sync.RWMutex{},
}

// GetPluginFormat retrieves the PluginFormat associated with the given format string from the lookup.
func (pfl *PluginFormatLookup) GetPluginFormat(format string) PluginFormat {
	pfl.mu.RLock()
	defer pfl.mu.RUnlock()
	return pfl.formats[format]
}

// IsValidFormat checks if the provided format string exists as a key in the PluginFormatLookup map.
// Returns true if valid.
func (pfl *PluginFormatLookup) IsValidFormat(format string) bool {
	_, ok := pfl.formats[format]
	return ok
}

/**
 * Plugin Language Types
**/

// PluginLanguage represents a language in which plugins can be implemented. It is represented as an integer type.
type PluginLanguage int

// Go represents the PluginLanguage value for the Go programming language.
// Python represents the PluginLanguage value for the Python programming language.
// Swift represents the PluginLanguage value for the Swift programming language.
// Ruby represents the PluginLanguage value for the Ruby programming language.
// CPP represents the PluginLanguage value for the C++ programming language.
// Java represents the PluginLanguage value for the Java programming language.
// Kotlin represents the PluginLanguage value for the Kotlin programming language.
// Node represents the PluginLanguage value for the Node.js runtime environment.
// Dart represents the PluginLanguage value for the Dart programming language.
// CSharp represents the PluginLanguage value for the C# programming language.
// ObjC represents the PluginLanguage value for the Objective-C programming language.
// PHP represents the PluginLanguage value for the PHP programming language.
const (
	Go PluginLanguage = iota
	Python
	Swift
	Ruby
	CPP
	Java
	Kotlin
	Node
	Dart
	CSharp
	ObjC
	PHP
)

// PluginLanguages is a mapping of PluginLanguage constants to their corresponding string values.
var PluginLanguages = map[PluginLanguage]string{
	Go:     "go",
	Python: "python",
	Swift:  "swift",
	Ruby:   "ruby",
	CPP:    "c++",
	Java:   "java",
	Kotlin: "kotlin",
	Node:   "node",
	Dart:   "dart",
	CSharp: "c#",
	ObjC:   "objc",
	PHP:    "php",
}

// String returns the string representation of the PluginLanguage using the PluginLanguages map.
func (l PluginLanguage) String() string {
	return PluginLanguages[l]
}

// PluginLanguageLookup maps string identifiers to their corresponding PluginLanguage constants.
var PluginLanguageLookup = map[string]PluginLanguage{
	"go":     Go,
	"python": Python,
	"swift":  Swift,
	"ruby":   Ruby,
	"c++":    CPP,
	"java":   Java,
	"kotlin": Kotlin,
	"node":   Node,
	"dart":   Dart,
	"c#":     CSharp,
	"objc":   ObjC,
	"php":    PHP,
}

// IsValidLanguage checks if the provided language string exists in the PluginLanguageLookup map.
func IsValidLanguage(lang string) bool {
	_, ok := PluginLanguageLookup[lang]
	return ok
}
