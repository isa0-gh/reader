const BASE = "/api/v1";
function authHeader() {
    const token = localStorage.getItem("token");
    return token ? { Authorization: `Bearer ${token}` } : {};
}
async function request(path, init) {
    const res = await fetch(BASE + path, {
        headers: { "Content-Type": "application/json", ...authHeader() },
        ...init,
    });
    if (!res.ok)
        throw new Error(await res.text());
    return res.json();
}
export const api = {
    register: (email, password, name) => request("/auth/register", { method: "POST", body: JSON.stringify({ email, password, name }) }),
    login: (email, password) => request("/auth/login", { method: "POST", body: JSON.stringify({ email, password }) }),
    getSeries: (id) => request(`/series/${id}`),
    getChapter: (id) => request(`/chapters/${id}`),
};
