package adapter

import (
	"bytes"
	"encoding/binary"
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"strings"

	"github.com/HugoSmits86/nativewebp"
	"golang.org/x/image/webp"
)

const maximumBytes = 64 << 20

func dimensions(width, height, dimension, pixels, decoded int) error {
	if dimension < 1 || dimension > 4096 || pixels < 0 || pixels > 1048576 || decoded < 0 || decoded > 4194304 {
		return errors.New("invalid image limits")
	}
	if width < 1 || height < 1 || width > dimension || height > dimension || width > pixels/height || width > decoded/4/height {
		return errors.New("image dimensions or decoded byte limit exceeded")
	}
	return nil
}

func recoverDecode(data *string, width, height *int, failure *string) {
	if recover() != nil {
		*data, *width, *height, *failure = "", 0, 0, "invalid image: codec rejected input"
	}
}

func recoverEncode(data *string, failure *string) {
	if recover() != nil {
		*data, *failure = "", "image codec failed"
	}
}

func decode(data string, maxInput, maxDimension, maxPixels, maxDecoded int, isWebP bool) (string, int, int, string) {
	if maxInput < 0 || maxInput > maximumBytes || len(data) > maxInput {
		return "", 0, 0, "image input byte limit exceeded"
	}
	var config image.Config
	var err error
	if isWebP {
		if err = validateWebP(data); err == nil {
			config, err = webp.DecodeConfig(strings.NewReader(data))
		}
	} else {
		config, err = jpeg.DecodeConfig(strings.NewReader(data))
	}
	if err != nil {
		return "", 0, 0, err.Error()
	}
	if err = dimensions(config.Width, config.Height, maxDimension, maxPixels, maxDecoded); err != nil {
		return "", 0, 0, err.Error()
	}
	var decoded image.Image
	if isWebP {
		decoded, err = webp.Decode(strings.NewReader(data))
	} else {
		decoded, err = jpeg.Decode(strings.NewReader(data))
	}
	if err != nil {
		return "", 0, 0, err.Error()
	}
	bounds := decoded.Bounds()
	if bounds.Dx() != config.Width || bounds.Dy() != config.Height {
		return "", 0, 0, "decoded dimensions differ from image header"
	}
	pixels := make([]byte, 0, bounds.Dx()*bounds.Dy()*4)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			sample := color.NRGBAModel.Convert(decoded.At(x, y)).(color.NRGBA)
			alpha := uint32(sample.A)
			pixels = append(pixels, byte((uint32(sample.R)*alpha+127)/255),
				byte((uint32(sample.G)*alpha+127)/255), byte((uint32(sample.B)*alpha+127)/255), sample.A)
		}
	}
	return string(pixels), bounds.Dx(), bounds.Dy(), ""
}

func DecodeJPEG(data string, maxInput, maxDimension, maxPixels, maxDecoded int) (pixels string, width, height int, failure string) {
	defer recoverDecode(&pixels, &width, &height, &failure)
	return decode(data, maxInput, maxDimension, maxPixels, maxDecoded, false)
}

func DecodeWebP(data string, maxInput, maxDimension, maxPixels, maxDecoded int) (pixels string, width, height int, failure string) {
	defer recoverDecode(&pixels, &width, &height, &failure)
	return decode(data, maxInput, maxDimension, maxPixels, maxDecoded, true)
}

func validateWebP(data string) error {
	invalid := errors.New("invalid or animated WebP container")
	if len(data) < 12 || data[:4] != "RIFF" || data[8:12] != "WEBP" || uint64(binary.LittleEndian.Uint32([]byte(data[4:8])))+8 != uint64(len(data)) {
		return invalid
	}
	frames := 0
	for offset := 12; offset < len(data); {
		if len(data)-offset < 8 {
			return invalid
		}
		kind := data[offset : offset+4]
		size := uint64(binary.LittleEndian.Uint32([]byte(data[offset+4 : offset+8])))
		end := uint64(offset) + 8 + size
		padded := end + size%2
		if padded > uint64(len(data)) {
			return invalid
		}
		if size%2 != 0 && data[int(end)] != 0 {
			return invalid
		}
		if kind == "ANIM" || kind == "ANMF" {
			return invalid
		}
		if kind == "VP8X" && (size != 10 || data[offset+8]&2 != 0) {
			return invalid
		}
		if kind == "VP8 " || kind == "VP8L" {
			frames++
		}
		offset = int(padded)
	}
	if frames != 1 {
		return invalid
	}
	return nil
}

