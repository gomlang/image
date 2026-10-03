# jpeg

`ecosystem::image::jpeg` decodes JPEG and encodes baseline JPEG through Go's
[`image/jpeg`](https://pkg.go.dev/image/jpeg) backend. Baseline and progressive
JPEG decoding is supported; arithmetic-coded, lossless and higher-bit-depth
JPEG variants follow the backend's unsupported-format errors.

`decode(bytes)`, `decode_with_limits(bytes, codecs::Limits)` and `encode(image)`
return the parent image package's `Result` and `Error` types. `JpegCodec`
implements `image::Codec`. `JpegCodec::standard()` uses quality 85 and white as
the alpha background. `JpegCodec::new(limits, quality, background)` accepts
quality 1 through 100 and requires an opaque `color::Rgba8` background.
`limits()` and `quality()` expose the selected settings.

JPEG cannot carry alpha. Before encoding, premultiplied pixels are composited
against the selected background in encoded sRGB byte space, without linear-light
conversion. Decoded pixels are opaque. Source rectangles may have any supported
origin; decoded images always start at `(0, 0)`. Empty images cannot be encoded.
JPEG encoding is lossy, including at quality 100.

The shared [`codecs::Limits`](../codecs/README.md) checks the encoded input and
header dimensions before decoding or allocating the output pixel buffer. Output
writes are bounded. EXIF orientation is not applied; EXIF, ICC, comments and other
metadata are not preserved. Color management is not performed. The backend can
accept trailing data after a completed JPEG; use a container-level framing check
when a protocol requires an exact datastream boundary.

See the root README for native adapter setup; Go 1.26 is required, with no C
compiler or system codec library requirement.
