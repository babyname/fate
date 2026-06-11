package ent

func (cc *CharacterCreate) SetCharacter(input *Character) *CharacterCreate {
	cc.SetChar(input.Char)
	cc.SetIsSimplified(input.IsSimplified)
	cc.SetIsTraditional(input.IsTraditional)
	cc.SetIsKangxi(input.IsKangxi)
	cc.SetIsVariant(input.IsVariant)
	cc.SetIsAncient(input.IsAncient)
	cc.SetRegular(input.Regular)
	cc.SetNameable(input.Nameable)
	return cc
}

func (cc *CharacterCreate) SetCharacterWithOptional(input *Character) *CharacterCreate {
	cc.SetChar(input.Char)
	cc.SetUnicode(input.Unicode)
	cc.SetIsSimplified(input.IsSimplified)
	cc.SetIsTraditional(input.IsTraditional)
	cc.SetIsKangxi(input.IsKangxi)
	cc.SetIsVariant(input.IsVariant)
	cc.SetIsAncient(input.IsAncient)
	cc.SetPinyin(input.Pinyin)
	cc.SetRadical(input.Radical)
	cc.SetRadicalStroke(input.RadicalStroke)
	cc.SetSimplifiedStroke(input.SimplifiedStroke)
	cc.SetTraditionalStroke(input.TraditionalStroke)
	cc.SetKangxiStroke(input.KangxiStroke)
	cc.SetScienceStroke(input.ScienceStroke)
	cc.SetWuXing(input.WuXing)
	cc.SetRegular(input.Regular)
	cc.SetCommonLevel(input.CommonLevel)
	cc.SetGenderHint(input.GenderHint)
	cc.SetNameable(input.Nameable)
	cc.SetMeaning(input.Meaning)
	cc.SetSource(input.Source)
	cc.SetSourceConfidence(input.SourceConfidence)
	cc.SetComment(input.Comment)
	return cc
}

func (cuo *CharacterUpdateOne) SetCharacter(input *Character) *CharacterUpdateOne {
	cuo.SetChar(input.Char)
	cuo.SetIsSimplified(input.IsSimplified)
	cuo.SetIsTraditional(input.IsTraditional)
	cuo.SetIsKangxi(input.IsKangxi)
	cuo.SetIsVariant(input.IsVariant)
	cuo.SetIsAncient(input.IsAncient)
	cuo.SetRegular(input.Regular)
	cuo.SetNameable(input.Nameable)
	return cuo
}

func (cu *CharacterUpdate) SetCharacter(input *Character) *CharacterUpdate {
	cu.SetChar(input.Char)
	cu.SetIsSimplified(input.IsSimplified)
	cu.SetIsTraditional(input.IsTraditional)
	cu.SetIsKangxi(input.IsKangxi)
	cu.SetIsVariant(input.IsVariant)
	cu.SetIsAncient(input.IsAncient)
	cu.SetRegular(input.Regular)
	cu.SetNameable(input.Nameable)
	return cu
}

func (cu *CharacterUpdate) SetCharacterWithOptional(input *Character) *CharacterUpdate {
	cu.SetChar(input.Char)
	cu.SetUnicode(input.Unicode)
	cu.SetIsSimplified(input.IsSimplified)
	cu.SetIsTraditional(input.IsTraditional)
	cu.SetIsKangxi(input.IsKangxi)
	cu.SetIsVariant(input.IsVariant)
	cu.SetIsAncient(input.IsAncient)
	cu.SetPinyin(input.Pinyin)
	cu.SetRadical(input.Radical)
	cu.SetRadicalStroke(input.RadicalStroke)
	cu.SetSimplifiedStroke(input.SimplifiedStroke)
	cu.SetTraditionalStroke(input.TraditionalStroke)
	cu.SetKangxiStroke(input.KangxiStroke)
	cu.SetScienceStroke(input.ScienceStroke)
	cu.SetWuXing(input.WuXing)
	cu.SetRegular(input.Regular)
	cu.SetCommonLevel(input.CommonLevel)
	cu.SetGenderHint(input.GenderHint)
	cu.SetNameable(input.Nameable)
	cu.SetMeaning(input.Meaning)
	cu.SetSource(input.Source)
	cu.SetSourceConfidence(input.SourceConfidence)
	cu.SetComment(input.Comment)
	return cu
}

func (vc *VersionCreate) SetVersion(input *Version) *VersionCreate {
	vc.SetCurrentVersion(input.CurrentVersion)
	vc.SetUpdatedUnix(input.UpdatedUnix)
	return vc
}

func (vc *VersionCreate) SetVersionWithOptional(input *Version) *VersionCreate {
	vc.SetCurrentVersion(input.CurrentVersion)
	vc.SetUpdatedUnix(input.UpdatedUnix)
	return vc
}

func (vuo *VersionUpdateOne) SetVersion(input *Version) *VersionUpdateOne {
	vuo.SetCurrentVersion(input.CurrentVersion)
	vuo.SetUpdatedUnix(input.UpdatedUnix)
	return vuo
}

func (vu *VersionUpdate) SetVersion(input *Version) *VersionUpdate {
	vu.SetCurrentVersion(input.CurrentVersion)
	vu.SetUpdatedUnix(input.UpdatedUnix)
	return vu
}

func (vu *VersionUpdate) SetVersionWithOptional(input *Version) *VersionUpdate {
	vu.SetCurrentVersion(input.CurrentVersion)
	vu.SetUpdatedUnix(input.UpdatedUnix)
	return vu
}
