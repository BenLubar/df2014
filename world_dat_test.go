package df2014

import (
	"io"
	"os"
	"os/exec"
	"reflect"
	"testing"
)

func TestWorldDat_21_93_19a(t *testing.T) {
	testWorldDat(t, "work/df_21_93_19a/data/save/region8.dat")
}
func TestWorldDat_21_93_19c(t *testing.T) {
	testWorldDat(t, "work/df_21_93_19c/data/save/region4.dat")
}
func TestWorldDat_21_95_19a(t *testing.T) {
	testWorldDat(t, "work/df_21_95_19a/data/save/region1.dat")
}
func TestWorldDat_21_95_19b(t *testing.T) {
	testWorldDat(t, "work/df_21_95_19b/data/save/region2.dat")
}
func TestWorldDat_21_95_19c(t *testing.T) {
	testWorldDat(t, "work/df_21_95_19c/data/save/region4.dat")
}
func TestWorldDat_21_100_19a(t *testing.T) {
	testWorldDat(t, "work/df_21_100_19a/data/save/region1.dat")
}
func TestWorldDat_21_101_19a(t *testing.T) {
	testWorldDat(t, "work/df_21_101_19a/data/save/region1.dat")
}
func TestWorldDat_21_101_19d(t *testing.T) {
	testWorldDat(t, "work/df_21_101_19d/data/save/region1.dat")
}
func TestWorldDat_21_102_19a(t *testing.T) {
	testWorldDat(t, "work/df_21_102_19a/data/save/region1.dat")
}
func TestWorldDat_21_104_19b(t *testing.T) {
	testWorldDat(t, "work/df_21_104_19b/data/save/region1.dat")
}
func TestWorldDat_21_104_19d(t *testing.T) {
	testWorldDat(t, "work/df_21_104_19d/data/save/region1.dat")
}
func TestWorldDat_21_104_21a(t *testing.T) {
	testWorldDat(t, "work/df_21_104_21a/data/save/region1.dat")
}
func TestWorldDat_21_104_21b(t *testing.T) {
	testWorldDat(t, "work/df_21_104_21b/data/save/region1.dat")
}
func TestWorldDat_21_105_21a(t *testing.T) {
	testWorldDat(t, "work/df_21_105_21a/data/save/region1.dat")
}
func TestWorldDat_22_107_21a(t *testing.T) {
	testWorldDat(t, "work/df_22_107_21a/data/save/region1.dat")
}
func TestWorldDat_22_110_22e(t *testing.T) {
	testWorldDat(t, "work/df_22_110_22e/data/save/region1.dat")
}
func TestWorldDat_22_110_22f(t *testing.T) {
	testWorldDat(t, "work/df_22_110_22f/data/save/region1.dat")
}
func TestWorldDat_22_110_23a(t *testing.T) {
	testWorldDat(t, "work/df_22_110_23a/data/save/region1.dat")
}
func TestWorldDat_22_120_23a(t *testing.T) {
	testWorldDat(t, "work/df_22_120_23a/data/save/region1.dat")
}
func TestWorldDat_22_121_23b(t *testing.T) {
	testWorldDat(t, "work/df_22_121_23b/data/save/region1.dat")
}
func TestWorldDat_22_123_23a(t *testing.T) {
	testWorldDat(t, "work/df_22_123_23a/data/save/region1.dat")
}
func TestWorldDat_23_130_23a(t *testing.T) {
	testWorldDat(t, "work/df_23_130_23a/data/save/region1.dat")
}

