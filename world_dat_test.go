package df2014

import (
	"github.com/kr/pretty"
	"io"
	"os"
	"os/exec"
	"testing"
)

func TestRegion1WorldDat(t *testing.T) {
	testWorldDat(t, "work/df_linux/data/save/region1/world.dat")
}

func TestRegion2WorldDat(t *testing.T) {
	testWorldDat(t, "work/df_linux/data/save/region2/world.dat")
}

func TestRegion3WorldDat(t *testing.T) {
	testWorldDat(t, "work/df_linux/data/save/region3/world.dat")
}

func TestRegion4WorldDat(t *testing.T) {
	testWorldDat(t, "work/df_linux/data/save/region4/world.dat")
}

func TestRegion5WorldDat(t *testing.T) {
	testWorldDat(t, "work/df_linux/data/save/region5/world.dat")
}

func testWorldDat(t *testing.T, fn string) {
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
	t.Logf("%# v", pretty.Formatter(w))

	t.Logf("%d books", len(w.Books))

	xxd := exec.Command("xxd")
	xxd.Stdin = io.LimitReader(r, 1<<10)
	b, err := xxd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	if len(b) != 0 {
		t.Fatal("Unparsed data:\n" + string(b))
	}
}
