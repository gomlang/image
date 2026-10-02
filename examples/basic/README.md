# image example

An example exercising the public `ecosystem::image` and
`ecosystem::image::draw` API. It blends premultiplied pixels, serializes an
image through both the `RawRgbaCodec` and PNG `Codec` implementations, decodes
them and checks the results.

This example shares the library root manifest and its dependencies. From the library root, run `goml verify --example basic` to build and test it as an independent downstream module.
