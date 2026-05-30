import { Badge } from '@/components/ui/badge';
import { useAppStore } from '@/store/app';
import type { NameResult } from '@/types/api';
import { ArrowUpDown, Star } from 'lucide-react';
import { useState, useMemo } from 'react';

interface CandidateTableProps {
  results: NameResult[];
}

type SortKey = 'score' | 'name' | 'san_cai';

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
          cmp = a.full_name.localeCompare(b.full_name, 'zh');
          break;
        case 'san_cai':
          cmp = a.san_cai.localeCompare(b.san_cai, 'zh');
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
      <div className="overflow-x-auto">
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
              <th className="px-4 py-3 text-left text-muted-foreground font-medium">笔画</th>
              <th className="px-4 py-3 text-left text-muted-foreground font-medium">五行</th>
              <th className="px-4 py-3 text-left text-muted-foreground font-medium">
                <button
                  onClick={() => handleSort('score')}
                  className="flex items-center gap-1 hover:text-foreground transition-colors"
                >
                  评分
                  <ArrowUpDown className="h-3 w-3" />
                </button>
              </th>
              <th className="px-4 py-3 text-left text-muted-foreground font-medium">等级</th>
              <th className="px-4 py-3 text-left text-muted-foreground font-medium">三才</th>
            </tr>
          </thead>
          <tbody>
            {sorted.map((result) => (
              <tr
                key={result.full_name}
                onClick={() => {
                  setSelectedName(result);
                  setDetailModalOpen(true);
                }}
                className="border-b border-white/5 hover:bg-white/5 cursor-pointer transition-colors"
              >
                <td className="px-4 py-3">
                  <span className="font-serif font-bold text-foreground">{result.full_name}</span>
                </td>
                <td className="px-4 py-3 font-mono text-xs text-muted-foreground">
                  {result.char1.pinyin} {result.char2.pinyin}
                </td>
                <td className="px-4 py-3 font-mono text-xs text-muted-foreground">
                  {result.strokes}
                </td>
                <td className="px-4 py-3">
                  <div className="flex gap-1">
                    <Badge variant="celestial" className="text-xs py-0">
                      {result.char1.wu_xing}
                    </Badge>
                    <Badge variant="celestial" className="text-xs py-0">
                      {result.char2.wu_xing}
                    </Badge>
                  </div>
                </td>
                <td className="px-4 py-3">
                  <div className="flex items-center gap-1">
                    <Star
                      className={`h-3 w-3 ${result.score >= 90 ? 'text-amber-400' : result.score >= 75 ? 'text-blue-400' : 'text-muted-foreground'}`}
                    />
                    <span className="font-mono font-medium">{result.score.toFixed(0)}</span>
                  </div>
                </td>
                <td className="px-4 py-3">
                  <span className="font-mono text-xs">{result.grade}</span>
                </td>
                <td className="px-4 py-3">
                  <Badge
                    variant={result.san_cai_luck === '吉' || result.san_cai_luck === '大吉' ? 'celestial' : 'outline'}
                    className="text-xs"
                  >
                    {result.san_cai}
                  </Badge>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}