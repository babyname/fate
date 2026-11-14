export const sameNameSources = {
  nationwide: {
    scope: '中国大陆全国',
    shortScope: '全国',
    agency: '公安部政务服务平台',
    url: 'https://ywtb.mps.gov.cn/',
    instructions: '在「便民应用」点「查询同名人数」，选择「全国」，按提示登录和查询。如可选择性别，请选「全部」。',
  },
  jiangsu: {
    scope: '江苏省',
    shortScope: '江苏',
    agency: '江苏公安',
    url: 'https://ythpt.jsga.gov.cn/jsgawx/tongming/tmcx.html',
    instructions: '粘贴全名查询江苏省内人数。也可在「苏服办」搜索「重名查询」。',
  },
  shandong: {
    scope: '山东省',
    shortScope: '山东',
    agency: '山东公安',
    url: 'https://msjw.gat.shandong.gov.cn/populationclient/search.shtml',
    instructions: '粘贴全名并填写官方页面的验证码，查询山东省内人数。',
  },
};

export type SameNameSource = keyof typeof sameNameSources;

export interface SameNameRecord {
  name: string;
  source: SameNameSource;
  count: number;
  under: boolean;
  recordedAt: string;
}

export function parseSameNameCount(input: string, under: boolean): number | null {
  const value = input.trim().normalize('NFKC');
  if (!/^(?:[0-9]+|[0-9]{1,3}(?:,[0-9]{3})+)$/.test(value)) return null;
  const count = Number(value.replace(/,/g, ''));
  return Number.isSafeInteger(count) && (!under || count > 0) ? count : null;
}

export function formatSameNameCount({ count, under }: SameNameRecord): string {
  return `${under ? '少于 ' : ''}${count.toLocaleString('zh-CN')} 人`;
}
