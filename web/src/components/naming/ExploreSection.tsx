import { useState, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { ChevronDown, ChevronUp, Loader2, BookOpen, Filter } from 'lucide-react';
import type { ExcellentEntry } from '@/types/api';
import { useAppStore } from '@/store/app';
import { api } from '@/lib/api';

interface ExploreSectionProps {
  taskId: string;
}

export function ExploreSection({ taskId }: ExploreSectionProps) {
  const [entries, setEntries] = useState<ExcellentEntry[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [expanded, setExpanded] = useState(false);
  const [poetryOnly, setPoetryOnly] = useState(false);
  const setSelectedName = useAppStore((s) => s.setSelectedName);
  const setDetailModalOpen = useAppStore((s) => s.setDetailModalOpen);

  const fetchExplore = (hasPoetry?: boolean) => {
    setLoading(true);
    api.explore(taskId, 10, hasPoetry).then((res) => {
      setEntries(res.names);
      setTotal(res.total);
      setLoading(false);
    }).catch(() => {
      setLoading(false);
    });
  };

  useEffect(() => {
    if (!taskId) return;
    fetchExplore();
  }, [taskId]);

  const handlePoetryFilter = () => {
    const next = !poetryOnly;
    setPoetryOnly(next);
    fetchExplore(next);
  };

  const handleLoadMore = async () => {
    setExpanded(true);
    setLoading(true);
    try {
      const res = await api.explore(taskId, 20, poetryOnly);
      setEntries(res.names);
    } catch {
    } finally {
      setLoading(false);
    }
  };

  const handleNameClick = async (entry: ExcellentEntry) => {
    try {
      const res = await api.getNameDetail(taskId, entry.char1, entry.char2);
      setSelectedName(res.name_result);
      setDetailModalOpen(true);
    } catch {
    }
  };

  const displayEntries = expanded ? entries : entries.slice(0, 6);

  if (loading && entries.length === 0) {
    return (
      <div className="flex items-center justify-center py-6">
        <Loader2 className="h-5 w-5 animate-spin text-blue-400" />
      </div>
    );
  }

  if (entries.length === 0) return null;

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h2 className="text-lg font-semibold text-foreground">
          探索名字
        </h2>
        <div className="flex items-center gap-2">
          <Button
            variant={poetryOnly ? 'celestial' : 'outline'}
            size="sm"
            onClick={handlePoetryFilter}
            className="gap-1.5"
          >
            <BookOpen className="h-3.5 w-3.5" />
            {poetryOnly ? '诗词优先' : '全部'}
          </Button>
          <Badge variant="secondary" className="font-mono">
            {total} 个结果
          </Badge>
        </div>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
        {displayEntries.map((entry) => (
          <button
            key={`${entry.char1}${entry.char2}`}
            onClick={() => handleNameClick(entry)}
            className="glass-card-hover p-3 text-left transition-all duration-300 group"
          >
            <div className="flex items-center justify-between mb-1">
              <span className="text-lg font-serif font-bold text-foreground group-hover:text-blue-300 transition-colors">
                {entry.char1}{entry.char2}
              </span>
              <span className="text-sm font-mono text-muted-foreground">
                {entry.score.toFixed(0)}分
              </span>
            </div>
            <div className="flex items-center gap-1.5">
              <Badge variant="celestial" className="text-xs py-0">{entry.wu_xing1}</Badge>
              <Badge variant="celestial" className="text-xs py-0">{entry.wu_xing2}</Badge>
              <span className="text-xs text-muted-foreground font-mono">{entry.grade}</span>
              {entry.has_poetry && (
                <Badge variant="stardust" className="text-xs py-0 gap-0.5">
                  <BookOpen className="h-2.5 w-2.5" />
                  诗
                </Badge>
              )}
            </div>
          </button>
        ))}
      </div>

      <div className="flex justify-center gap-2">
        {entries.length > 6 && !expanded && (
          <Button
            variant="ghost"
            size="sm"
            onClick={handleLoadMore}
            disabled={loading}
            className="gap-1"
          >
            {loading ? (
              <Loader2 className="h-4 w-4 animate-spin" />
            ) : (
              <>
                查看更多 ({entries.length - 6}) <ChevronDown className="h-4 w-4" />
              </>
            )}
          </Button>
        )}
        {expanded && entries.length > 6 && (
          <Button
            variant="ghost"
            size="sm"
            onClick={() => setExpanded(false)}
            className="gap-1"
          >
            收起 <ChevronUp className="h-4 w-4" />
          </Button>
        )}
      </div>
    </div>
  );
}