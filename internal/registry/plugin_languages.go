package registry

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
