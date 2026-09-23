# draw

`draw(destination, region, source, source_origin, op)` maps the source origin
to the minimum corner of the destination region. It clips against destination
and source bounds and returns the number of pixels written. `copy` is the
`Op::Src` shortcut; `over` is the `Op::Over` shortcut. Pixels outside the source
are skipped, leaving the destination unchanged. An empty clipped region is a
successful zero-pixel draw. Source origins outside ±1,000,000 are rejected.

`Src` replaces destination pixels, including fully transparent source pixels.
`Over` performs Porter-Duff source-over on **premultiplied encoded sRGB**
channels: each output channel is `source + round(destination * (255 -
source_alpha) / 255)`, with halfway values rounded upward. This matches the
library's explicit 8-bit pixel representation; it is not a linear-light
composite. The source region is read into bounded temporary storage before
any writes, so overlapping self-copy and self-over use original source pixels.
