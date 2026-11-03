export interface Name {
  id: number;
  name: string;
  created_at: string;
}

export interface CreateNameResponse {
  message: string;
  name: string;
}

export interface SearchNameResponse {
  name: string;
  exists: boolean;
}

export interface ApiError {
  error: string;
}
