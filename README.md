# GusT

## Intro
`GusT` enables development of modern web apps and services. The GusT framework is an implementation of the "gust stack", which is comprised of Go, Rust, and TypeScript. These modern languages are each very good at specific tasks and can all be useful in developing a modern web product. We want to enable developers to use these technologies together, by providing libraries that allow them to work seamlessly with each other. 

## Purpose
The GusT framework aims to make life easier for developers by eliminating boilerplate, using secure defaults, providing plug-in points, and offering common pieces needed for web apps. GusT is fairly opinionated, but aims to provide flexibility in the right places. 

## Background
We see Go as the best language to build web APIs and rich backend services, and so GusT's Go components are all focused on building those things. 

Rust is incredibly memory and compute efficient, and is therefore well suited for realtime, high-throughput, and batch data tasks, so GusT's Rust components help in building those kinds of applications, while integrating with GusT's Go backend components. 

TypeScript is a well-loved and powerful language for building web apps, and so GusT's TypeScript components allow you to tie your webapps into your GusT-build backend Go and Rust services to make your frontend code responsive and easily maintainable.

GusT consists of components that can be used to help you build your web apps and services. GusT components can be used alone or together. Below is a list of in-development and planned components.

### In development:

**GusT API**

The `gapi` component will help quickly build API services with Go. It includes secure-by-default settings such as built-in LetsEncrypt, lots of customizability, and helpers galore. It will integrate with GusT's Rust components to allow performing more complex and performance-oriented work.

**GusT Functions**

`gfn` will aid in the development of functions that can be called in an RPC manner, or scheduled in a job-like manner. Functions will be written in Rust, and can be triggered from the companion Go or TypeScript client libraries. `gfn` will include secure communication and encryption by default.

**GusT Logger**

`glog` will be a low-effort logging package that will allow for structured or text-based logging, that will easily work with popular third-party logging systems.

### Planned:

**GusT Authentication**

The `gauth` component will provide an authentication library for service-service authentication (such as between a `gapi` service and a `gfn` function) as well as client-server authentication that can be extended to fit any need.

