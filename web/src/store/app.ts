import { create } from 'zustand';
import type { NameResult, ExcellentEntry } from '@/types/api';

interface AppState {
  surname: string;
  born: string;
  sex: string;
  avoidChars: string[];
  requireChars: string[];
  poetryMode: number;
  count: number;
  taskId: string | null;
  topNames: NameResult[];
  top10: ExcellentEntry[];
  total: number;
  isGenerating: boolean;
  selectedName: NameResult | null;
  detailModalOpen: boolean;

  setSurname: (v: string) => void;
  setBorn: (v: string) => void;
  setSex: (v: string) => void;
  setAvoidChars: (v: string[]) => void;
  setRequireChars: (v: string[]) => void;
  setPoetryMode: (v: number) => void;
  setCount: (v: number) => void;
  setTaskId: (v: string | null) => void;
  setTopNames: (v: NameResult[]) => void;
  setTop10: (v: ExcellentEntry[]) => void;
  setTotal: (v: number) => void;
  setIsGenerating: (v: boolean) => void;
  setSelectedName: (v: NameResult | null) => void;
  setDetailModalOpen: (v: boolean) => void;
  reset: () => void;
}

const initialState = {
  surname: '',
  born: '',
  sex: 'boy' as const,
  avoidChars: [] as string[],
  requireChars: [] as string[],
  poetryMode: 0,
  count: 20,
  taskId: null as string | null,
  topNames: [] as NameResult[],
  top10: [] as ExcellentEntry[],
  total: 0,
  isGenerating: false,
  selectedName: null as NameResult | null,
  detailModalOpen: false,
};

export const useAppStore = create<AppState>((set) => ({
  ...initialState,
  setSurname: (v) => set({ surname: v }),
  setBorn: (v) => set({ born: v }),
  setSex: (v) => set({ sex: v }),
  setAvoidChars: (v) => set({ avoidChars: v }),
  setRequireChars: (v) => set({ requireChars: v }),
  setPoetryMode: (v) => set({ poetryMode: v }),
  setCount: (v) => set({ count: v }),
  setTaskId: (v) => set({ taskId: v }),
  setTopNames: (v) => set({ topNames: v }),
  setTop10: (v) => set({ top10: v }),
  setTotal: (v) => set({ total: v }),
  setIsGenerating: (v) => set({ isGenerating: v }),
  setSelectedName: (v) => set({ selectedName: v }),
  setDetailModalOpen: (v) => set({ detailModalOpen: v }),
  reset: () => set(initialState),
}));