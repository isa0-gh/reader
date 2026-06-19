import { useState, useEffect, useCallback } from 'react';
import { api } from './api';

export interface Favorite {
  id: number;
  user_id: number;
  series_id: number;
  created_at: string;
  series?: {
    id: number;
    title: string;
    slug: string;
    cover_image?: {
      key: string;
      bucket: string;
    };
  };
}

export function useFavorites() {
  const [favorites, setFavorites] = useState<Favorite[]>([]);
  const [loading, setLoading] = useState(false);

  const fetchFavorites = useCallback(async () => {
    setLoading(true);
    try {
      const res = await api.get('/favorites');
      setFavorites(res.data || []);
    } catch (err) {
      console.error('Failed to fetch favorites:', err);
      setFavorites([]);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchFavorites();
  }, [fetchFavorites]);

  const addFavorite = async (seriesID: number) => {
    await api.post(`/favorites/${seriesID}`);
    await fetchFavorites();
  };

  const removeFavorite = async (seriesID: number) => {
    await api.delete(`/favorites/${seriesID}`);
    await fetchFavorites();
  };

  const toggleFavorite = async (seriesID: number, isFavorite: boolean) => {
    if (isFavorite) {
      await removeFavorite(seriesID);
    } else {
      await addFavorite(seriesID);
    }
  };

  return {
    favorites,
    loading,
    addFavorite,
    removeFavorite,
    toggleFavorite,
    refetch: fetchFavorites,
  };
}

export function useFavoriteCheck(seriesID: number | null) {
  const [isFavorite, setIsFavorite] = useState(false);
  const [totalFavorites, setTotalFavorites] = useState(0);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (!seriesID) return;

    setLoading(true);
    api
      .get(`/favorites/${seriesID}/check`)
      .then((res) => {
        setIsFavorite(res.data.is_favorite);
        setTotalFavorites(res.data.total_favorites || 0);
      })
      .catch(() => {
        setIsFavorite(false);
        setTotalFavorites(0);
      })
      .finally(() => setLoading(false));
  }, [seriesID]);

  const toggle = async () => {
    if (!seriesID) return;

    try {
      if (isFavorite) {
        await api.delete(`/favorites/${seriesID}`);
        setIsFavorite(false);
        setTotalFavorites((prev) => Math.max(0, prev - 1));
      } else {
        await api.post(`/favorites/${seriesID}`);
        setIsFavorite(true);
        setTotalFavorites((prev) => prev + 1);
      }
    } catch (err) {
      console.error('Toggle favorite failed:', err);
    }
  };

  return { isFavorite, totalFavorites, loading, toggle };
}
