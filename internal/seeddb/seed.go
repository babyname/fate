package seeddb

// SeedCharacter 种子汉字数据，用于导入到目标数据库。
type SeedCharacter struct {
	Char              string   `json:"char"`
	Unicode           string   `json:"unicode,omitempty"`
	IsSimplified      bool     `json:"is_simplified"`
	IsTraditional     bool     `json:"is_traditional"`
	IsKangxi          bool     `json:"is_kangxi"`
	IsVariant         bool     `json:"is_variant"`
	IsAncient         bool     `json:"is_ancient"`
	Pinyin            []string `json:"pinyin,omitempty"`
	Radical           string   `json:"radical,omitempty"`
	RadicalStroke     int      `json:"radical_stroke,omitempty"`
	SimplifiedStroke  int      `json:"simplified_stroke,omitempty"`
	TraditionalStroke int      `json:"traditional_stroke,omitempty"`
	KangxiStroke      int      `json:"kangxi_stroke,omitempty"`
	ScienceStroke     int      `json:"science_stroke,omitempty"`
	WuXing            string   `json:"wu_xing,omitempty"`
	Regular           bool     `json:"regular"`
	CommonLevel       int      `json:"common_level,omitempty"`
	GenderHint        string   `json:"gender_hint,omitempty"`
	Nameable          bool     `json:"nameable"`
	Meaning           string   `json:"meaning,omitempty"`
	Source            string   `json:"source,omitempty"`
	SourceConfidence  float64  `json:"source_confidence,omitempty"`
	Comment           string   `json:"comment,omitempty"`
	SimplifiedOfChar  string   `json:"simplified_of_char,omitempty"`
	VariantOfChar     string   `json:"variant_of_char,omitempty"`
}

// SeedWuGeLucky 种子五格吉凶数据。
type SeedWuGeLucky struct {
	LastStroke1  int    `json:"last_stroke_1"`
	LastStroke2  int    `json:"last_stroke_2"`
	FirstStroke1 int    `json:"first_stroke_1"`
	FirstStroke2 int    `json:"first_stroke_2"`
	TianGe       int    `json:"tian_ge"`
	TianDaYan    string `json:"tian_da_yan"`
	RenGe        int    `json:"ren_ge"`
	RenDaYan     string `json:"ren_da_yan"`
	DiGe         int    `json:"di_ge"`
	DiDaYan      string `json:"di_da_yan"`
	WaiGe        int    `json:"wai_ge"`
	WaiDaYan     string `json:"wai_da_yan"`
	ZongGe       int    `json:"zong_ge"`
	ZongDaYan    string `json:"zong_da_yan"`
	ZongLucky    bool   `json:"zong_lucky"`
	ZongSex      bool   `json:"zong_sex"`
	ZongMax      bool   `json:"zong_max"`
}

// SeedWuXing 种子五行数据。
type SeedWuXing struct {
	ID      int    `json:"id"`
	First   string `json:"first"`
	Second  string `json:"second"`
	Third   string `json:"third"`
	Fortune string `json:"fortune"`
}

// FieldChange 记录字段变更详情。
type FieldChange struct {
	Char     string `json:"char"`
	Field    string `json:"field"`
	OldValue string `json:"old_value"`
	NewValue string `json:"new_value"`
	Reason   string `json:"reason"`
	Source   string `json:"source"`
}

// DataReport 数据质量报告。
type DataReport struct {
	Characters CharacterReport `json:"characters"`
	WuGeLucky  WuGeLuckyReport `json:"wu_ge_lucky"`
	WuXing     WuXingReport    `json:"wu_xing"`
}

// CharacterReport 汉字数据统计报告。
type CharacterReport struct {
	Total            int            `json:"total"`
	WithWuXing       int            `json:"with_wu_xing"`
	WithoutWuXing    int            `json:"without_wu_xing"`
	WuXingCoverage   float64        `json:"wu_xing_coverage"`
	WithPinyin       int            `json:"with_pinyin"`
	PinyinCoverage   float64        `json:"pinyin_coverage"`
	RegularCount     int            `json:"regular_count"`
	NameableCount    int            `json:"nameable_count"`
	SimplifiedCount  int            `json:"simplified_count"`
	TraditionalCount int            `json:"traditional_count"`
	KangxiCount      int            `json:"kangxi_count"`
	VariantCount     int            `json:"variant_count"`
	WuXingDist       map[string]int `json:"wu_xing_distribution"`
	StrokeIssues     []StrokeIssue  `json:"stroke_issues,omitempty"`
}

// StrokeIssue 笔画数据问题记录。
type StrokeIssue struct {
	Char    string `json:"char"`
	Field   string `json:"field"`
	Value   int    `json:"value"`
	Message string `json:"message"`
}

// WuGeLuckyReport 五格吉凶统计报告。
type WuGeLuckyReport struct {
	Total      int     `json:"total"`
	LuckyCount int     `json:"lucky_count"`
	LuckyRate  float64 `json:"lucky_rate"`
	MaxCount   int     `json:"max_count"`
	SexCount   int     `json:"sex_count"`
}

// WuXingReport 五行统计报告。
type WuXingReport struct {
	Total        int            `json:"total"`
	LuckyCount   int            `json:"lucky_count"`
	UnluckyCount int            `json:"unlucky_count"`
	FortuneDist  map[string]int `json:"fortune_distribution"`
}

// Exporter 从源数据库导出数据为种子 JSON 文件的导出器。
type Exporter struct {
	dbPath     string
	seedDir    string
	rawDataDir string
	changes    []FieldChange

	pinyinMap    map[string][]string
	totalStrokes map[string]int
	definitions  map[string]string
	wuxingMap    map[string]string
	rsUnicode    map[string]int
}

// Importer 将种子 JSON 数据导入到目标数据库的导入器。
type Importer struct {
	seedDir string
	cfg     DBConfig
}

// Reporter 生成数据质量报告的报告器。
type Reporter struct {
	seedDir string
}

// DBConfig 数据库连接配置。
type DBConfig struct {
	Driver string
	DSN    string
	Host   string
	Port   string
	User   string
	Pwd    string
	Name   string
}

// NewExporter 创建数据导出器实例。
func NewExporter(dbPath, seedDir string, rawDataDirs ...string) *Exporter {
	rawDataDir := "data/raw"
	if len(rawDataDirs) > 0 {
		rawDataDir = rawDataDirs[0]
	}
	return &Exporter{
		dbPath:       dbPath,
		seedDir:      seedDir,
		rawDataDir:   rawDataDir,
		pinyinMap:    make(map[string][]string),
		totalStrokes: make(map[string]int),
		definitions:  make(map[string]string),
		wuxingMap:    make(map[string]string),
		rsUnicode:    make(map[string]int),
	}
}

// NewImporter 创建数据导入器实例。
func NewImporter(seedDir string, cfg DBConfig) *Importer {
	return &Importer{seedDir: seedDir, cfg: cfg}
}

// NewReporter 创建数据报告器实例。
func NewReporter(seedDir string) *Reporter {
	return &Reporter{seedDir: seedDir}
}

func (e *Exporter) recordChange(char, field, oldVal, newVal, reason, source string) {
	e.changes = append(e.changes, FieldChange{
		Char:     char,
		Field:    field,
		OldValue: oldVal,
		NewValue: newVal,
		Reason:   reason,
		Source:   source,
	})
}
