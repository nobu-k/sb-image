package image

import (
	"bytes"
	"fmt"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"image"
	"image/color"
	"image/jpeg"
)

// ToJPEG converts a raw image to a jpeg image.
func ToJPEG(width, height int, colorModel string, img data.Blob, quality int) (data.Blob, error) {
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
