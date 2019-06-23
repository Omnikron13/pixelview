package pixelview

import (
    "testing"
    "io/ioutil"
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

// TODO: test FromReader() & FromImage() (especially things like sub-images?)

