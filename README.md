# xeno

Xeno is a collection of modular libraries that can be used to quickly develop production-grade web services.

The Xeno platform aims to make life easier for developers by eliminating boilerplate, using secure defaults, providing plug-in points, and offering common tools needed for web services.

Go is the language being targeted for development of Xeno, but all protocols will be portable to other langiages.

Xeno consists of modules, or "cogs" that can be used to help you build your web services. Cogs can be used alone or together. Below is a list of in-development and planned cogs.

### In development:

**Xeno API** ( [start here](/x/README.md) )

The `x` cog will help quickly build API services. It includes secure-by default settings, lots of customizability, and helpers galore.


### Planned:

**Xeno RPC**

xrpc will make development of RPC services easier by creating a wrapper around gRPC. xrpc will allow for gRPC development without needing to maintain protobuf definitions. xrpc will be encrypted and authenticated by default.

**Xeno Auth Hub**

The Xeno Auth Hub and xauth cog will provide an authentication authority for services within a system.

**Xeno Worker**

xwork will allow for the easy development of workers that can handle background tasks and eventually, scheduled tasks

**Xeno Static**

Xeno static will be a fileserver application that allows for quick, easy, and secure static file serving.