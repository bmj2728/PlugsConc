package cat

import "fmt"

type Cat struct {
	pluginName    string
	pluginType    string
	pluginVersion string
}

func (c Cat) Name() string {
	return c.pluginName
}

func (c Cat) Version() string {
	return c.pluginVersion
}

func (c Cat) Type() string {
	return c.pluginType
}

func (c Cat) Speak() {
	fmt.Println("Meow")
}
