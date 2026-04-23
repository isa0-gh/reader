const BASE = "/api/v1";

function authHeader(): Record<string, string> {
  const token = localStorage.getItem("token");
  return token ? { Authorization: `Bearer ${token}` } : {};
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(BASE + path, {
    headers: { "Content-Type": "application/json", ...authHeader() },
    ...init,
  });
  if (!res.ok) throw new Error(await res.text());
  return res.json();
}

export const api = {
  register: (email: string, password: string, name: string) =>
    request("/auth/register", { method: "POST", body: JSON.stringify({ email, password, name }) }),

  login: (email: string, password: string) =>
    request<{ token: string; user: { id: number; name: string; role: string } }>(
      "/auth/login",
      { method: "POST", body: JSON.stringify({ email, password }) }
    ),

  getSeries: (id: number) => request<Series>(`/series/${id}`),
  getChapter: (id: number) => request<Chapter>(`/chapters/${id}`),
};

export interface Series {
  id: number;
  title: string;
  slug: string;
  description: string;
  cover_image: string;
  author: string;
  artist: string;
  status: string;
  chapters: Chapter[];
}

export interface Chapter {
  id: number;
  series_id: number;
  number: number;
  title: string;
  pages: Page[];
}

export interface Page {
  id: number;
  key: string;
  page_number: number;
}
