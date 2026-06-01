import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Label } from '@/components/ui/label';
import { Badge } from '@/components/ui/badge';
import { Loader2, Search, BookOpen } from 'lucide-react';
import { api } from '@/lib/api';
import type { PoetrySearchResult } from '@/types/api';

export function PoetrySearchTab() {
  const [keyword, setKeyword] = useState('');
  const [results, setResults] = useState<PoetrySearchResult[]>([]);
  const [loading, setLoading] = useState(false);
  const [searched, setSearched] = useState(false);

  const handleSearch = async () => {
    if (!keyword.trim()) return;
    setLoading(true);
    setSearched(true);
    try {
      const res = await api.searchPoetry(keyword);
      setResults(res.results || []);
    } catch {
      setResults([]);
    } finally {
      setLoading(false);
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      handleSearch();
    }
  };

  const typeColorMap: Record<string, string> = {
    '诗': 'bg-red-100 text-red-800',
    '词': 'bg-blue-100 text-blue-800',
    '曲': 'bg-orange-100 text-orange-800',
    '赋': 'bg-purple-100 text-purple-800',
  };

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>诗词搜索</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex gap-2">
            <Input
              placeholder="输入关键词搜索诗词，例如：清风明月"
              value={keyword}
              onChange={(e) => setKeyword(e.target.value)}
              onKeyDown={handleKeyDown}
              className="flex-1"
            />
            <Button
              onClick={handleSearch}
              disabled={loading}
            >
              {loading ? (
                <Loader2 className="h-4 w-4 animate-spin mr-2" />
              ) : (
                <Search className="h-4 w-4 mr-2" />
              )}
              搜索
            </Button>
          </div>
        </CardContent>
      </Card>

      {loading && (
        <div className="flex flex-col items-center justify-center py-12 text-muted-foreground">
          <Loader2 className="mb-2 h-8 w-8 animate-spin" />
          <span className="text-sm">正在搜索诗词...</span>
        </div>
      )}

      {!loading && searched && results.length === 0 && (
        <div className="flex flex-col items-center justify-center py-12 text-muted-foreground">
          <BookOpen className="mb-2 h-10 w-10 opacity-40" />
          <span className="text-sm">未找到相关诗词，换个关键词试试</span>
        </div>
      )}

      {!loading && results.length > 0 && (
        <div className="space-y-3">
          {results.map((item, index) => (
            <Card key={index}>
              <CardContent className="pt-6">
                <div className="mb-2 flex items-center gap-2">
                  <span className="font-medium text-foreground" style={{ fontFamily: 'STKaiti, KaiTi, serif' }}>
                    {item.title}
                  </span>
                  <span className="text-xs text-muted-foreground">
                    {item.dynasty}·{item.author}
                  </span>
                  <Badge
                    variant="secondary"
                    className={typeColorMap[item.type] || 'bg-gray-100 text-gray-800'}
                  >
                    {item.type}
                  </Badge>
                </div>
                <p className="text-sm leading-relaxed text-foreground" style={{ fontFamily: 'STKaiti, KaiTi, serif' }}>
                  {item.sentence}
                </p>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