func TestWorldDat_27_169_32a(t *testing.T) {
	testWorldDat(t, "work/df_27_169_32a/data/save/region1/world.dat")
}
func TestWorldDat_27_169_33a(t *testing.T) {
	testWorldDat(t, "work/df_27_169_33a/data/save/region1/world.dat")
}
func TestWorldDat_27_169_33b(t *testing.T) {
	testWorldDat(t, "work/df_27_169_33b/data/save/region1/world.dat")
}
func TestWorldDat_27_169_33c(t *testing.T) {
	testWorldDat(t, "work/df_27_169_33c/data/save/region1/world.dat")
}
func TestWorldDat_27_169_33d(t *testing.T) {
	testWorldDat(t, "work/df_27_169_33d/data/save/region1/world.dat")
}
func TestWorldDat_27_169_33e(t *testing.T) {
	testWorldDat(t, "work/df_27_169_33e/data/save/region1/world.dat")
}
func TestWorldDat_27_169_33f(t *testing.T) {
	testWorldDat(t, "work/df_27_169_33f/data/save/region1/world.dat")
}
func TestWorldDat_27_169_33g(t *testing.T) {
	testWorldDat(t, "work/df_27_169_33g/data/save/region1/world.dat")
}
func TestWorldDat_27_173_38a(t *testing.T) {
	testWorldDat(t, "work/df_27_173_38a/data/save/region1/world.dat")
}
func TestWorldDat_27_176_38a(t *testing.T) {
	testWorldDat(t, "work/df_27_176_38a/data/save/region1/world.dat")
}
func TestWorldDat_27_176_38b(t *testing.T) {
	testWorldDat(t, "work/df_27_176_38b/data/save/region1/world.dat")
}
func TestWorldDat_27_176_38c(t *testing.T) {
	testWorldDat(t, "work/df_27_176_38c/data/save/region1/world.dat")
}
func TestWorldDat_28_181_39a(t *testing.T) {
	testWorldDat(t, "work/df_28_181_39a/data/save/region1/world.dat")
}
func TestWorldDat_28_181_39b(t *testing.T) {
	testWorldDat(t, "work/df_28_181_39b/data/save/region1/world.dat")
}
func TestWorldDat_28_181_39c(t *testing.T) {
	testWorldDat(t, "work/df_28_181_39c/data/save/region1/world.dat")
}
func TestWorldDat_28_181_39d(t *testing.T) {
	testWorldDat(t, "work/df_28_181_39d/data/save/region1/world.dat")
}
func TestWorldDat_28_181_39e(t *testing.T) {
	testWorldDat(t, "work/df_28_181_39e/data/save/region1/world.dat")
}
func TestWorldDat_28_181_39f(t *testing.T) {
	testWorldDat(t, "work/df_28_181_39f/data/save/region1/world.dat")
}
func TestWorldDat_28_181_40a(t *testing.T) {
	testWorldDat(t, "work/df_28_181_40a/data/save/region1/world.dat")
}
func TestWorldDat_28_181_40b(t *testing.T) {
	testWorldDat(t, "work/df_28_181_40b/data/save/region1/world.dat")
}
func TestWorldDat_28_181_40c(t *testing.T) {
	testWorldDat(t, "work/df_28_181_40c/data/save/region1/world.dat")
}
func TestWorldDat_28_181_40d(t *testing.T) {
	testWorldDat(t, "work/df_28_181_40d/data/save/region1/world.dat")
}

func TestWorldDat_31_01(t *testing.T) {
	testWorldDat(t, "work/df_31_01/data/save/region1/world.dat")
}
func TestWorldDat_31_02(t *testing.T) {
	testWorldDat(t, "work/df_31_02/data/save/region1/world.dat")
}
func TestWorldDat_31_03(t *testing.T) {
	testWorldDat(t, "work/df_31_03/data/save/region1/world.dat")
}
func TestWorldDat_31_04(t *testing.T) {
	testWorldDat(t, "work/df_31_04/data/save/region1/world.dat")
}
func TestWorldDat_31_05(t *testing.T) {
	testWorldDat(t, "work/df_31_05/data/save/region1/world.dat")
}
func TestWorldDat_31_06(t *testing.T) {
	testWorldDat(t, "work/df_31_06/data/save/region1/world.dat")
}
func TestWorldDat_31_08(t *testing.T) {
	testWorldDat(t, "work/df_31_08/data/save/region1/world.dat")
}
func TestWorldDat_31_09(t *testing.T) {
	testWorldDat(t, "work/df_31_09/data/save/region1/world.dat")
}
func TestWorldDat_31_10(t *testing.T) {
	testWorldDat(t, "work/df_31_10/data/save/region1/world.dat")
}
func TestWorldDat_31_11(t *testing.T) {
	testWorldDat(t, "work/df_31_11/data/save/region1/world.dat")
}
func TestWorldDat_31_12(t *testing.T) {
	testWorldDat(t, "work/df_31_12/data/save/region1/world.dat")
}
func TestWorldDat_31_13(t *testing.T) {
	testWorldDat(t, "work/df_31_13/data/save/region1/world.dat")
}
func TestWorldDat_31_14(t *testing.T) {
	testWorldDat(t, "work/df_31_14/data/save/region1/world.dat")
}
func TestWorldDat_31_15(t *testing.T) {
	testWorldDat(t, "work/df_31_15/data/save/region1/world.dat")
}
func TestWorldDat_31_16(t *testing.T) {
	testWorldDat(t, "work/df_31_16/data/save/region1/world.dat")
}
func TestWorldDat_31_17(t *testing.T) {
	testWorldDat(t, "work/df_31_17/data/save/region1/world.dat")
}
func TestWorldDat_31_18(t *testing.T) {
	testWorldDat(t, "work/df_31_18/data/save/region1/world.dat")
}
func TestWorldDat_31_19(t *testing.T) {
	testWorldDat(t, "work/df_31_19/data/save/region1/world.dat")
}
func TestWorldDat_31_21(t *testing.T) {
	testWorldDat(t, "work/df_31_21/data/save/region1/world.dat")
}
func TestWorldDat_31_22(t *testing.T) {
	testWorldDat(t, "work/df_31_22/data/save/region1/world.dat")
}
func TestWorldDat_31_23(t *testing.T) {
	testWorldDat(t, "work/df_31_23/data/save/region1/world.dat")
}
func TestWorldDat_31_24(t *testing.T) {
	testWorldDat(t, "work/df_31_24/data/save/region1/world.dat")
}
func TestWorldDat_31_25(t *testing.T) {
	testWorldDat(t, "work/df_31_25/data/save/region1/world.dat")
}

