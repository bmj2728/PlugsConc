package registry

// PluginType represents a category or type of plugin, managed using integer values.
type PluginType int

// AnimalPlugin represents a standard plugin type for animals.
// AnimalGRPCPlugin represents a gRPC plugin type for animals.
const (
	AnimalPlugin PluginType = iota
	AnimalGRPCPlugin
)

// PluginTypes maps PluginType constants to their string representations for better readability
// and debugging convenience.
var PluginTypes = map[PluginType]string{
	AnimalPlugin:     "AnimalPlugin",
	AnimalGRPCPlugin: "AnimalGRPCPlugin",
}

// String returns the string representation of the PluginType by looking it up in the PluginTypes map.
func (t PluginType) String() string {
	return PluginTypes[t]
}

// PluginFormat represents the format type for plugins, defined as an integer enumeration.
type PluginFormat int

// GRPC represents the gRPC plugin format in the PluginFormat enumeration.
// RPC represents the RPC plugin format in the PluginFormat enumeration.
const (
	GRPC PluginFormat = iota
	RPC
)

// PluginFormats maps PluginFormat constants to their corresponding string representations.
var PluginFormats = map[PluginFormat]string{
	GRPC: "grpc",
	RPC:  "rpc",
}

// String returns the string representation of the PluginFormat.
func (f PluginFormat) String() string {
	return PluginFormats[f]
}

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

// PluginLanguages maps PluginLanguage enum values to their corresponding string representations.
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
	ObjC:   "obj-c",
	PHP:    "php",
}

// String returns the string representation of the PluginLanguage using the PluginLanguages map.
func (l PluginLanguage) String() string {
	return PluginLanguages[l]
}
