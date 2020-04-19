![logo_transparent_small](https://user-images.githubusercontent.com/5942370/79701451-37f2ec00-826b-11ea-8e37-b6f89111a5d7.png)

## Intro
`Vektor` enables development of modern web services. Vektor is a Go webserver framework designed to simplify the development of web APIs by eliminating boilerplate, using secure defaults, providing plug-in points, and offering common pieces needed for web apps. Vektor is fairly opinionated, but aims to provide flexibility in the right places.

## Background
We see Go as the best language to build web APIs and rich backend services, and so Vektor's Go components are all focused on building those things. 

Vektor consists of components that can be used to help you build your web apps and services. Vektor components can be used alone or together. Below is a list of in-development and planned components.

### In development:

**Vektor API**

The `vk` component is central to Vektor. It helps to quickly build production-ready API services with Go. It includes secure-by-default settings such as built-in LetsEncrypt, lots of customizability, and helpers galore. It will soon integrate with SubOrbital's Hive job scheduler to allow performing more complex and performance-oriented work. `vk` enables minimal-boilerplate servers with an intuitive wrapper around the most performant HTTP router, `httprouter`.

**Vektor Logger**

`vlog` will be a low-effort logging package that will allow for structured or text-based logging, that will easily work with popular third-party logging systems.

### Planned:

**Vektor Authentication**

The `vauth` component will provide an authentication library for service-service authentication (such as between `vk` services) as well as client-server authentication that can be extended to fit any need.