type boundedWriter struct {
	buffer bytes.Buffer
	limit  int
	err    error
}

func (w *boundedWriter) Write(data []byte) (int, error) {
	if w.err != nil {
		return 0, w.err
	}
	if len(data) > w.limit-w.buffer.Len() {
		w.err = errors.New("image output byte limit exceeded")
		return 0, w.err
	}
	return w.buffer.Write(data)
}

func sourceImage(data string, width, height, maxDimension, maxPixels, maxDecoded int) (*image.RGBA, error) {
	if err := dimensions(width, height, maxDimension, maxPixels, maxDecoded); err != nil {
		return nil, err
	}
	if len(data) != width*height*4 {
		return nil, errors.New("invalid image pixel count")
	}
	for i := 0; i < len(data); i += 4 {
		if data[i] > data[i+3] || data[i+1] > data[i+3] || data[i+2] > data[i+3] {
			return nil, errors.New("invalid premultiplied pixel")
		}
	}
	return &image.RGBA{Pix: []byte(data), Stride: 4 * width, Rect: image.Rect(0, 0, width, height)}, nil
}

func EncodeJPEG(data string, width, height, maxDimension, maxPixels, maxDecoded, maxOutput, quality int, background uint32) (encoded string, failure string) {
	defer recoverEncode(&encoded, &failure)
	if maxOutput < 0 || maxOutput > maximumBytes || quality < 1 || quality > 100 || background&255 != 255 {
		return "", "invalid JPEG encoding options"
	}
	img, err := sourceImage(data, width, height, maxDimension, maxPixels, maxDecoded)
	if err != nil {
		return "", err.Error()
	}
	bg := color.RGBA{R: byte(background >> 24), G: byte(background >> 16), B: byte(background >> 8), A: 255}
	for i := 0; i < len(img.Pix); i += 4 {
		inverse := uint32(255 - img.Pix[i+3])
		img.Pix[i] += byte((uint32(bg.R)*inverse + 127) / 255)
		img.Pix[i+1] += byte((uint32(bg.G)*inverse + 127) / 255)
		img.Pix[i+2] += byte((uint32(bg.B)*inverse + 127) / 255)
		img.Pix[i+3] = 255
	}
	output := &boundedWriter{limit: maxOutput}
	if err = jpeg.Encode(output, img, &jpeg.Options{Quality: quality}); err != nil {
		return "", err.Error()
	}
	if output.err != nil {
		return "", output.err.Error()
	}
	return output.buffer.String(), ""
}

func EncodeWebP(data string, width, height, maxDimension, maxPixels, maxDecoded, maxOutput, effort int) (encoded string, failure string) {
	defer recoverEncode(&encoded, &failure)
	if maxOutput < 0 || maxOutput > maximumBytes || effort < 0 || effort > 6 {
		return "", "invalid WebP encoding options"
	}
	img, err := sourceImage(data, width, height, maxDimension, maxPixels, maxDecoded)
	if err != nil {
		return "", err.Error()
	}
	straight := image.NewNRGBA(img.Bounds())
	for i := 0; i < len(img.Pix); i += 4 {
		alpha := uint32(img.Pix[i+3])
		straight.Pix[i+3] = byte(alpha)
		if alpha != 0 {
			for channel := 0; channel < 3; channel++ {
				straight.Pix[i+channel] = byte((uint32(img.Pix[i+channel])*255 + alpha/2) / alpha)
			}
		}
	}
	output := &boundedWriter{limit: maxOutput}
	if err = nativewebp.Encode(output, straight, &nativewebp.Options{CompressionLevel: nativewebp.CompressionLevel(effort)}); err != nil {
		return "", err.Error()
	}
	if output.err != nil {
		return "", output.err.Error()
	}
	if output.buffer.Len() == 0 {
		return "", io.ErrUnexpectedEOF.Error()
	}
	return output.buffer.String(), ""
}
