---
title: Binaries
---

# Running **GoosyMock** from binaries

## Downloading pre-built binaries

**GoosyMock** offers pre-build binaries for following operating systems and
architectures:

| Operating system | OS codename | Architectures |
|-|-|-|
| Linux | `linux` | `amd64`, `386`, `arm`, `arm64`, `ppc64le` |
| Windows | `windows` | `amd64`, `386` |
| MacOS | `darwin` | `arm64` |

Properly named binaries are attached as release assets and can be downloaded
from [latest tag's attachments](https://git.sr.ht/~icikowski/goosymock/refs).

## Building binaries manually

!!! info "Prerequisites"
    - **Git** 2.34+
    - **Go** 1.19+
    - **Taskfile** 3.17+

### Cloning the repository

=== "Clone via HTTPS"

    ```bash
    git clone https://git.sr.ht/~icikowski/goosymock.git
    ```

=== "Clone via SSH"

    ```bash
    git clone git@github.com:Icikowski/GoosyMock.git
    ```

### Building the sources

!!! tip "Using _Taskfile_"
    This project utilizes _[Taskfile](https://taskfile.dev)_ for build 
    automatization. If you are willing to use `task` command for building
    binaries, please install _Taskfile_ as described in 
    [official documentation](https://taskfile.dev/installation).

Binary will be built for current OS & architecture and placed in
`target/binaries` directory.

=== "Building static binary"

    ```bash
    task build:static
    ```

=== "Building dynamic binary"

    ```bash
    task build:dynamic
    ```
