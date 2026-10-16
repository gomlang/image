# webp

`ecosystem::image::webp` supports still WebP images: lossy VP8 and lossless VP8L
decoding, and lossless VP8L encoding. Extended VP8X still-image containers are
accepted, including alpha. Animated WebP is rejected instead of silently taking
one frame. Container lengths, padding and the presence of exactly one image
chunk are checked; trailing bytes are rejected.

The codecs are implemented in GoML. The decoder follows the Go Authors' VP8,
VP8L and WebP implementations; attribution is in
[third-party notices](../THIRD_PARTY_NOTICES.md). No native adapter or external
codec dependency is required.

`decode(bytes)`, `decode_with_limits(bytes, codecs::Limits)` and `encode(image)`
return the parent `Result` and `Error` types. `WebpCodec` implements `image::Codec`.
`WebpCodec::standard()` selects encoding effort 4. `WebpCodec::new(limits, effort)`
accepts effort 0 through 6. Effort 0 prioritizes speed, 4 uses spatial prediction, and 6 searches more
backwards references and predictor modes for compression. Values 1, 2, 3 and 5
select the same algorithm as 4. All values
encode losslessly. `limits()` and `effort()` expose the requested settings. Lossy
WebP encoding is not implemented.

WebP straight-alpha colors are converted to the common premultiplied `Rgba8`
model on decode. Hidden RGB at zero alpha is discarded. Encoding selects straight
color values that reconstruct the input premultiplied bytes, preserving every
stored channel and alpha through lossless encode/decode. Bounds origins are not
stored and decode returns origin `(0, 0)`. Empty images cannot be encoded.

The shared [`codecs::Limits`](../codecs/README.md) checks input length and header
dimensions before full decode and limits the encoded output writer. The
decoder checks extended-canvas and frame dimensions agree before decoding a
VP8L frame, preventing a small canvas header from disguising a large allocation.
Metadata, EXIF orientation and ICC color management are not applied or preserved.
The decoder converts VP8's limited-range BT.601 samples explicitly, following
[RFC 6386](https://www.rfc-editor.org/rfc/rfc6386.html#section-9.2), instead of
applying Go's JPEG full-range conversion. The 14-bit integer coefficients and
rounding preserve the original codec's exact output. Chroma sampling uses the decoded
4:2:0 pixel's nearest chroma sample; edge colors can differ from libwebp's
interpolated chroma upsampling. Use lossless encoding when exact stored pixels matter.
