package fate

import (
	"testing"
	"time"

	"github.com/babyname/fate/config"
	_ "github.com/sqlite3ent/sqlite3"
)

func TestNew(t *testing.T) {
	type args struct {
		cfg *config.Config
	}
	tests := []struct {
		name    string
		args    args
		nowant  Fate
		wantErr bool
	}{
		{
			name: "",
			args: args{
				cfg: config.DefaultConfig(),
			},
			nowant:  nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == tt.nowant {
				t.Errorf("New() got = %v, nowant %v", got, tt.nowant)
			}
			session := got.NewSessionWithFilter(NewFilter(FilterOption{
				CharacterFilter:     true,
				CharacterFilterType: 0,
				MinCharacter:        3,
				MaxCharacter:        18,
				SexFilter:           true,
			}))
			err = session.Start(&Input{
				Last: [2]string{"章"},
				Born: time.Now(),
				Sex:  1,
			})
			if err != nil {
				t.Fatal(err)
			}

		})
	}
}
