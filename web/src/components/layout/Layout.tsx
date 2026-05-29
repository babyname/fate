import { Header } from './Header';
import { NameDetailModal } from '@/components/naming/NameDetailModal';
import { useAppStore } from '@/store/app';

interface LayoutProps {
  children: React.ReactNode;
}

export function Layout({ children }: LayoutProps) {
  const selectedName = useAppStore((s) => s.selectedName);
  const detailModalOpen = useAppStore((s) => s.detailModalOpen);
  const setDetailModalOpen = useAppStore((s) => s.setDetailModalOpen);
  const setSelectedName = useAppStore((s) => s.setSelectedName);

  return (
    <div className="min-h-screen bg-background text-foreground">
      <Header />
      <main className="mx-auto max-w-6xl px-4 py-6">
        {children}
      </main>
      <NameDetailModal
        name={selectedName}
        open={detailModalOpen}
        onOpenChange={(open) => {
          setDetailModalOpen(open);
          if (!open) setSelectedName(null);
        }}
      />
    </div>
  );
}
