package srw

import "testing"

func TestMain(m *testing.M) {
	f := NewFractal(6, []int{})
	f.SetRauzy()
	f.SetBdd()
	f.Save("test.png")
}
