package pixelview

import (
    "testing"
    "flag"
    "os"
    "path/filepath"
    "io/ioutil"
    "image"
    "image/color"
    _ "image/png"
)


var update = flag.Bool("update", false, "update .golden files")


// This test is actually pretty pointless, given that the specific tests
// below are the actual logic, and that FromFile() is an incredibly light
// convenience function.
func TestFromFile(t *testing.T) {
    golden := filepath.Join("testdata", t.Name()+".golden")

    buf, err := ioutil.ReadFile(golden)
    if err != nil {
        t.Error("Golden file could not be read")
    }
    reference := string(buf)

    s, err := FromFile(filepath.Join("testdata", "pixelview.png"))
    if err != nil {
        panic(err)
    }

    if s != reference {
        t.Error("Output did not match reference")
    }

    if *update {
        ioutil.WriteFile(golden, []byte(s), 0644)
    }
}


func TestFromImageGeneric(t *testing.T) {
    golden := filepath.Join("testdata", t.Name()+".golden")

    buf, err := ioutil.ReadFile(golden)
    if err != nil {
        t.Error("Golden file could not be read")
    }
    reference := string(buf)

    f, err := os.Open(filepath.Join("testdata", "pixelview.png"))
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

    if *update {
        ioutil.WriteFile(golden, []byte(s), 0644)
    }
}


func TestFromPaletted(t *testing.T) {
    golden := filepath.Join("testdata", t.Name()+".golden")

    buf, err := ioutil.ReadFile(golden)
    if err != nil {
        t.Error("Golden file could not be read")
    }
    reference := string(buf)

    f, err := os.Open(filepath.Join("testdata", "pixelview.png"))
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
    s, err := fromPaletted(paletted)
    if err != nil {
        t.Errorf("Error encountered during execution: %s", err)
    }
    if s != reference {
        t.Error("Output did not match reference")
    }

    if *update {
        ioutil.WriteFile(golden, []byte(s), 0644)
    }
}


func TestEncode(t *testing.T) {
    var prevfg, prevbg color.Color
    fg := &color.RGBA {
        R: 0x12,
        G: 0x34,
        B: 0x56,
    }
    bg := &color.RGBA {
        R: 0xAB,
        G: 0xCD,
        B: 0xEF,
    }

    t.Run("No RLE", func(t *testing.T) {
        ref := "[#123456:#abcdef]▀"
        s := encode(fg, bg, &prevfg, &prevbg)
        if s != ref {
            t.Errorf("Output (%s) did not match reference (%s)", s, ref)
        }
    })

    t.Run("Full RLE", func(t *testing.T) {
        ref := "▀"
        s := encode(fg, bg, &prevfg, &prevbg)
        if s != ref {
            t.Errorf("Output (%s) did not match reference (%s)", s, ref)
        }
    })

    t.Run("FG RLE", func(t *testing.T) {
        ref := "[:#abcdef]▀"
        prevbg = nil
        s := encode(fg, bg, &prevfg, &prevbg)
        if s != ref {
            t.Errorf("Output (%s) did not match reference (%s)", s, ref)
        }
    })

    t.Run("BG RLE", func(t *testing.T) {
        ref := "[#123456:]▀"
        prevfg = nil
        s := encode(fg, bg, &prevfg, &prevbg)
        if s != ref {
            t.Errorf("Output (%s) did not match reference (%s)", s, ref)
        }
    })
}


// TODO: test FromReader() & FromImage() (especially things like sub-images?)


func BenchmarkFromImageGeneric(b *testing.B) {
    f, err := os.Open(filepath.Join("testdata", "pixelview.png"))
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
    f, err := os.Open(filepath.Join("testdata", "pixelview.png"))
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
        fromPaletted(paletted)
    }
}

