package main

import (
    "fmt"
    "testing"
    . "github.com/smartystreets/goconvey/convey"
    "github.com/Sam-Izdat/kee"
)

func TestUUIDNew(t *testing.T) {
    kee.UUID.Options.Cache = false
    kee.UUID.Options.PadB64 = true
    kee.UUID.Options.PadB32 = true
    kee.UUID.Options.WrapA85 = false

    Convey("When a V4 UUID is generated randomly", t, func() {

        id := kee.UUID.New()

        Convey("The type returned should be 'kee.uuid'", func() {
            So(fmt.Sprintf("%T",id), ShouldEqual, "kee.uuid")
        })

        Convey("Encoding it", func() {

            Convey("Should not return an empty string as Hex (shorthand)", func() {
                So(id.String(), ShouldNotEqual, "")
            })

            Convey("Should not return an empty string as Hex (explicit)", func() {
                So(id.Hex(), ShouldNotEqual, "")
            })

            Convey("Should not return an empty string as ASCII 85", func() {
                So(id.A85(), ShouldNotEqual, "")
            })
            Convey("Should not return an empty string as Base 64", func() {
                So(id.B64(), ShouldNotEqual, "")
            })
            Convey("Should not return an empty string as Base 64 URL-safe", func() {
                So(id.URL64(), ShouldNotEqual, "")
            })

            Convey("Should not return an empty string as Base 32", func() {
                So(id.B32(), ShouldNotEqual, "")
            })

            Convey("Should not return an empty string as Base 32 URL-safe", func() {
                So(id.URL32(), ShouldNotEqual, "")
            })

            Convey("Should not return an empty string as URN", func() {
                So(id.URN(), ShouldNotEqual, "")
            })

            Convey("Should not return a nil slice", func() {
                So(id.Slc(), ShouldNotEqual, nil)
            })

        })

    })
}
