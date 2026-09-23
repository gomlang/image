# png

`ecosystem::image::png` decodes bounded PNG datastreams into the parent
module's premultiplied `Rgba8` images. It implements the parent `Codec` trait
through `PngCodec`; `decode(bytes)` and `encode(image)` use standard limits and
compression level 6. `PngCodec::new(limits, level)` accepts levels 0 through 9.

Decode supports all standard color types and their legal bit depths:

| Color type | Samples | Bit depths |
| --- | --- | --- |
| 0 | grayscale | 1, 2, 4, 8, 16 |
| 2 | RGB | 8, 16 |
| 3 | indexed palette | 1, 2, 4, 8 |
| 4 | grayscale and alpha | 8, 16 |
| 6 | RGBA | 8, 16 |

It handles `PLTE`, `tRNS`, all five PNG row filters and both noninterlaced and
Adam7 data. `IDAT` may span consecutive chunks. Every chunk's CRC32 is checked
using `std::hash::crc32`, and its compressed data is decoded through
`ecosystem::compress::zlib`. Unknown critical chunks, invalid chunk names,
duplicate or misplaced required chunks, nonconsecutive `IDAT`, trailing bytes,
invalid palette indexes and malformed filter rows are rejected. Unknown
ancillary chunks are CRC-checked and ignored.

`Limits` bounds input bytes, pixels, chunk count, aggregate compressed bytes,
inflated scanline bytes and zlib work. `standard()` allows 8 MiB input and
compressed bytes, 16 MiB inflated bytes, 1,048,576 pixels, 4,096 chunks and
100,000,000 zlib work units. Image dimensions are also capped at 4,096 per
axis by `Image`. Dimensions and expected inflated size are checked before
zlib allocation; exact inflated length is checked after decoding.

PNG stores **straight** alpha. The decoder rescales each sample to eight bits
and then premultiplies through `Nrgba8.to_rgba8()`. Thus 16-bit precision and
hidden RGB at zero alpha are lost in the resulting `Image`. For `tRNS`, the
unscaled source sample is compared with its transparency key before rescaling.
Metadata chunks such as `gAMA`, `iCCP`, `sRGB`, text and animation extensions
are ignored; color values are treated as encoded sRGB bytes without color
management. The encoder writes one noninterlaced 8-bit RGBA `IDAT` stream with
filter type 0 for each row. It unpremultiplies pixels before writing and
requires a nonempty image at origin `(0, 0)`. This conversion can be lossy for
translucent pixels.

The tests use deterministic independently assembled PNG fixtures under
[`tests/data`](tests/data/README.md), covering every accepted bit-depth/color
type combination, filters, Adam7, transparency, chunk placement, CRC failures
and budgets. `just ecosystem-test image` also exercises the published codec
from an independent consumer.

Format and decoding rules follow the [W3C PNG Specification, Third
Edition](https://www.w3.org/TR/png-3/).
