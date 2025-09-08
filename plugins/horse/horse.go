package horse

import "fmt"

type Horse struct {
	pluginName    string
	pluginType    string
	pluginVersion string
}

func (h Horse) Name() string {
	return h.pluginName
}

func (h Horse) Version() string {
	return h.pluginVersion
}

func (h Horse) Type() string {
	return h.pluginType
}

func (h Horse) Speak() {
	fmt.Println("Woof")
}
