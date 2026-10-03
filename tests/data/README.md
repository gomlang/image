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
