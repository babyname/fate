//go:build integration

package fate_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/babyname/fate/v4"
	"github.com/babyname/fate/v4/config"
	"github.com/babyname/fate/v4/ent"
	"github.com/babyname/fate/v4/ranking"
)

func TestScoreNameMatchesGeneration(t *testing.T) {
	engine, err := fate.New(config.DefaultConfig())
	if err != nil {
		t.Fatalf("open the local character database: %v", err)
	}
	weight, err := ranking.NewWeight(65)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	req := fate.ScoreNameRequest{
		LastName:             "张",
		Born:                 time.Date(2024, time.June, 15, 10, 30, 0, 0, time.UTC),
		Sex:                  fate.SexBoy,
		FilterType:           fate.CharacterFilterTypeKangxi,
		XiYongMethod:         "geju",
		RecommendationWeight: &weight,
	}
	filter := fate.DefaultFilterOption()
	filter.CharacterFilterType = req.FilterType
	filter.XiYongMethod = req.XiYongMethod
	generated, err := engine.Generate(ctx, fate.GenerateRequest{
		LastName:             req.LastName,
		Born:                 req.Born,
		Sex:                  req.Sex,
		Filter:               &filter,
		RecommendationWeight: req.RecommendationWeight,
	})
	if err != nil || len(generated.TopNames) == 0 {
		t.Fatalf("generate weighted comparison names: %v", err)
	}
	for _, want := range generated.TopNames {
		req.FirstName = want.FirstName
		got, err := fate.ScoreName(ctx, engine, req)
		if err != nil {
			t.Fatalf("score %s: %v", want.FullName, err)
		}
		if got.Rank != 0 || got.Ranking == nil {
			t.Fatalf("%s: got rank=%d, weighted scores=%+v", got.FullName, got.Rank, got.Ranking)
		}
		want.Rank = 0
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("direct result for %s differs from generated result\ndirect: %+v\ngenerated: %+v", got.FullName, got, want)
		}
	}
	t.Logf("direct scores and full analyses matched %d weighted catalogue names", len(generated.TopNames))

	req.LastName = "欧阳"
	req.FirstName = "安安"
	req.RecommendationWeight = nil
	got, err := fate.ScoreName(ctx, engine, req)
	if err != nil || got.FullName != "欧阳安安" || got.Ranking != nil || got.WuGe == nil {
		t.Fatalf("score repeated characters and a compound surname with the default weights: %+v, %v", got, err)
	}
	t.Logf("default scoring returned %.1f and five-grid details for %s", got.Score, got.FullName)

	req.FirstName = "安"
	if _, err := fate.ScoreName(ctx, engine, req); err == nil {
		t.Error("a one-character given name unexpectedly passed validation")
	}
	req.FirstName = "安\U00031350"
	if _, err := fate.ScoreName(ctx, engine, req); !ent.IsNotFound(err) {
		t.Errorf("unknown given name character error = %v, want character not found", err)
	}
}
