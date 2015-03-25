package core

import (
	"reflect"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_countUsable(t *testing.T) {
	Convey("计算使用次数", t, func() {
		So(_countUsable(5, 8), ShouldEqual, 33)
	})
}

func Test_parseItem(t *testing.T) {
	str := "DEPOSIT_ITEMS ({\"100313130_5_8\",\"100102496_8_8\",\"100101132_8_8\",\"0_0_0\",\"100112073_8_8\",})"
	ret := parseItem(strings.FieldsFunc(str, splitInit))
	item := []Item{{100313130, 33}, {100102496, 36}, {100101132, 36}, {100112073, 36}}

	Convey("分割字符串", t, func() {
		So(reflect.DeepEqual(ret, item), ShouldBeTrue)
	})
}

func Test_loadTargetFile(t *testing.T) {
	filePath, err := _loadTargetFile("../../data/upinfo.tmp")
	file := NewUserFileSet()
	file.Add("100/d1hgb2238.storage.o")
	file.Add("110/nsveaia33358.storage.o")
	file.Add("50/2012yihan.storage.o")
	file.Add("104/hw63158527.storage.o")
	Convey("读取upinfo.tmp", t, func() {
		So(err, ShouldBeNil)
		So(filePath.Len(), ShouldEqual, 4)
		So(reflect.DeepEqual(filePath, file), ShouldBeTrue)
	})
}

func Test_loadStorageFile(t *testing.T) {
	user, err := _loadStorageFile("../../data/userdata/104/hw63158527.storage.o")
	Convey("读取仓库", t, func() {
		So(err, ShouldBeNil)
		So(user.Account, ShouldEqual, "hw63158527")
		So(len(user.ComStorage), ShouldEqual, 102)
		So(len(user.VipStorage), ShouldEqual, 128)
	})
}
