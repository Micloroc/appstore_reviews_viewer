import { Review } from '../types';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

const apiRequest = async <T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> => {
  const url = `${API_BASE_URL}${endpoint}`;
  
  const defaultOptions: RequestInit = {
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
    ...options,
  };

  const response = await fetch(url, defaultOptions);
  
  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }
  
  const contentType = response.headers.get('content-type');
  const contentLength = response.headers.get('content-length');
  
  if (!contentType || contentLength === '0' || response.status === 204) {
    return {} as T;
  }
  
  try {
    const text = await response.text();
    if (!text || text.trim() === '') {
      return {} as T;
    }
    return JSON.parse(text);
  } catch (error) {
    console.warn('Failed to parse JSON response:', error);
    return {} as T;
  }
};

export const reviewsApi = {
  getByAppId: async (appId: string): Promise<{ reviews: Review[] }> => {
    return apiRequest<{ reviews: Review[] }>(`/api/v1/app/${appId}/reviews/recent`);
  },
};

export const appsApi = {
  add: async (appId: string): Promise<void> => {
    return apiRequest<void>('/api/v1/app', {
      method: 'POST',
      body: JSON.stringify({ appId }),
    });
  },
};
