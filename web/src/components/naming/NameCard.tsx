import { Badge } from '@/components/ui/badge';
import { WuGeLuck } from '@/components/naming/WuGeLuck';
import { recommendationScore } from '@/lib/recommendation-scores';
import { useAppStore } from '@/store/app';
import { useScoringPreferences } from '@/store/scoring-preferences';
import type { NameResult } from '@/types/api';
import { Star } from 'lucide-react';

const wuxingColors: Record<string, string> = {
  金: 'text-amber-300',
  木: 'text-green-400',
  水: 'text-blue-400',
  火: 'text-red-400',
  土: 'text-yellow-600',
};

interface NameCardProps {
  name: NameResult;
  showWuGeLuck?: boolean;
}

export function NameCard({ name, showWuGeLuck }: NameCardProps) {
  const weight = useScoringPreferences((s) => s.traditionalWeight);
  const score = name.ranking ? recommendationScore(name.ranking, 'total', weight) : name.score;
  const setSelectedName = useAppStore((s) => s.setSelectedName);
  const setDetailModalOpen = useAppStore((s) => s.setDetailModalOpen);

  const handleClick = () => {
    setSelectedName(name);
    setDetailModalOpen(true);
  };

  const scoreColor =
    score >= 90
      ? 'text-gradient-stardust'
      : score >= 75
        ? 'text-gradient-celestial'
        : 'text-foreground';

  return (
    <button
      onClick={handleClick}
      className="glass-card-hover w-full text-left p-4 transition-all duration-300 group"
    >
      <div className="flex items-start justify-between mb-3">
        <div>
          <h3 className="text-2xl font-serif font-bold text-foreground group-hover:text-blue-300 transition-colors">
            {name.full_name}
          </h3>
          <p className="text-sm text-muted-foreground font-mono mt-0.5">
            {name.full_pinyin || `${name.char1.pinyin} ${name.char2.pinyin}`}
          </p>
        </div>
        <div className="flex items-center gap-1">
          <Star className={`h-4 w-4 ${score >= 90 ? 'text-amber-400' : 'text-muted-foreground'}`} />
          <span className={`text-lg font-bold font-mono ${scoreColor}`}>
            {score.toFixed(name.ranking ? 1 : 0)}
          </span>
        </div>
      </div>

      <p className="text-sm text-foreground/70 line-clamp-2 mb-3 leading-relaxed">
        {name.recommendation || [name.char1.meaning, name.char2.meaning].filter(Boolean).join('，')}
      </p>

      <div className="flex flex-wrap gap-1.5 mb-2">
        <Badge variant="celestial" className={`text-xs ${wuxingColors[name.char1.wu_xing] || ''}`}>
          {name.char1.wu_xing}
        </Badge>
        <Badge variant="celestial" className={`text-xs ${wuxingColors[name.char2.wu_xing] || ''}`}>
          {name.char2.wu_xing}
        </Badge>
        {name.san_cai_luck && (
          <Badge variant="outline" className="text-xs">
            {name.san_cai_luck}
          </Badge>
        )}
      </div>

      <div className="flex items-center gap-2 text-xs text-muted-foreground">
        {name.ranking ? (
          <span className="tabular-nums">传统数理 {name.ranking.traditional.toFixed(1)} · 音义 {name.ranking.sound_and_meaning.toFixed(1)}</span>
        ) : <span className="font-mono">{name.grade}</span>}
        <span>三才: {name.san_cai}</span>
        {name.has_poetry && (
          <Badge variant="stardust" className="text-xs py-0">诗</Badge>
        )}
      </div>
      {showWuGeLuck && (
        <div className="mt-3 border-t border-white/10 pt-3">
          <WuGeLuck wuGe={name.wu_ge} />
        </div>
      )}
    </button>
  );
}
