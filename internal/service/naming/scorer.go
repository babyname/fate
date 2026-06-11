package naming

import (
	"context"

	"github.com/babyname/fate/v4/ent"
	"github.com/babyname/fate/v4/internal/filter"
	"github.com/babyname/fate/v4/internal/log"
	"github.com/babyname/fate/v4/internal/rating"
	"github.com/babyname/fate/v4/internal/wuge"
)

// runScoring evaluates all candidate name pairs derived from lucky wuge
// combinations and pushes qualifying entries into the excellent table.
//
// For each lucky combination, it iterates through characters with matching
// stroke counts and rates each pair. Poetry mode 2 requires both characters
// to appear in the poetry set.
//
// The context is checked on every character-pair iteration to support
// cancellation.
func runScoring(
	ctx context.Context,
	lucky []wuge.WuGeResult,
	rater *rating.Rater,
	table *ExcellentTable,
	poetrySet map[string]bool,
	chars map[int][]*ent.Character,
	flt filter.Filter,
) {
	poetryMode := flt.PoetryMode()

	for i := range lucky {
		tmp := &lucky[i]
		if flt.CheckSkipStrokeNumberScope(tmp.FirstStroke1, tmp.FirstStroke2) {
			continue
		}
		if flt.CheckSkipSexFilter(tmp) {
			continue
		}
		if flt.CheckSkipDaYanFilter(tmp) {
			continue
		}
		if flt.CheckSkipWuXingFilter(tmp.TianGe, tmp.RenGe, tmp.DiGe) {
			continue
		}

		f1s := chars[tmp.FirstStroke1]
		f2s := chars[tmp.FirstStroke2]

		for i1 := range f1s {
			c1 := f1s[i1]
			if poetryMode == 2 && !poetrySet[c1.Char] {
				continue
			}
			for i2 := range f2s {
				select {
				case <-ctx.Done():
					return
				default:
				}

				c2 := f2s[i2]
				if poetryMode == 2 && !poetrySet[c2.Char] {
					continue
				}

				nr := rater.RateName("", c1, c2)

				hasPoetry := poetrySet[c1.Char] || poetrySet[c2.Char]

				table.TryPush(ExcellentEntry{
					Char1:     c1.Char,
					Char2:     c2.Char,
					Score:     nr.TotalScore,
					Grade:     nr.Grade,
					WuXing1:   c1.WuXing,
					WuXing2:   c2.WuXing,
					HasPoetry: hasPoetry,
				})
			}
		}
	}
	log.Info("score candidates done", "heap_size", table.HeapLen())
}
