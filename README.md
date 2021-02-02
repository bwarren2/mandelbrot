# Intro

This project makes a mandelbrot image, or a gif of zooming in/out to a mandelbrot.

It is a toy go project for learning.  My major learning objectives were:
1. General go technical familiarity.
2. Project layout
3. Using Cobra for CLI generation

# Test
`go test ./...`

If you are making changes and want to rewrite the golden test files, then add a flag:

`go test -write-file`
# Run

```
go build cmd/main.go && ./main  new --filename foo.png
```

# Gif
```
go build cmd/main.go && ./main gif  --height 45 --width 45 --maxIterations 10 --frames 10 --scaleIn .98 --x -1.5 --filename out.gif
```
