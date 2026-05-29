import { Badge } from '@/components/ui/badge';
import { useAppStore } from '@/store/app';
import type { NameResult } from '@/types/api';
import { ArrowUpDown, Star } from 'lucide-react';
import { useState, useMemo } from 'react';

interface CandidateTableProps {
  results: NameResult[];
}

type SortKey = 'score' | 'name' | 'wuxing';

export function CandidateTable({ results }: CandidateTableProps) {
  const [sortKey, setSortKey] = useState<SortKey>('score');
  const [sortDesc, setSortDesc] = useState(true);
  const setSelectedName = useAppStore((s) => s.setSelectedName);
  const setDetailModalOpen = useAppStore((s) => s.setDetailModalOpen);

  const sorted = useMemo(() => {
    return [...results].sort((a, b) => {
      let cmp = 0;
      switch (sortKey) {
        case 'score':
          cmp = a.score - b.score;
          break;
        case 'name':
          cmp = a.name.localeCompare(b.name, 'zh');
          break;
        case 'wuxing':
          cmp = a.wuxing.balance - b.wuxing.balance;
          break;
      }
      return sortDesc ? -cmp : cmp;
    });
  }, [results, sortKey, sortDesc]);

  const handleSort = (key: SortKey) => {
    if (sortKey === key) {
      setSortDesc(!sortDesc);
    } else {
      setSortKey(key);
      setSortDesc(true);
    }
  };

  if (results.length === 0) return null;

  return (
    <div className="glass-card overflow-hidden">
      <div className="overflow-x-auto scrollbar-thin">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-white/10">
              <th className="px-4 py-3 text-left text-muted-foreground font-medium">
                <button
                  onClick={() => handleSort('name')}
                  className="flex items-center gap-1 hover:text-foreground transition-colors"
                >
                  名字
                  <ArrowUpDown className="h-3 w-3" />
                </button>
              </th>
              <th className="px-4 py-3 text-left text-muted-foreground font-medium">拼音</th>
              <th className="px-4 py-3 text-left text-muted-foreground font-medium">
                <button
                  onClick={() => handleSort('score')}
                  className="flex items-center gap-1 hover:text-foreground transition-colors"
                >
                  评分
                  <ArrowUpDown className="h-3 w-3" />
                </button>
              </th>
              <th className="px-4 py-3 text-left text-muted-foreground font-medium">五行</th>
              <th className="px-4 py-3 text-left text-muted-foreground font-medium">出处</th>
            </tr>
          </thead>
          <tbody>
            {sorted.map((result) => (
              <tr
                key={result.name}
                onClick={() => {
                  setSelectedName(result);
                  setDetailModalOpen(true);
                }}
                className="border-b border-white/5 hover:bg-white/5 cursor-pointer transition-colors"
              >
                <td className="px-4 py-3">
                  <span className="font-serif font-bold text-foreground">{result.name}</span>
                </td>
                <td className="px-4 py-3 font-mono text-xs text-muted-foreground">
                  {result.pronunciation.pinyin}
                </td>
                <td className="px-4 py-3">
                  <div className="flex items-center gap-1">
                    <Star
                      className={`h-3 w-3 ${result.score >= 90 ? 'text-amber-400' : result.score >= 75 ? 'text-blue-400' : 'text-muted-foreground'}`}
                    />
                    <span className="font-mono font-medium">{result.score}</span>
                  </div>
                </td>
                <td className="px-4 py-3">
                  <div className="flex gap-1">
                    {result.wuxing.elements.map((el) => (
                      <Badge key={el} variant="celestial" className="text-xs py-0">
                        {el}
                      </Badge>
                    ))}
                  </div>
                </td>
                <td className="px-4 py-3">
                  {result.poetry.length > 0 ? (
                    <span className="text-xs text-muted-foreground font-serif truncate block max-w-[200px]">
                      《{result.poetry[0].title}》
                    </span>
                  ) : (
                    <span className="text-xs text-muted-foreground">—</span>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
