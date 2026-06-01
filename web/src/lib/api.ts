import type {
  GenerateRequest,
  GenerateResponse,
  TaskStatusResponse,
  TaskResultResponse,
  ExploreResponse,
  NameDetailResponse,
  NamesResponse,
  NameScoreRequest,
  PoetrySearchResponse,
} from '@/types/api';

const BASE = '/api';

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE}${path}`, {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error(err.error || `Request failed: ${res.status}`);
  }
  return res.json();
}

export const api = {
  generate(data: GenerateRequest): Promise<GenerateResponse> {
    return request('/generate', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  },

  getTaskStatus(taskId: string): Promise<TaskStatusResponse> {
    return request(`/generate/status?task_id=${taskId}`);
  },

  getTaskResult(taskId: string): Promise<TaskResultResponse> {
    return request(`/generate/result?task_id=${taskId}`);
  },

  getNames(taskId: string, page: number = 1, size: number = 20): Promise<NamesResponse> {
    const params = new URLSearchParams();
    params.set('task_id', taskId);
    params.set('page', String(page));
    params.set('size', String(size));
    return request(`/generate/names?${params.toString()}`);
  },

  explore(taskId: string, count?: number, hasPoetry?: boolean, wuxing?: string): Promise<ExploreResponse> {
    const params = new URLSearchParams();
    params.set('task_id', taskId);
    if (count) params.set('count', String(count));
    if (hasPoetry) params.set('has_poetry', 'true');
    if (wuxing) params.set('wuxing', wuxing);
    return request(`/generate/explore?${params.toString()}`);
  },

  getNameDetail(taskId: string, char1: string, char2: string): Promise<NameDetailResponse> {
    return request(`/name-detail?task_id=${taskId}&char1=${char1}&char2=${char2}`);
  },

  nameScore(data: NameScoreRequest): Promise<NameDetailResponse> {
    return request('/name-score', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  },

  searchPoetry(keyword: string): Promise<PoetrySearchResponse> {
    const params = new URLSearchParams();
    params.set('keyword', keyword);
    return request(`/poetry/search?${params.toString()}`);
  },
};