# codecs

`ecosystem::image::codecs::Limits` provides shared allocation and input/output
limits for the native JPEG and WebP adapters. PNG retains its existing limits.

| Field | Standard value | Maximum accepted value |
| --- | --- | --- |
| `max_input_bytes` | 8 MiB | 64 MiB |
| `max_output_bytes` | 8 MiB | 64 MiB |
| `max_dimension` | 4,096 | 4,096 |
| `max_pixels` | 1,048,576 | 1,048,576 |
| `max_decoded_bytes` | 4,194,304 | 4,194,304 |

All limits may be reduced; byte and pixel limits may be zero. A dimension limit
must be at least one. Invalid limits are recoverable errors. Decoded bytes count
the final packed RGBA8 representation, four bytes per pixel. Header dimensions
are checked against dimension, pixel and decoded-byte limits before full decoding.

These limits bound input, final output and pixel counts, not total process memory
or elapsed CPU time. Backends allocate their own bounded-by-image-size working
buffers, and the Go/GoML bridge copies pixel buffers. In particular the WebP
encoder creates an internal compressed stream before the bounded output writer
receives it. There is no cancellation or exact workspace-byte budget. Applications
processing untrusted images should select small limits and constrain concurrent
codec calls. Unexpected backend panics are converted to recoverable codec errors;
process-level out-of-memory failures are outside that guarantee.

The `pack`, `unpack` and `encoded` functions implement the raw bridge used by the
codec packages; consumers normally use `jpeg` or `webp` rather than these helpers.
