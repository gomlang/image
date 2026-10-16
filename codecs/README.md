# codecs

`ecosystem::image::codecs::Limits` provides shared allocation and input/output
limits for the GoML JPEG and WebP codecs. PNG retains its existing limits.

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
or elapsed CPU time. The codecs allocate image-sized working buffers and bounded
entropy tables. Encoders check their output budget while writing compressed
bytes. There is no cancellation or exact workspace-byte budget. Applications
processing untrusted images should select small limits and constrain concurrent
codec calls. Invalid streams return recoverable codec errors; process-level
out-of-memory failures are outside that guarantee.

The legacy `pack`, `unpack` and `encoded` helpers retain their public signatures
using `std::ffi::String` for compatibility. They convert raw pixel and error
buffers; the JPEG and WebP implementations do not use them or external bindings.
