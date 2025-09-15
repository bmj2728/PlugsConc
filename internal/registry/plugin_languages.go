package registry

import "sync"

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

// PluginLanguages defines a thread-safe structure to manage a mapping of PluginLanguage constants to their string
// names.
type PluginLanguages struct {
	mu        sync.RWMutex
	languages map[PluginLanguage]string
}

// AvailablePluginLanguages is a mapping of PluginLanguage constants to their corresponding string values.
var AvailablePluginLanguages = PluginLanguages{
	mu: sync.RWMutex{},
	languages: map[PluginLanguage]string{
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
	},
}

// Get retrieves the string representation of the specified PluginLanguage from the thread-safe languages map.
func (l *PluginLanguages) Get(language PluginLanguage) string {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.languages[language]
}

// GetByString retrieves the string value of a PluginLanguage based on its corresponding string identifier.
func (l *PluginLanguages) GetByString(language string) string {
	l.mu.RLock()
	defer l.mu.RUnlock()
	AvailablePluginLanguageLookup.mu.RLock()
	defer AvailablePluginLanguageLookup.mu.RUnlock()
	return l.Get(AvailablePluginLanguageLookup.GetLanguage(language))
}

// PluginLanguageLookup is a thread-safe map for storing and retrieving programming languages for plugins.
// It uses a sync.RWMutex to manage concurrent access to its language mappings.
// The languages field holds a map associating language names (as strings) with PluginLanguage values.
type PluginLanguageLookup struct {
	mu        sync.RWMutex
	languages map[string]PluginLanguage
}

// AvailablePluginLanguageLookup provides thread-safe access to a lookup table of supported plugin languages.
var AvailablePluginLanguageLookup = PluginLanguageLookup{
	mu: sync.RWMutex{},
	languages: map[string]PluginLanguage{
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
	},
}

// GetLanguage retrieves the PluginLanguage corresponding to the provided language string from the lookup map.
func (l *PluginLanguageLookup) GetLanguage(language string) PluginLanguage {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.languages[language]
}

// IsValidLanguage checks if the given language is present in the available plugin language lookup map.
func IsValidLanguage(lang string) bool {
	AvailablePluginLanguageLookup.mu.RLock()
	defer AvailablePluginLanguageLookup.mu.RUnlock()
	_, ok := AvailablePluginLanguageLookup.languages[lang]
	return ok
}
