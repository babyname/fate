package fate

// DefaultFilterOption returns the bounded filter option used by Generate and
// GenerateStream when GenerateRequest.Filter is nil.
//
// The shared DefaultFilter() is deliberately permissive: its stroke range is
// effectively unbounded, which makes the candidate cross-product in the
// scoring pipeline intractable. A measured run with DefaultFilter() did not
// converge after 10 minutes (>2000s CPU, ~1.1 GiB RSS); the bounded option
// below converges in 11-19s. Generate therefore never falls back to
// DefaultFilter() -- doing so would make a trivially valid request unusable.
//
// The bounds match the Web form and internal/service/naming pipeline tests,
// so the library, the HTTP handler and the examples all score the same
// candidate space.
//
// Use it as the base for custom options:
//
//	opt := fate.DefaultFilterOption().WithAvoidCharacters("病", "死")
//	result, err := engine.Generate(ctx, fate.GenerateRequest{Filter: &opt, ...})
func DefaultFilterOption() FilterOption {
	return FilterOption{
		CharacterFilter:     true,
		CharacterFilterType: CharacterFilterTypeDefault,
		MinStroke:           3,
		MaxStroke:           18,
		RegularFilter:       true,
		DaYanFilter:         true,
		WuXingFilter:        true,
	}
}

// resolveFilterOption is the single decision point for the nil-Filter policy.
// Generate and GenerateStream both go through it so they cannot drift apart.
func resolveFilterOption(opt *FilterOption) FilterOption {
	if opt == nil {
		return DefaultFilterOption()
	}
	return *opt
}
