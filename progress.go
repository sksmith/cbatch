package cbatch

import (
	"fmt"
	"os"
)

type progressBar struct {
	percent int64  // progress percentage
	cur     int64  // current progress
	total   int64  // total value for progress
	rate    string // the actual progress bar to be printed
	graph   string // the fill value for progress bar
}

func (bar *progressBar) New(start, total int64) {
	bar.cur = start
	bar.total = total
	if bar.graph == "" {
		bar.graph = "█"
	}
	bar.percent = bar.getPercent()
	for i := 0; i < int(bar.percent); i += 2 {
		bar.rate += bar.graph // initial progress position
	}
}

func (bar *progressBar) getPercent() int64 {
	return int64(float32(bar.cur) / float32(bar.total) * 100)
}

func (bar *progressBar) NewOptionWithGraph(start, total int64, graph string) {
	bar.graph = graph
	bar.New(start, total)
}

func (bar *progressBar) Play(cur int64) {
	bar.cur = cur
	last := bar.percent
	bar.percent = bar.getPercent()
	if bar.percent != last {
		bar.rate = ""
		for i := int64(0); i < bar.percent; i += 2 {
			bar.rate += bar.graph
		}
	}
	fmt.Fprintf(os.Stderr, "\r[%-50s]%3d%% %8d/%d", bar.rate, bar.percent, bar.cur, bar.total)
}

func (bar *progressBar) Finish() {
	fmt.Fprintf(os.Stderr, "\n")
}
