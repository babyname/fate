import { useParams, useNavigate } from 'react-router-dom';
import { useState, useEffect, useCallback, useRef } from 'react';
import type { NameResult } from '@/types/api';
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { NameCard } from '@/components/naming/NameCard';
import { ExploreSection } from '@/components/naming/ExploreSection';
import { CandidateTable } from '@/components/naming/CandidateTable';
import { SameNameLookupButton } from '@/components/naming/SameNameLookup';
import { useAppStore } from '@/store/app';
import { api } from '@/lib/api';
import { ArrowLeft, Grid3X3, List, Search } from 'lucide-react';

const PAGE_SIZE = 20;

export function ResultsPage() {
  const { taskId } = useParams<{ taskId: string }>();
  const navigate = useNavigate();
  const [viewMode, setViewMode] = useState<'table' | 'card'>('table');
  const [exploreVisited, setExploreVisited] = useState(false);
  const [currentPage, setCurrentPage] = useState(1);
  const [pagedNames, setPagedNames] = useState<NameResult[]>([]);
  const [totalCount, setTotalCount] = useState(0);
  const [loadingPagination, setLoadingPagination] = useState(false);
  const [polling, setPolling] = useState(true);
  const [pollError, setPollError] = useState('');
  const pollRef = useRef<ReturnType<typeof setInterval> | null>(null);

  const setStoreTopNames = useAppStore((s) => s.setTopNames);
  const setStoreTop10 = useAppStore((s) => s.setTop10);
  const setStoreTotal = useAppStore((s) => s.setTotal);
  const setIsGenerating = useAppStore((s) => s.setIsGenerating);

  const loadPage = useCallback(async (page: number) => {
    if (!taskId) return;
    setLoadingPagination(true);
    try {
      const res = await api.getNames(taskId, page, PAGE_SIZE);
      const names = res.names ?? [];
      setPagedNames(names);
      setTotalCount(res.total);
      setCurrentPage(res.page);
      setStoreTopNames(names);
      setStoreTotal(res.total);
    } catch (error) {
      setPollError(error instanceof Error ? error.message : '读取候选名失败。');
    } finally {
      setLoadingPagination(false);
    }
  }, [taskId, setStoreTopNames, setStoreTotal]);

  const onGenerationDone = useCallback(async () => {
    await loadPage(1);
    setPolling(false);
    setIsGenerating(false);
  }, [loadPage, setIsGenerating]);

  useEffect(() => {
    if (!taskId) return;
    setPolling(true);
    setIsGenerating(true);

    pollRef.current = setInterval(async () => {
      try {
        const status = await api.getTaskStatus(taskId);
        if (status.state === 'done') {
          if (pollRef.current) clearInterval(pollRef.current);
          await onGenerationDone();
        } else if (status.state === 'failed') {
          if (pollRef.current) clearInterval(pollRef.current);
          setPolling(false);
          setIsGenerating(false);
          setPollError(status.error || '生成失败');
        }
      } catch (error) {
        if (pollRef.current) clearInterval(pollRef.current);
        setPolling(false);
        setIsGenerating(false);
        setPollError(error instanceof Error ? error.message : '读取命名任务失败。');
      }
    }, 500);

    return () => {
      if (pollRef.current) clearInterval(pollRef.current);
    };
  }, [taskId, onGenerationDone, setIsGenerating]);

  if (polling) {
    return (
      <div className="flex flex-col items-center justify-center py-24 space-y-4">
        <div className="h-10 w-10 animate-spin rounded-full border-4 border-outline-variant border-t-primary" />
        <p className="text-lg font-medium text-foreground">正在生成名字...</p>
        <p className="text-sm text-muted-foreground font-mono">Task: {taskId}</p>
      </div>
    );
  }

  if (pollError) {
    return (
      <div className="flex flex-col items-center justify-center py-24 space-y-4">
        <p className="text-lg text-destructive">{pollError}</p>
        <Button variant="outline" onClick={() => navigate('/')}>
          <ArrowLeft className="h-4 w-4" />
          返回首页
        </Button>
      </div>
    );
  }

  if (pagedNames.length === 0 && !loadingPagination) {
    return (
      <div className="flex flex-col items-center justify-center py-24 space-y-4">
        <p className="text-lg text-muted-foreground">暂无结果</p>
        <Button variant="outline" onClick={() => navigate('/')}>
          <ArrowLeft className="h-4 w-4" />
          返回首页
        </Button>
      </div>
    );
  }

  const top3 = pagedNames.slice(0, 3);
  const moreNames = pagedNames.slice(3, 7);
  const totalPages = Math.ceil(totalCount / PAGE_SIZE);

  return (
    <div className="space-y-8">
      <div className="flex items-center gap-3">
        <Button variant="ghost" size="icon" onClick={() => navigate('/')}>
          <ArrowLeft className="h-5 w-5" />
        </Button>
        <div>
          <h1 className="font-serif text-xl font-semibold text-foreground">命名结果</h1>
          <p className="text-xs text-muted-foreground font-mono">Task: {taskId}</p>
        </div>
        <Badge variant="celestial">{totalCount} 个候选</Badge>
        <div className="h-px flex-1 bg-gradient-to-r from-outline-variant to-transparent" />
      </div>

      <Tabs defaultValue="main" onValueChange={(value) => {
        if (value === 'explore') setExploreVisited(true);
      }}>
        <TabsList className="w-full sm:w-auto">
          <TabsTrigger value="main">命名结果</TabsTrigger>
          <TabsTrigger value="explore" className="gap-2">
            <Search className="h-4 w-4" />
            探索更多
          </TabsTrigger>
        </TabsList>

        <TabsContent value="main" className="mt-6 space-y-6">
          <div>
            <h2 className="font-serif text-lg font-bold text-foreground mb-4">精选好名</h2>
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">
              {top3.map((item, index) => (
                <NameCard key={`${item.full_name}-${index}`} name={item} />
              ))}
            </div>
          </div>

          {moreNames.length > 0 && (
            <div>
              <h3 className="mb-3 font-serif text-base font-semibold text-foreground">更多推荐</h3>
              <div className="flex gap-3 overflow-x-auto pb-2">
                {moreNames.map((item, index) => (
                  <div key={`${item.full_name}-more-${index}`} className="min-w-[200px] max-w-[240px] flex-shrink-0">
                    <NameCard name={item} />
                  </div>
                ))}
              </div>
            </div>
          )}

          <div className="h-px bg-gradient-to-r from-outline-variant to-transparent" />

          <div>
            <div className="flex items-center justify-between mb-3">
              <h3 className="font-serif text-base font-semibold text-foreground">候选名册</h3>
              <div className="flex items-center gap-2 border rounded-md p-1 bg-card">
                <Button
                  variant={viewMode === 'table' ? 'default' : 'ghost'}
                  size="sm"
                  className="gap-1.5"
                  onClick={() => setViewMode('table')}
                >
                  <List className="h-4 w-4" />
                  表格
                </Button>
                <Button
                  variant={viewMode === 'card' ? 'default' : 'ghost'}
                  size="sm"
                  className="gap-1.5"
                  onClick={() => setViewMode('card')}
                >
                  <Grid3X3 className="h-4 w-4" />
                  卡片
                </Button>
              </div>
            </div>
            {viewMode === 'table' ? (
              <CandidateTable results={pagedNames} />
            ) : (
              <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
                {pagedNames.map((result) => (
                  <div key={result.full_name}>
                    <NameCard name={result} showWuGeLuck />
                    <div className="px-2 pt-2"><SameNameLookupButton fullName={result.full_name} compact /></div>
                  </div>
                ))}
              </div>
            )}

            {totalPages > 1 && (
              <div className="flex items-center justify-center gap-2 mt-6">
                <Button
                  variant="outline"
                  size="sm"
                  disabled={currentPage <= 1 || loadingPagination}
                  onClick={() => loadPage(currentPage - 1)}
                >
                  上一页
                </Button>
                <span className="text-sm text-muted-foreground">
                  第 {currentPage} 页 / 共 {totalPages} 页
                </span>
                <Button
                  variant="outline"
                  size="sm"
                  disabled={currentPage >= totalPages || loadingPagination}
                  onClick={() => loadPage(currentPage + 1)}
                >
                  下一页
                </Button>
              </div>
            )}
          </div>
        </TabsContent>

        <TabsContent value="explore" forceMount={exploreVisited || undefined} className="mt-6 data-[state=inactive]:hidden">
          {taskId && exploreVisited && <ExploreSection key={taskId} taskId={taskId} surname={pagedNames[0]?.surname || ''} />}
        </TabsContent>
      </Tabs>
    </div>
  );
}
