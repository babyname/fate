package chronosfate

import (
	"github.com/godcong/chronos/v2"
)

func GetFateData(input FateInput) (*FateData, error) {
	if input.Calendar == nil {
		return nil, FateError{Code: ErrCodeInputInvalid, Message: "calendar is nil"}
	}

	lunar := input.Calendar.Lunar()
	solar := input.Calendar.Solar()

	baziInfo, err := calculateBazi(lunar, solar, input.Gender)
	if err != nil {
		return nil, err
	}

	wuxingXiji := calculateWuxingXiji(lunar)
	wuxingStrength := calculateWuxingStrength(lunar)

	data := &FateData{
		BaziInfo:       baziInfo,
		WuxingXijiInfo: wuxingXiji,
		WuxingStrength: wuxingStrength,
	}

	geJuInfo := determineGeJu(baziInfo, wuxingStrength)
	data.GeJuInfo = geJuInfo

	switch input.XiYongMethod {
	case XiYongMethodGeJu:
		if geJuInfo != nil {
			data.XiYongJiChou = XiYongJiChou{
				Xi:   geJuInfo.XiShen,
				Yong: geJuInfo.YongShen,
				Ji:   geJuInfo.JiShen,
				Chou: geJuInfo.ChouShen,
			}
		}
	default:
		data.XiYongJiChou = balanceXiYongJi(wuxingStrength, baziInfo)
	}

	return data, nil
}

func calculateBazi(lunar chronos.Lunar, solar chronos.Solar, gender int) (BaziInfo, error) {
	eightChar := lunar.GetEightChar()
	zodiacObj := lunar.GetZodiac()
	constellationObj := solar.GetConstellation()

	return BaziInfo{
		SiZhu:        eightChar.FourPillars(),
		WuXing:       eightChar.FiveElements(),
		NaYin:        eightChar.NaYin(),
		ShiShenGan:   eightChar.TenGodsStems(),
		ShiShenZhi:   eightChar.TenGodsBranches(),
		CangGan:      eightChar.HiddenStems(),
		DaYun:        eightChar.DaYun(gender),
		Zodiac:       zodiacObj.Chinese(),
		Constellation: constellationObj.Chinese(),
	}, nil
}
