### Notes

## Protocol Buffers
- New .proto files should be nested in a folder within the proto dir.
  - These should be formatted as the .proto's name + pb
  - e.g. proto/animalpb/animal.proto, proto/dbpb/db.proto

## Plugins
- Pluggable interfaces must provide concrete rpc/grpc implementations, and plugin types that provide clients/servers
- Plugins comm via ipc(rpc/grpc)
- Plugins can be written in numerous languages
- Plugins must provide a handshake for integrity NOT security

## os.Root
- Need better understanding of how paths are affected by this