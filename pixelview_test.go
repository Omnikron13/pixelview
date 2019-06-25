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
    // And this seems even more pointless as it is essentially
    // testing if os.Open() works, but without this test apparently
    // there is missing code coverage...
    t.Run("Missing File", func(t *testing.T) {
        _, err := FromFile(filepath.Join("testdata", "404.png"))
        if err == nil {
            t.Error("Didn't error on non-existent file")
        }
    })

    golden, reference := getGolden(t)
    t.Run("Output", func(t *testing.T) {
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
    })
}


// Ideally fromImageGeneric() shouldn't ever be called, as more efficient
// and specific functions can process the common specific types of images.
// It obviously needs to work properly still though, as implementing specific
// functions for all types would be onerous.
func TestFromImageGeneric(t *testing.T) {
    golden, reference := getGolden(t)

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
    golden, reference := getGolden(t)

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


func TestFromNRGBA(t *testing.T) {
    golden, reference := getGolden(t)

    f, err := os.Open(filepath.Join("testdata", "nrgba.png"))
    if err != nil {
        panic(err)
    }
    img, _, err := image.Decode(f)
    if err != nil {
        panic(err)
    }
    nrgba, ok := img.(*image.NRGBA)
    if !ok {
        panic("Type assertion failed before test could be run")
    }

    s, err := fromNRGBA(nrgba)
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


// This helper function loads golden files into strings for comparison
// and also returns the (relative) path to them, which is required if
// -update is passed to go test so new versions can be written.
func getGolden(t *testing.T) (golden, reference string) {
    golden = filepath.Join("testdata", t.Name()+".golden")

    buf, err := ioutil.ReadFile(golden)
    if err != nil {
        t.Error("Golden file could not be read")
    }
    reference = string(buf)
    return
}


// TODO: test FromReader() & FromImage() (especially things like sub-images?)


func BenchmarkFromImageGeneric(b *testing.B) {
    benchmarks := []string{
        "paletted.png",
        "nrgba.png",
    }
    for _, s := range benchmarks {
        b.Run(s, func(b *testing.B) {
            img := loadTestImage(b, s)
            for n := 0; n < b.N; n++ {
                fromImageGeneric(img)
            }
        })
    }
}


func BenchmarkFromPaletted(b *testing.B) {
    img := loadTestImage(b, "paletted.png")
    paletted, ok := img.(*image.Paletted)
    if !ok {
        panic("Type assertion failed")
    }
    for n := 0; n < b.N; n++ {
        fromPaletted(paletted)
    }
}


func BenchmarkFromNRGBA(b *testing.B) {
    img := loadTestImage(b, "nrgba.png")
    nrgba, ok := img.(*image.NRGBA)
    if !ok {
        panic("Type assertion failed")
    }
    for n := 0; n < b.N; n++ {
        fromNRGBA(nrgba)
    }
}


func loadTestImage(tb testing.TB, filename string) (img image.Image) {
    f, err := os.Open(filepath.Join("testdata", filename))
    if err != nil {
        tb.Fatal("Test image could not be opened")
    }
    img, _, err = image.Decode(f)
    if err != nil {
        tb.Fatal("Test image could not be decoded")
    }
    return
}

