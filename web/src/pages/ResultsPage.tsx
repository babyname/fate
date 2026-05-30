import { useState, useEffect, useCallback } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { NameCard } from '@/components/naming/NameCard';
import { ExploreSection } from '@/components/naming/ExploreSection';
import { CandidateTable } from '@/components/naming/CandidateTable';
import { useAppStore } from '@/store/app';
import { api } from '@/lib/api';
import type { NameResult, ExcellentEntry } from '@/types/api';
import { ArrowLeft, Loader2, Grid3X3, List, RefreshCw } from 'lucide-react';

export function ResultsPage() {
  const { taskId } = useParams<{ taskId: string }>();
  const navigate = useNavigate();
  const [topNames, setTopNames] = useState<NameResult[]>([]);
  const [top10, setTop10] = useState<ExcellentEntry[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [progress, setProgress] = useState(0);
  const [viewMode, setViewMode] = useState<'card' | 'table'>('card');

  const storeTopNames = useAppStore((s) => s.topNames);
  const storeTop10 = useAppStore((s) => s.top10);
  const setStoreTopNames = useAppStore((s) => s.setTopNames);
  const setStoreTop10 = useAppStore((s) => s.setTop10);
  const setStoreTotal = useAppStore((s) => s.setTotal);

  const fetchResult = useCallback(async () => {
    if (!taskId) return;
    try {
      const res = await api.getTaskResult(taskId);
      setTopNames(res.top_names);
      setTop10(res.top10);
      setTotal(res.total);
      setStoreTopNames(res.top_names);
      setStoreTop10(res.top10);
      setStoreTotal(res.total);
      setLoading(false);
      setProgress(100);
    } catch {
      setError('获取结果失败');
      setLoading(false);
    }
  }, [taskId, setStoreTopNames, setStoreTop10, setStoreTotal]);

  const pollStatus = useCallback(async () => {
    if (!taskId) return;
    try {
      const res = await api.getTaskStatus(taskId);
      if (res.state === 'done') {
        fetchResult();
      } else if (res.state === 'failed') {
        setError(res.error || '生成失败');
        setLoading(false);
      } else {
        setProgress(res.total ? Math.min((res.total / 100) * 100, 99) : 30);
      }
    } catch {
      setError('获取状态失败');
      setLoading(false);
    }
  }, [taskId, fetchResult]);

  useEffect(() => {
    if (storeTopNames.length > 0) {
      setTopNames(storeTopNames);
      setTop10(storeTop10);
      setLoading(false);
      setProgress(100);
      return;
    }
    pollStatus();
  }, [storeTopNames, storeTop10, pollStatus]);

  useEffect(() => {
    if (!loading || !taskId) return;
    const interval = setInterval(pollStatus, 2000);
    return () => clearInterval(interval);
  }, [loading, taskId, pollStatus]);

  if (loading) {
    return (
      <div className="flex flex-col items-center justify-center py-24 space-y-4">
        <Loader2 className="h-12 w-12 animate-spin text-blue-400" />
        <div className="text-center space-y-2">
          <p className="text-lg font-medium text-foreground">正在生成名字...</p>
          <p className="text-sm text-muted-foreground font-mono">
            任务ID: {taskId}
          </p>
          {progress > 0 && (
            <div className="w-64 h-2 rounded-full bg-white/5 overflow-hidden mt-3">
              <div
                className="h-full rounded-full bg-gradient-to-r from-blue-500 to-violet-500 transition-all duration-500"
                style={{ width: `${progress}%` }}
              />
            </div>
          )}
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex flex-col items-center justify-center py-24 space-y-4">
        <p className="text-lg text-destructive">{error}</p>
        <Button variant="outline" onClick={() => navigate('/')}>
          <ArrowLeft className="h-4 w-4" />
          返回首页
        </Button>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Button variant="ghost" size="icon" onClick={() => navigate('/')}>
            <ArrowLeft className="h-5 w-5" />
          </Button>
          <div>
            <h1 className="text-xl font-semibold text-foreground">命名结果</h1>
            <p className="text-xs text-muted-foreground font-mono">Task: {taskId}</p>
          </div>
        </div>
        <div className="flex items-center gap-2">
          <Badge variant="celestial">{total} 个名字</Badge>
          <Button
            variant="ghost"
            size="icon"
            onClick={() => {
              setLoading(true);
              setTopNames([]);
              setTop10([]);
              setStoreTopNames([]);
              setStoreTop10([]);
              pollStatus();
            }}
          >
            <RefreshCw className="h-4 w-4" />
          </Button>
        </div>
      </div>

      <Tabs defaultValue="card" onValueChange={(v) => setViewMode(v as 'card' | 'table')}>
        <div className="flex items-center justify-between">
          <TabsList>
            <TabsTrigger value="card" className="gap-1.5">
              <Grid3X3 className="h-4 w-4" />
              卡片
            </TabsTrigger>
            <TabsTrigger value="table" className="gap-1.5">
              <List className="h-4 w-4" />
              表格
            </TabsTrigger>
          </TabsList>
        </div>

        <TabsContent value="card" className="mt-4 space-y-6">
          <ExploreSection taskId={taskId!} />

          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
            {topNames.map((result) => (
              <NameCard key={result.full_name} name={result} />
            ))}
          </div>
        </TabsContent>

        <TabsContent value="table" className="mt-4 space-y-4">
          <CandidateTable results={topNames} />
        </TabsContent>
      </Tabs>
    </div>
  );
}