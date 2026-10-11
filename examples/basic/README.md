# image example

An example exercising the public `ecosystem::image` and
`ecosystem::image::draw` API. It blends premultiplied pixels, serializes an
image through both the `RawRgbaCodec` and PNG `Codec` implementations, decodes
them and checks the results.

This example shares the library root manifest and its dependencies. From the library root, run `goml test` to build and test the library and its examples.

`thumbnail(source)` demonstrates strict cropping and bilinear resize. The example
then encodes JPEG (white alpha background) and lossless WebP. Its example
test loads a checked-in Pillow/libwebp fixture, checks corner pixels,
and transcodes the resized output into PNG and JPEG. The library-level native
adapter setup is described in the root README. The example shares the root
Go module.
