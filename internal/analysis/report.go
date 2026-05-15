// Package analysis provides name analysis and report generation.
package analysis

import (
	"time"

	v2 "github.com/godcong/chronos/v2"
)

// NewReport creates a new FateReport populated with the given surname, birth date,
// sex, fate data from chronos, and the total number of candidate names.
// It extracts bazi and wuxing sections from fateData when available.
func NewReport(surname, born, sex string, fateData *v2.FateData, totalNames int) *FateReport {
	report := &FateReport{
		GeneratedAt: time.Now().Format("2006-01-02 15:04:05"),
		Surname:     surname,
		Born:        born,
		Sex:         sex,
		TotalNames:  totalNames,
	}

	if fateData != nil {
		if fateData.Bazi != nil {
			report.Bazi = &BaziSection{
				Sizhu:         fateData.Bazi.Sizhu,
				Wuxing:        fateData.Bazi.Wuxing,
				Nayin:         fateData.Bazi.Nayin,
				Zodiac:        fateData.Bazi.Zodiac,
				Constellation: fateData.Bazi.Constellation,
			}
		}
		if fateData.WuxingXiji != nil {
			report.WuXing = &WuXingSection{
				DayGan:     fateData.WuxingXiji.DayGan,
				DayWuxing:  fateData.WuxingXiji.DayWuxing,
				QiangRuo:   fateData.WuxingXiji.QiangRuo,
				XiWuxing:   fateData.WuxingXiji.XiWuxing,
				YongWuxing: fateData.WuxingXiji.YongWuxing,
				JiWuxing:   fateData.WuxingXiji.JiWuxing,
				ChouWuxing: fateData.WuxingXiji.ChouWuxing,
				XianWuxing: fateData.WuxingXiji.XianWuxing,
				Method:     fateData.WuxingXiji.MethodName,
				MethodName: fateData.WuxingXiji.MethodName,
				Analysis:   fateData.WuxingXiji.Analysis,
			}
			if fateData.WuxingXiji.GeJu != nil {
				report.WuXing.GeJuName = fateData.WuxingXiji.GeJu.Name
			}
		}
	}

	return report
}
