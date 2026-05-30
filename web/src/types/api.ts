export interface GenerateRequest {
  surname: string
  born: string
  sex: string
  poetry_mode?: number
  xiyong_method?: string
  avoid_chars?: string[]
  require_chars?: string[]
  filter_type?: number
}

export interface GenerateResponse {
  task_id: string
  state: string
}

export interface TaskStatusResponse {
  task_id: string
  state: string
  total?: number
  error?: string
}

export interface TaskResultResponse {
  task_id: string
  top_names: NameResult[]
  top3: NameResult[]
  top10: ExcellentEntry[]
  total: number
}

export interface ExploreResponse {
  task_id: string
  names: ExcellentEntry[]
  shown_count: number
  total: number
}

export interface NameDetailResponse {
  name_result: NameResult
}

export interface CharInfo {
  char: string
  traditional_char: string
  pinyin: string
  wu_xing: string
  simplified_stroke: number
  traditional_stroke: number
  science_stroke: number
  kangxi_stroke: number
  radical: string
  meaning: string
  is_xi_yong: boolean
}

export interface GeItem {
  name: string
  stroke: number
  lucky: string
  da_yan: string
  sky_nine: string
  yin_yang_wu_xing: string
  analysis: string
}

export interface WuGeResult {
  tian_ge: GeItem
  ren_ge: GeItem
  di_ge: GeItem
  wai_ge: GeItem
  zong_ge: GeItem
}

export interface WuXingSection {
  xi: string
  ji: string
  ri_zhu_qiang_ruo: string
  tiao_hou_shen: string
  yong_shen: string
  chou_shen: string
  geju_name: string
  analysis: string
}

export interface ZhouYiBasic {
  ben_gua_name: string
  ben_gua_desc: string
  ben_gua_ji_xiong: string
  bian_gua_name: string
  dong_yao_desc: string
  dong_yao_ji_xiong: string
  da_xiang: string
  score: number
}

export interface BaziBasic {
  four_pillars: string[]
  five_elements: string[]
  na_yin: string[]
  zodiac: string
  constellation: string
}

export interface ScoreDetail {
  wen_hua_yin_xiang: number
  wu_xing_ba_zi: number
  sheng_xiao: number
  wu_ge_shu_li: number
  yin_yun: number
}

export interface NameResult {
  rank: number
  full_name: string
  surname: string
  first_name: string
  strokes: string
  char1: CharInfo
  char2: CharInfo
  wu_ge: WuGeResult
  san_cai: string
  san_cai_luck: string
  san_cai_detail: string
  ji_chu_yun: string
  cheng_gong_yun: string
  ren_ji_guan_xi: string
  zhou_yi: ZhouYiBasic
  bazi?: BaziBasic
  wu_xing?: WuXingSection
  score: number
  grade: string
  score_detail: ScoreDetail
  interpret: string
  has_poetry: boolean
}

export interface ExcellentEntry {
  char1: string
  char2: string
  score: number
  grade: string
  wu_xing1: string
  wu_xing2: string
  has_poetry: boolean
}