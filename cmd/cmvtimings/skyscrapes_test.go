package main

import "fmt"

func Example() {
	check(process("../../testdata/dffd_0006190/skyscrapesEnd1.cmv"))
	check(process("../../testdata/dffd_0006190/skyscrapesEnd2.cmv"))
	check(process("../../testdata/dffd_0006190/SkyscrapesEnd3.cmv"))

	// Output:
	// 1m43.86s
	// You have discovered an eerie cavern. The air above the dark stone floor is alive with vortices of purple light and dark, boiling clouds. Seemingly bottomless glowing pits mark the surface.
	//
	// 4m6.56s
	// Horrifying screams come from the darkness below!
	//
	// 28m35.7s
	// 26m44.76s
	// error: unexpected EOF (near frame 55300)
}

func check(err error) {
	if err != nil {
		fmt.Println("error:", err)
	}
}
