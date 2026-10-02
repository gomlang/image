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
management. The encoder writes one noninterlaced 8-bit RGBA `IDAT` stream. For
each row it tries all five PNG filters and chooses the smallest sum of absolute
signed-byte residuals, breaking ties by the lowest filter number. This
deterministic heuristic improves compression of smooth and repeated rows without
changing decoded pixels; it does not guarantee the smallest compressed stream.
Filtering uses bounded row storage and a fixed five passes per row. The
`max_work` budget applies to zlib, while dimensions bound filter work.
It unpremultiplies pixels before writing and
requires a nonempty image at origin `(0, 0)`. This conversion can be lossy for
translucent pixels.

`encode_with_options(image, limits, EncodeOptions)` selects compression `level`,
`filter` and `encoding`. `EncodeOptions::standard()` preserves level 6, adaptive
filtering and RGBA8. `Encoding::{Rgba8, Rgb8, Grayscale8, GrayscaleAlpha8}` maps
to PNG color types 6, 2, 0 and 4 with eight-bit samples. RGB and grayscale without
alpha require every pixel to be opaque. Both grayscale modes require equal
unpremultiplied RGB channels; incompatible pixels return a codec error instead of
dropping alpha or converting colors. Existing premultiplied-to-straight rounding
still applies. Scanline budgets use the chosen channel count.

`Filter::{None, Sub, Up, Average, Paeth}` uses one fixed PNG filter on every row,
providing a cheaper alternative to the five-pass `Adaptive` heuristic. A fixed
filter can produce larger or smaller files depending on the image. Predictor
distances follow the encoding's channel count. `PngCodec::with_filter` and
`with_encoding` select the same policies for the `Codec` trait; decoding remains
independent of encoding preferences. Palette, low-bit-depth, 16-bit and interlaced
encoding are not implemented.

The tests use deterministic independently assembled PNG fixtures under
[`tests/data`](tests/data/README.md), covering every accepted bit-depth/color
type combination, filters, Adam7, transparency, chunk placement, CRC failures
and budgets. From the library root, `(cd ../verification && just ecosystem-test image)`
also exercises the public codec through the example and independent downstream checks.

Format and decoding rules follow the [W3C PNG Specification, Third
Edition](https://www.w3.org/TR/png-3/).
