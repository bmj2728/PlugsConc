### Notes

## Protocol Buffers
- New .proto files should be nested in a folder within the proto dir.
  - These should be formatted as the .proto's name + pb
  - e.g. proto/animalpb/animal.proto, proto/dbpb/db.proto
  - This is to avoid name collisions with other .proto files and add consistency for plugin developers
  - buf build to validate, buf generate to generate code
  - need to learn how to use buf linter

## Plugins
- Pluggable interfaces must provide concrete rpc/grpc implementations, and plugin types that provide clients/servers
- Plugins comm via ipc(rpc/grpc)
- Plugins can be written in numerous languages
- Plugins must provide a handshake for integrity NOT security
- Plugins are loaded from a specified directory
- Each plugin should have a folder with the same name as the plugin
- the folder should contain the plugin binary and a manifest.yaml file
- The initial load of plugins will just load the manifests and not the binaries
- the manifests contain the data necessary to load the plugin and establish a connection
- manifest schema can be found in manifest.example.yaml
- manifests are hashed on load to quickly identify changes
- manifests are stored in a map keyed by the directory path
- future state will allow plugins to be loaded dynamically using fsnotify

## Concurrency

- Research concurrent file system ops using os.Root.FS()'s fs.FS object
- Manifests provides a thread-safe map for concurrent access
- we currently create a new os.Root for each plugin while processing the manifest 
- I think this should ensure thread-safety
- the map can then be processed concurrently with no concerns as the data is in the manifest objects