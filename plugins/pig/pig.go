package pig

import "fmt"

type Pig struct {
	pluginName    string
	pluginType    string
	pluginVersion string
}

func (p Pig) Name() string {
	return p.pluginName
}

func (p Pig) Version() string {
	return p.pluginVersion
}

func (p Pig) Type() string {
	return p.pluginType
}

func (p Pig) Speak() {
	fmt.Println("Woof")
}
