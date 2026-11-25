package ranking

import "testing"

func TestCombineUsesPublishedComponents(t *testing.T) {
	for _, tc := range []struct {
		weight int
		want   float64
	}{{0, 88.5}, {40, 80.9}, {100, 69.6}} {
		weight, err := NewWeight(tc.weight)
		if err != nil {
			t.Fatal(err)
		}
		result := weight.Combine(69.56, 88.48)
		if result.Total != tc.want || result.Traditional != 69.6 || result.SoundAndMeaning != 88.5 || result.TraditionalWeight != tc.weight {
			t.Errorf("weight %d: result = %+v, want total %.1f and rounded components", tc.weight, result, tc.want)
		}
	}
}

func TestNewWeightRejectsOutsidePercentageRange(t *testing.T) {
	for _, weight := range []int{-1, 101} {
		if _, err := NewWeight(weight); err == nil {
			t.Errorf("NewWeight(%d) succeeded with an invalid percentage", weight)
		}
	}
}
