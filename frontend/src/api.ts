const BASE = "/api/v1";

function authHeader(): Record<string, string> {
  const token = localStorage.getItem("token");
  return token ? { Authorization: `Bearer ${token}` } : {};
}

async function request<T = void>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(BASE + path, {
    headers: { "Content-Type": "application/json", ...authHeader() },
    ...init,
  });
  const text = await res.text();
  if (!res.ok) {
    let message = text;
    try {
      const parsed = JSON.parse(text) as { error?: string };
      if (parsed.error) message = parsed.error;
    } catch {
      // error body wasn't JSON, fall back to raw text
    }
    throw new Error(message);
  }
  return text ? JSON.parse(text) : (undefined as T);
}

export const api = {
  register: (email: string, password: string, name: string) =>
    request("/auth/register", { method: "POST", body: JSON.stringify({ email, password, name }) }),

  login: (email: string, password: string) =>
    request<{ token: string; user: User }>(
      "/auth/login",
      { method: "POST", body: JSON.stringify({ email, password }) }
    ),

  changePassword: (currentPassword: string, newPassword: string) =>
    request("/users/me/password", {
      method: "PATCH",
      body: JSON.stringify({ current_password: currentPassword, new_password: newPassword }),
    }),

  changeEmail: (currentPassword: string, newEmail: string) =>
    request<User>("/users/me/email", {
      method: "PATCH",
      body: JSON.stringify({ current_password: currentPassword, new_email: newEmail }),
    }),

  listOrphanedObjects: () =>
    request<S3Object[]>("/admin/s3/orphaned"),

  purgeOrphanedObjects: () =>
    request<{ deleted: string[]; failed: string[] }>("/admin/s3/orphaned", { method: "DELETE" }),

  listUsers: (params?: { limit?: number; after?: number; before?: number }) => {
    const q = new URLSearchParams();
    if (params?.limit) q.set("limit", String(params.limit));
    if (params?.after) q.set("after", String(params.after));
    if (params?.before) q.set("before", String(params.before));
    return request<User[]>(`/users?${q}`);
  },

  updateUserRole: (id: number, role: string) =>
    request(`/users/${id}/role`, { method: "PATCH", body: JSON.stringify({ role }) }),

  deleteUser: (id: number) =>
    request(`/users/${id}`, { method: "DELETE" }),

  listSeries: (params?: { limit?: number; offset?: number; q?: string }) => {
    const qs = new URLSearchParams();
    if (params?.limit) qs.set("limit", String(params.limit));
    if (params?.offset) qs.set("offset", String(params.offset));
    if (params?.q) qs.set("q", params.q);
    return request<{ items: Series[]; total: number }>(`/series?${qs}`);
  },
  getSeries: (id: number) => request<Series>(`/series/${id}`),
  getChapter: (id: number) => request<Chapter>(`/chapters/${id}`),

  deleteSeries: (id: number) =>
    request(`/series/${id}`, { method: "DELETE" }),

  deleteChapter: (id: number) =>
    request(`/chapters/${id}`, { method: "DELETE" }),

  createSeries: (data: Partial<Series> & { cover_key?: string; cover_bucket?: string }) =>
    request<Series>("/series", { method: "POST", body: JSON.stringify(data) }),

  updateSeries: (id: number, data: Partial<Series> & { cover_key?: string; cover_bucket?: string }) =>
    request<Series>(`/series/${id}`, { method: "PUT", body: JSON.stringify(data) }),

  createChapter: (data: { series_id: number; number: number; title: string }) =>
    request<Chapter>("/chapters", { method: "POST", body: JSON.stringify(data) }),

  updateChapter: (id: number, data: { number: number; title: string }) =>
    request<Chapter>(`/chapters/${id}`, { method: "PUT", body: JSON.stringify(data) }),

  uploadChapterPages: (chapterId: number, pages: { key: string; bucket: string; page_number: number }[]) =>
    request(`/chapters/${chapterId}/pages`, { method: "POST", body: JSON.stringify({ pages }) }),

  reorderChapterPages: (chapterId: number, pages: { id: number; page_number: number }[]) =>
    request(`/chapters/${chapterId}/pages`, { method: "PUT", body: JSON.stringify({ pages }) }),

  deleteChapterPage: (chapterId: number, pageId: number) =>
    request(`/chapters/${chapterId}/pages/${pageId}`, { method: "DELETE" }),

  presign: (filename: string, prefix: string) =>
    request<{ upload_url: string; key: string; public_url: string; bucket: string }>("/upload/presign", {
      method: "POST",
      body: JSON.stringify({ filename, prefix }),
    }),

  listComments: (chapterId: number, params?: { limit?: number; offset?: number }) => {
    const qs = new URLSearchParams();
    if (params?.limit) qs.set("limit", String(params.limit));
    if (params?.offset) qs.set("offset", String(params.offset));
    return request<{ items: Comment[]; total: number }>(`/chapters/${chapterId}/comments?${qs}`);
  },

  createComment: (chapterId: number, body: string) =>
    request<Comment>(`/chapters/${chapterId}/comments`, { method: "POST", body: JSON.stringify({ body }) }),

  deleteComment: (id: number) =>
    request(`/comments/${id}`, { method: "DELETE" }),

  suspendComments: (userId: number, duration: string) =>
    request<User>(`/users/${userId}/comment-suspension`, { method: "PATCH", body: JSON.stringify({ duration }) }),

  updateProfile: (name: string, bio: string) =>
    request<User>("/users/me/profile", { method: "PATCH", body: JSON.stringify({ name, bio }) }),

  setAvatar: (key: string, bucket: string) =>
    request<User>("/users/me/avatar", { method: "PATCH", body: JSON.stringify({ key, bucket }) }),

  clearAvatar: () =>
    request<User>("/users/me/avatar", { method: "DELETE" }),

  listMyComments: (params?: { limit?: number; offset?: number }) => {
    const qs = new URLSearchParams();
    if (params?.limit) qs.set("limit", String(params.limit));
    if (params?.offset) qs.set("offset", String(params.offset));
    return request<{ items: Comment[]; total: number }>(`/users/me/comments?${qs}`);
  },

  listFavorites: (params?: { limit?: number; offset?: number }) => {
    const qs = new URLSearchParams();
    if (params?.limit) qs.set("limit", String(params.limit));
    if (params?.offset) qs.set("offset", String(params.offset));
    return request<{ items: Favorite[]; total: number }>(`/favorites?${qs}`);
  },

  addFavorite: (seriesId: number) =>
    request<Favorite>("/favorites", { method: "POST", body: JSON.stringify({ series_id: seriesId }) }),

  removeFavorite: (seriesId: number) =>
    request(`/favorites/${seriesId}`, { method: "DELETE" }),

  updateFavoriteProgress: (seriesId: number, chapterId: number) =>
    request(`/favorites/${seriesId}/progress`, { method: "PATCH", body: JSON.stringify({ chapter_id: chapterId }) }),

  listBackgrounds: () => request<Background[]>("/backgrounds"),

  addBackground: (key: string, bucket: string) =>
    request<Background>("/admin/backgrounds", { method: "POST", body: JSON.stringify({ key, bucket }) }),

  deleteBackground: (id: number) =>
    request(`/admin/backgrounds/${id}`, { method: "DELETE" }),
};

