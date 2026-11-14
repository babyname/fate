import { useState, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { ChevronDown, ChevronUp, Loader2, BookOpen } from 'lucide-react';
import type { ExcellentEntry } from '@/types/api';
import { useAppStore } from '@/store/app';
import { api } from '@/lib/api';

interface ExploreSectionProps {
  taskId: string;
  surname: string;
}

export function ExploreSection({ taskId, surname }: ExploreSectionProps) {
  const [entries, setEntries] = useState<ExcellentEntry[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [expanded, setExpanded] = useState(false);
  const [poetryOnly, setPoetryOnly] = useState(false);
  const [exhausted, setExhausted] = useState(false);
  const [error, setError] = useState('');
  const [retryCount, setRetryCount] = useState(0);
  const setSelectedName = useAppStore((s) => s.setSelectedName);
  const setDetailModalOpen = useAppStore((s) => s.setDetailModalOpen);

  useEffect(() => {
    if (!taskId) return;
    let active = true;
    setLoading(true);
    setError('');
    setEntries([]);
    setExpanded(false);
    setExhausted(false);
    api.explore(taskId, 6, poetryOnly).then((res) => {
      if (!active) return;
      const names = res.names ?? [];
      setEntries(names);
      setTotal(res.total);
      setExhausted(names.length === 0);
    }).catch((err) => {
      if (active) setError(err instanceof Error ? err.message : '加载名字失败，请重试');
    }).finally(() => {
      if (active) setLoading(false);
    });
    return () => { active = false; };
  }, [taskId, poetryOnly, retryCount]);

  const handlePoetryFilter = () => {
    setPoetryOnly((value) => !value);
  };

  const handleLoadMore = async () => {
    if (!expanded && entries.length > 6) {
      setExpanded(true);
      return;
    }
    setLoading(true);
    setError('');
    try {
      const res = await api.explore(taskId, 20, poetryOnly);
      const names = res.names ?? [];
      setEntries((previous) => [...previous, ...names]);
      setTotal(res.total);
      setExhausted(names.length === 0);
      setExpanded(true);
    } catch (err) {
      setError(err instanceof Error ? err.message : '加载名字失败，请重试');
    } finally {
      setLoading(false);
    }
  };

  const handleNameClick = async (entry: ExcellentEntry) => {
    try {
      const res = await api.getNameDetail(taskId, entry.char1, entry.char2);
      setSelectedName(res.name_result);
      setDetailModalOpen(true);
    } catch (err) {
      setError(err instanceof Error ? err.message : '加载名字详情失败，请重试');
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
            disabled={loading}
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

      {error && (
        <div role="alert" className="space-y-2 text-sm text-destructive">
          <p>{error}</p>
          {entries.length === 0 && (
            <Button variant="outline" size="sm" onClick={() => setRetryCount((value) => value + 1)}>
              重试
            </Button>
          )}
        </div>
      )}

      {!error && entries.length === 0 && (
        <p role="status" className="glass-card p-6 text-center text-sm text-muted-foreground">
          暂无可探索的名字，请调整筛选或返回首页重新生成。
        </p>
      )}

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
        {displayEntries.map((entry) => (
          <button
            key={`${entry.char1}${entry.char2}`}
            onClick={() => handleNameClick(entry)}
            className="glass-card-hover p-3 text-left transition-all duration-300 group"
          >
            <div className="flex items-center justify-between mb-1">
              <span className="text-lg font-serif font-bold text-foreground group-hover:text-blue-300 transition-colors">
                {surname}{entry.char1}{entry.char2}
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

      {exhausted && entries.length > 0 && (
        <p role="status" className="text-center text-sm text-muted-foreground">
          本轮可探索的名字已全部展示，可返回命名结果或重新生成。
        </p>
      )}

      <div className="flex justify-center gap-2">
        {entries.length > 0 && (!exhausted || (!expanded && entries.length > 6)) && (
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
                {expanded ? '继续探索' : '查看更多'} <ChevronDown className="h-4 w-4" />
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
