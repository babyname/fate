import type {
  GenerateRequest,
  GenerateResponse,
  TaskStatusResponse,
  PoetrySearchRequest,
  PoetrySearchResponse,
  ReportRequest,
  ReportResponse,
} from '@/types/api';

const BASE = '/api';

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE}${path}`, {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({ message: res.statusText }));
    throw new Error(err.message || `Request failed: ${res.status}`);
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
    return request(`/tasks/${taskId}`);
  },

  searchPoetry(data: PoetrySearchRequest): Promise<PoetrySearchResponse> {
    const params = new URLSearchParams();
    params.set('keyword', data.keyword);
    if (data.dynasty) params.set('dynasty', data.dynasty);
    if (data.author) params.set('author', data.author);
    if (data.page) params.set('page', String(data.page));
    if (data.page_size) params.set('page_size', String(data.page_size));
    return request(`/poetry/search?${params.toString()}`);
  },

  downloadReport(data: ReportRequest): Promise<ReportResponse> {
    return request('/report', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  },
};
