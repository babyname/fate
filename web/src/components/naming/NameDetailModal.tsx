import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from '@/components/ui/dialog';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Printer, Star, Flame, Sparkles, Hexagon } from 'lucide-react';
import type { NameResult } from '@/types/api';

interface NameDetailModalProps {
  name: NameResult | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

const wuxingColors: Record<string, string> = {
  金: 'text-amber-300',
  木: 'text-green-400',
  水: 'text-blue-400',
  火: 'text-red-400',
  土: 'text-yellow-600',
};

const luckyColor = (s: string) => {
  if (s === '大吉' || s === '吉') return 'text-green-400';
  if (s === '凶' || s === '大凶') return 'text-red-400';
  return 'text-muted-foreground';
};

export function NameDetailModal({ name, open, onOpenChange }: NameDetailModalProps) {
  if (!name) return null;

  const handlePrint = () => {
    const printContent = `<!DOCTYPE html>
<html>
<head><meta charset="utf-8"><title>${name.full_name} - FATE</title>
<style>
body{font-family:"Microsoft YaHei",serif;max-width:600px;margin:40px auto;padding:20px;color:#333}
h1{text-align:center;font-size:2.5em;margin-bottom:0}
h2{border-bottom:1px solid #ddd;padding-bottom:4px;margin-top:24px}
table{width:100%;border-collapse:collapse;margin:8px 0}
td,th{border:1px solid #ddd;padding:6px 10px;text-align:center;font-size:13px}
th{background:#f5f5f5}
.score{text-align:center;font-size:3em;font-weight:bold;color:#e5a100}
</style></head>
<body>
<h1>${name.full_name}</h1>
<p style="text-align:center;color:#666">${name.char1.pinyin} ${name.char2.pinyin} | ${name.strokes}画</p>
<div class="score">${name.score.toFixed(0)}<span style="font-size:16px;color:#666"> / ${name.grade}</span></div>
<h2>释义</h2>
<p>${name.char1.meaning || '无'}；${name.char2.meaning || '无'}</p>
<h2>五格分析</h2>
<table>
<tr><th>格</th><th>笔画</th><th>吉凶</th><th>五行</th></tr>
<tr><td>天格</td><td>${name.wu_ge.tian_ge.stroke}</td><td>${name.wu_ge.tian_ge.lucky}</td><td>${name.wu_ge.tian_ge.yin_yang_wu_xing}</td></tr>
<tr><td>人格</td><td>${name.wu_ge.ren_ge.stroke}</td><td>${name.wu_ge.ren_ge.lucky}</td><td>${name.wu_ge.ren_ge.yin_yang_wu_xing}</td></tr>
<tr><td>地格</td><td>${name.wu_ge.di_ge.stroke}</td><td>${name.wu_ge.di_ge.lucky}</td><td>${name.wu_ge.di_ge.yin_yang_wu_xing}</td></tr>
<tr><td>外格</td><td>${name.wu_ge.wai_ge.stroke}</td><td>${name.wu_ge.wai_ge.lucky}</td><td>${name.wu_ge.wai_ge.yin_yang_wu_xing}</td></tr>
<tr><td>总格</td><td>${name.wu_ge.zong_ge.stroke}</td><td>${name.wu_ge.zong_ge.lucky}</td><td>${name.wu_ge.zong_ge.yin_yang_wu_xing}</td></tr>
</table>
<h2>三才</h2><p>${name.san_cai} — ${name.san_cai_luck}</p><p>${name.san_cai_detail}</p>
<h2>周易</h2><p>${name.zhou_yi.ben_gua_name} — ${name.zhou_yi.ben_gua_ji_xiong}</p><p>${name.zhou_yi.da_xiang}</p>
${name.bazi ? `<h2>八字</h2><p>四柱: ${name.bazi.four_pillars.join(' ')}</p><p>纳音: ${name.bazi.na_yin.join(' ')}</p><p>生肖: ${name.bazi.zodiac} | 星座: ${name.bazi.constellation}</p>` : ''}
${name.wu_xing ? `<h2>五行</h2><p>喜: ${name.wu_xing.xi} | 忌: ${name.wu_xing.ji}</p><p>${name.wu_xing.analysis}</p>` : ''}
<h2>评分明细</h2>
<table>
<tr><th>维度</th><th>得分</th></tr>
<tr><td>文化印象</td><td>${name.score_detail.wen_hua_yin_xiang}</td></tr>
<tr><td>五行八字</td><td>${name.score_detail.wu_xing_ba_zi}</td></tr>
<tr><td>生肖</td><td>${name.score_detail.sheng_xiao}</td></tr>
<tr><td>五格数理</td><td>${name.score_detail.wu_ge_shu_li}</td></tr>
<tr><td>音韵</td><td>${name.score_detail.yin_yun}</td></tr>
</table>
</body></html>`;
    const win = window.open('', '_blank');
    if (win) {
      win.document.write(printContent);
      win.document.close();
      win.print();
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-2xl max-h-[85vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-3">
            <span className="text-3xl font-serif font-bold text-gradient-celestial">
              {name.full_name}
            </span>
            <span className="text-sm font-mono text-muted-foreground">
              {name.char1.pinyin} {name.char2.pinyin}
            </span>
          </DialogTitle>
          <DialogDescription className="sr-only">
            {name.full_name}的详细命名分析报告
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-6">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <Star className="h-5 w-5 text-amber-400" />
              <span className="text-2xl font-bold font-mono text-gradient-stardust">
                {name.score.toFixed(0)}
              </span>
              <span className="text-sm text-muted-foreground">综合评分</span>
              <Badge variant="celestial">{name.grade}</Badge>
            </div>
            <Button variant="outline" size="sm" onClick={handlePrint}>
              <Printer className="h-4 w-4" />
              打印
            </Button>
          </div>

          <div className="glass-card p-4 space-y-2">
            <h3 className="text-sm font-medium text-foreground flex items-center gap-2">
              <Sparkles className="h-4 w-4 text-blue-400" />
              释义
            </h3>
            <div className="space-y-1">
              <p className="text-sm leading-relaxed text-foreground/80">
                <span className="font-serif text-lg">{name.char1.char}</span> — {name.char1.meaning || '无释义'}
              </p>
              <p className="text-sm leading-relaxed text-foreground/80">
                <span className="font-serif text-lg">{name.char2.char}</span> — {name.char2.meaning || '无释义'}
              </p>
            </div>
          </div>

          <div className="glass-card p-4 space-y-3">
            <h3 className="text-sm font-medium text-foreground flex items-center gap-2">
              <Hexagon className="h-4 w-4 text-violet-400" />
              五格分析
            </h3>
            <div className="grid grid-cols-5 gap-2 text-center">
              {(['tian_ge', 'ren_ge', 'di_ge', 'wai_ge', 'zong_ge'] as const).map((key) => {
                const ge = name.wu_ge[key];
                const labels: Record<string, string> = { tian_ge: '天', ren_ge: '人', di_ge: '地', wai_ge: '外', zong_ge: '总' };
                return (
                  <div key={key} className="space-y-1">
                    <p className="text-xs text-muted-foreground">{labels[key]}格</p>
                    <p className="text-lg font-mono font-bold text-foreground">{ge.stroke}</p>
                    <p className={`text-xs ${luckyColor(ge.lucky)}`}>{ge.lucky}</p>
                    <p className="text-xs text-muted-foreground">{ge.yin_yang_wu_xing}</p>
                  </div>
                );
              })}
            </div>
          </div>

          <div className="glass-card p-4 space-y-2">
            <h3 className="text-sm font-medium text-foreground">三才配置</h3>
            <div className="flex items-center gap-3">
              <span className="font-mono text-lg text-foreground">{name.san_cai}</span>
              <Badge variant={name.san_cai_luck === '吉' || name.san_cai_luck === '大吉' ? 'celestial' : 'destructive'}>
                {name.san_cai_luck}
              </Badge>
            </div>
            <p className="text-xs text-muted-foreground">{name.san_cai_detail}</p>
            <div className="grid grid-cols-2 gap-2 text-xs text-muted-foreground">
              <p>基础运: {name.ji_chu_yun}</p>
              <p>成功运: {name.cheng_gong_yun}</p>
              <p className="col-span-2">人际关系: {name.ren_ji_guan_xi}</p>
            </div>
          </div>

          {name.wu_xing && (
            <div className="glass-card p-4 space-y-3">
              <h3 className="text-sm font-medium text-foreground flex items-center gap-2">
                <Flame className="h-4 w-4 text-orange-400" />
                五行分析
              </h3>
              <div className="grid grid-cols-2 gap-2 text-sm">
                <div><span className="text-muted-foreground">喜: </span><span className={wuxingColors[name.wu_xing.xi] || ''}>{name.wu_xing.xi}</span></div>
                <div><span className="text-muted-foreground">忌: </span><span className={wuxingColors[name.wu_xing.ji] || ''}>{name.wu_xing.ji}</span></div>
                <div><span className="text-muted-foreground">用神: </span><span>{name.wu_xing.yong_shen}</span></div>
                <div><span className="text-muted-foreground">格局: </span><span>{name.wu_xing.geju_name}</span></div>
              </div>
              <p className="text-xs text-muted-foreground">{name.wu_xing.analysis}</p>
            </div>
          )}

          {name.bazi && (
            <div className="glass-card p-4 space-y-2">
              <h3 className="text-sm font-medium text-foreground">八字</h3>
              <div className="grid grid-cols-4 gap-1 text-center text-sm">
                {name.bazi.four_pillars.map((p, i) => (
                  <div key={i}>
                    <p className="font-mono text-foreground">{p}</p>
                    <p className="text-xs text-muted-foreground">{name.bazi!.five_elements[i]}</p>
                    <p className="text-xs text-muted-foreground">{name.bazi!.na_yin[i]}</p>
                  </div>
                ))}
              </div>
              <div className="flex gap-2 text-xs text-muted-foreground">
                <span>生肖: {name.bazi.zodiac}</span>
                <span>星座: {name.bazi.constellation}</span>
              </div>
            </div>
          )}

          <div className="glass-card p-4 space-y-3">
            <h3 className="text-sm font-medium text-foreground flex items-center gap-2">
              <Star className="h-4 w-4 text-amber-400" />
              评分明细
            </h3>
            {([
              { key: 'wen_hua_yin_xiang', label: '文化印象' },
              { key: 'wu_xing_ba_zi', label: '五行八字' },
              { key: 'sheng_xiao', label: '生肖' },
              { key: 'wu_ge_shu_li', label: '五格数理' },
              { key: 'yin_yun', label: '音韵' },
            ] as const).map(({ key, label }) => {
              const score = name.score_detail[key as keyof typeof name.score_detail];
              return (
                <div key={key} className="flex items-center gap-3">
                  <span className="w-16 text-xs text-muted-foreground">{label}</span>
                  <div className="flex-1 h-2 rounded-full bg-white/5 overflow-hidden">
                    <div
                      className="h-full rounded-full bg-gradient-to-r from-blue-500 to-violet-500 transition-all"
                      style={{ width: `${score}%` }}
                    />
                  </div>
                  <span className="text-xs font-mono text-muted-foreground w-8 text-right">{score}</span>
                </div>
              );
            })}
          </div>

          {name.zhou_yi && (
            <div className="glass-card p-4 space-y-2">
              <h3 className="text-sm font-medium text-foreground">周易卦象</h3>
              <p className="text-sm text-foreground">
                <span className="font-serif">{name.zhou_yi.ben_gua_name}</span>
                <Badge
                  variant={name.zhou_yi.ben_gua_ji_xiong === '吉' ? 'celestial' : 'outline'}
                  className="ml-2 text-xs"
                >
                  {name.zhou_yi.ben_gua_ji_xiong}
                </Badge>
              </p>
              <p className="text-xs text-muted-foreground">{name.zhou_yi.da_xiang}</p>
              {name.zhou_yi.bian_gua_name && (
                <p className="text-xs text-muted-foreground">
                  变卦: {name.zhou_yi.bian_gua_name} | 动爻: {name.zhou_yi.dong_yao_desc}
                </p>
              )}
            </div>
          )}

          {name.interpret && (
            <div className="glass-card p-4 space-y-2">
              <h3 className="text-sm font-medium text-foreground">综合解读</h3>
              <p className="text-sm leading-relaxed text-foreground/80 font-serif">
                {name.interpret}
              </p>
            </div>
          )}
        </div>
      </DialogContent>
    </Dialog>
  );
}