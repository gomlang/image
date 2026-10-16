# Third-party notices

The JPEG implementation in `jpeg/` is adapted from Go 1.26 `image/jpeg`.
Copyright 2009, 2011, 2012, 2025 The Go Authors. All rights reserved.

The VP8, VP8L and WebP decoders in `webp/` are adapted from
`golang.org/x/image` v0.46.0 (`vp8`, `vp8l`, `webp`).
Copyright 2011, 2014 The Go Authors. All rights reserved.

The original copyright notices and redistribution conditions are preserved in
[the Go Authors BSD license](licenses/GO-AUTHORS-BSD.txt). These implementations
were translated to GoML; no Go source is linked into production codecs.

The lossless WebP encoder follows the WebP lossless bitstream specification and
adapts coding techniques from `github.com/HugoSmits86/nativewebp` v1.3.0.
Copyright (c) 2024 Hugo Smits. Its MIT terms are preserved in
[the nativewebp license](licenses/NATIVEWEBP-MIT.txt).
The GoML encoder uses canonical Huffman coding, LZ77 references and spatial
prediction.
