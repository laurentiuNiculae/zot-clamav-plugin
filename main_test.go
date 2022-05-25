package main_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/laurentiuNiculae/zot-clamav-plugin/utils"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	// EicarSignature is used as a test value for virus scanners. The virus scanner
	// should find this string and call the file infected.
	EicarSignature = `X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*`
)

func TestScanImage(t *testing.T) {
	Convey("Scan filesystem that has test infected file", t, func() {
		infectedDir, err := createInfectedTestDirectory()
		defer os.RemoveAll(infectedDir)
		So(err, ShouldBeNil)
		So(infectedDir, ShouldNotBeEmpty)

		result, err := utils.ScanImage("testImg", infectedDir)
		So(err, ShouldNotBeNil)
		So(result, ShouldNotBeEmpty)
	})

	Convey("Scan clean filesystem ", t, func() {
		infectedDir, err := createCleanTestDirectory()
		defer os.RemoveAll(infectedDir)
		So(err, ShouldBeNil)
		So(infectedDir, ShouldNotBeEmpty)

		result, err := utils.ScanImage("testImg", infectedDir)
		So(err, ShouldBeNil)
		So(result, ShouldNotBeEmpty)
	})
}

func createInfectedTestDirectory() (infectedDir string, err error) {
	infectedDir, err = ioutil.TempDir("", "clamav-test-infected-dir*")
	So(err, ShouldBeNil)

	f1, err := ioutil.TempFile(infectedDir, "infected-file*.txt")
	So(err, ShouldBeNil)
	defer f1.Close()
	f1.WriteString(EicarSignature)

	f2, err := ioutil.TempFile(infectedDir, "uninfected-file*.txt")
	So(err, ShouldBeNil)
	defer f2.Close()
	f2.WriteString("This text is fine, there is no problem with it.")

	f3, err := ioutil.TempFile(infectedDir, "infected-file*.txt")
	So(err, ShouldBeNil)
	defer f3.Close()
	f3.WriteString(EicarSignature)

	return infectedDir, err
}

func createCleanTestDirectory() (infectedDir string, err error) {
	infectedDir, err = ioutil.TempDir("", "clamav-test-infected-dir*")
	So(err, ShouldBeNil)

	f1, err := ioutil.TempFile(infectedDir, "infected-file*.txt")
	So(err, ShouldBeNil)
	defer f1.Close()
	f1.WriteString("This text is fine, there is no problem with it.")

	f2, err := ioutil.TempFile(infectedDir, "uninfected-file*.txt")
	So(err, ShouldBeNil)
	defer f2.Close()
	f2.WriteString("This text is fine, there is no problem with it.")

	return infectedDir, err
}
