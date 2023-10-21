# Self-replicating words (and Rauzy fractals)

This is a Go implementation of the self-replicating words and Rauzy fractals.
Details on mathematical settings can be found [here](https://arXiv:2308.10213).


## How to use

### Self-replicating words

```go
package main

import "github.com/hyosangkang/srw"

furn main() {
	f := NewFractal(6, []int{0, 1, 0, 2})
	f.Save("test.png")
}
```

`NewFractal(n, seq)` takes two arguments. `n` is the number of layers (`int`) and `seq` is the sequence of letters (`[]int`) in the initial word.
To visualize the fractal, use `Save(fn)` where `fn` is the file name (`string`) to save the image.

### Rauzy fractal

```go
package main

import "github.com/hyosangkang/srw"

furn main() {
	f := NewFractal(6, []int{})
	f.SetRauzy() 
    f.SetBdd() // to show the boundary only
	f.Save("test.png")
}
```

`SetRauzy()` forces the fractal to be the Rauzy fractal. 
With `SetBdd()`, only the boundary points in the Rauzy fractal are shown in the image.
