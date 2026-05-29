import { useState, useCallback } from 'react';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { X } from 'lucide-react';
import { cn } from '@/lib/utils';

interface FixedCharInputProps {
  surname: string;
  value: { position: number; char: string }[];
  onChange: (value: { position: number; char: string }[]) => void;
}

export function FixedCharInput({ surname, value, onChange }: FixedCharInputProps) {
  const [inputValue, setInputValue] = useState('');
  const totalLength = surname.length + 2;

  const slots = Array.from({ length: totalLength }, (_, i) => {
    if (i < surname.length) {
      return { type: 'surname' as const, char: surname[i], position: i };
    }
    return { type: 'given' as const, char: '', position: i };
  });

  const fixedPositions = new Set(value.map((v) => v.position));

  const handleSlotClick = useCallback(
    (position: number) => {
      if (position < surname.length) return;
      const existing = value.find((v) => v.position === position);
      if (existing) {
        onChange(value.filter((v) => v.position !== position));
      } else if (inputValue.trim()) {
        onChange([...value, { position, char: inputValue.trim() }]);
        setInputValue('');
      }
    },
    [value, onChange, inputValue, surname.length],
  );

  const handleRemove = useCallback(
    (position: number) => {
      onChange(value.filter((v) => v.position !== position));
    },
    [value, onChange],
  );

  return (
    <div className="space-y-3">
      <label className="text-sm font-medium text-foreground">固定用字</label>
      <Input
        placeholder="输入要固定的字，然后点击下方位置"
        value={inputValue}
        onChange={(e) => setInputValue(e.target.value)}
        className="font-serif text-base"
      />
      <div className="flex items-center gap-2">
        {slots.map((slot) => (
          <button
            key={slot.position}
            onClick={() => handleSlotClick(slot.position)}
            className={cn(
              'flex h-12 w-12 items-center justify-center rounded-md border text-lg font-serif transition-all duration-200',
              slot.type === 'surname'
                ? 'border-border bg-white/5 text-muted-foreground cursor-default'
                : fixedPositions.has(slot.position)
                  ? 'border-blue-500/50 bg-blue-500/10 text-blue-300 cursor-pointer hover:bg-blue-500/20'
                  : 'border-dashed border-border bg-transparent text-muted-foreground cursor-pointer hover:border-white/30 hover:bg-white/5',
            )}
            disabled={slot.type === 'surname'}
          >
            {slot.char || (slot.type === 'given' ? '·' : slot.char)}
          </button>
        ))}
      </div>
      {value.length > 0 && (
        <div className="flex flex-wrap gap-2">
          {value.map((v) => (
            <Badge key={v.position} variant="celestial" className="gap-1">
              位置{v.position + 1}: {v.char}
              <button onClick={() => handleRemove(v.position)} className="ml-1 hover:text-white">
                <X className="h-3 w-3" />
              </button>
            </Badge>
          ))}
        </div>
      )}
    </div>
  );
}
