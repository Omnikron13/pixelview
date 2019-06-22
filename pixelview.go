package pixelview

import (
    "os"
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

    img, _, err := image.Decode(f)
    if err != nil {
        return
    }

    if (img.Bounds().Max.Y - img.Bounds().Min.Y) % 2 != 0 {
        errors.New("pixelview: Can't process image with uneven height")
    }

    for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y += 2 {
        for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
            encoded += fmt.Sprintf(
                "[%s:%s]▀",
                hexColour(img.At(x, y)),
                hexColour(img.At(x, y + 1)),
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