export async function uploadFile(file: File, prefix: string): Promise<{ key: string; bucket: string; public_url: string }> {
  const { upload_url, key, public_url, bucket } = await api.presign(file.name, prefix);
  const res = await fetch(upload_url, { method: "PUT", body: file });
  if (!res.ok) throw new Error("Upload failed");
  return { key, bucket, public_url };
}

export interface S3Object {
  id: number;
  chapter_id?: number;
  bucket: string;
  key: string;
  page_number: number;
  created_at: string;
}

export interface User {
  id: number;
  email: string;
  name: string;
  bio?: string;
  avatar?: S3Object | null;
  role: string;
  created_at: string;
  comment_suspended_until?: string | null;
}

export interface Comment {
  id: number;
  chapter_id: number;
  chapter?: { id: number; series_id: number; number: number; title: string } | null;
  user_id: number;
  user?: { id: number; name: string; avatar?: S3Object | null; comment_suspended_until?: string | null };
  body: string;
  created_at: string;
}

export interface Series {
  id: number;
  title: string;
  slug: string;
  description: string;
  cover_image: S3Object | null;
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

export interface Background {
  id: number;
  image_id: number;
  image: S3Object;
  created_at: string;
}

export interface Favorite {
  id: number;
  user_id: number;
  series_id: number;
  series?: Series;
  last_read_chapter_id?: number | null;
  last_read_chapter?: Chapter | null;
  created_at: string;
  updated_at: string;
}
