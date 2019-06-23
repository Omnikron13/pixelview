package pixelview

import (
    "testing"
    "flag"
    "os"
    "path/filepath"
    "io/ioutil"
    "image"
    _ "image/png"
)


var update = flag.Bool("update", false, "update .golden files")


func TestFromImageGeneric(t *testing.T) {
    golden := filepath.Join("testdata", t.Name()+".golden")
    buf, err := ioutil.ReadFile(golden)
    if err != nil {
        panic(err)
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
        panic(err)
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
    s, err := FromPaletted(paletted)
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
        FromPaletted(paletted)
    }
}

