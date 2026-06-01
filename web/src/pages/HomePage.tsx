import { useState, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Badge } from '@/components/ui/badge';
import { useAppStore } from '@/store/app';
import { api } from '@/lib/api';
import { Sparkles, Settings, Loader2, X, Plus, BookOpen } from 'lucide-react';

export function HomePage() {
  const navigate = useNavigate();

  const surname = useAppStore((s) => s.surname);
  const born = useAppStore((s) => s.born);
  const sex = useAppStore((s) => s.sex);
  const avoidChars = useAppStore((s) => s.avoidChars);
  const requireChars = useAppStore((s) => s.requireChars);
  const poetryMode = useAppStore((s) => s.poetryMode);
  const isGenerating = useAppStore((s) => s.isGenerating);
  const setSurname = useAppStore((s) => s.setSurname);
  const setBorn = useAppStore((s) => s.setBorn);
  const setSex = useAppStore((s) => s.setSex);
  const setAvoidChars = useAppStore((s) => s.setAvoidChars);
  const setRequireChars = useAppStore((s) => s.setRequireChars);
  const setPoetryMode = useAppStore((s) => s.setPoetryMode);
  const setTaskId = useAppStore((s) => s.setTaskId);
  const setIsGenerating = useAppStore((s) => s.setIsGenerating);

  const [showAdvanced, setShowAdvanced] = useState(false);
  const [avoidInput, setAvoidInput] = useState('');
  const [requireInput, setRequireInput] = useState('');
  const [error, setError] = useState<string | null>(null);

  const handleGenerate = useCallback(async () => {
    if (!surname.trim() || !born.trim()) {
      setError('请填写姓氏和出生日期');
      return;
    }
    setError(null);
    setIsGenerating(true);
    try {
      const bornFormatted = born.includes(':') ? born : `${born} 12:00`;
      const bornFinal = bornFormatted.replace(/-/g, '/');
      const res = await api.generate({
        surname: surname.trim(),
        born: bornFinal,
        sex,
        poetry_mode: poetryMode,
        avoid_chars: avoidChars.length > 0 ? avoidChars : undefined,
        require_chars: requireChars.length > 0 ? requireChars : undefined,
      });
      setTaskId(res.task_id);
      navigate(`/results/${res.task_id}`);
    } catch (err) {
      setError(err instanceof Error ? err.message : '生成失败');
    } finally {
      setIsGenerating(false);
    }
  }, [surname, born, sex, poetryMode, avoidChars, requireChars, navigate, setTaskId, setIsGenerating]);

  const handleAddAvoid = useCallback(() => {
    const chars = avoidInput.trim().split('').filter((c) => c.trim() && !avoidChars.includes(c));
    if (chars.length > 0) setAvoidChars([...avoidChars, ...chars]);
    setAvoidInput('');
  }, [avoidInput, avoidChars, setAvoidChars]);

  const handleAddRequire = useCallback(() => {
    const chars = requireInput.trim().split('').filter((c) => c.trim() && !requireChars.includes(c));
    if (chars.length > 0) setRequireChars([...requireChars, ...chars]);
    setRequireInput('');
  }, [requireInput, requireChars, setRequireChars]);

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
              <Select value={sex} onValueChange={setSex}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="boy">男</SelectItem>
                  <SelectItem value="girl">女</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          <div className="space-y-2">
            <label className="text-sm font-medium text-foreground">出生日期</label>
            <Input
              type="date"
              value={born.includes('T') ? born.split('T')[0] : born}
              onChange={(e) => setBorn(e.target.value)}
              className="w-48"
            />
            <p className="text-xs text-muted-foreground">格式: YYYY/MM/DD（用于八字五行分析）</p>
          </div>

          <div className="space-y-2">
            <label className="text-sm font-medium text-foreground flex items-center gap-2">
              <BookOpen className="h-4 w-4 text-amber-400" />
              诗词模式
            </label>
            <Select value={String(poetryMode)} onValueChange={(v) => setPoetryMode(Number(v))}>
              <SelectTrigger className="w-64">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="0">不限 — 所有字均可</SelectItem>
                <SelectItem value="2">诗词优先 — 仅选有诗词出处的字</SelectItem>
              </SelectContent>
            </Select>
            <p className="text-xs text-muted-foreground">
              诗词优先模式下，名字中的字均出自唐诗宋词诗经，更具文化底蕴
            </p>
          </div>

          {error && (
            <p className="text-sm text-destructive">{error}</p>
          )}

          <Button
            variant="celestial"
            className="w-full h-12 text-base"
            onClick={handleGenerate}
            disabled={isGenerating || !surname.trim() || !born.trim()}
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
          <CardContent className="space-y-6">
            <div className="space-y-3">
              <label className="text-sm font-medium text-foreground">避用字</label>
              <div className="flex gap-2">
                <Input
                  placeholder="输入要避开的字"
                  value={avoidInput}
                  onChange={(e) => setAvoidInput(e.target.value)}
                  onKeyDown={(e) => { if (e.key === 'Enter') { e.preventDefault(); handleAddAvoid(); } }}
                  className="font-serif"
                />
                <Button variant="outline" size="icon" onClick={handleAddAvoid}>
                  <Plus className="h-4 w-4" />
                </Button>
              </div>
              {avoidChars.length > 0 && (
                <div className="flex flex-wrap gap-2">
                  {avoidChars.map((c) => (
                    <Badge key={c} variant="destructive" className="gap-1">
                      {c}
                      <button onClick={() => setAvoidChars(avoidChars.filter((x) => x !== c))} className="ml-1 hover:text-white">
                        <X className="h-3 w-3" />
                      </button>
                    </Badge>
                  ))}
                </div>
              )}
            </div>

            <div className="space-y-3">
              <label className="text-sm font-medium text-foreground">必须包含的字</label>
              <div className="flex gap-2">
                <Input
                  placeholder="输入必须包含的字"
                  value={requireInput}
                  onChange={(e) => setRequireInput(e.target.value)}
                  onKeyDown={(e) => { if (e.key === 'Enter') { e.preventDefault(); handleAddRequire(); } }}
                  className="font-serif"
                />
                <Button variant="outline" size="icon" onClick={handleAddRequire}>
                  <Plus className="h-4 w-4" />
                </Button>
              </div>
              {requireChars.length > 0 && (
                <div className="flex flex-wrap gap-2">
                  {requireChars.map((c) => (
                    <Badge key={c} variant="celestial" className="gap-1">
                      {c}
                      <button onClick={() => setRequireChars(requireChars.filter((x) => x !== c))} className="ml-1 hover:text-white">
                        <X className="h-3 w-3" />
                      </button>
                    </Badge>
                  ))}
                </div>
              )}
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}