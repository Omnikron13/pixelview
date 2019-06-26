// Package pixelview is a simple package which converts images to text formatted for tview.
// It uses coloured unicode half-block characters (▀) to represent pixels.
package pixelview

import (
    "os"
    "io"
    "fmt"
    "image"
    "image/color"
    "github.com/pkg/errors"
)


// FromFile is a convenience function that converts a file on disk to a formatted string.
// See FromImage() for more details.
func FromFile(filename string) (encoded string, err error) {
    f, err := os.Open(filename)
    if err != nil {
        return
    }
    defer f.Close()
    return FromReader(io.Reader(f))
}


// FromReader is a convenience function that converts an io.Reader to a formatted string.
// See FromImage() for more details.
func FromReader(reader io.Reader) (encoded string, err error) {
    img, _, err := image.Decode(reader)
    if err != nil {
        return
    }
    return FromImage(img)
}


// FromImage is the primary function of this package,
// It takes an image.Image and converts it to a string formatted for tview.
// The unicode half-block character (▀) with a fg & bg colour set will represent
// pixels in the returned string.
// Because each character represents two pixels, it is not possible to convert an
// image if its height is uneven. Attempts to do so will return an error.
func FromImage(img image.Image) (encoded string, err error) {
    if (img.Bounds().Max.Y - img.Bounds().Min.Y) % 2 != 0 {
        err = errors.New("pixelview: Can't process image with uneven height")
        return
    }

    switch v := img.(type) {
    default:
        return fromImageGeneric(img)
    case *image.Paletted:
        return fromPaletted(v)
    case *image.NRGBA:
        return fromNRGBA(v)
    }
}


// fromImageGeneric is the fallback function for processing images.
// It will be used for more exotic image formats than png or gif.
func fromImageGeneric(img image.Image) (encoded string, err error) {
    for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y += 2 {
        var prevfg, prevbg color.Color
        for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
            fg := img.At(x, y)
            bg := img.At(x, y + 1)
            encoded += encode(fg, bg, &prevfg, &prevbg)
        }
        encoded += "\n"
    }
    return
}


// fromPaletted saves a few μs when working with paletted images.
// These are what PNG8 images are decoded as.
func fromPaletted(img *image.Paletted) (encoded string, err error) {
    for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y += 2 {
        var prevfg, prevbg color.Color
        for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
            i := (y - img.Rect.Min.Y) * img.Stride + (x - img.Rect.Min.X)
            fg := img.Palette[img.Pix[i]]
            bg := img.Palette[img.Pix[i + img.Stride]]
            encoded += encode(fg, bg, &prevfg, &prevbg)
        }
        encoded += "\n"
    }
    return
}


// fromNRGBA saves a handful of μs when working with NRGBA images.
// These are what PNG24 images are decoded as.
func fromNRGBA(img *image.NRGBA) (encoded string, err error) {
    for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y += 2 {
        var prevfg, prevbg color.Color
        for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
            i := (y - img.Rect.Min.Y) * img.Stride + (x - img.Rect.Min.X) * 4
            fg := color.NRGBA{img.Pix[i], img.Pix[i+1], img.Pix[i+2], img.Pix[i+3]}
            i += img.Stride
            bg := color.NRGBA{img.Pix[i], img.Pix[i+1], img.Pix[i+2], img.Pix[i+3]}
            encoded += encode(fg, bg, &prevfg, &prevbg)
        }
        encoded += "\n"
    }
    return
}


// encode converts a fg & bg colour into a formatted pair of 'pixels',
// using the prevfg & prevbg colours to perform something akin to run-length encoding
func encode(fg, bg color.Color, prevfg, prevbg *color.Color) (encoded string) {
    if fg == *prevfg && bg == *prevbg {
        encoded = "▀"
        return
    }
    if fg == *prevfg {
        encoded = fmt.Sprintf(
            "[:%s]▀",
            hexColour(bg),
        )
        *prevbg = bg
        return
    }
    if bg == *prevbg {
        encoded = fmt.Sprintf(
            "[%s:]▀",
            hexColour(fg),
        )
        *prevfg = fg
        return
    }
    encoded = fmt.Sprintf(
        "[%s:%s]▀",
        hexColour(fg),
        hexColour(bg),
    )
    *prevfg = fg
    *prevbg = bg
    return
}


func hexColour(c color.Color) string {
    r, g, b, _ := c.RGBA()
    return fmt.Sprintf("#%.2x%.2x%.2x", r >> 8, g >> 8, b >> 8)
}

