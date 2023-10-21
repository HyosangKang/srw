package srw

import (
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
)

type fractal struct {
	numL       int // number of layers, number of bdd points
	seq, bseq  []int
	wvec       map[int][][3]int
	bvec       map[int][][2]float64
	pts, bpts  map[int][]*pt
	bdd, rauzy bool
}

func NewFractal(n int, seq []int) *fractal {
	f := &fractal{
		numL: n,
		seq:  seq,
	}
	// By default, boundary sequence is same as the morphic sequence
	f.bseq = []int{}
	for i := 0; i < len(seq); i++ {
		f.bseq = append(f.bseq, i)
	}
	// The first points is the origin.
	f.pts = make(map[int][]*pt)
	o := &pt{
		xy:   [2]float64{0, 0},
		par:  nil,
		word: [3]int{0, 0, 0},
	}
	f.pts[-1] = []*pt{o}
	f.bpts = make(map[int][]*pt)
	f.bpts[-1] = []*pt{o}
	return f
}

func (f *fractal) SetRauzy() {
	f.rauzy = true
	f.seq = []int{0, 1, 0, 2, 0, 1, 0}
	f.bseq = []int{0, 2, 1, 5, 3, 4, 0}
}

func (f *fractal) SetBdd() {
	f.bdd = true
}

func (f *fractal) Save(fn string) {
	f.init(f.numL)
	p := plot.New()
	noAxis(p)
	f.draw(p)
	axisEqual(p)
	p.Save(800, 800, fn)
}

func (f *fractal) init(n int) {
	f.wvec = make(map[int][][3]int)
	f.wvec[-1] = [][3]int{}
	for _, i := range f.seq {
		vi := [3]int{}
		vi[i] = 1
		f.wvec[-1] = append(f.wvec[-1], vi)
	}
	var bvec [][2]float64
	for i := 0; i < len(f.seq); i++ {
		bvec = append(bvec, [2]float64{0, 0})
	}
	f.bvec = make(map[int][][2]float64)
	f.bvec[-1] = bvec
	for i := 0; i < 20; i++ {
		f.morph(i)
	}
	wvec := f.wvec[19][0]
	ev := nmlz([3]float64{float64(wvec[0]), float64(wvec[1]), float64(wvec[2])})
	var basis [][2]float64
	e1 := nmlz(oproj([3]float64{1, 0, 0}, ev))
	e2 := oproj([3]float64{0, 1, 0}, ev)
	e2 = nmlz(oproj(e2, e1))
	for i := 0; i < 3; i++ {
		v := [3]float64{}
		v[i] = 1
		basis = append(basis, [2]float64{dot(v, e1), dot(v, e2)})
	}
	bvec = [][2]float64{}
	for _, i := range f.seq {
		bvec = append(bvec, basis[i])
	}
	f.bvec[-1] = bvec
	f.numL = n
	for i := 0; i < n; i++ {
		f.layer(i)
		f.morph(i)
	}
}

func (f *fractal) layer(n int) {
	f.pts[n] = []*pt{}
	for _, p := range f.pts[n-1] {
		var ps []*pt
		xy := p.xy
		wv := p.word
		for i := 0; i < len(f.seq); i++ {
			for j := 0; j < 2; j++ {
				xy[j] += f.bvec[n-1][i][j]
			}
			for j := 0; j < 3; j++ {
				wv[j] += f.wvec[n-1][i][j]
			}
			np := &pt{
				xy:   xy,
				par:  []*pt{p},
				word: wv,
			}
			var id int
			f.pts[n], id = add(f.pts[n], np)
			if id > -1 {
				np = f.pts[n][id]
				np.par = append(np.par, p)
			}
			ps = append(ps, np)
		}
		p.chd = []*pt{}
		for _, i := range f.bseq {
			p.chd = append(p.chd, ps[i])
		}
	}
	if n == 0 {
		f.bpts[0] = f.bpts[-1][0].chd[:len(f.bseq)-1]
	} else {
		f.bpts[n] = []*pt{}
		for _, bp := range f.bpts[n-1] {
			for _, cp := range bp.chd {
				if !in(cp, n) {
					f.bpts[n], _ = add(f.bpts[n], cp)
				}
			}
		}
	}
}

func (f *fractal) morph(n int) {
	f.bvec[n] = [][2]float64{}
	f.wvec[n] = [][3]int{}
	for _, i := range f.seq {
		var bvec [2]float64
		var wvec [3]int
		m := i
		if f.rauzy && i == 2 {
			m = 3
		}
		for k := 0; k < len(f.seq)-m; k++ {
			for j := 0; j < 2; j++ {
				bvec[j] += f.bvec[n-1][k][j]
			}
			for j := 0; j < 3; j++ {
				wvec[j] += f.wvec[n-1][k][j]
			}
		}
		f.bvec[n] = append(f.bvec[n], bvec)
		f.wvec[n] = append(f.wvec[n], wvec)
	}
}

func (f *fractal) draw(p *plot.Plot) plotter.XYs {
	var xys plotter.XYs
	if f.bdd && f.rauzy {
		for _, point := range f.bpts[f.numL-1] {
			xy := point.plotter()
			xys = append(xys, xy)
		}
	} else {
		for i := 0; i < f.numL; i++ {
			for _, point := range f.pts[i] {
				xy := point.plotter()
				xys = append(xys, xy)
			}
		}
	}
	s, _ := plotter.NewScatter(xys)
	s.GlyphStyle.Shape = draw.CircleGlyph{}
	p.Add(s)

	return xys
}

func noAxis(p *plot.Plot) {
	p.X.LineStyle.Width = 0
	p.X.Tick.LineStyle.Width = 0
	p.Y.LineStyle.Width = 0
	p.Y.Tick.LineStyle.Width = 0
	p.X.Tick.Label.Font.Size = 0
	p.Y.Tick.Label.Font.Size = 0
}

func axisEqual(p *plot.Plot) {
	xl := p.X.Max - p.X.Min
	yl := p.Y.Max - p.Y.Min
	if xl > yl {
		p.Y.Max += (xl - yl) / 2
		p.Y.Min -= (xl - yl) / 2
	} else {
		p.X.Max += (yl - xl) / 2
		p.X.Min -= (yl - xl) / 2
	}
}

func oproj(v, w [3]float64) [3]float64 {
	a := [3]float64{}
	for i, vv := range projection(v, w) {
		a[i] = v[i] - vv
	}
	return a
}

func projection(v, w [3]float64) []float64 {
	a := []float64{}
	k := dot(v, w) / (norm(w) * norm(w))
	for _, ww := range w {
		a = append(a, k*ww)
	}
	return a
}

func dot(v, w [3]float64) float64 {
	a := 0.0
	for i := 0; i < 3; i++ {
		a += v[i] * w[i]
	}
	return a
}

func nmlz(v [3]float64) [3]float64 {
	w := [3]float64{}
	k := norm(v)
	for i, vv := range v {
		w[i] = vv / k
	}
	return w
}

func norm(v [3]float64) float64 {
	a := 0.0
	for _, vv := range v {
		a += float64(vv * vv)
	}
	return math.Sqrt(a)
}
