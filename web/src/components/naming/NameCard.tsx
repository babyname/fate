import { Badge } from '@/components/ui/badge';
import { useAppStore } from '@/store/app';
import type { NameResult } from '@/types/api';
import { Star, BookOpen } from 'lucide-react';

interface NameCardProps {
  name: NameResult;
}

export function NameCard({ name }: NameCardProps) {
  const setSelectedName = useAppStore((s) => s.setSelectedName);
  const setDetailModalOpen = useAppStore((s) => s.setDetailModalOpen);

  const handleClick = () => {
    setSelectedName(name);
    setDetailModalOpen(true);
  };

  const scoreColor =
    name.score >= 90
      ? 'text-gradient-stardust'
      : name.score >= 75
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
            {name.name}
          </h3>
          <p className="text-sm text-muted-foreground font-mono mt-0.5">
            {name.pronunciation.pinyin}
          </p>
        </div>
        <div className="flex items-center gap-1">
          <Star className={`h-4 w-4 ${name.score >= 90 ? 'text-amber-400' : 'text-muted-foreground'}`} />
          <span className={`text-lg font-bold font-mono ${scoreColor}`}>
            {name.score}
          </span>
        </div>
      </div>

      <p className="text-sm text-foreground/70 line-clamp-2 mb-3 leading-relaxed">
        {name.meaning}
      </p>

      <div className="flex flex-wrap gap-1.5 mb-2">
        {name.wuxing.elements.map((el) => (
          <Badge key={el} variant="celestial" className="text-xs">
            {el}
          </Badge>
        ))}
        {name.wuxing.missing.length > 0 && (
          <Badge variant="outline" className="text-xs text-muted-foreground">
            缺{name.wuxing.missing.join('')}
          </Badge>
        )}
      </div>

      {name.poetry.length > 0 && (
        <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
          <BookOpen className="h-3 w-3" />
          <span className="truncate font-serif">
            「{name.poetry[0].content.slice(0, 20)}...」— {name.poetry[0].author}
          </span>
        </div>
      )}
    </button>
  );
}
