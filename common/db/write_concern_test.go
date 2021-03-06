package db

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestBuildWriteConcern(t *testing.T) {
	Convey("Given a write concern string value, and a boolean indicating if the "+
		"write concern is to be used on a replica set, on calling BuildWriteConcern...", t, func() {
		Convey("no error should be returned if the write concern is valid", func() {
			writeConcern, err := BuildWriteConcern(`{w:34}`, true)
			So(err, ShouldBeNil)
			So(writeConcern.W, ShouldEqual, 34)
			writeConcern, err = BuildWriteConcern(`{w:"majority"}`, true)
			So(err, ShouldBeNil)
			So(writeConcern.WMode, ShouldEqual, "majority")
			writeConcern, err = BuildWriteConcern(`majority`, true)
			So(err, ShouldBeNil)
			So(writeConcern.WMode, ShouldEqual, "majority")
			writeConcern, err = BuildWriteConcern(`tagset`, true)
			So(err, ShouldBeNil)
			So(writeConcern.WMode, ShouldEqual, "tagset")
		})
		Convey("on replica sets, only a write concern of 1 or 0 should be returned", func() {
			writeConcern, err := BuildWriteConcern(`{w:34}`, false)
			So(err, ShouldBeNil)
			So(writeConcern.W, ShouldEqual, 1)
			writeConcern, err = BuildWriteConcern(`{w:"majority"}`, false)
			So(err, ShouldBeNil)
			So(writeConcern.W, ShouldEqual, 1)
			writeConcern, err = BuildWriteConcern(`tagset`, false)
			So(err, ShouldBeNil)
			So(writeConcern.W, ShouldEqual, 1)
		})
		Convey("with a w value of 0, without j set, a nil write concern should be returned", func() {
			writeConcern, err := BuildWriteConcern(`{w:0}`, false)
			So(err, ShouldBeNil)
			So(writeConcern, ShouldBeNil)
		})
		Convey("with a negative w value, an error should be returned", func() {
			_, err := BuildWriteConcern(`{w:-1}`, false)
			So(err, ShouldNotBeNil)
			_, err = BuildWriteConcern(`{w:-2}`, false)
			So(err, ShouldNotBeNil)
		})
		Convey("with a w value of 0, with j set, a non-nil write concern should be returned", func() {
			writeConcern, err := BuildWriteConcern(`{w:0, j:true}`, false)
			So(err, ShouldBeNil)
			So(writeConcern.J, ShouldBeTrue)
		})
	})
}

func TestConstructWCObject(t *testing.T) {
	Convey("Given a write concern string value, on calling constructWCObject...", t, func() {

		Convey("non-JSON string values should be assigned to the 'WMode' "+
			"field in their entirety", func() {
			writeConcernString := "majority"
			writeConcern, err := constructWCObject(writeConcernString)
			So(err, ShouldBeNil)
			So(writeConcern.WMode, ShouldEqual, writeConcernString)
		})

		Convey("non-JSON int values should be assigned to the 'w' field "+
			"in their entirety", func() {
			writeConcernString := `{w: 4}`
			writeConcern, err := constructWCObject(writeConcernString)
			So(err, ShouldBeNil)
			So(writeConcern.W, ShouldEqual, 4)
		})

		Convey("JSON strings with valid j, wtimeout, fsync and w, should be "+
			"assigned accordingly", func() {
			writeConcernString := `{w: 3, j: true, fsync: false, wtimeout: 43}`
			expectedW := 3
			expectedWTimeout := 43
			writeConcern, err := constructWCObject(writeConcernString)
			So(err, ShouldBeNil)
			So(writeConcern.W, ShouldEqual, expectedW)
			So(writeConcern.J, ShouldBeTrue)
			So(writeConcern.FSync, ShouldBeFalse)
			So(writeConcern.WTimeout, ShouldEqual, expectedWTimeout)
		})

		Convey("JSON strings with an argument for j that is not false should set j true", func() {
			writeConcernString := `{w: 3, j: "rue"}`
			writeConcern, err := constructWCObject(writeConcernString)
			So(err, ShouldBeNil)
			So(writeConcern.W, ShouldEqual, 3)
			So(writeConcern.J, ShouldBeTrue)
		})

		Convey("JSON strings with an argument for fsync that is not false should set fsync true", func() {
			writeConcernString := `{w: 3, fsync: "rue"}`
			writeConcern, err := constructWCObject(writeConcernString)
			So(err, ShouldBeNil)
			So(writeConcern.W, ShouldEqual, 3)
			So(writeConcern.FSync, ShouldBeTrue)
		})

		Convey("JSON strings with an invalid wtimeout argument should error out", func() {
			writeConcernString := `{w: 3, wtimeout: "rue"}`
			_, err := constructWCObject(writeConcernString)
			So(err, ShouldNotBeNil)
			writeConcernString = `{w: 3, wtimeout: "43"}`
			_, err = constructWCObject(writeConcernString)
			So(err, ShouldNotBeNil)
		})

		Convey("JSON strings with any non-false j argument should not error out", func() {
			writeConcernString := `{w: 3, j: "t"}`
			writeConcern, err := constructWCObject(writeConcernString)
			So(err, ShouldBeNil)
			So(writeConcern.J, ShouldBeTrue)
			writeConcernString = `{w: 3, j: "f"}`
			writeConcern, err = constructWCObject(writeConcernString)
			So(err, ShouldBeNil)
			So(writeConcern.J, ShouldBeTrue)
			writeConcernString = `{w: 3, j: false}`
			writeConcern, err = constructWCObject(writeConcernString)
			So(err, ShouldBeNil)
			So(writeConcern.J, ShouldBeFalse)
			writeConcernString = `{w: 3, j: 0}`
			writeConcern, err = constructWCObject(writeConcernString)
			So(err, ShouldBeNil)
			So(writeConcern.J, ShouldBeFalse)
		})

		Convey("JSON strings with a shorthand fsync argument should not error out", func() {
			writeConcernString := `{w: 3, fsync: "t"}`
			writeConcern, err := constructWCObject(writeConcernString)
			So(err, ShouldBeNil)
			So(writeConcern.FSync, ShouldBeTrue)
			writeConcernString = `{w: "3", fsync: "f"}`
			writeConcern, err = constructWCObject(writeConcernString)
			So(err, ShouldBeNil)
			So(writeConcern.FSync, ShouldBeTrue)
			writeConcernString = `{w: "3", fsync: false}`
			writeConcern, err = constructWCObject(writeConcernString)
			So(err, ShouldBeNil)
			So(writeConcern.FSync, ShouldBeFalse)
			writeConcernString = `{w: "3", fsync: 0}`
			writeConcern, err = constructWCObject(writeConcernString)
			So(err, ShouldBeNil)
			So(writeConcern.FSync, ShouldBeFalse)
		})

		Convey("Unacknowledge write concern strings should return a nil object "+
			"if journaling is not required", func() {
			writeConcernString := `{w: 0}`
			writeConcern, err := constructWCObject(writeConcernString)
			So(err, ShouldBeNil)
			So(writeConcern, ShouldBeNil)
			writeConcernString = `{w: 0}`
			writeConcern, err = constructWCObject(writeConcernString)
			So(err, ShouldBeNil)
			So(writeConcern, ShouldBeNil)
			writeConcernString = `0`
			writeConcern, err = constructWCObject(writeConcernString)
			So(err, ShouldBeNil)
			So(writeConcern, ShouldBeNil)
		})
	})
}
