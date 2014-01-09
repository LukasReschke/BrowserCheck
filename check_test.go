package browsercheck

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

// Test empty UA
func TestEmpty(t *testing.T) {
	var x []Application

	Convey("Given no UA", t, func() {
		UA := ""

		Convey("When version is tested", func() {
			x = Check(UA)

			Convey("No outdated application should be found", func() {
				So(len(x), ShouldEqual, 0)
			})
		})
	})
}

// Test actual versions
func TestActual(t *testing.T) {
	var x []Application
	var oldWindows string

	Convey("Given an actual Windows 8.1 with Internet Explorer 11", t, func() {
		oldWindows = "Mozilla/5.0 (Windows NT 6.3; Trident/7.0; rv:11.0) like Gecko"

		Convey("When version is tested", func() {
			x = Check(oldWindows)

			Convey("No outdated application should be found", func() {
				So(len(x), ShouldEqual, 0)
			})
		})
	})
}

// Test outdated versions
func TestOutdated(t *testing.T) {
	var x []Application

	Convey("Given an outdated Windows XP with Internet Explorer 7", t, func() {
		outdatedString := "Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; .NET CLR 1.1.4322)"

		Convey("When version is tested", func() {
			x = Check(outdatedString)

			Convey("One outdated application should be found", func() {
				So(len(x), ShouldEqual, 1)
			})
		})
	})

	Convey("Given an old OSX Version", t, func() {
		outdatedString := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/32.0.1700.68 Safari/537.36"

		Convey("When version is tested", func() {
			x = Check(outdatedString)

			Convey("One outdated application should be found", func() {
				So(len(x), ShouldEqual, 1)
			})
		})
	})

	Convey("Given an old OSX Version with outdated Flash", t, func() {
		outdatedString := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/32.0.1700.68 Safari/537.36 Shockwave Flash 10.1 r100"

		Convey("When version is tested", func() {
			x = Check(outdatedString)

			Convey("Two outdated applications should be found", func() {
				So(len(x), ShouldEqual, 2)
			})
		})
	})
}

// Benchmark the check function
func BenchmarkCheck(b *testing.B) {
	for i := 0; i < b.N; i++ {
		outdatedString := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/32.0.1700.68 Safari/537.36 Shockwave Flash 10.1 r100"
		Check(outdatedString)
	}
}
