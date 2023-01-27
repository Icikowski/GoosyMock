---
title: "About project"
hide: ["navigation"]
---

# GoosyMock

![GoosyMock icon](assets/img/logo.png#only-light){ align=right .smallImg }
![GoosyMock icon](assets/img/logo-inv.png#only-dark){ align=right .smallImg }

**GoosyMock** is a configurable test service for mocking HTTP responses,
featuring SSL support, dedicated administration API and custom payloads
(binary files that can be served on particular routes). It's also prepared
for Kubernetes deployments (Helm chart), making it easy to adapt and use.

This project is successor to [**GPTS**](https://github.com/Icikowski/GPTS)
and was initially meant to be an improvement over existing solution by
introducing new internal architecture, better approach for serving custom
content and overall improvements to both code's and repository's structure.
Changes turned out to be hard to apply, which led to rewriting whole source
code from scratch.

## Features

- [X] Dynamic default response for non-configured paths
- [X] Declarative configuration
    - [X] Support for `YAML` and `JSON` configuration formats
    - [X] Support for method-specific responses per each route (`GET`, `POST`, `PUT`, `PATCH`, `DELETE`)
    - [X] Support for wildcards (`*`) in routes' paths 
    - [X] Support for custom payloads (that can be uploaded, reused, updated and deleted)
    - [X] Well-documented administration API (OpenAPI specification included)
- [X] Docker support
    - [X] Based on latest [Google's `distroless`](https://github.com/GoogleContainerTools/distroless) image
    - [X] Small size (≈ 20MB)
    - [X] Running in rootless mode
- [X] Kubernetes support
    - [X] Helm chart available
    - [X] Rootless container image
    - [X] All settings configurable via chart values
    - [X] Support for Ingress controllers
