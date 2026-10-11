# image

`ecosystem::image` provides bounded, checked 8-bit premultiplied RGBA images.
It depends on `ecosystem::color` for `Rgba8`. The nested
[`draw`](draw/README.md) package supplies clipped copy and source-over,
[`png`](png/README.md) provides a bounded PNG codec,
[`jpeg`](jpeg/README.md) adds JPEG decoding/encoding,
[`webp`](webp/README.md) adds still WebP decoding and lossless encoding, and
[`transform`](transform/README.md) provides crop, nearest/bilinear resize, and lossless orientation.

```toml
[dependencies]
"ecosystem::image" = true
"ecosystem::color" = true
```

`Rect::new(min_x, min_y, max_x, max_y)` uses half-open coordinates. Each
coordinate must lie within ±1,000,000, width and height must be nonnegative
and at most 4,096. Empty rectangles are valid. `intersection`, `contains`,
`is_empty`, dimension and coordinate accessors work without allocation.

`Image::new(bounds)` initializes transparent pixels; `filled(bounds, pixel)`
and `from_pixels(bounds, pixels)` construct other images. The pixel count is
limited to 1,048,576. `from_pixels` checks the exact count and copies caller
storage. `at(x, y)` returns `None` outside bounds; `set(x, y, pixel)` returns
`OutOfBounds`. `pixels()` and `snapshot()` return independent storage.
`subimage(region)` intersects its rectangle with the image and returns a
**copied** image retaining the intersected coordinates. Editing either image
afterward does not alter the other.

All stored pixels are `Rgba8`: encoded sRGB red, green and blue already
multiplied by 8-bit alpha, with each channel at most alpha. The constructor in
`ecosystem::color` enforces that invariant. These pixel operations do not
convert to linear light; use `Srgb.composite_over` for linear-light blending.

`Codec` is a public trait with `decode(Bytes) -> Result[Image, Error]` and
`encode(Image) -> Result[Bytes, Error]`, so applications can provide format
codecs without changing the image core. `RawRgbaCodec` is a deterministic
uncompressed example. Its `R8PM` format has a 4-byte ASCII magic, little-endian
32-bit width and height, then row-major premultiplied RGBA8 bytes. It requires
origin `(0, 0)` on encode. Decode rejects invalid dimensions, size mismatches,
trailing data and channels above alpha before constructing an image. This is
an internal interchange format, not PNG, JPEG or a registered media type.

`(cd ../workflows && just ecosystem-test image)` runs library and example format,
build, tests, cached-build and executable checks.

## Native adapter setup

JPEG uses the Go standard `image/jpeg` implementation. WebP uses the pinned
`golang.org/x/image` v0.46.0 decoder and `github.com/HugoSmits86/nativewebp` v1.3.0
lossless encoder. Go 1.26 or newer is required; cgo and system image libraries
are not required. The native adapter is declared at module level, so **all
consuming modules, including consumers using only PNG or transforms, must have
a minimal root `go.mod`**:

```go
module example.com/my-image-app

go 1.26.0
```

GoML resolves the adapter from the selected image module and manages its Go
requirement and source replacement. It does not modify the consumer's `go.mod`.
External Go dependencies must be downloaded before offline FFI validation. For
this repository, run `go mod download all`. CI prepares native dependencies
automatically; the local ecosystem verifier expects the cache to be populated. A new consumer can prime the shared Go
module cache using:

```sh
go mod download github.com/HugoSmits86/nativewebp@v1.3.0 golang.org/x/image@v0.46.0
```

The example shares the root manifest and native Go module.
Generated FFI bindings are owned by `bindings.json`; regenerate them with
`goml bind-go bindings.json`. Do not edit generated files.

The new codecs share [`codecs::Limits`](codecs/README.md). They preserve the
existing `Image`, `Rect`, `Error`, `Codec`, draw and PNG APIs. JPEG explicitly
flattens alpha against a configurable background; WebP preserves stored
premultiplied bytes with lossless encoding. Metadata and color management are
outside the current API.

## Development and examples

Requires the [current GoML toolchain](https://github.com/gomlang/workflows/blob/main/ci/toolchain.json) with unversioned registry support. The `examples/basic/` example shares the root manifest and its dependencies. From the library root, run:

```sh
goml run --example basic
goml test
go test ./adapter
(cd ../workflows && just ecosystem-test image)
```

`goml test` builds the example and runs its tests. `(cd ../workflows && just ecosystem-test image)` runs the library-specific smoke and compatibility checks.
