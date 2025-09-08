package cow

import "fmt"

type Cow struct {
	pluginName    string
	pluginType    string
	pluginVersion string
}

func (c Cow) Name() string {
	return c.pluginName
}

func (c Cow) Version() string {
	return c.pluginVersion
}

func (c Cow) Type() string {
	return c.pluginType
}

func (c Cow) Speak() {
	fmt.Println("Moo")
}
