export interface Review {
  id: string;
  appId: string;
  author: string;
  content: string;
  score: number;
  submittedAt: string;
}

export interface App {
  id: string;
  name?: string;
}
