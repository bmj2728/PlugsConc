package logger

import (
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

// The plugin writer and proxy sink are used to proxy logs to a plugin
// the proxy gets assigned to sync or async interceptor
// the proxy sink uses the plugin client to write to the ShipLog func

type PluginWriter struct {
	client *plugin.Client
}

func (w *PluginWriter) Write(p []byte) (n int, err error) {
	// We'll use this write to call the gRPC function
	// TODO update func once proto buffer and generated code available
	return 0, nil
}

func (w *PluginWriter) Close() error {
	w.client.Kill()
	return nil
}

func PluginProxySink(name string,
	client *plugin.Client,
	level hclog.Level,
	color hclog.ColorOption,
	includeLocation bool,
	isJSON bool,
) hclog.SinkAdapter {
	w := &PluginWriter{
		client: client,
	}
	opts := hclog.LoggerOptions{
		Name:            name,
		Level:           level,
		Output:          w,
		Color:           color,
		IncludeLocation: includeLocation,
		JSONFormat:      isJSON,
	}
	return hclog.NewSinkAdapter(&opts)
}

// TODO helpers for parsing log json into proto buffer
// proto has section for known k-v pairs and then a map<string, well-known.Value>
