package scripts

import (
	"fmt"
	"log/slog"
	"strings"

	"golang.org/x/net/context"

	"github.com/babyname/fate/ent"
	"github.com/babyname/fate/ent/ncharacter"
)

const perLimit = 500

type Fix struct {
	ID      int      `json:"id"`
	Char    string   `json:"char"`
	FixType []string `json:"fix"`
}

func (f Fix) String() string {
	return fmt.Sprintf("id: %v, char: %v, fix: [%v]", f.ID, f.Char, strings.Join(f.FixType, ","))
}

func NeedFix(ctx context.Context, client *ent.Client, hook func(Fix) bool) error {
	count, err := client.NCharacter.Query().Count(ctx)
	if err != nil {
		return err
	}

	var ncs []*ent.NCharacter
	for i := 0; i < count; i += perLimit {
		slog.Info("update character", "offset", i)
		ncs, err = client.NCharacter.Query().Offset(i).Limit(perLimit).All(ctx)
		if err != nil {
			slog.Info("found error on", "offset", i, "limit", perLimit, "error", err)
			continue
		}

		for _, nc := range ncs {
			fix := Fix{
				ID:   nc.ID,
				Char: nc.Char,
			}
			if len(nc.PinYin) == 0 {
				fix.FixType = append(fix.FixType, ncharacter.FieldPinYin)
			}
			if nc.CharStroke == 0 {
				fix.FixType = append(fix.FixType, ncharacter.FieldCharStroke)
			}
			if len(nc.WuXing) == 0 {
				fix.FixType = append(fix.FixType, ncharacter.FieldWuXing)
			}
			if len(nc.Lucky) == 0 {
				fix.FixType = append(fix.FixType, ncharacter.FieldLucky)
			}
			if len(fix.FixType) != 0 {
				if !hook(fix) {
					return nil
				}
			}
		}
	}
	return nil
}
