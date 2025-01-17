# Image Service

Image Service provides a way to download Talos boot and install images generated with specific customizations.

The list of provided assets:

* ISO
* kernel/initramfs/kernel command line
* UKI
* disk images in various formats
* `installer` container images

Supported frontends:

* HTTP
* PXE service
* Container Registry

Official Image Service is available at [https://imager.talos.dev](https://imager.talos.dev).

## HTTP Frontend API

### `POST /flavors`

Create a new image flavor.

The request body is a YAML (JSON) encoded flavor description:

```yaml
customization:
    extraKernelArgs: # optional
        - vga=791
    systemExtensions: # optional
      officialExtensions: # optional
        - siderolabs/gvisor
        - siderolabs/amd-ucode
```

Output is a JSON-encoded flavor ID:

```json
{"id":"2a63b6e7dab90ec9d44f213339b9545bd39c6499b22a14cf575c1ca4b6e39ff8"}
```

This ID can be used to download images with this flavor.

Well-known flavor IDs:

* `376567988ad370138ad8b2698212367b8edcb69b5fd68c80be1f2ec7d603b4ba` - default flavor (without any customizations)

### `GET /image/:flavor/:version/:path`

Download a Talos boot image with the specified flavor and Talos version.

* `:flavor` is a flavor ID returned by `POST /flavor`
* `:version` is a Talos version, e.g. `v1.5.0`
* `:path` is a specific image path, details below

Common used parameters:

* `<arch>` image architecture: `amd64` or `arm64`
* `<platform>` Talos platform, e.g. `metal`, `aws`, `gcp`, etc.
* `<board>` is a board name (only for `arm64` `metal` platform), e.g. `rpi_generic`
* `-secureboot` identifies a Secure Boot asset

Supported image paths:

* `kernel-<arch>` (e.g. `kernel-amd64`) - raw kernel image
* `cmdline-<platform>[-<board>]-<arch>[-secureboot]` (e.g. `cmdline-metal-amd64`) - kernel command line
* `initramfs-<arch>.xz` (e.g. `initramfs-amd64.xz`) - initramfs image (including system extensions if configured)
* `<platform>-<arch>[-secureboot].iso` (e.g. `metal-amd64.iso`) - ISO image
* `<platform>-<arch>-secureboot-uki.efi` (e.g. `metal-amd64-secureboot-uki.efi) UEFI UKI image (Secure Boot compatible)
* `installer-<arch>[-secureboot].tar` (e.g. `installer-amd64.tar`) is a custom Talos installer image (including system extensions if configured)
* disk images in different formats (see Talos documentation for a full list):
  * `metal-<arch>[-secureboot].raw.xz` (e.g. `metal-amd64.raw.xz`) - raw disk image for metal platform
  * `aws-<arch>.raw.xz` (e.g. `aws-amd64.raw.xz`) - raw disk image for AWS platform, that can be imported as an AMI
  * `gcp-<arch>.raw.tar.gz` (e.g. `gcp-amd64.raw.tar.gz`) - raw disk image for GCP platform, that can be imported as a GCE image
  * ... other support image types

### `GET /versions`

Returns a list of Talos versions available for image generation.

```json
["v1.5.0","v1.5.1", "v1.5.2"]
```

### `GET /version/:version/extensions/official`

Returns a list of official system extensions available for the specified Talos version.

```json
[
  {
    "name": "siderolabs/amd-ucode",
    "ref": "ghcr.io/siderolabs/amd-ucode:20230804",
    "digest": "sha256:761a5290a4bae9ceca11468d2ba8ca7b0f94e6e3a107ede2349ae26520682832",
  },

]
```

## PXE Frontend API

PXE frontend provides an [iPXE script](https://ipxe.org/scripting) which automatically downloads and boots Talos.
The bare metal machine should be configured to boot from the URL provided by this API, e.g.:

```text
#!ipxe
chain --replace --autofree https://image.service/pxe/376567988ad370138ad8b2698212367b8edcb69b5fd68c80be1f2ec7d603b4ba/v1.5.0/metal-${buildarch}
```

### `GET /pxe/:flavor/:version/:path`

Returns an iPXE script which downloads and boots Talos with the specified flavor and Talos version, architecture and platform.

* `:flavor` is a flavor ID returned by `POST /flavor`
* `:version` is a Talos version, e.g. `v1.5.0`
* `:path` is a `<platform>-<arch>[-secureboot]` path, e.g. `metal-amd64`

In non-SecureBoot flavor, the following iPXE script is returned:

```text
#!ipxe
kernel https://image.service/image/:flavor/:version/kernel-<arch> <kernel-cmdline>
initrd https://image.service/image/:flavor/:version/initramfs-<arch>.xz
boot
```

For SecureBoot flavor, the following iPXE script is returned:

```text
#!ipxe
kernel https://image.service/image/:flavor/:version/<platform>-<arch>-secureboot.uki.efi
boot
```

## OCI Registry Frontend API

The Talos `installer` image used both for the initial install and upgrade can be pulled from the Image Service OCI registry.
If the image hasn't been created yet, it will be built on demand automatically.

### `docker pull <registry>/installer[-secureboot]/<flavor>:<version>`

Example: `docker pull imager.talos.dev/installer/376567988ad370138ad8b2698212367b8edcb69b5fd68c80be1f2ec7d603b4ba:v1.5.0`

Pulls the Talos `installer` image with the specified flavor and Talos version.
The image platform (architecture) will be determined by the architecture of the Talos Linux machine.

## Development

Run integration tests in local mode, with registry mirrors:

```bash
make integration TEST_FLAGS="-test.image-registry=127.0.0.1:5004 -test.flavor-service-repository=127.0.0.1:5005/image-service/flavor -test.installer-external-repository=127.0.0.1:5005/test -test.installer-internal-repository=127.0.0.1:5005/test" REGISTRY=127.0.0.1:5005
```
