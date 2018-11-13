# NetServer
[![Build Status](http://drone.giacomofurlan.name/api/badges/MicroLayers/net-server/status.svg?branch=master)](http://drone.giacomofurlan.name/MicroLayers/net-server)

This server aims to simplify the creation of micro services listening on TCP or
Unix socket, using JSON or Protobuf as message protocol, letting the developer
write only the handlers via a go plugin.

**IMPORTANT**: Windows is currently not supported, Being the go plugin system a
Unix-only feature as of Go 1.11.2.

## Features

- Unified configuration file: use one YAML configuration file for both the common
NetServer settings and the module ones.
- Listening on TCP and Unix sockets in just a configuration edit.
- You have to write only a module: create a struct compliant with
[module/NetServerModule](https://github.com/MicroLayers/net-server/blob/master/module/NetServerModule.go),
export it and compile it as a go plugin: the connections are managed by the main server.

## Configuration

The basic configuration file, which can be found [here](https://github.com/MicroLayers/net-server/blob/master/dist/configuration.yml),
has a few sections:

- **module**: the path of the module to load (either relative or absolute)
- **listen -> unix**: Unix socket listening options
- **listen -> tcp**: TCP socket listening options
- listen -> http: (not implemented yet)

## Building a module

1. Create a struct compliant with [module/NetServerModule](https://github.com/MicroLayers/net-server/blob/master/module/NetServerModule.go)
2. export a new instance of that struct with the symbol name `NetServerModule`
3. build the package as a plugin, i.e. `build -buildmode=plugin -o my-module.so my-module-main-file.go`
4. build NetServer **with the same build arguments of the module** (except for
buildmode), otherwise the application will complain about a different version of
the binaries.

### Examples

You can find an echo server module [here](https://github.com/MicroLayers/net-server/blob/master/module/examples/echo/EchoModule.go).

## Running a server

1. Copy the [configuration](https://github.com/MicroLayers/net-server/blob/master/dist/configuration.yml)
file within the same directory of the net-server binary and the module static object.
2. Edit the configuration accordingly, being sure to point to the module file.
3. run `net-server` with the optional flag `-config path/to/configuration.yml` if
the configuration file is not in the same folder of the executable or has a different name.
