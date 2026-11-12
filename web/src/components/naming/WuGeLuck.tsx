import { wuGeFields, wuGeLuckColor } from '@/lib/wu-ge';
import type { WuGeResult } from '@/types/api';

export function WuGeLuck({ wuGe }: { wuGe: WuGeResult }) {
  return (
    <dl aria-label="五格吉凶" className="grid min-w-[10rem] grid-cols-5 gap-x-1 text-center text-xs whitespace-nowrap">
      {wuGeFields.map(({ key, label }) => {
        const luck = wuGe?.[key]?.lucky || '—';
        return (
          <div key={key} className="space-y-1">
            <dt className="text-muted-foreground">{label}</dt>
            <dd className={`font-medium ${wuGeLuckColor(luck)}`}>{luck}</dd>
          </div>
        );
      })}
    </dl>
  );
}
