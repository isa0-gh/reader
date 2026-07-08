import { createContext, useContext, useEffect, useState, useCallback, ReactNode } from "react";
import { useAuth } from "./AuthContext";
import { api, Favorite, Series, Chapter } from "./api";

interface FavoritesCtx {
  favorites: Favorite[];
  isFavorite(seriesId: number): boolean;
  toggleFavorite(seriesId: number, series?: Series): Promise<void>;
  pingProgress(seriesId: number, chapter: Chapter): void;
}

const Ctx = createContext<FavoritesCtx>(null!);

export function FavoritesProvider({ children }: { children: ReactNode }) {
  const { user } = useAuth();
  const [favorites, setFavorites] = useState<Favorite[]>([]);

  const load = useCallback(() => {
    if (!user) { setFavorites([]); return; }
    api.listFavorites({ limit: 200 }).then((res) => setFavorites(res.items)).catch(() => {});
  }, [user]);

  useEffect(() => { load(); }, [load]);

  function isFavorite(seriesId: number) {
    return favorites.some((f) => f.series_id === seriesId);
  }

  async function toggleFavorite(seriesId: number, series?: Series) {
    if (isFavorite(seriesId)) {
      await api.removeFavorite(seriesId);
      setFavorites((fs) => fs.filter((f) => f.series_id !== seriesId));
    } else {
      const fav = await api.addFavorite(seriesId);
      setFavorites((fs) => [{ ...fav, series: series ?? fav.series }, ...fs]);
    }
  }

  // Best-effort: the backend silently no-ops if the series isn't favorited,
  // but skip the call entirely client-side to avoid pinging on every single
  // chapter view for non-favorited series.
  function pingProgress(seriesId: number, chapter: Chapter) {
    if (!isFavorite(seriesId)) return;
    api.updateFavoriteProgress(seriesId, chapter.id).catch(() => {});
    setFavorites((fs) => fs.map((f) =>
      f.series_id === seriesId ? { ...f, last_read_chapter_id: chapter.id, last_read_chapter: chapter } : f
    ));
  }

  return (
    <Ctx.Provider value={{ favorites, isFavorite, toggleFavorite, pingProgress }}>
      {children}
    </Ctx.Provider>
  );
}

export const useFavorites = () => useContext(Ctx);
