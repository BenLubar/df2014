package df2014

import (
	"io"
	"os"
	"os/exec"
	"testing"
)

func TestWorldDat(t *testing.T) {
	f, err := os.Open("work/df_linux/data/save/region1/world.dat")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	r := &Reader{f}

	w, err := r.WorldDat()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", w)

	xxd := exec.Command("xxd")
	xxd.Stdin = io.LimitReader(r, 128)
	b, err := xxd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	if len(b) != 0 {
		t.Fatal("Unparsed data:\n" + string(b))
	}
}
