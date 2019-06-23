package pixelview

import (
    "os"
    "io"
    "fmt"
    "image"
    "image/color"
    "github.com/pkg/errors"
)


func FromFile(filename string) (encoded string, err error) {
    f, err := os.Open(filename)
    if err != nil {
        return
    }
    return FromReader(io.Reader(f))
}


func FromReader(reader io.Reader) (encoded string, err error) {
    img, _, err := image.Decode(reader)
    if err != nil {
        return
    }
    return FromImage(img)
}


func FromImage(img image.Image) (encoded string, err error) {
    if (img.Bounds().Max.Y - img.Bounds().Min.Y) % 2 != 0 {
        err = errors.New("pixelview: Can't process image with uneven height")
        return
    }

    for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y += 2 {
        var prevfg, prevbg color.Color
        for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
            fg := img.At(x, y)
            bg := img.At(x, y + 1)
            if fg == prevfg && bg == prevbg {
                encoded += "▀"
                continue
            }
            if fg == prevfg {
                encoded += fmt.Sprintf(
                    "[:%s]▀",
                    hexColour(bg),
                )
                prevbg = bg
                continue
            }
            if bg == prevbg {
                encoded += fmt.Sprintf(
                    "[%s:]▀",
                    hexColour(fg),
                )
                prevfg = fg
                continue
            }
            encoded += fmt.Sprintf(
                "[%s:%s]▀",
                hexColour(fg),
                hexColour(bg),
            )
            prevfg = fg
            prevbg = bg
        }
        encoded += "\n"
    }
    return
}


func FromPaletted(img *image.Paletted) (encoded string, err error) {
    if (img.Bounds().Max.Y - img.Bounds().Min.Y) % 2 != 0 {
        err = errors.New("pixelview: Can't process image with uneven height")
        return
    }

    for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y += 2 {
        for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
            i := (y - img.Rect.Min.Y) * img.Stride + (x - img.Rect.Min.X)
            fg := img.Pix[i]
            bg := img.Pix[i + img.Stride]
            // TODO: RLE
            encoded += fmt.Sprintf(
                "[%s:%s]▀",
                hexColour(img.Palette[fg]),
                hexColour(img.Palette[bg]),
            )
        }
        encoded += "\n"
    }
    return
}


func hexColour(c color.Color) string {
    r, g, b, _ := c.RGBA()
    return fmt.Sprintf("#%.2x%.2x%.2x", r >> 8, g >> 8, b >> 8)
}

