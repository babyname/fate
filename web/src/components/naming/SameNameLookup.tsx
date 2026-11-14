import { useId, useRef, useState } from 'react';
import { Copy, ExternalLink, Search } from 'lucide-react';
import { Button } from '@/components/ui/button';
import {
  Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger,
} from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { formatSameNameCount, parseSameNameCount, sameNameSources } from '@/lib/same-name';
import type { SameNameSource } from '@/lib/same-name';
import { useSameNameRecords } from '@/store/same-name-records';

const selectClass = 'h-10 rounded-md border border-border bg-background px-2 text-sm text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring';
const dateFormat = new Intl.DateTimeFormat('zh-CN', { dateStyle: 'medium' });

function SameNameRecordForm({ name, source }: { name: string; source: SameNameSource }) {
  const id = useId();
  const record = useSameNameRecords((s) => s.records.find((item) => item.name === name && item.source === source));
  const save = useSameNameRecords((s) => s.save);
  const remove = useSameNameRecords((s) => s.remove);
  const [under, setUnder] = useState(record?.under ?? false);
  const [countInput, setCountInput] = useState(record?.count.toString() ?? '');
  const [error, setError] = useState('');

  return (
    <form className="space-y-3 border-t border-white/10 pt-4" onSubmit={(event) => {
      event.preventDefault();
      const count = parseSameNameCount(countInput, under);
      if (count === null) {
        setError(under ? '请输入大于 0 的整数作为人数上限。' : '请输入 0 或正整数。');
        return;
      }
      save({ name, source, count, under, recordedAt: new Date().toISOString() });
      setError('');
    }}>
      <div>
        <h3 className="text-sm font-medium">记录查询结果</h3>
        <p className="mt-1 text-xs text-muted-foreground">查完可手动记录，结果保存在当前浏览器。</p>
      </div>
      <div className="flex flex-wrap items-end gap-2">
        <label className="flex flex-col gap-1.5 text-sm">
          结果类型
          <select value={under ? 'under' : 'exact'} onChange={(event) => {
            setUnder(event.target.value === 'under');
            setError('');
          }} className={selectClass}>
            <option value="exact">具体人数</option>
            <option value="under">少于</option>
          </select>
        </label>
        <label className="flex min-w-[7rem] flex-1 flex-col gap-1.5 text-sm" htmlFor={id}>
          {under ? '人数上限' : '同名人数'}
          <Input
            id={id}
            inputMode="numeric"
            value={countInput}
            maxLength={20}
            aria-invalid={Boolean(error)}
            onChange={(event) => {
              setCountInput(event.target.value);
              setError('');
            }}
          />
        </label>
        <Button type="submit" disabled={!countInput.trim()}>保存</Button>
      </div>
      {error && <p role="alert" className="text-sm text-destructive">{error}</p>}
      {record && (
        <div className="flex flex-wrap items-center justify-between gap-x-3 gap-y-1">
          <p role="status" className="text-xs text-muted-foreground">
            已记录 · {sameNameSources[source].scope} · {formatSameNameCount(record)} · {dateFormat.format(new Date(record.recordedAt))}
          </p>
          <Button type="button" variant="ghost" size="sm" onClick={() => {
            remove(name, source);
            setCountInput('');
            setUnder(false);
            setError('');
          }}>删除记录</Button>
        </div>
      )}
    </form>
  );
}

function SameNameLookupForm({ fullName, initialSource }: { fullName?: string; initialSource: SameNameSource }) {
  const inputRef = useRef<HTMLInputElement>(null);
  const [name, setName] = useState(fullName ?? '');
  const [source, setSource] = useState(initialSource);
  const [copyStatus, setCopyStatus] = useState('');
  const value = name.trim();
  const service = sameNameSources[source];

  const copyName = async () => {
    try {
      await navigator.clipboard.writeText(value);
      setCopyStatus('已复制全名');
    } catch {
      inputRef.current?.focus();
      inputRef.current?.select();
      setCopyStatus('已选中姓名，可手动复制。');
    }
  };

  return (
    <div className="space-y-4">
      <div className="flex flex-wrap items-end gap-2">
        <label className="flex min-w-[10rem] flex-1 flex-col gap-1.5 text-sm">
          全名
          <Input
            ref={inputRef}
            value={name}
            onChange={(event) => {
              setName(event.target.value);
              setCopyStatus('');
            }}
            placeholder="输入姓氏和名字"
            maxLength={60}
            autoComplete="off"
          />
        </label>
        <Button type="button" variant="outline" disabled={!value} onClick={() => void copyName()}>
          <Copy className="h-4 w-4" aria-hidden="true" />
          复制全名
        </Button>
      </div>
      {copyStatus && <p role="status" className="text-xs text-muted-foreground">{copyStatus}</p>}
      <div className="space-y-2">
        <div className="flex flex-wrap items-end gap-2">
          <label className="flex min-w-[10rem] flex-1 flex-col gap-1.5 text-sm">
            查询范围
            <select value={source} onChange={(event) => setSource(event.target.value as SameNameSource)} className={selectClass}>
              {Object.entries(sameNameSources).map(([key, item]) => <option key={key} value={key}>{item.scope}</option>)}
            </select>
          </label>
          <Button asChild>
            <a href={service.url} target="_blank" rel="noopener noreferrer">
              打开官方查询
              <ExternalLink className="h-4 w-4" aria-hidden="true" />
            </a>
          </Button>
        </div>
        <p className="text-xs text-muted-foreground">{service.agency} · {service.instructions}</p>
      </div>
      {value && <SameNameRecordForm key={JSON.stringify([value, source])} name={value} source={source} />}
      <a
        href="https://www.ris.gov.tw/documents/data/5/2/112namestat.pdf"
        target="_blank"
        rel="noopener noreferrer"
        className="inline-flex items-center gap-1.5 text-xs text-muted-foreground underline-offset-4 hover:text-foreground hover:underline"
      >
        台湾地区姓名统计（历史榜单）
        <ExternalLink className="h-3 w-3" aria-hidden="true" />
      </a>
    </div>
  );
}

export function SameNameLookupButton({ fullName, compact }: {
  fullName?: string;
  compact?: boolean;
}) {
  const [open, setOpen] = useState(false);
  const record = useSameNameRecords((s) => fullName ? s.records.find((item) => item.name === fullName) : undefined);
  const summary = compact && record ? `${sameNameSources[record.source].shortScope} · ${formatSameNameCount(record)}` : '';

  return (
    <span className="inline-flex" onClick={(event) => event.stopPropagation()}>
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogTrigger asChild>
          <Button
            type="button"
            variant={compact ? 'ghost' : 'outline'}
            size="sm"
            aria-label={fullName ? ['查同名', fullName, summary].filter(Boolean).join(' ') : undefined}
            className={compact ? 'h-7 gap-1 px-1.5 text-xs text-muted-foreground' : undefined}
          >
            <Search className={compact ? 'h-3.5 w-3.5' : 'h-4 w-4'} aria-hidden="true" />
            查同名
            {summary && <span className="tabular-nums">{summary}</span>}
          </Button>
        </DialogTrigger>
        <DialogContent className="max-h-[90vh] w-[calc(100%-2rem)] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>查询同名人数</DialogTitle>
            <DialogDescription>复制全名，在所选地区的官方页面查询。</DialogDescription>
          </DialogHeader>
          {open && <SameNameLookupForm fullName={fullName} initialSource={record?.source ?? 'nationwide'} />}
        </DialogContent>
      </Dialog>
    </span>
  );
}
