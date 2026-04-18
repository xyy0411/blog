const trimTrailingSlash = (value?: string) => value?.replace(/\/$/, '') ?? '';

export const apiBaseUrl = trimTrailingSlash(import.meta.env.VITE_API_BASE_URL) || '';

export const apiUrl = (path: string) => {
  if (!path.startsWith('/')) {
    path = `/${path}`;
  }

  return `${apiBaseUrl}${path}`;
};

const base = {
  url: apiBaseUrl,
  login: '/api/auth/login',
  register: '/api/auth/register',
  setName: '/api/auth/setName',
  showArticle: '/api/article/:article_id',
  publishArticleComment: '/api/article/comment/:article_id',
  deleteArticleComment: '/api/article/:article_id/comment/:comment_id',
  createArticle: '/api/article',
  articleList: '/api/article/all',
  matchingToday: '/api/matching/record/today',
  matchingWeek: '/api/matching/record/week',
  matchingAll: '/api/matching/record/all',
} as const;

export default base;
