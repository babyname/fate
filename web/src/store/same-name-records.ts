import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { SameNameRecord, SameNameSource } from '@/lib/same-name';

interface SameNameRecords {
  records: SameNameRecord[];
  save: (record: SameNameRecord) => void;
  remove: (name: string, source: SameNameSource) => void;
}

export const useSameNameRecords = create<SameNameRecords>()(persist(
  (set) => ({
    records: [],
    save: (record) => set(({ records }) => ({
      records: [
        record,
        ...records.filter((previous) => previous.name !== record.name || previous.source !== record.source),
      ],
    })),
    remove: (name, source) => set(({ records }) => ({
      records: records.filter((record) => record.name !== name || record.source !== source),
    })),
  }),
  {
    name: 'fate-same-name-records',
    partialize: ({ records }) => ({ records }),
  },
));
