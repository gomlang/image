# transform

`ecosystem::image::transform` operates on the existing premultiplied image model
without native code.

`crop(source, region)` requires the complete half-open rectangle to be contained
in the source bounds. It returns an independent copy with origin `(0, 0)` and
returns `OutOfBounds` for an outside region. An empty rectangle contained within
the source, including at its upper boundary, produces an empty image. The existing
`Image.subimage(region)` remains the clipped alternative and retains source
coordinates.

`resize(source, width, height, Filter::{Nearest, Bilinear})` returns an independent
image at origin `(0, 0)`. Negative dimensions, dimensions above 4,096 and pixel
counts above 1,048,576 fail before allocation. An empty destination is valid; an
empty source cannot produce a nonempty destination.

Both filters map destination pixel centers into source coordinates. Nearest
selects the containing source pixel, breaking exact midpoint ties toward the
larger coordinate. Bilinear interpolates four neighbors, clamps samples at source
edges and rounds each final channel to the nearest byte, with halves upward.
Integer weights avoid floating-point nondeterminism. Red, green, blue and alpha
are interpolated in their premultiplied form so transparency does not introduce
dark color fringes. All calculations use encoded sRGB values, without linear-light
conversion. Bilinear resize is not an area or antialiasing filter for substantial
downsampling; aspect ratio is selected by the caller.

## Lossless orientation

`orient(source, Orientation)` returns an independent image at `(0, 0)` without
resampling or changing any premultiplied RGBA channel. Rotations are clockwise
in image coordinates (x grows right, y grows down). The eight operations cover
all quarter-turn and mirror combinations:

| Orientation | Result for rows `AB / CD / EF` |
| --- | --- |
| `Identity` | `AB / CD / EF` |
| `FlipHorizontal` | `BA / DC / FE` |
| `Rotate180` | `FE / DC / BA` |
| `FlipVertical` | `EF / CD / AB` |
| `Transpose` | `ACE / BDF` |
| `Rotate90` | `ECA / FDB` |
| `Transverse` | `FDB / ECA` |
| `Rotate270` | `BDF / ACE` |

The last four swap width and height; this also preserves the shape of zero-width
or zero-height images after swapping. Source coordinate offsets are normalized,
and even `Identity` copies storage. Work and additional storage are O(pixels).
Existing image dimension/pixel bounds apply unchanged because the operation
preserves pixel count. The function does not parse or apply EXIF metadata;
callers choose the orientation explicitly.
