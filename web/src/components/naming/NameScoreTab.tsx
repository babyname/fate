import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Loader2, Star, BookOpen } from 'lucide-react';
import { api } from '@/lib/api';
import type { NameResult } from '@/types/api';
import { useAppStore } from '@/store/app';

export function NameScoreTab() {
  const [surname, setSurname] = useState('');
  const [name1, setName1] = useState('');
  const [name2, setName2] = useState('');
  const [born, setBorn] = useState('');
  const [sex, setSex] = useState('1');
  const [result, setResult] = useState<NameResult | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const setSelectedName = useAppStore((s) => s.setSelectedName);
  const setDetailModalOpen = useAppStore((s) => s.setDetailModalOpen);

  const handleScore = async () => {
    if (!surname || !name1 || !name2) {
      setError('请输入完整的姓名');
      return;
    }
    setLoading(true);
    setError('');
    try {
      const res = await api.nameScore({
        surname,
        name1,
        name2,
        born,
        sex,
      });
      setResult(res.name_result);
    } catch (err) {
      setError(err instanceof Error ? err.message : '测分失败');
    } finally {
      setLoading(false);
    }
  };

  const handleViewDetail = () => {
    if (result) {
      setSelectedName(result);
      setDetailModalOpen(true);
    }
  };

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>姓名测分</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-2">
              <span className="text-sm text-muted-foreground">姓氏</span>
              <Input
                placeholder="例如：李"
                value={surname}
                onChange={(e) => setSurname(e.target.value)}
              />
            </div>
            <div className="grid grid-cols-2 gap-2">
              <div className="space-y-2">
                <span className="text-sm text-muted-foreground">名一</span>
                <Input
                  placeholder="例如：明"
                  value={name1}
                  onChange={(e) => setName1(e.target.value)}
                />
              </div>
              <div className="space-y-2">
                <span className="text-sm text-muted-foreground">名二</span>
                <Input
                  placeholder="例如：华"
                  value={name2}
                  onChange={(e) => setName2(e.target.value)}
                />
              </div>
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-2">
              <span className="text-sm text-muted-foreground">出生日期 (可选)</span>
              <Input
                type="datetime-local"
                value={born}
                onChange={(e) => setBorn(e.target.value)}
              />
            </div>
            <div className="space-y-2">
              <span className="text-sm text-muted-foreground">性别 (可选)</span>
              <div className="flex gap-2">
                <Button
                  variant={sex === '1' ? 'default' : 'outline'}
                  size="sm"
                  onClick={() => setSex('1')}
                >
                  男
                </Button>
                <Button
                  variant={sex === '2' ? 'default' : 'outline'}
                  size="sm"
                  onClick={() => setSex('2')}
                >
                  女
                </Button>
              </div>
            </div>
          </div>

          {error && (
            <p className="text-red-500 text-sm">{error}</p>
          )}

          <Button
            onClick={handleScore}
            disabled={loading}
            className="w-full"
          >
            {loading ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                测算中...
              </>
            ) : (
                '立即测分'
            )}
          </Button>
        </CardContent>
      </Card>

      {result && (
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle>{result.full_name}</CardTitle>
              <Button size="sm" onClick={handleViewDetail}>查看详情</Button>
            </div>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center gap-2">
              <div className="flex items-center">
                <Star className={result.score >= 90 ? 'text-yellow-400' : result.score >= 75 ? 'text-blue-400' : 'text-gray-400'} />
                <span className="text-3xl font-bold ml-2">{result.score}</span>
              </div>
              <Badge variant={result.score >= 90 ? 'celestial' : 'secondary'} className="text-lg px-3 py-1">
                {result.grade}
              </Badge>
            </div>

            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              <div className="text-center p-3 bg-surface rounded-lg">
                <p className="text-xs text-muted-foreground mb-1">文化印象</p>
                <p className="font-bold">{result.score_detail.wen_hua_yin_xiang}</p>
              </div>
              <div className="text-center p-3 bg-surface rounded-lg">
                <p className="text-xs text-muted-foreground mb-1">五行八字</p>
                <p className="font-bold">{result.score_detail.wu_xing_ba_zi}</p>
              </div>
              <div className="text-center p-3 bg-surface rounded-lg">
                <p className="text-xs text-muted-foreground mb-1">生肖</p>
                <p className="font-bold">{result.score_detail.sheng_xiao}</p>
              </div>
              <div className="text-center p-3 bg-surface rounded-lg">
                <p className="text-xs text-muted-foreground mb-1">五格数理</p>
                <p className="font-bold">{result.score_detail.wu_ge_shu_li}</p>
              </div>
            </div>

            <div className="flex items-center gap-2">
              <Badge variant="celestial">{result.char1.wu_xing}</Badge>
              <Badge variant="celestial">{result.char2.wu_xing}</Badge>
              {result.has_poetry && (
                <Badge variant="stardust" className="gap-1">
                  <BookOpen className="h-3 w-3" />
                  有诗词出处
                </Badge>
              )}
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
