package dog

import "fmt"

type Dog struct {
	pluginName    string
	pluginType    string
	pluginVersion string
}

func (d Dog) Name() string {
	return d.pluginName
}

func (d Dog) Version() string {
	return d.pluginVersion
}

func (d Dog) Type() string {
	return d.pluginType
}

func (d Dog) Speak() {
	fmt.Println("Woof")
}
