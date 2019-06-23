package pixelview

import (
    "testing"
    "os"
    "io/ioutil"
    "image"
    _ "image/png"
)


func TestFromFile(t *testing.T) {
    buf, err := ioutil.ReadFile("pixelview.raw")
    if err != nil {
        panic(err)
    }
    reference := string(buf)

    s, err := FromFile("pixelview.png")
    if err != nil {
        t.Errorf("Error encountered during execution: %s", err)
    }
    if s != reference {
        t.Error("Output did not match reference")
    }

    s, err = FromFile("pixelview_uneven.png")
    if err == nil {
        t.Error("Allowed uneven height image without error")
    }
    if s != "" {
        t.Errorf("Returned non-zero result when given uneven image: %s", s)
    }
}


func TestFromImageGeneric(t *testing.T) {
    buf, err := ioutil.ReadFile("pixelview.raw")
    if err != nil {
        panic(err)
    }
    reference := string(buf)

    f, err := os.Open("pixelview.png")
    if err != nil {
        panic(err)
    }
    img, _, err := image.Decode(f)
    if err != nil {
        panic(err)
    }
    s, err := fromImageGeneric(img)
    if err != nil {
        t.Errorf("Error encountered during execution: %s", err)
    }
    if s != reference {
        t.Error("Output did not match reference")
    }
}


func TestFromPaletted(t *testing.T) {
    buf, err := ioutil.ReadFile("pixelview.raw")
    if err != nil {
        panic(err)
    }
    reference := string(buf)

    f, err := os.Open("pixelview.png")
    if err != nil {
        panic(err)
    }
    img, _, err := image.Decode(f)
    if err != nil {
        panic(err)
    }
    paletted, ok := img.(*image.Paletted)
    if !ok {
        panic("Type assertion failed before test could be run")
    }
    s, err := FromPaletted(paletted)
    if err != nil {
        t.Errorf("Error encountered during execution: %s", err)
    }
    if s != reference {
        t.Error("Output did not match reference")
    }
}


// TODO: test FromReader() & FromImage() (especially things like sub-images?)


func BenchmarkFromImageGeneric(b *testing.B) {
    f, err := os.Open("pixelview.png")
    if err != nil {
        panic(err)
    }
    img, _, err := image.Decode(f)
    if err != nil {
        panic(err)
    }
    for n := 0; n < b.N; n++ {
        fromImageGeneric(img)
    }
}


func BenchmarkFromPaletted(b *testing.B) {
    f, err := os.Open("pixelview.png")
    if err != nil {
        panic(err)
    }
    img, _, err := image.Decode(f)
    if err != nil {
        panic(err)
    }
    paletted, ok := img.(*image.Paletted)
    if !ok {
        panic("Type assertion failed")
    }
    for n := 0; n < b.N; n++ {
        fromImageGeneric(paletted)
    }
}