func TestWorldDat_34_01(t *testing.T) {
	testWorldDat(t, "work/df_34_01/data/save/region1/world.dat")
}
func TestWorldDat_34_02(t *testing.T) {
	testWorldDat(t, "work/df_34_02/data/save/region1/world.dat")
}
func TestWorldDat_34_03(t *testing.T) {
	testWorldDat(t, "work/df_34_03/data/save/region1/world.dat")
}
func TestWorldDat_34_04(t *testing.T) {
	testWorldDat(t, "work/df_34_04/data/save/region1/world.dat")
}
func TestWorldDat_34_05(t *testing.T) {
	testWorldDat(t, "work/df_34_05/data/save/region1/world.dat")
}
func TestWorldDat_34_06(t *testing.T) {
	testWorldDat(t, "work/df_34_06/data/save/region1/world.dat")
}
func TestWorldDat_34_07(t *testing.T) {
	testWorldDat(t, "work/df_34_07/data/save/region1/world.dat")
}
func TestWorldDat_34_08(t *testing.T) {
	testWorldDat(t, "work/df_34_08/data/save/region1/world.dat")
}
func TestWorldDat_34_09(t *testing.T) {
	testWorldDat(t, "work/df_34_09/data/save/region1/world.dat")
}
func TestWorldDat_34_10(t *testing.T) {
	testWorldDat(t, "work/df_34_10/data/save/region1/world.dat")
}
func TestWorldDat_34_11(t *testing.T) {
	testWorldDat(t, "work/df_34_11/data/save/region1/world.dat")
}

func TestWorldDat_40_01(t *testing.T) {
	testWorldDat(t, "work/df_40_01/data/save/region1/world.dat")
}
func TestWorldDat_40_02(t *testing.T) {
	testWorldDat(t, "work/df_40_02/data/save/region1/world.dat")
}
func TestWorldDat_40_03(t *testing.T) {
	testWorldDat(t, "work/df_40_03/data/save/region1/world.dat")
}
func TestWorldDat_40_04(t *testing.T) {
	testWorldDat(t, "work/df_40_04/data/save/region1/world.dat")
}
func TestWorldDat_40_05(t *testing.T) {
	testWorldDat(t, "work/df_40_05/data/save/region1/world.dat")
}
func TestWorldDat_40_06(t *testing.T) {
	testWorldDat(t, "work/df_40_06/data/save/region1/world.dat")
}
func TestWorldDat_40_07(t *testing.T) {
	testWorldDat(t, "work/df_40_07/data/save/region1/world.dat")
}
func TestWorldDat_40_08(t *testing.T) {
	testWorldDat(t, "work/df_40_08/data/save/region1/world.dat")
}
func TestWorldDat_40_09(t *testing.T) {
	testWorldDat(t, "work/df_40_09/data/save/region1/world.dat")
}
func TestWorldDat_40_10(t *testing.T) {
	testWorldDat(t, "work/df_40_10/data/save/region1/world.dat")
}
func TestWorldDat_40_11(t *testing.T) {
	testWorldDat(t, "work/df_40_11/data/save/region1/world.dat")
}
func TestWorldDat_40_12(t *testing.T) {
	testWorldDat(t, "work/df_40_12/data/save/region1/world.dat")
}
func TestWorldDat_40_13(t *testing.T) {
	testWorldDat(t, "work/df_40_13/data/save/region1/world.dat")
}
func TestWorldDat_40_14(t *testing.T) {
	testWorldDat(t, "work/df_40_14/data/save/region1/world.dat")
}
func TestWorldDat_40_15(t *testing.T) {
	testWorldDat(t, "work/df_40_15/data/save/region1/world.dat")
}
func TestWorldDat_40_16(t *testing.T) {
	testWorldDat(t, "work/df_40_16/data/save/region1/world.dat")
}
func TestWorldDat_40_17(t *testing.T) {
	testWorldDat(t, "work/df_40_17/data/save/region1/world.dat")
}
func TestWorldDat_40_18(t *testing.T) {
	testWorldDat(t, "work/df_40_18/data/save/region1/world.dat")
}
func TestWorldDat_40_19(t *testing.T) {
	testWorldDat(t, "work/df_40_19/data/save/region1/world.dat")
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
			t.Error("Unparsed data:\n" + string(b))
		}
	}
}
