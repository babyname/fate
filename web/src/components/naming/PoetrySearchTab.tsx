import { useState, useCallback } from 'react';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Search, BookOpen } from 'lucide-react';
import { api } from '@/lib/api';
import type { PoetryItem } from '@/types/api';

export function PoetrySearchTab() {
  const [keyword, setKeyword] = useState('');
  const [dynasty, setDynasty] = useState('');
  const [results, setResults] = useState<PoetryItem[]>([]);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);

  const handleSearch = useCallback(
    async (p = 1) => {
      if (!keyword.trim()) return;
      setLoading(true);
      try {
        const res = await api.searchPoetry({
          keyword: keyword.trim(),
          dynasty: dynasty || undefined,
          page: p,
          page_size: 10,
        });
        setResults(res.poems);
        setTotal(res.total);
        setPage(p);
      } catch {
        setResults([]);
        setTotal(0);
      } finally {
        setLoading(false);
      }
    },
    [keyword, dynasty],
  );

  return (
    <div className="space-y-4">
      <div className="flex gap-2">
        <Input
          placeholder="搜索诗词关键字..."
          value={keyword}
          onChange={(e) => setKeyword(e.target.value)}
          onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
          className="font-serif"
        />
        <Input
          placeholder="朝代（可选）"
          value={dynasty}
          onChange={(e) => setDynasty(e.target.value)}
          className="w-28"
        />
        <Button variant="celestial" onClick={() => handleSearch()} disabled={loading || !keyword.trim()}>
          <Search className="h-4 w-4" />
        </Button>
      </div>

      {total > 0 && (
        <p className="text-sm text-muted-foreground">
          找到 <span className="text-foreground font-medium">{total}</span> 首相关诗词
        </p>
      )}

      <div className="space-y-3">
        {results.map((poem, i) => (
          <Card key={i} className="glass-card-hover">
            <CardHeader className="pb-2">
              <div className="flex items-center justify-between">
                <CardTitle className="font-serif text-base">{poem.title}</CardTitle>
                <div className="flex items-center gap-2">
                  <Badge variant="stardust" className="text-xs">
                    {poem.dynasty}
                  </Badge>
                  <span className="text-sm text-muted-foreground">{poem.author}</span>
                </div>
              </div>
            </CardHeader>
            <CardContent>
              <p className="font-serif text-sm leading-relaxed text-foreground/80 whitespace-pre-line">
                {poem.content}
              </p>
              {poem.matched_chars && poem.matched_chars.length > 0 && (
                <div className="mt-2 flex gap-1">
                  {poem.matched_chars.map((c) => (
                    <Badge key={c} variant="celestial" className="text-xs">
                      {c}
                    </Badge>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>
        ))}
      </div>

      {total > 10 && (
        <div className="flex justify-center gap-2">
          <Button
            variant="outline"
            size="sm"
            disabled={page <= 1}
            onClick={() => handleSearch(page - 1)}
          >
            上一页
          </Button>
          <span className="flex items-center text-sm text-muted-foreground">
            {page} / {Math.ceil(total / 10)}
          </span>
          <Button
            variant="outline"
            size="sm"
            disabled={page >= Math.ceil(total / 10)}
            onClick={() => handleSearch(page + 1)}
          >
            下一页
          </Button>
        </div>
      )}

      {results.length === 0 && !loading && keyword && (
        <div className="flex flex-col items-center justify-center py-12 text-muted-foreground">
          <BookOpen className="h-12 w-12 mb-3 opacity-30" />
          <p className="text-sm">输入关键字搜索古诗词</p>
        </div>
      )}
    </div>
  );
}
