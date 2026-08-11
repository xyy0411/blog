// 文章模型（对齐后端 models.Article，字段名按 JSON tag 用 snake_case）
export interface Article {
  id?: number;
  user_id: number;          // 关联用户
  title: string;            // 标题
  abstract: string;         // 简介（当前后端 CreateArticle 不自动生成，留作前端列表兜底截取 content）
  content: string;          // 内容（markdown 字符串，由后端 Contents 切片拼接而来）
  cover: string;            // 封面
  look_count: number;       // 浏览量
  likes: number;            // 点赞数
  comment_count: number;    // 评论数
  collect_count: number;    // 收藏数
  open_comment: boolean;    // 评论开关
  created_at?: string;
  updated_at?: string;
}

// 评论模型（对齐后端 models.Comment）
export interface Comment {
  id?: number;
  user_id: number;
  article_id: number;
  content: string;
  ip?: string;
  created_at?: string;
  updated_at?: string;
}

// 文章 + 评论组合（详情页用，对齐后端 CommentWithArticle）
export interface CommentWithArticle {
  article: Article;
  comments: Comment[];
}

// 文章内容切片（编辑器多段输入，对齐后端 ArticleContent）
// 后端 ArticleContent.String() 会用 "\n" 把切片拼接成单个字符串存到 Article.content
export interface ArticleContent {
  content: string[];
}

// 提交文章时用的结构（对齐后端 CommitArticle）
export interface CommitArticle {
  title: string;
  contents: ArticleContent;
  cover: string;
  open_comment: boolean;
}
