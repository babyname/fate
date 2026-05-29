import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Download, FileText, FileCode } from 'lucide-react';
import { api } from '@/lib/api';
import { cn } from '@/lib/utils';

interface ReportDownloadProps {
  taskId: string;
  names?: string[];
}

type ReportFormat = 'text' | 'markdown';

const formats: { value: ReportFormat; label: string; icon: typeof FileText; ext: string }[] = [
  { value: 'text', label: '纯文本', icon: FileText, ext: '.txt' },
  { value: 'markdown', label: 'Markdown', icon: FileCode, ext: '.md' },
];

export function ReportDownload({ taskId, names }: ReportDownloadProps) {
  const [selectedFormat, setSelectedFormat] = useState<ReportFormat>('markdown');
  const [downloading, setDownloading] = useState(false);

  const handleDownload = async () => {
    setDownloading(true);
    try {
      const res = await api.downloadReport({
        task_id: taskId,
        format: selectedFormat,
        names,
      });

      const blob = new Blob([res.content], {
        type: selectedFormat === 'markdown' ? 'text/markdown' : 'text/plain',
      });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = res.filename || `fate-report-${taskId}${formats.find((f) => f.value === selectedFormat)?.ext}`;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      URL.revokeObjectURL(url);
    } catch (err) {
      console.error('Download failed:', err);
    } finally {
      setDownloading(false);
    }
  };

  return (
    <div className="glass-card p-4 space-y-3">
      <h3 className="text-sm font-medium text-foreground">下载报告</h3>
      <div className="flex gap-2">
        {formats.map((fmt) => {
          const Icon = fmt.icon;
          return (
            <button
              key={fmt.value}
              onClick={() => setSelectedFormat(fmt.value)}
              className={cn(
                'flex items-center gap-2 rounded-md px-3 py-2 text-sm transition-all duration-200 border',
                selectedFormat === fmt.value
                  ? 'border-blue-500/50 bg-blue-500/10 text-blue-300'
                  : 'border-border bg-transparent text-muted-foreground hover:bg-white/5',
              )}
            >
              <Icon className="h-4 w-4" />
              {fmt.label}
            </button>
          );
        })}
      </div>
      <Button
        variant="celestial"
        size="sm"
        className="w-full"
        onClick={handleDownload}
        disabled={downloading}
      >
        <Download className="h-4 w-4" />
        {downloading ? '生成中...' : '下载报告'}
      </Button>
    </div>
  );
}
