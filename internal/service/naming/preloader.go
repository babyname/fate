package naming

import (
	"context"

	"github.com/babyname/fate/v4/ent"
	"github.com/babyname/fate/v4/internal/filter"
	"github.com/babyname/fate/v4/internal/log"
	"github.com/babyname/fate/v4/internal/repository"
	"github.com/babyname/fate/v4/internal/wuge"
)

// preloadChars loads characters for all stroke values needed by the lucky
// wuge combinations. Each unique stroke value is queried once and cached.
// Strokes that fail the filter's stroke-number scope are skipped.
func preloadChars(
	ctx context.Context,
	db *repository.Repository,
	lucky []wuge.WuGeResult,
	flt filter.Filter,
) map[int][]*ent.Character {
	strokesNeeded := make(map[int]struct{})
	for i := range lucky {
		tmp := &lucky[i]
		if !flt.CheckSkipStrokeNumberScope(tmp.FirstStroke1, tmp.FirstStroke2) &&
			!flt.CheckSkipSexFilter(tmp) &&
			!flt.CheckSkipDaYanFilter(tmp) &&
			!flt.CheckSkipWuXingFilter(tmp.TianGe, tmp.RenGe, tmp.DiGe) {
			strokesNeeded[tmp.FirstStroke1] = struct{}{}
			strokesNeeded[tmp.FirstStroke2] = struct{}{}
		}
	}

	chars := make(map[int][]*ent.Character, len(strokesNeeded))
	for stroke := range strokesNeeded {
		cs, err := db.GetCharactersCached(ctx, stroke, flt.QueryStrokeFilter(stroke), flt.QueryRegularFilter)
		if err != nil {
			log.Error("preload characters", err)
			continue
		}
		chars[stroke] = cs
	}
	return chars
}
