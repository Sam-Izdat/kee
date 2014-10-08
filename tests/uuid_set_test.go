package main

import (
    "testing"
    "reflect"
    . "github.com/smartystreets/goconvey/convey"
    "github.com/Sam-Izdat/kee"
)

func TestUUIDSet(t *testing.T) {
    kee.UUID.Options.Cache = false
    kee.UUID.Options.PadB64 = true
    kee.UUID.Options.PadB32 = true
    kee.UUID.Options.WrapA85 = false

    Convey("When the UUID is manually set with valid test value A", t, func() {
        testVal := [16]byte{23, 22, 217, 229, 211, 95, 75, 134, 139, 156, 156, 34, 97, 225, 43, 143}
        id := kee.UUID.Set(testVal)

        Convey("The resulting array should match the given array", func() {
            So(reflect.DeepEqual(id.Arr(), testVal), ShouldEqual, true)
        })

        Convey("Encoding it", func() {

            Convey("Should match expected A value as Hex (shorthand)", func() {
                So(id.String(), ShouldEqual, "1716d9e5-d35f-4b86-8b9c-9c2261e12b8f")
            })

            Convey("Should match expected A value as Hex (explicit)", func() {
                So(id.Hex(), ShouldEqual, "1716d9e5-d35f-4b86-8b9c-9c2261e12b8f")
            })

            Convey("Should match expected A value as ASCII 85", func() {
                So(id.A85(), ShouldEqual, `(Db]cdpGb&Mk$:]@Gr_t`)
            })

            Convey("Should match expected A value as Base 64", func() {
                So(id.B64(), ShouldEqual, "FxbZ5dNfS4aLnJwiYeErjw==")         
            })

            Convey("Should match expected A value as Base 64 URL-safe", func() {
                So(id.URL64(), ShouldEqual, "FxbZ5dNfS4aLnJwiYeErjw")
            })
            Convey("Should match expected A value as Base 32", func() {
                So(id.B32(), ShouldEqual, "C4LNTZOTL5FYNC44TQRGDYJLR4======")
            })

            Convey("Should match expected A value as Base 32 URL-safe", func() {
                So(id.URL32(), ShouldEqual, "C4LN-TZOT-L5FY-NC44-TQRG-DYJL-R4")
            })

            Convey("Should match expected A value as URN", func() {
                So(id.URN(), ShouldEqual, "urn:uuid:1716d9e5-d35f-4b86-8b9c-9c2261e12b8f")
            })

        })

    })

    Convey("When the UUID is manually set with valid test value B", t, func() {

        testVal := [16]byte{131, 156, 130, 220, 79, 63, 71, 254, 159, 187, 154, 25, 249, 59, 62, 227}
        id := kee.UUID.Set(testVal)

        Convey("The resulting array should match the given array", func() {
            So(reflect.DeepEqual(id.Arr(), testVal), ShouldEqual, true)
        })

        Convey("Encoding it", func() {

            Convey("Should match expected B value as Hex (shorthand)", func() {
                So(id.String(), ShouldEqual, "839c82dc-4f3f-47fe-9fbb-9a19f93b3ee3")
            })

            Convey("Should match expected B value as Hex (explicit)", func() {
                So(id.Hex(), ShouldEqual, "839c82dc-4f3f-47fe-9fbb-9a19f93b3ee3")
            })

            Convey("Should match expected B value as ASCII 85", func() {
                So(id.A85(), ShouldEqual, `K:IPK:HqAKT=^O0q)^e#`)
            })

            Convey("Should match expected B value as Base 64", func() {
                So(id.B64(), ShouldEqual, "g5yC3E8/R/6fu5oZ+Ts+4w==")         
            })

            Convey("Should match expected B value as Base 64 URL-safe", func() {
                So(id.URL64(), ShouldEqual, "g5yC3E8_R_6fu5oZ-Ts-4w")
            })
            Convey("Should match expected B value as Base 32", func() {
                So(id.B32(), ShouldEqual, "QOOIFXCPH5D75H53TIM7SOZ64M======")
            })

            Convey("Should match expected B value as Base 32 URL-safe", func() {
                So(id.URL32(), ShouldEqual, "QOOI-FXCP-H5D7-5H53-TIM7-SOZ6-4M")
            })

            Convey("Should match expected B value as URN", func() {
                So(id.URN(), ShouldEqual, "urn:uuid:839c82dc-4f3f-47fe-9fbb-9a19f93b3ee3")
            })

        })

    })
}