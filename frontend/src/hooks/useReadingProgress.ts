import { useState, useEffect } from 'react';
import { api } from './api';

export interface ReadingProgress {
  id: number;
  user_id: number;
  chapter_id: number;
  page_index: number;
  completed: boolean;
  created_at: string;
  updated_at: string;
  chapter?: {
    id: number;
    title: string;
    number: number;
    series_id: number;
  };
}

export function useReadingProgress(chapterID: number | null) {
  const [progress, setProgress] = useState<ReadingProgress | null>(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (!chapterID) return;
    
    setLoading(true);
    api.get(`/progress/chapters/${chapterID}`)
      .then(res => setProgress(res.data))
      .catch(() => setProgress(null))
      .finally(() => setLoading(false));
  }, [chapterID]);

  const updateProgress = async (pageIndex: number, completed: boolean) => {
    if (!chapterID) return;
    
    const res = await api.put(`/progress/chapters/${chapterID}`, {
      page_index: pageIndex,
      completed
    });
    setProgress(res.data);
  };

  return { progress, loading, updateProgress };
}

export function useMyProgress() {
  const [progressList, setProgressList] = useState<ReadingProgress[]>([]);
  const [loading, setLoading] = useState(false);

  const fetch = () => {
    setLoading(true);
    api.get('/progress/me')
      .then(res => setProgressList(res.data || []))
      .catch(() => setProgressList([]))
      .finally(() => setLoading(false));
  };

  useEffect(() => {
    fetch();
  }, []);

  return { progressList, loading, refetch: fetch };
}

export function useSeriesProgress(seriesID: number | null) {
  const [progressList, setProgressList] = useState<ReadingProgress[]>([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (!seriesID) return;
    
    setLoading(true);
    api.get(`/progress/series/${seriesID}`)
      .then(res => setProgressList(res.data || []))
      .catch(() => setProgressList([]))
      .finally(() => setLoading(false));
  }, [seriesID]);

  return { progressList, loading };
}
