import axios from "axios";

const api = axios.create({
  baseURL: "http://localhost:8000",
});

// Intercepteur pour ajouter le token JWT à chaque requête
api.interceptors.request.use((config) => {
  const token = localStorage.getItem("token");
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Auth
export const login = (data) => api.post("/auth/login", data);
export const register = (data) => api.post("/auth/register", data);

// Resources
export const getResources = () => api.get("/resources");

// Reservations (user)
export const createReservation = (data) =>
  api.post("/reservations", {
    ...data,
    start_at: new Date(data.start_at).toISOString(),
    end_at: new Date(data.end_at).toISOString(),
  });

export const getUserReservations = (userId) =>
  api.get(`/reservations?userId=${userId}`);

// Notifications (user)
export const getUserNotifications = (userId) =>
  api.get(`/notifications?userId=${userId}`);

// Admin - Notifications
export const getAdminNotifications = () => api.get("/admin/notifications");
export const markNotificationRead = (id) =>
  api.put(`/admin/notifications/${id}/read`);

// Admin - Resources
export const createAdminResource = (data) => api.post("/admin/resources", data);
export const deleteAdminResource = (id) => api.delete(`/admin/resources/${id}`);

// Admin - Reservations
export const getAdminReservations = () => api.get("/admin/reservations");
export const approveReservation = (id) =>
  api.put(`/admin/reservations/${id}/approve`);
export const rejectReservation = (id) =>
  api.put(`/admin/reservations/${id}/reject`);

// Admin - Users
export const getAdminUsers = () => api.get("/admin/users");
export const deleteAdminUser = (id) => api.delete(`/admin/user/${id}`);

export default api;
