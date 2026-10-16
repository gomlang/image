# JPEG reference fixtures

The `video-*.jpeg` files are copied unchanged from Go 1.26's
`src/image/testdata`. They cover baseline and progressive coding, all six Go
chroma sampling ratios, separate DC progression, partial progressive images,
CMYK, and grayscale frames with both conventional and redundant sampling factors.
They are distributed under the [Go Authors BSD license](../../../licenses/GO-AUTHORS-BSD.txt).
Copyright 2009–2025 The Go Authors. All rights reserved.

The expected RGBA CRC-32 values in `../jpeg_test.goml` were calculated by decoding
these files with Go 1.26 `image/jpeg`, visiting pixels in row order, and taking
the high byte of each channel returned by `Color.RGBA()`. Baseline/progressive
pairs intentionally share expected checksums. The implementation also passed
pixel-by-pixel comparison against the complete 26-file Go JPEG corpus.

`generate.py` uses only Python's standard library to create original eight-bit
JPEG streams with constant DCT blocks and deliberately simple Huffman tables:

- `rgb.jpg` stores direct RGB `(32, 96, 160)` with Adobe transform 0.
- `cmyk.jpg` stores inverted CMYK `(40, 80, 160, 200)`, yielding RGB `(31, 62, 125)`.
- `ycck.jpg` stores YCCK `(96, 128, 128, 200)`, yielding RGB `(125, 125, 125)`.
- `quantization-change.jpg` redefines table 0 between separate component scans.
- `quantization-wide.jpg` uses the same scans with sixteen-bit quantization entries.
- `restart.jpg` resets the grayscale DC predictor between two MCUs.

The quantization and restart fixtures decode to uniform gray 64. All original
fixtures have also been checked with Go 1.26's decoder. Regenerate them with
`python3 jpeg/tests/data/generate.py` from the module root.

`encoded-{1,25,85,100}.jpg` are Go 1.26 encoder outputs for the original 19 by 17
premultiplied RGBA pattern defined by `encoding_source` in `../jpeg_test.goml`,
composited onto RGB `(32, 96, 160)` using `(background * (255-alpha) + 127) / 255`.
The four filenames select the JPEG quality. Tests require byte-for-byte matching,
including 4:2:0 sampling, clipped edge MCUs, quantization, and entropy coding.

The original generated fixtures and generator use this repository's MIT license.
No external codec, Python installation, or Go source is required to run the tests.
