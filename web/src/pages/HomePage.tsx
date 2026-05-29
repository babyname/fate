import { useState, useCallback, useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { FixedCharInput } from '@/components/naming/FixedCharInput';
import { BlacklistPanel } from '@/components/naming/BlacklistPanel';
import { PoetrySearchTab } from '@/components/naming/PoetrySearchTab';
import { useAppStore } from '@/store/app';
import { api } from '@/lib/api';
import { Sparkles, Settings, BookOpen, Loader2 } from 'lucide-react';

export function HomePage() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const defaultTab = searchParams.get('tab') === 'poetry' ? 'poetry' : 'generator';

  const surname = useAppStore((s) => s.surname);
  const gender = useAppStore((s) => s.gender);
  const fixedChars = useAppStore((s) => s.fixedChars);
  const blacklist = useAppStore((s) => s.blacklist);
  const count = useAppStore((s) => s.count);
  const isGenerating = useAppStore((s) => s.isGenerating);
  const setSurname = useAppStore((s) => s.setSurname);
  const setGender = useAppStore((s) => s.setGender);
  const setFixedChars = useAppStore((s) => s.setFixedChars);
  const setBlacklist = useAppStore((s) => s.setBlacklist);
  const setCount = useAppStore((s) => s.setCount);
  const setTaskId = useAppStore((s) => s.setTaskId);
  const setIsGenerating = useAppStore((s) => s.setIsGenerating);

  const [showAdvanced, setShowAdvanced] = useState(false);

  const handleGenerate = useCallback(async () => {
    if (!surname.trim()) return;
    setIsGenerating(true);
    try {
      const res = await api.generate({
        surname: surname.trim(),
        gender,
        fixed_chars: fixedChars,
        blacklist,
        count,
      });
      setTaskId(res.task_id);
      navigate(`/results/${res.task_id}`);
    } catch (err) {
      console.error('Generate failed:', err);
    } finally {
      setIsGenerating(false);
    }
  }, [surname, gender, fixedChars, blacklist, count, navigate, setTaskId, setIsGenerating]);

  return (
    <div className="space-y-6">
      <div className="text-center py-8">
        <h1 className="text-4xl font-serif font-bold text-gradient-celestial mb-2">
          天命之名
        </h1>
        <p className="text-muted-foreground font-mono text-sm tracking-wider">
          FATE · Celestial Naming Engine
        </p>
      </div>

      <Tabs defaultValue={defaultTab} className="w-full">
        <TabsList className="w-full justify-center">
          <TabsTrigger value="generator" className="gap-2 flex-1">
            <Sparkles className="h-4 w-4" />
            Generator
          </TabsTrigger>
          <TabsTrigger value="poetry" className="gap-2 flex-1">
            <BookOpen className="h-4 w-4" />
            Poetry
          </TabsTrigger>
        </TabsList>

        <TabsContent value="generator" className="space-y-4 mt-4">
          <Card className="glass-card">
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Sparkles className="h-5 w-5 text-blue-400" />
                基本配置
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div className="space-y-2">
                  <label className="text-sm font-medium text-foreground">姓氏</label>
                  <Input
                    placeholder="输入姓氏"
                    value={surname}
                    onChange={(e) => setSurname(e.target.value)}
                    className="font-serif text-lg"
                  />
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-medium text-foreground">性别</label>
                  <Select value={gender} onValueChange={(v) => setGender(v as 'male' | 'female' | 'neutral')}>
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="male">男</SelectItem>
                      <SelectItem value="female">女</SelectItem>
                      <SelectItem value="neutral">中性</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </div>

              <div className="space-y-2">
                <label className="text-sm font-medium text-foreground">生成数量</label>
                <Select value={String(count)} onValueChange={(v) => setCount(Number(v))}>
                  <SelectTrigger className="w-32">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="10">10个</SelectItem>
                    <SelectItem value="20">20个</SelectItem>
                    <SelectItem value="30">30个</SelectItem>
                    <SelectItem value="50">50个</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              <Button
                variant="celestial"
                className="w-full h-12 text-base"
                onClick={handleGenerate}
                disabled={isGenerating || !surname.trim()}
              >
                {isGenerating ? (
                  <>
                    <Loader2 className="h-5 w-5 animate-spin" />
                    生成中...
                  </>
                ) : (
                  <>
                    <Sparkles className="h-5 w-5" />
                    开始命名
                  </>
                )}
              </Button>
            </CardContent>
          </Card>

          <div className="flex items-center justify-center">
            <Button
              variant="ghost"
              size="sm"
              onClick={() => setShowAdvanced(!showAdvanced)}
              className="gap-2"
            >
              <Settings className="h-4 w-4" />
              {showAdvanced ? '收起高级配置' : '高级配置'}
            </Button>
          </div>

          {showAdvanced && (
            <Card className="glass-card">
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Settings className="h-5 w-5 text-violet-400" />
                  高级配置
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <FixedCharInput
                  surname={surname}
                  value={fixedChars}
                  onChange={setFixedChars}
                />
                <BlacklistPanel value={blacklist} onChange={setBlacklist} />
              </CardContent>
            </Card>
          )}
        </TabsContent>

        <TabsContent value="poetry" className="mt-4">
          <Card className="glass-card">
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <BookOpen className="h-5 w-5 text-amber-400" />
                诗词查字
              </CardTitle>
            </CardHeader>
            <CardContent>
              <PoetrySearchTab />
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  );
}
