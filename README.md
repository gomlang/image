# image

`ecosystem::image` provides bounded, checked 8-bit premultiplied RGBA images.
It depends on `ecosystem::color` for `Rgba8`. The nested
[`draw`](draw/README.md) package supplies clipped copy and source-over,
and [`png`](png/README.md) provides a bounded PNG codec.

```toml
[dependencies]
"ecosystem::image" = "0.1.0"
"ecosystem::color" = "0.1.0"
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

`just ecosystem-test image` runs library and independent consumer format,
build, tests, cached-build and executable checks.
