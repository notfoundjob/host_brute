package common

import (
	"fmt"
	"strings"
)

func getPercent(cur, total int64) int64 {
	return int64(float32(cur) / float32(total) * 100)
}

func (bar *Bar) NewOption(cur, total int64) {
	bar.total = total
}

func (bar *Bar) Play(cur int64) {
	percent := getPercent(cur, bar.total)
	fmt.Printf("\r[%-50s]%3d%% %8d/%d", strings.Repeat("█", int(percent/2)), percent, cur, bar.total)
}
