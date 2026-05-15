// Package config 提供命运起名系统的配置加载、校验与序列化功能。
package config

import (
	"errors"
	"fmt"
	"math"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 总配置
type Config struct {
	Database DBConfig     `yaml:"database"`
	Rate     RateConfig   `yaml:"rate"`
	Filter   FilterConfig `yaml:"filter"`
	Output   OutputConfig `yaml:"output"`
}

// DatabaseConfig 数据库配置 (deprecated: use DBConfig)
type DatabaseConfig = DBConfig

// RateConfig 评分权重与阈值配置
type RateConfig struct {
	// 评分权重（三者之和须为 1.0）
	WuxingWeight float64 `yaml:"wuxing_weight"`
	BihuaWeight  float64 `yaml:"bihua_weight"`
	YinyunWeight float64 `yaml:"yinyun_weight"`

	// 强弱判断阈值
	StrengthThresholdLow  float64 `yaml:"strength_threshold_low"`
	StrengthThresholdHigh float64 `yaml:"strength_threshold_high"`

	// 权重参数
	TongleiWeight  float64 `yaml:"tonglei_weight"`
	ShengwoWeight  float64 `yaml:"shengwo_weight"`
	KewoWeight     float64 `yaml:"kewo_weight"`
	XiewoWeight    float64 `yaml:"xiewo_weight"`
	HaowoWeight    float64 `yaml:"haowo_weight"`
	CangganWeight  float64 `yaml:"canggan_weight"`
	LingChenWeight float64 `yaml:"lingchen_weight"`
}

// FilterConfig 筛选配置
type FilterConfig struct {
	// 常用字筛选
	EnableCommonFilter bool `yaml:"enable_common_filter"`
	CommonLevel        int  `yaml:"common_level"`

	// 笔画范围
	MinStroke int `yaml:"min_stroke"`
	MaxStroke int `yaml:"max_stroke"`

	// 其他筛选
	EnableSexFilter     bool `yaml:"enable_sex_filter"`
	EnableDayanFilter   bool `yaml:"enable_dayan_filter"`
	EnableRegularFilter bool `yaml:"enable_regular_filter"`

	// StrokeMode 笔画计算模式: "science"(科学, 默认), "simplified"(简体), "traditional"(繁体), "kangxi"(康熙)
	StrokeMode string `yaml:"stroke_mode"`

	// NameStyle 输出姓名风格: "simplified"(简体, 默认), "traditional"(繁体)
	NameStyle string `yaml:"name_style"`

	// GenderFilter 性别筛选: "male", "female", "neutral", ""(全部)
	GenderFilter string `yaml:"gender_filter"`

	// PreferredWuXing 偏好的五行元素列表
	PreferredWuXing []string `yaml:"preferred_wuxing"`

	// AvoidWuXing 需要避开的五行元素列表
	AvoidWuXing []string `yaml:"avoid_wuxing"`

	// AvoidCharacters 需要避开的汉字列表
	AvoidCharacters []string `yaml:"avoid_characters"`

	// PreferredRadicals 偏好的部首列表
	PreferredRadicals []string `yaml:"preferred_radicals"`

	// AvoidRadicals 需要避开的部首列表
	AvoidRadicals []string `yaml:"avoid_radicals"`
}

// OutputConfig 输出格式与数量配置
type OutputConfig struct {
	Format          string `yaml:"format"`
	MaxResults      int    `yaml:"max_results"`
	IncludeAnalysis bool   `yaml:"include_analysis"`
}

// DefaultConfig 返回带默认值的配置实例
func DefaultConfig() *Config {
	return &Config{
		Database: DBConfig{
			Name:   "fate",
			Driver: "sqlite3",
		},
		Rate: RateConfig{
			WuxingWeight:          0.4,
			BihuaWeight:           0.3,
			YinyunWeight:          0.3,
			StrengthThresholdLow:  0.4,
			StrengthThresholdHigh: 0.6,
			TongleiWeight:         1.0,
			ShengwoWeight:         0.8,
			KewoWeight:            1.0,
			XiewoWeight:           0.8,
			HaowoWeight:           0.7,
			CangganWeight:         0.5,
			LingChenWeight:        1.5,
		},
		Filter: FilterConfig{
			EnableCommonFilter:  true,
			CommonLevel:         3,
			MinStroke:           5,
			MaxStroke:           25,
			EnableSexFilter:     false,
			EnableDayanFilter:   true,
			EnableRegularFilter: true,
			StrokeMode:          "science",
			NameStyle:           "simplified",
			GenderFilter:        "",
			PreferredWuXing:     nil,
			AvoidWuXing:         nil,
			AvoidCharacters:     nil,
			PreferredRadicals:   nil,
			AvoidRadicals:       nil,
		},
		Output: OutputConfig{
			Format:          "text",
			MaxResults:      100,
			IncludeAnalysis: true,
		},
	}
}

// LoadConfig 从指定路径加载配置，文件不存在则返回默认配置
func LoadConfig(path string) (*Config, error) {
	cfg := DefaultConfig()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return cfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	if err := ValidateConfig(cfg); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	return cfg, nil
}

// ValidateConfig 校验配置项的合法性
func ValidateConfig(cfg *Config) error {
	if cfg == nil {
		return errors.New("config is nil")
	}

	if math.Abs(cfg.Rate.WuxingWeight+cfg.Rate.BihuaWeight+cfg.Rate.YinyunWeight-1.0) > 0.001 {
		return errors.New("rate weights sum must be 1.0")
	}

	if cfg.Rate.StrengthThresholdLow <= 0 || cfg.Rate.StrengthThresholdLow >= 1 {
		return errors.New("strength_threshold_low must be between 0 and 1")
	}

	if cfg.Rate.StrengthThresholdHigh <= 0 || cfg.Rate.StrengthThresholdHigh >= 1 {
		return errors.New("strength_threshold_high must be between 0 and 1")
	}

	if cfg.Rate.StrengthThresholdLow >= cfg.Rate.StrengthThresholdHigh {
		return errors.New("strength_threshold_low must be less than strength_threshold_high")
	}

	if cfg.Filter.MinStroke < 0 {
		return errors.New("min_stroke must be >= 0")
	}

	if cfg.Filter.MaxStroke < cfg.Filter.MinStroke {
		return errors.New("max_stroke must be >= min_stroke")
	}

	validStrokeModes := map[string]bool{"science": true, "simplified": true, "traditional": true, "kangxi": true}
	if !validStrokeModes[cfg.Filter.StrokeMode] {
		return fmt.Errorf("stroke_mode must be one of science, simplified, traditional, kangxi, got %q", cfg.Filter.StrokeMode)
	}

	validNameStyles := map[string]bool{"simplified": true, "traditional": true}
	if !validNameStyles[cfg.Filter.NameStyle] {
		return fmt.Errorf("name_style must be one of simplified, traditional, got %q", cfg.Filter.NameStyle)
	}

	validGenderFilters := map[string]bool{"male": true, "female": true, "neutral": true, "": true}
	if !validGenderFilters[cfg.Filter.GenderFilter] {
		return fmt.Errorf("gender_filter must be one of male, female, neutral, or empty, got %q", cfg.Filter.GenderFilter)
	}

	if cfg.Output.MaxResults <= 0 {
		return errors.New("max_results must be > 0")
	}

	return nil
}

// SaveConfig 将配置序列化并写入指定路径
func SaveConfig(cfg *Config, path string) error {
	if cfg == nil {
		return errors.New("config is nil")
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write config file: %w", err)
	}

	return nil
}
