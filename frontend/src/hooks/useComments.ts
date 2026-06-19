import { useState, useEffect, useCallback } from 'react';
import { api } from './api';

export interface Comment {
  id: number;
  chapter_id: number;
  user_id: number;
  parent_id?: number;
  content: string;
  edited: boolean;
  deleted: boolean;
  like_count: number;
  reply_count: number;
  created_at: string;
  updated_at: string;
  user?: {
    id: number;
    username: string;
    email: string;
  };
}

export interface CommentListResponse {
  comments: Comment[];
  total: number;
  page: number;
  page_size: number;
}

export function useChapterComments(chapterID: number | null) {
  const [comments, setComments] = useState<Comment[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(false);

  const fetchComments = useCallback(
    async (pageNum: number = 1) => {
      if (!chapterID) return;

      setLoading(true);
      try {
        const res = await api.get<CommentListResponse>(
          `/chapters/${chapterID}/comments?page=${pageNum}&page_size=20`
        );
        setComments(res.data.comments || []);
        setTotal(res.data.total);
        setPage(pageNum);
      } catch (err) {
        console.error('Failed to fetch comments:', err);
        setComments([]);
      } finally {
        setLoading(false);
      }
    },
    [chapterID]
  );

  useEffect(() => {
    fetchComments(1);
  }, [fetchComments]);

  const createComment = async (content: string, parentID?: number) => {
    if (!chapterID) return;

    await api.post(`/chapters/${chapterID}/comments`, {
      content,
      parent_id: parentID,
    });
    await fetchComments(page);
  };

  const updateComment = async (commentID: number, content: string) => {
    await api.patch(`/comments/${commentID}`, { content });
    await fetchComments(page);
  };

  const deleteComment = async (commentID: number) => {
    await api.delete(`/comments/${commentID}`);
    await fetchComments(page);
  };

  const likeComment = async (commentID: number) => {
    await api.post(`/comments/${commentID}/like`);
    // Optimistically update local state
    setComments((prev) =>
      prev.map((c) =>
        c.id === commentID ? { ...c, like_count: c.like_count + 1 } : c
      )
    );
  };

  const unlikeComment = async (commentID: number) => {
    await api.delete(`/comments/${commentID}/like`);
    // Optimistically update local state
    setComments((prev) =>
      prev.map((c) =>
        c.id === commentID
          ? { ...c, like_count: Math.max(0, c.like_count - 1) }
          : c
      )
    );
  };

  return {
    comments,
    total,
    page,
    loading,
    fetchComments,
    createComment,
    updateComment,
    deleteComment,
    likeComment,
    unlikeComment,
  };
}

export function useCommentReplies(commentID: number | null) {
  const [replies, setReplies] = useState<Comment[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);

  const fetchReplies = useCallback(async () => {
    if (!commentID) return;

    setLoading(true);
    try {
      const res = await api.get<CommentListResponse>(
        `/comments/${commentID}/replies?page=1&page_size=50`
      );
      setReplies(res.data.comments || []);
      setTotal(res.data.total);
    } catch (err) {
      console.error('Failed to fetch replies:', err);
      setReplies([]);
    } finally {
      setLoading(false);
    }
  }, [commentID]);

  useEffect(() => {
    fetchReplies();
  }, [fetchReplies]);

  return {
    replies,
    total,
    loading,
    refetch: fetchReplies,
  };
}
