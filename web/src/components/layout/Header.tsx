import { useNavigate, useLocation } from 'react-router-dom';
import { Sparkles, BookOpen } from 'lucide-react';
import { cn } from '@/lib/utils';

const navItems = [
  { label: 'Generator', path: '/', icon: Sparkles },
  { label: 'Poetry', path: '/?tab=poetry', icon: BookOpen },
];

export function Header() {
  const navigate = useNavigate();
  const location = useLocation();

  return (
    <header className="sticky top-0 z-40 border-b border-white/10 bg-background/80 backdrop-blur-xl">
      <div className="mx-auto flex h-16 max-w-6xl items-center justify-between px-4">
        <div className="flex items-center gap-3">
          <h1 className="text-2xl font-bold font-serif tracking-wider text-gradient-celestial">
            FATE
          </h1>
          <span className="hidden sm:inline-block text-xs text-muted-foreground font-mono tracking-widest uppercase">
            Celestial Naming Engine
          </span>
        </div>
        <nav className="flex items-center gap-1">
          {navItems.map((item) => {
            const Icon = item.icon;
            const isActive =
              item.path === '/'
                ? location.pathname === '/' && !new URLSearchParams(location.search).get('tab')
                : location.search.includes('tab=poetry');
            return (
              <button
                key={item.label}
                onClick={() => navigate(item.path)}
                className={cn(
                  'flex items-center gap-1.5 rounded-md px-3 py-1.5 text-sm font-medium transition-all duration-200',
                  isActive
                    ? 'bg-white/10 text-foreground border border-white/10'
                    : 'text-muted-foreground hover:text-foreground hover:bg-white/5',
                )}
              >
                <Icon className="h-4 w-4" />
                {item.label}
              </button>
            );
          })}
        </nav>
      </div>
    </header>
  );
}
