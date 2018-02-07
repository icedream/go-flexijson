package example

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_SomeStruct(t *testing.T) {
	Convey("MarshalJSON should work", t, func() {
		data := new(SomeStruct)
		data.A = "Hello"
		data.B = "world"
		data.Sub = new(SomeStruct)
		data.Sub.A = "hi"
		data.Sub.Fields = map[string]interface{}{
			"G": "bye",
		}
		data.Fields = map[string]interface{}{
			"C": "and",
			"D": "everyone",
			"E": "else",
			"abc": struct {
				F float64
			}{
				F: 0.1234,
			},
		}

		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(data)
		So(strings.TrimSpace(buf.String()), ShouldEqual, `{"A":"Hello","B":"world","C":"and","D":"everyone","E":"else","abc":{"F":0.1234},"sub":{"A":"hi","G":"bye"}}`)
	})
	Convey("UnmarshalJSON should work", t, func() {
		data := new(SomeStruct)

		r := strings.NewReader(`{"A":"Hello","B":"world","C":"and","D":"everyone","E":"else","abc":{"F":1337},"sub":{"A":"hi","G":"bye"}}`)

		json.NewDecoder(r).Decode(data)
		So(data.A, ShouldEqual, "Hello")
		So(data.B, ShouldEqual, "world")
		So(data.Fields, ShouldNotBeNil)
		So(len(data.Fields), ShouldBeGreaterThanOrEqualTo, 4)
		So(data.Fields["C"], ShouldEqual, "and")
		So(data.Fields["D"], ShouldEqual, "everyone")
		So(data.Fields["E"], ShouldEqual, "else")
		So(data.Fields["abc"], ShouldNotBeNil)
		abc, ok := data.Fields["abc"].(map[string]interface{})
		So(ok, ShouldBeTrue)
		So(len(abc), ShouldEqual, 1)
		So(abc["F"], ShouldEqual, 1337)
		So(data.Sub, ShouldNotBeNil)
		So(data.Sub.A, ShouldEqual, "hi")
		So(data.Sub.Fields, ShouldNotBeNil)
		So(len(data.Fields), ShouldBeGreaterThanOrEqualTo, 1)
		So(data.Sub.Fields["G"], ShouldEqual, "bye")
	})
}
