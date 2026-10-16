# Independent JPEG and WebP fixtures

`generate.py` creates these original, tiny fixtures through Pillow 12.3.0,
using its independent libjpeg backend (reported JPEG API version 6.2) and libwebp
1.6.0. Pillow is used only to regenerate fixtures; tests read the checked-in files
without Python or system codec dependencies.

- `solid.jpg`: 8 by 8 RGB `(32, 96, 160)`, quality 100, no chroma subsampling.
- `progressive.jpg`: the same pixels in a progressive JPEG.
- `gray.webp`: 8 by 8 gray 128, lossy WebP quality 100.
- `color.webp`, `black.webp`, `white.webp`: lossy 8 by 8 images checking
  BT.601 limited-range conversion against libwebp colors `(32,96,160)`, black
  and white. `color-alpha.webp` uses the same color with alpha 128.
- `alpha.webp`: 2 by 2 lossless pixels with straight RGBA values
  `(255,0,0,255)`, `(0,255,0,128)`, `(0,0,255,64)`, `(255,255,255,0)`.

JPEG fixture comparisons allow two byte values of codec rounding. The lossy
gray fixture also allows two; the lossless fixture checks exact premultiplied
pixels including discarded hidden color at alpha zero. A copy of `alpha.webp`
lives in the independent consumer example. These fixtures are covered by the
repository MIT license.

To regenerate in an environment with Pillow 12.3.0 installed:

```sh
python3 tests/data/generate.py
```

`generate_lossless.py` creates original patterned RGBA images and encodes them
with FFmpeg's independent libwebp lossless encoder. `lossless-gradient.webp`
(37×29), `lossless-palette.webp` (31×25), `lossless-noise.webp` (19×17) and
`lossless-predictor.webp` (65×51) exercise transform and entropy choices with
odd dimensions, alpha, color indexing and nonuniform pixels. Matching `.rgba`
files contain exact premultiplied reference pixels, computed from the original
straight RGBA source without a decoder. Tests check these files and all seven
encoding efforts without external processes. Regeneration requires `ffmpeg`
with the `libwebp` encoder; runtime tests do not.

`lossy-rounding.webp` is the original 17×19 pattern also embedded in the internal
VP8 regression test. libwebp's `WebPEncodeRGB` encodes quality 65 from RGB
`((17*x+3*y)%256, (5*x+19*y)%256, (13*x+11*y)%256)`.
`lossy-rounding-alpha.webp` wraps the same VP8 payload in a VP8X container with
uncompressed alpha `(11*x+13*y)%256`. Their `.rgba` files are fixed results from
the original Go adapter using `golang.org/x/image` v0.46.0, including its final
14-bit BT.601 conversion, nearest 4:2:0 chroma and premultiplication. Pixel `(9,0)`
in the opaque image is exactly `(161,50,108,255)`; lower precision 8-bit
conversion incorrectly yields red 162. Tests compare every output channel,
including the alpha case, without executing an external codec.
