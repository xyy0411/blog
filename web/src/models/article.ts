export interface Article {
  userId?: number;        // 关联用户
  title: string;         // 标题
  abstract: string;      // 简介
  content: string;       // 内容
  cover?: string;         // 封面
  lookCount?: number;     // 浏览量
  likes?: number;         // 点赞数
  commentCount?: number;  // 评论数
  collectCount?: number;  // 收藏数
  openComment?: boolean;   // 文章评论开关
  status: number;        // 状态 草稿 审核中 已发布
}
