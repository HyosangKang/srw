package srw

import (
	"gonum.org/v1/plot/plotter"
)

type pt struct {
	word [3]int
	xy   [2]float64
	par  []*pt
	chd  []*pt
}

func (p *pt) plotter() plotter.XY {
	return plotter.XY{
		X: p.xy[0], Y: p.xy[1],
	}
}

func (p *pt) Id() int {
	return p.word[0] + p.word[1] + p.word[2] - 1
}

func (cp *pt) ansc() []*pt {
	if cp.par == nil {
		return []*pt{}
	}
	cps := []*pt{}
	cps = append(cps, cp.par...)
	for _, pa := range cp.par {
		cps = append(cps, pa.ansc()...)
	}
	return cps
}

func in(p *pt, n int) bool {
	cps := p.ansc()
	for _, cp := range cps[1:] {
		for j := 1; j < len(cp.chd); j++ {
			if tridom(cp.xy, cp.chd[j-1].xy, cp.chd[j].xy, p.xy) {
				return true
			}
		}
	}
	return false
}

func tridom(c, p0, p1, p [2]float64) bool {
	for i := 0; i < 2; i++ {
		p[i] -= c[i]
		p0[i] -= c[i]
		p1[i] -= c[i]
	}
	m := [2][2]float64{p0, p1}
	d := m[0][0]*m[1][1] - m[0][1]*m[1][0]
	a1 := (m[1][1]*p[0] - m[1][0]*p[1]) / d
	a2 := (-m[0][1]*p[0] + m[0][0]*p[1]) / d
	a := a1 + a2
	if 0 <= a1 && 0 <= a2 && a <= 1 {
		return true
	}
	return false
}

func add(cps []*pt, p *pt) ([]*pt, int) {
	for i, cp := range cps {
		if cp.Id() == p.Id() {
			return cps, i
		}
	}
	return append(cps, p), -1
}
