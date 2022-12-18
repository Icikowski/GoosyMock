---
title: "About project"
hide: ["navigation"]
---

# GoosyMock

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

- [X] Default response for 
- [X] Declarative configuration
    - [X] Support `YAML` and `JSON` configuration formats
    - [X] Support per-method response definitions (`GET`, `POST`, `PUT`, `PATCH`, `DELETE`)
    - [X] Support default response definition (for not configured methods)
    - [X] Support sub paths handling
    - [X] Well-documented administration API
- [X] Docker support
    - [X] Based on latest [Google's "distroless"](https://github.com/GoogleContainerTools/distroless) image
    - [X] Small size (≈ 20MB)
    - [X] Running in rootless mode
- [X] Kubernetes support
    - [X] Helm chart available
    - [X] Image running as non-root
    - [X] All settings can be configured via chart values
    - [X] Support Ingress controllers
