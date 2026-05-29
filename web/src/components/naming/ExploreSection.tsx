import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { ChevronDown, ChevronUp } from 'lucide-react';
import type { NameResult } from '@/types/api';
import { useAppStore } from '@/store/app';

interface ExploreSectionProps {
  results: NameResult[];
}

export function ExploreSection({ results }: ExploreSectionProps) {
  const [expanded, setExpanded] = useState(false);
  const setSelectedName = useAppStore((s) => s.setSelectedName);
  const setDetailModalOpen = useAppStore((s) => s.setDetailModalOpen);

  const displayResults = expanded ? results : results.slice(0, 6);

  if (results.length === 0) return null;

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h2 className="text-lg font-semibold text-foreground">
          探索名字
        </h2>
        <Badge variant="secondary" className="font-mono">
          {results.length} 个结果
        </Badge>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
        {displayResults.map((result) => (
          <button
            key={result.name}
            onClick={() => {
              setSelectedName(result);
              setDetailModalOpen(true);
            }}
            className="glass-card-hover p-3 text-left transition-all duration-300 group"
          >
            <div className="flex items-center justify-between mb-1">
              <span className="text-lg font-serif font-bold text-foreground group-hover:text-blue-300 transition-colors">
                {result.name}
              </span>
              <span className="text-sm font-mono text-muted-foreground">
                {result.score}分
              </span>
            </div>
            <p className="text-xs text-muted-foreground line-clamp-1 font-serif">
              {result.meaning}
            </p>
          </button>
        ))}
      </div>

      {results.length > 6 && (
        <div className="flex justify-center">
          <Button
            variant="ghost"
            size="sm"
            onClick={() => setExpanded(!expanded)}
            className="gap-1"
          >
            {expanded ? (
              <>
                收起 <ChevronUp className="h-4 w-4" />
              </>
            ) : (
              <>
                查看更多 ({results.length - 6}) <ChevronDown className="h-4 w-4" />
              </>
            )}
          </Button>
        </div>
      )}
    </div>
  );
}
