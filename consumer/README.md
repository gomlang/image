# image consumer

An independent module exercising the published `ecosystem::image` and
`ecosystem::image::draw` API. It blends premultiplied pixels, serializes an
image through both the `RawRgbaCodec` and PNG `Codec` implementations, decodes
them and checks the results.
