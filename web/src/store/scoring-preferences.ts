import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { DEFAULT_TRADITIONAL_WEIGHT } from '@/lib/recommendation-scores';

interface ScoringPreferences {
  traditionalWeight: number;
  setTraditionalWeight: (weight: number) => void;
}

function validWeight(weight: unknown): weight is number {
  return typeof weight === 'number' && Number.isInteger(weight) && weight >= 0 && weight <= 100;
}

export const useScoringPreferences = create<ScoringPreferences>()(persist(
  (set) => ({
    traditionalWeight: DEFAULT_TRADITIONAL_WEIGHT,
    setTraditionalWeight: (weight) => {
      if (validWeight(weight)) set({ traditionalWeight: weight });
    },
  }),
  {
    name: 'fate-scoring-preferences',
    partialize: ({ traditionalWeight }) => ({ traditionalWeight }),
    merge: (saved, current) => {
      const weight = (saved as Partial<ScoringPreferences> | undefined)?.traditionalWeight;
      return validWeight(weight) ? { ...current, traditionalWeight: weight } : current;
    },
  },
));
