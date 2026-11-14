import type { WuGeResult } from '@/types/api';

export const wuGeFields: { key: keyof WuGeResult; label: string }[] = [
  { key: 'tian_ge', label: '天格' },
  { key: 'ren_ge', label: '人格' },
  { key: 'di_ge', label: '地格' },
  { key: 'wai_ge', label: '外格' },
  { key: 'zong_ge', label: '总格' },
];

export function wuGeLuckColor(luck: string) {
  switch (luck) {
    case '吉':
    case '大吉':
      return 'text-green-400';
    case '半吉':
      return 'text-amber-300';
    case '凶':
    case '大凶':
      return 'text-red-400';
    default:
      return 'text-muted-foreground';
  }
}
