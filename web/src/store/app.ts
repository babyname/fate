import { create } from 'zustand';
import type { NameResult, FixedChar } from '@/types/api';

interface AppState {
  surname: string;
  gender: 'male' | 'female' | 'neutral';
  fixedChars: FixedChar[];
  blacklist: string[];
  count: number;
  taskId: string | null;
  results: NameResult[];
  isGenerating: boolean;
  selectedName: NameResult | null;
  detailModalOpen: boolean;

  setSurname: (v: string) => void;
  setGender: (v: 'male' | 'female' | 'neutral') => void;
  setFixedChars: (v: FixedChar[]) => void;
  setBlacklist: (v: string[]) => void;
  setCount: (v: number) => void;
  setTaskId: (v: string | null) => void;
  setResults: (v: NameResult[]) => void;
  setIsGenerating: (v: boolean) => void;
  setSelectedName: (v: NameResult | null) => void;
  setDetailModalOpen: (v: boolean) => void;
  reset: () => void;
}

const initialState = {
  surname: '',
  gender: 'neutral' as const,
  fixedChars: [],
  blacklist: [],
  count: 20,
  taskId: null,
  results: [],
  isGenerating: false,
  selectedName: null,
  detailModalOpen: false,
};

export const useAppStore = create<AppState>((set) => ({
  ...initialState,
  setSurname: (v) => set({ surname: v }),
  setGender: (v) => set({ gender: v }),
  setFixedChars: (v) => set({ fixedChars: v }),
  setBlacklist: (v) => set({ blacklist: v }),
  setCount: (v) => set({ count: v }),
  setTaskId: (v) => set({ taskId: v }),
  setResults: (v) => set({ results: v }),
  setIsGenerating: (v) => set({ isGenerating: v }),
  setSelectedName: (v) => set({ selectedName: v }),
  setDetailModalOpen: (v) => set({ detailModalOpen: v }),
  reset: () => set(initialState),
}));
