import { useId, useState } from 'react';
import { SlidersHorizontal } from 'lucide-react';
import { Button } from '@/components/ui/button';
import {
  Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger,
} from '@/components/ui/dialog';
import { DEFAULT_TRADITIONAL_WEIGHT } from '@/lib/recommendation-scores';
import { useScoringPreferences } from '@/store/scoring-preferences';

const presets = [
  { label: '偏重传统数理', weight: 70 },
  { label: '各占一半', weight: 50 },
  { label: '偏重音义', weight: 30 },
];

export function ScoreWeightControl({ onApply }: { onApply?: () => void }) {
  const traditionalWeight = useScoringPreferences((s) => s.traditionalWeight);
  const setTraditionalWeight = useScoringPreferences((s) => s.setTraditionalWeight);
  const [open, setOpen] = useState(false);
  const [draft, setDraft] = useState(traditionalWeight);
  const sliderId = useId();

  const apply = () => {
    if (draft !== traditionalWeight) {
      setTraditionalWeight(draft);
      onApply?.();
    }
    setOpen(false);
  };

  return (
    <div className="flex flex-wrap items-center gap-x-3 gap-y-1">
      <Dialog open={open} onOpenChange={(next) => {
        if (next) setDraft(traditionalWeight);
        setOpen(next);
      }}>
        <DialogTrigger asChild>
          <Button type="button" variant="outline" size="sm">
            <SlidersHorizontal aria-hidden="true" className="h-4 w-4" />
            评分权重
          </Button>
        </DialogTrigger>
        <DialogContent className="w-[calc(100%-2rem)]">
          <DialogHeader>
            <DialogTitle>调整评分权重</DialogTitle>
            <DialogDescription>按你的偏好分配权重，两项合计 100%。</DialogDescription>
          </DialogHeader>
          <div className="space-y-3">
            <div className="flex justify-between gap-2 text-sm font-medium">
              <label htmlFor={sliderId}>传统数理 {draft}%</label>
              <span>音义 {100 - draft}%</span>
            </div>
            <input
              id={sliderId}
              type="range"
              min={0}
              max={100}
              step={1}
              value={draft}
              aria-label="传统数理权重"
              aria-valuetext={`传统数理 ${draft}%，音义 ${100 - draft}%`}
              onChange={(event) => setDraft(event.currentTarget.valueAsNumber)}
              className="w-full cursor-pointer accent-primary"
            />
            <div className="flex flex-wrap gap-2">
              {presets.map(({ label, weight }) => (
                <Button
                  key={weight}
                  type="button"
                  variant={draft === weight ? 'secondary' : 'outline'}
                  size="sm"
                  aria-pressed={draft === weight}
                  onClick={() => setDraft(weight)}
                >
                  {label}
                </Button>
              ))}
            </div>
          </div>
          <div className="flex flex-wrap items-center justify-between gap-2">
            <Button type="button" variant="ghost" size="sm" onClick={() => setDraft(DEFAULT_TRADITIONAL_WEIGHT)}>
              恢复默认
            </Button>
            <Button type="button" onClick={apply}>应用权重</Button>
          </div>
        </DialogContent>
      </Dialog>
      <span className="text-xs text-muted-foreground tabular-nums">
        传统数理 {traditionalWeight}% · 音义 {100 - traditionalWeight}%
      </span>
    </div>
  );
}
