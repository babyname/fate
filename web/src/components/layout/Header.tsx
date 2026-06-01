import { useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { Sparkles, Star, BookOpen } from 'lucide-react';
import { cn } from '@/lib/utils';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';
import { NameScoreTab } from '@/components/naming/NameScoreTab';
import { PoetrySearchTab } from '@/components/naming/PoetrySearchTab';

export function Header() {
  const navigate = useNavigate();
  const location = useLocation();
  const [scoreOpen, setScoreOpen] = useState(false);
  const [poetryOpen, setPoetryOpen] = useState(false);

  const isHome = location.pathname === '/';

  return (
    <header className="sticky top-0 z-40 border-b border-white/10 bg-background/80 backdrop-blur-xl">
      <div className="mx-auto flex h-16 max-w-6xl items-center justify-between px-4">
        <div className="flex items-center gap-3">
          <button
            onClick={() => navigate('/')}
            className="text-2xl font-bold font-serif tracking-wider text-gradient-celestial hover:opacity-80 transition-opacity"
          >
            FATE
          </button>
          <span className="hidden sm:inline-block text-xs text-muted-foreground font-mono tracking-widest uppercase">
            Celestial Naming Engine
          </span>
        </div>
        <nav className="flex items-center gap-1">
          <button
            onClick={() => navigate('/')}
            className={cn(
              'flex items-center gap-1.5 rounded-md px-3 py-1.5 text-sm font-medium transition-all duration-200',
              isHome
                ? 'bg-white/10 text-foreground border border-white/10'
                : 'text-muted-foreground hover:text-foreground hover:bg-white/5',
            )}
          >
            <Sparkles className="h-4 w-4" />
            命名
          </button>

          <Dialog open={scoreOpen} onOpenChange={setScoreOpen}>
            <DialogTrigger asChild>
              <button className="flex items-center gap-1.5 rounded-md px-3 py-1.5 text-sm font-medium text-muted-foreground hover:text-foreground hover:bg-white/5 transition-all duration-200">
                <Star className="h-4 w-4" />
                姓名测分
              </button>
            </DialogTrigger>
            <DialogContent className="max-w-2xl max-h-[85vh] overflow-y-auto">
              <DialogHeader>
                <DialogTitle className="font-serif">姓名测分</DialogTitle>
              </DialogHeader>
              <NameScoreTab />
            </DialogContent>
          </Dialog>

          <Dialog open={poetryOpen} onOpenChange={setPoetryOpen}>
            <DialogTrigger asChild>
              <button className="flex items-center gap-1.5 rounded-md px-3 py-1.5 text-sm font-medium text-muted-foreground hover:text-foreground hover:bg-white/5 transition-all duration-200">
                <BookOpen className="h-4 w-4" />
                诗词搜索
              </button>
            </DialogTrigger>
            <DialogContent className="max-w-2xl max-h-[85vh] overflow-y-auto">
              <DialogHeader>
                <DialogTitle className="font-serif">诗词搜索</DialogTitle>
              </DialogHeader>
              <PoetrySearchTab />
            </DialogContent>
          </Dialog>
        </nav>
      </div>
    </header>
  );
}