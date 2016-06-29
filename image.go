package image

import (
	"bytes"
	"fmt"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"image"
	"image/color"
	"image/jpeg"
)

// EncodeJPEG converts a raw image to a jpeg image.
func EncodeJPEG(width, height int, colorModel string, img data.Blob, quality int) (data.Blob, error) {
	r := rawImage{
		w: width,
		h: height,
		i: img,
	}

	var i image.Image
	switch colorModel {
	case "rgb":
		i = &rgbRawImage{r}
	case "bgr":
		i = &bgrRawImage{r}
	default:
		return nil, fmt.Errorf("unknown colorModel: %v", colorModel)
	}

	w := bytes.NewBuffer(nil)
	if err := jpeg.Encode(w, i, &jpeg.Options{
		Quality: quality,
	}); err != nil {
		return nil, err
	}
	return data.Blob(w.Bytes()), nil
}

type rawImage struct {
	w, h int
	i    data.Blob
}

func (r *rawImage) ColorModel() color.Model {
	return color.NRGBAModel
}

func (r *rawImage) Bounds() image.Rectangle {
	return image.Rectangle{
		Min: image.Point{},
		Max: image.Pt(r.w, r.h),
	}
}

type rgbRawImage struct {
	rawImage
}

func (i *rgbRawImage) At(x, y int) color.Color {
	if x < 0 || i.w <= x || y < 0 || i.h <= y {
		return color.NRGBA{}
	}
	off := (y*i.w + x) * 3
	return color.NRGBA{
		R: i.i[off],
		G: i.i[off+1],
		B: i.i[off+2],
		A: 255,
	}
}

type bgrRawImage struct {
	rawImage
}

func (i *bgrRawImage) At(x, y int) color.Color {
	if x < 0 || i.w <= x || y < 0 || i.h <= y {
		return color.NRGBA{}
	}
	off := (y*i.w + x) * 3
	return color.NRGBA{
		R: i.i[off+2],
		G: i.i[off+1],
		B: i.i[off],
		A: 255,
	}
}

// DecodeJPEG decodes a JPEG image to a raw image. The color model is in bgr
// for OpenCV integration.
func DecodeJPEG(jpg data.Blob) (data.Map, error) {
	r := bytes.NewReader([]byte(jpg))
	img, err := jpeg.Decode(r)
	if err != nil {
		return nil, err
	}

	// TODO: provide a faster version for a specific Image type.
	b := img.Bounds()
	width, height := b.Size().X, b.Size().Y
	buf := make([]byte, width*height*3)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := img.At(x, y)
			r, g, b, _ := c.RGBA()
			off := (y*width + x) * 3
			buf[off+0] = byte(b)
			buf[off+1] = byte(g)
			buf[off+2] = byte(r)
		}
	}
	return data.Map{
		"width":       data.Int(width),
		"height":      data.Int(height),
		"format":      data.String("raw"),
		"color_model": data.String("bgr"),
		"image":       data.Blob(buf),
	}, nil
}
