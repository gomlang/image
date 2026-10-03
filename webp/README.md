# webp

`ecosystem::image::webp` supports still WebP images: lossy VP8 and lossless VP8L
decoding, and lossless VP8L encoding. Extended VP8X still-image containers are
accepted, including alpha. Animated WebP is rejected instead of silently taking
one frame. Container lengths, padding and the presence of exactly one image
chunk are checked; trailing bytes are rejected.

The decoder is [`golang.org/x/image/webp`](https://pkg.go.dev/golang.org/x/image/webp)
v0.46.0, and the encoder is
[`github.com/HugoSmits86/nativewebp`](https://github.com/HugoSmits86/nativewebp)
v1.3.0. Both use Go without cgo or a system libwebp installation.

`decode(bytes)`, `decode_with_limits(bytes, codecs::Limits)` and `encode(image)`
return the parent `Result` and `Error` types. `WebpCodec` implements `image::Codec`.
`WebpCodec::standard()` selects encoding effort 4. `WebpCodec::new(limits, effort)`
accepts effort 0 through 6, trading encoder work for compressed size; all values
encode losslessly. `limits()` and `effort()` expose the selected settings. Lossy
WebP encoding is not implemented.

WebP straight-alpha colors are converted to the common premultiplied `Rgba8`
model on decode. Hidden RGB at zero alpha is discarded. Encoding selects straight
color values that reconstruct the input premultiplied bytes, preserving every
stored channel and alpha through lossless encode/decode. Bounds origins are not
stored and decode returns origin `(0, 0)`. Empty images cannot be encoded.

The shared [`codecs::Limits`](../codecs/README.md) checks input length and header
dimensions before full decode and limits the encoded output writer. The selected
decoder checks extended-canvas and frame dimensions agree before decoding a
VP8L frame, preventing a small canvas header from disguising a large allocation.
Metadata, EXIF orientation and ICC color management are not applied or preserved.
The Go backend's VP8 YCbCr conversion is used; it can differ slightly from libwebp
for colored lossy inputs. Use lossless encoding when exact stored pixels matter.
