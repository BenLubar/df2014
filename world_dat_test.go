//go:generate go run download_test_data.go

package df2014

import (
	"io"
	"os"
	"os/exec"
	"reflect"
	"testing"
)

func TestWorldDat_dffd_0000004(t *testing.T) {
	testWorldDat(t, "testdata/dffd_0000004/world.dat")
}

func TestWorldDat_dffd_0000573(t *testing.T) {
	testWorldDat(t, "testdata/dffd_0000573/world.dat")
}

func TestWorldDat_dffd_0003810(t *testing.T) {
	testWorldDat(t, "testdata/dffd_0003810/world.dat")
}

func TestWorldDat_dffd_0005154(t *testing.T) {
	testWorldDat(t, "testdata/dffd_0005154/world.dat")
}

func TestWorldDat_dffd_0005574(t *testing.T) {
	testWorldDat(t, "testdata/dffd_0005574/world.dat")
}

func TestWorldDat_dffd_0005930(t *testing.T) {
	testWorldDat(t, "testdata/dffd_0005930/world.dat")
}

func TestWorldDat_dffd_0005994(t *testing.T) {
	testWorldDat(t, "testdata/dffd_0005994/world.dat")
}

func TestWorldDat_dffd_0006331(t *testing.T) {
	testWorldDat(t, "testdata/dffd_0006331/world.dat")
}

func TestWorldDat_dffd_0006808(t *testing.T) {
	testWorldDat(t, "testdata/dffd_0006808/world.dat")
}

func TestWorldDat_dffd_0007554(t *testing.T) {
	testWorldDat(t, "testdata/dffd_0007554/world.dat")
}

func TestWorldDat_dffd_0008345(t *testing.T) {
	testWorldDat(t, "testdata/dffd_0008345/world.dat")
}

func TestWorldDat_dffd_0010619(t *testing.T) {
	testWorldDat(t, "testdata/dffd_0010619/world.dat")
}

func TestWorldDat_dffd_0010759(t *testing.T) {
	testWorldDat(t, "testdata/dffd_0010759/world.dat")
}

func testWorldDat(t *testing.T, fn string) {
	t.Log(fn)

	f, err := os.Open(fn)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	r := &Reader{f}

	var w WorldDat

	err = r.Decode(&w)
	if err != nil {
		t.Error(err)
	}
	if testing.Verbose() {
		t.Log(string(prettyPrint(&w, reflect.ValueOf(w), nil, []byte{'\n'}, "")))
	}

	if err == nil || testing.Verbose() {
		xxd := exec.Command("xxd")
		xxd.Stdin = io.LimitReader(r, 1<<10)
		b, err := xxd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}
		if len(b) != 0 {
			t.Errorf("Unparsed data:\n%s", b)
		}
	}
}
