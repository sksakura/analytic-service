package profiler

import "testing"

func TestProfilerCycle(t *testing.T) {
	p := Profiler{}
	t.Run("TestProfilerCycle", func(t *testing.T) {
		p.On()
		if p.Status() != true {
			t.Error(p.Status(), true)
		}

		p.Off()
		if p.Status() != false {
			t.Error(p.Status(), false)
		}
	})
}
