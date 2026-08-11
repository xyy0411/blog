export interface Article {
  id?: number;
  user_id: number;
  title: string;
  abstract: string;
  content: string;
  cover: string;
  look_count: number;
  likes: number;
  comment_count: number;
  collect_count: number;
  open_comment: boolean;
  created_at?: string;
  updated_at?: string;
}

export interface Comment {
  id?: number;
  user_id: number;
  article_id: number;
  content: string;
  ip?: string;
  created_at?: string;
  updated_at?: string;
}

export interface CommentWithArticle {
  article: Article;
  comments: Comment[];
}

export interface CommitArticle {
  title: string;
  abstract: string;
  content: string;
  cover: string;
  open_comment: boolean;
}
