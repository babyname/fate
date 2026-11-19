import type { NameResult, RecommendationRanking, RecommendationScore, RecommendationSort } from '@/types/api';

export const DEFAULT_TRADITIONAL_WEIGHT = 40;

export function recommendationScore(
  ranking: RecommendationRanking,
  field: RecommendationScore,
  traditionalWeight: number,
): number {
  if (field !== 'total' || traditionalWeight === DEFAULT_TRADITIONAL_WEIGHT) return ranking[field];

  // Weight integer tenths to round half points consistently.
  const traditional = Math.round(ranking.traditional * 10);
  const soundAndMeaning = Math.round(ranking.sound_and_meaning * 10);
  return Math.round((traditional * traditionalWeight + soundAndMeaning * (100 - traditionalWeight)) / 100) / 10;
}

export function sortRecommendations(
  names: NameResult[],
  { field, descending }: RecommendationSort,
  traditionalWeight: number,
): NameResult[] {
  const direction = descending ? -1 : 1;
  return names
    .map((name) => ({
      name,
      score: name.ranking ? recommendationScore(name.ranking, field, traditionalWeight) : 0,
    }))
    .sort((a, b) => direction * (a.score - b.score) || a.name.rank - b.name.rank)
    .map(({ name }) => name);
}
