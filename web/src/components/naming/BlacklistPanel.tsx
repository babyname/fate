import { useState, useCallback } from 'react';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { X, Plus } from 'lucide-react';

interface BlacklistPanelProps {
  value: string[];
  onChange: (value: string[]) => void;
}

export function BlacklistPanel({ value, onChange }: BlacklistPanelProps) {
  const [inputValue, setInputValue] = useState('');

  const handleAdd = useCallback(() => {
    const chars = inputValue.trim().split('').filter((c) => c.trim());
    if (chars.length === 0) return;
    const newChars = chars.filter((c) => !value.includes(c));
    if (newChars.length > 0) {
      onChange([...value, ...newChars]);
    }
    setInputValue('');
  }, [inputValue, value, onChange]);

  const handleRemove = useCallback(
    (char: string) => {
      onChange(value.filter((c) => c !== char));
    },
    [value, onChange],
  );

  const handleKeyDown = useCallback(
    (e: React.KeyboardEvent) => {
      if (e.key === 'Enter') {
        e.preventDefault();
        handleAdd();
      }
    },
    [handleAdd],
  );

  return (
    <div className="space-y-3">
      <label className="text-sm font-medium text-foreground">避用字</label>
      <div className="flex gap-2">
        <Input
          placeholder="输入要避开的字"
          value={inputValue}
          onChange={(e) => setInputValue(e.target.value)}
          onKeyDown={handleKeyDown}
          className="font-serif"
        />
        <Button variant="outline" size="icon" onClick={handleAdd}>
          <Plus className="h-4 w-4" />
        </Button>
      </div>
      {value.length > 0 && (
        <div className="flex flex-wrap gap-2">
          {value.map((char) => (
            <Badge key={char} variant="destructive" className="gap-1">
              {char}
              <button onClick={() => handleRemove(char)} className="ml-1 hover:text-white">
                <X className="h-3 w-3" />
              </button>
            </Badge>
          ))}
        </div>
      )}
    </div>
  );
}
