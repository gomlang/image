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
