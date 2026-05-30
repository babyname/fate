import { useNavigate, useLocation } from 'react-router-dom';
import { Sparkles } from 'lucide-react';
import { cn } from '@/lib/utils';

export function Header() {
  const navigate = useNavigate();
  const location = useLocation();

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
        </nav>
      </div>
    </header>
  );
}