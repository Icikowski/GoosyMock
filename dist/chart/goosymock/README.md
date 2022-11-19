# GoosyMock

![GoosyMock banner](https://raw.githubusercontent.com/Icikowski/GoosyMock/master/images/banner.png)

**GoosyMock** is a configurable test service for mocking HTTP responses,
featuring SSL support, dedicated administration API and custom payloads
(binary files that can be served on particular routes). It's also prepared
for Kubernetes deployments (Helm chart), making it easy to adapt and use.

[**Explore documentation**](https://icikowski.github.io/GoosyMock) | [GitHub repository](https://github.com/Icikowski/GoosyMock)

## Easy installation

```bash
helm repo add icikowski https://charts.icikowski.pl
helm install goosymock icikowski/goosymock
```

## Parameters

### Global overrides

| Name | Description | Default value |
|-|-|-|
| `overrides.name` | Partially override deployment name | `""` |
| `overrides.fullName` | Fully override deployment name | `""` |
| `overrides.repository` | Override the image repository whose default is `ghcr.io` | `""` |
| `overrides.image` | Overrides the image name whose default is `icikowski/goosymock` | `""` |
| `overrides.tag` | Overrides the image tag whose default is the `.Chart.AppVersion` | `""` |

### GoosyMock parameters

#### Global configuration

| Name | Description | Default value |
|-|-|-|
| `goosyMock.logLevel` | Desired log level | `"info"` |
| `goosyMock.prettyLog` | Enable prettified and colored output | `false` |
| `goosyMock.maxPayloadSize` | Maximum accepted payload size (in megabytes) | `64` |

#### Admin API Service configuration

| Name | Description | Default value |
|-|-|-|
| `goosyMock.adminApi.port` | Container port on which Admin API Service will be exposed | `8081` |
| `goosyMock.adminApi.ssl.enabled` | Enable SSL for Admin API Service | `false` |
| `goosyMock.adminApi.ssl.port` | Container port on which secured Admin API Service will be exposed | `8444` |
| `goosyMock.adminApi.ssl.secretName` | Name of Kubernetes Secret containing TLS certificate & key (`tls.crt` & `tls.key` files) for Admin API Service SSL server; **required** if SSL is enabled for Admin API Service | `""` |

#### Content Service configuration

| Name | Description | Default value |
|-|-|-|
| `goosyMock.contentService.port` | Container port on which Content Service will be exposed | `8080` |
| `goosyMock.contentService.ssl.enabled` | Enable SSL for Content Service | `false` |
| `goosyMock.contentService.ssl.port` | Container port on which secured Content Service will be exposed | `8443` |
| `goosyMock.contentService.ssl.secretName` | Name of Kubernetes Secret containing TLS certificate & key (`tls.crt` & `tls.key` files) for Content Service SSL server; **required** if SSL is enabled for Content Service | `""` |

#### Health probes configuration

| Name | Description | Default value |
|-|-|-|
| `goosyMock.health.port` | Container port on which health probes (liveness & readiness) will be exposed | `8888` |

### Kubernetes Services

#### Admin API Service

| Name | Description | Default value |
|-|-|-|
| `service.adminApi.type` | Kubernetes Service type for Admin API Service | `"ClusterIP"` |
| `service.adminApi.port` | Kubernetes Service port for HTTP Admin API Service | `80` |
| `service.adminApi.securedPort` | Kubernetes Service port for HTTPS Admin API Service (used only if SSL is enabled for Admin API Service) | `443` |

#### Content Service

| Name | Description | Default value |
|-|-|-|
| `service.contentService.type` | Kubernetes Service type for Content Service | `"ClusterIP"` |
| `service.contentService.port` | Kubernetes Service port for HTTP Content Service | `80` |
| `service.contentService.securedPort` | Kubernetes Service port for HTTPS Content Service (used only if SSL is enabled for Content Service) | `443` |

### Kubernetes Ingresses

#### Admin API Service

| Name | Description | Default value |
|-|-|-|
| `ingress.adminApi.enabled` | Enable Kubernetes Ingress for Admin API Service | `false` |
| `ingress.adminApi.*` | Kubernetes Ingress configuration for Admin API Service | _N/A_ |

#### Content Service

| Name | Description | Default value |
|-|-|-|
| `ingress.contentService.enabled` | Enable Kubernetes Ingress for Content Service | `false` |
| `ingress.contentService.*` | Kubernetes Ingress configuration for Content Service | _N/A_ |
