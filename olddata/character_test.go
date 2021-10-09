package olddata

import (
	"fmt"
	"net/url"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/xormsharp/xorm"
)

const dsn = "%s:%s@tcp(%s)/%s?loc=%s&charset=utf8mb4&parseTime=true"

var dbEngine *xorm.Engine

func init() {
	var err error
	dbURL := fmt.Sprintf(dsn, "root", "root", "localhost", "fate", url.QueryEscape("Asia/Shanghai"))
	dbEngine, err = xorm.NewEngine("mysql", dbURL)
	if err != nil {
		panic(err)
	}
}

func TestRangeCharacters(t *testing.T) {
	type args struct {
		engine *xorm.Engine
		fn     func(c *Character) bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				engine: dbEngine,
				fn: func(c *Character) bool {
					fmt.Println("Char:", *c)
					return true
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := RangeCharacters(tt.args.engine, tt.args.fn); (err != nil) != tt.wantErr {
				t.Errorf("RangeCharacters() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
