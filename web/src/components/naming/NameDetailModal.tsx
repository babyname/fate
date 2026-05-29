import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from '@/components/ui/dialog';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Separator } from '@/components/ui/separator';
import { Printer, Star, BookOpen, Flame } from 'lucide-react';
import type { NameResult } from '@/types/api';

interface NameDetailModalProps {
  name: NameResult | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

const elementColors: Record<string, string> = {
  金: 'text-amber-300',
  木: 'text-green-400',
  水: 'text-blue-400',
  火: 'text-red-400',
  土: 'text-yellow-600',
};

export function NameDetailModal({ name, open, onOpenChange }: NameDetailModalProps) {
  if (!name) return null;

  const handlePrint = () => {
    const printContent = `
      <html>
        <head><title>${name.name} - FATE命名报告</title></head>
        <body style="font-family: serif; max-width: 600px; margin: 40px auto; padding: 20px;">
          <h1 style="text-align: center; font-size: 2em;">${name.name}</h1>
          <p style="text-align: center; color: #666;">${name.pronunciation.pinyin} | ${name.pronunciation.tone_pattern}</p>
          <hr/>
          <h2>释义</h2>
          <p>${name.meaning}</p>
          <h2>五行分析</h2>
          <p>五行: ${name.wuxing.elements.join(' ')} | 平衡度: ${name.wuxing.balance} | 缺: ${name.wuxing.missing.join('') || '无'}</p>
          <p>${name.wuxing.description}</p>
          ${name.poetry.length > 0 ? `<h2>诗词出处</h2>${name.poetry.map((p) => `<p>《${p.title}》— ${p.dynasty}·${p.author}<br/>${p.content}</p>`).join('')}` : ''}
        </body>
      </html>
    `;
    const win = window.open('', '_blank');
    if (win) {
      win.document.write(printContent);
      win.document.close();
      win.print();
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-2xl max-h-[85vh] overflow-y-auto scrollbar-thin">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-3">
            <span className="text-3xl font-serif font-bold text-gradient-celestial">
              {name.name}
            </span>
            <span className="text-sm font-mono text-muted-foreground">
              {name.pronunciation.pinyin}
            </span>
          </DialogTitle>
          <DialogDescription className="sr-only">
            {name.name}的详细命名分析报告
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-6">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <Star className="h-5 w-5 text-amber-400" />
              <span className="text-2xl font-bold font-mono text-gradient-stardust">
                {name.score}
              </span>
              <span className="text-sm text-muted-foreground">综合评分</span>
            </div>
            <Button variant="outline" size="sm" onClick={handlePrint}>
              <Printer className="h-4 w-4" />
              打印
            </Button>
          </div>

          <div className="glass-card p-4 space-y-2">
            <h3 className="text-sm font-medium text-foreground flex items-center gap-2">
              释义
            </h3>
            <p className="text-sm leading-relaxed text-foreground/80 font-serif">
              {name.meaning}
            </p>
          </div>

          <div className="glass-card p-4 space-y-3">
            <h3 className="text-sm font-medium text-foreground flex items-center gap-2">
              <Flame className="h-4 w-4 text-orange-400" />
              五行分析
            </h3>
            <div className="flex flex-wrap gap-2">
              {name.wuxing.elements.map((el) => (
                <Badge key={el} variant="celestial" className={`text-sm ${elementColors[el] || ''}`}>
                  {el}
                </Badge>
              ))}
              {name.wuxing.missing.length > 0 && (
                <Badge variant="outline" className="text-sm text-muted-foreground">
                  缺{name.wuxing.missing.join('')}
                </Badge>
              )}
            </div>
            <div className="flex items-center gap-2">
              <span className="text-xs text-muted-foreground">平衡度</span>
              <div className="flex-1 h-2 rounded-full bg-white/5 overflow-hidden">
                <div
                  className="h-full rounded-full bg-gradient-to-r from-blue-500 to-violet-500 transition-all"
                  style={{ width: `${name.wuxing.balance}%` }}
                />
              </div>
              <span className="text-xs font-mono text-muted-foreground">{name.wuxing.balance}%</span>
            </div>
            <p className="text-xs text-muted-foreground">{name.wuxing.description}</p>
          </div>

          {name.poetry.length > 0 && (
            <div className="glass-card p-4 space-y-3">
              <h3 className="text-sm font-medium text-foreground flex items-center gap-2">
                <BookOpen className="h-4 w-4 text-blue-400" />
                诗词出处
              </h3>
              {name.poetry.map((poem, i) => (
                <div key={i} className="space-y-1">
                  <div className="flex items-center gap-2">
                    <span className="font-serif font-medium text-foreground">
                      《{poem.title}》
                    </span>
                    <Badge variant="stardust" className="text-xs">
                      {poem.dynasty}
                    </Badge>
                    <span className="text-xs text-muted-foreground">{poem.author}</span>
                  </div>
                  <p className="font-serif text-sm leading-relaxed text-foreground/70 whitespace-pre-line">
                    {poem.content}
                  </p>
                  <p className="text-xs text-blue-400">
                    匹配字: {poem.matched_char}
                  </p>
                  {i < name.poetry.length - 1 && <Separator className="my-2 bg-white/5" />}
                </div>
              ))}
            </div>
          )}

          <div className="glass-card p-4 space-y-2">
            <h3 className="text-sm font-medium text-foreground">音律分析</h3>
            <div className="grid grid-cols-3 gap-4 text-center">
              <div>
                <p className="text-xs text-muted-foreground">拼音</p>
                <p className="font-mono text-sm text-foreground">{name.pronunciation.pinyin}</p>
              </div>
              <div>
                <p className="text-xs text-muted-foreground">声调</p>
                <p className="font-mono text-sm text-foreground">{name.pronunciation.tone_pattern}</p>
              </div>
              <div>
                <p className="text-xs text-muted-foreground">韵律</p>
                <p className="font-mono text-sm text-foreground">{name.pronunciation.rhythm}</p>
              </div>
            </div>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}
