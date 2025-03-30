const base = {
  url: `127.0.0.1:3000`,
  login: '/api/auth/login',
  register: '/api/auth/register',
  setName: '/api/auth/setName',
  showArticle: '/api/article/:article_id',
  publishArticleComment: '/api/article/comment/:article_id',
  deleteArticleComment: '/api/article/:article_id/comment/:comment_id',
  createArticle: '/api/article'
}

export default base

