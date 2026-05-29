export interface GenerateRequest {
  surname: string;
  gender: 'male' | 'female' | 'neutral';
  fixed_chars: FixedChar[];
  blacklist: string[];
  count: number;
  options?: GenerateOptions;
}

export interface FixedChar {
  position: number;
  char: string;
}

export interface GenerateOptions {
  wuxing_boost?: boolean;
  poetry_pref?: boolean;
  avoid_surname_homophone?: boolean;
}

export interface GenerateResponse {
  task_id: string;
  status: 'pending' | 'processing' | 'completed' | 'failed';
}

export interface TaskStatusResponse {
  task_id: string;
  status: 'pending' | 'processing' | 'completed' | 'failed';
  progress?: number;
  result?: NameResult[];
  error?: string;
}

export interface NameResult {
  name: string;
  surname: string;
  given_name: string;
  score: number;
  wuxing: WuxingAnalysis;
  poetry: PoetryReference[];
  meaning: string;
  pronunciation: PronunciationInfo;
}

export interface WuxingAnalysis {
  elements: string[];
  balance: number;
  missing: string[];
  description: string;
}

export interface PoetryReference {
  title: string;
  author: string;
  dynasty: string;
  content: string;
  matched_char: string;
}

export interface PronunciationInfo {
  pinyin: string;
  tone_pattern: string;
  rhythm: string;
}

export interface PoetrySearchRequest {
  keyword: string;
  dynasty?: string;
  author?: string;
  page?: number;
  page_size?: number;
}

export interface PoetrySearchResponse {
  total: number;
  page: number;
  page_size: number;
  poems: PoetryItem[];
}

export interface PoetryItem {
  title: string;
  author: string;
  dynasty: string;
  content: string;
  matched_chars?: string[];
}

export interface ReportRequest {
  task_id: string;
  format: 'text' | 'markdown';
  names?: string[];
}

export interface ReportResponse {
  content: string;
  filename: string;
  format: 'text' | 'markdown';
}
