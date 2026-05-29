import { useState, useEffect, useCallback } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { NameCard } from '@/components/naming/NameCard';
import { ExploreSection } from '@/components/naming/ExploreSection';
import { CandidateTable } from '@/components/naming/CandidateTable';
import { ReportDownload } from '@/components/naming/ReportDownload';
import { useAppStore } from '@/store/app';
import { api } from '@/lib/api';
import type { NameResult } from '@/types/api';
import { ArrowLeft, Loader2, Grid3X3, List, RefreshCw } from 'lucide-react';

export function ResultsPage() {
  const { taskId } = useParams<{ taskId: string }>();
  const navigate = useNavigate();
  const [results, setResults] = useState<NameResult[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [progress, setProgress] = useState(0);
  const [viewMode, setViewMode] = useState<'card' | 'table'>('card');

  const storeResults = useAppStore((s) => s.results);
  const setStoreResults = useAppStore((s) => s.setResults);

  const pollStatus = useCallback(async () => {
    if (!taskId) return;
    try {
      const res = await api.getTaskStatus(taskId);
      if (res.status === 'completed' && res.result) {
        setResults(res.result);
        setStoreResults(res.result);
        setLoading(false);
      } else if (res.status === 'failed') {
        setError(res.error || '生成失败');
        setLoading(false);
      } else {
        setProgress(res.progress || 0);
      }
    } catch {
      setError('获取结果失败');
      setLoading(false);
    }
  }, [taskId, setStoreResults]);

  useEffect(() => {
    if (storeResults.length > 0) {
      setResults(storeResults);
      setLoading(false);
      return;
    }
    pollStatus();
  }, [storeResults, pollStatus]);

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
          <Badge variant="celestial">{results.length} 个名字</Badge>
          <Button
            variant="ghost"
            size="icon"
            onClick={() => {
              setLoading(true);
              setResults([]);
              setStoreResults([]);
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
          <ExploreSection results={results} />

          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
            {results.map((result) => (
              <NameCard key={result.name} name={result} />
            ))}
          </div>
        </TabsContent>

        <TabsContent value="table" className="mt-4 space-y-4">
          <CandidateTable results={results} />
        </TabsContent>
      </Tabs>

      <ReportDownload taskId={taskId!} />
    </div>
  );
}
