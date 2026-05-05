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

// Companies
export const getCompanies = () => api.get("/companies");

// Places
export const getPlaces = (companyId) =>
  api.get(`/places${companyId ? `?company_id=${companyId}` : ""}`);

// Categories
export const getCategories = () => api.get("/categories");

// Resources
export const getResources = (placeId) =>
  api.get(`/resources${placeId ? `?place_id=${placeId}` : ""}`);

// Reservations (user)
export const createReservation = (data) =>
  api.post("/reservations", {
    ...data,
    start_datetime: new Date(data.start_datetime).toISOString(),
    end_datetime: new Date(data.end_datetime).toISOString(),
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

// Admin - Companies
export const createAdminCompany = (data) => api.post("/admin/companies", data);
export const deleteAdminCompany = (id) => api.delete(`/admin/companies/${id}`);

// Admin - Places
export const createAdminPlace = (data) => api.post("/admin/places", data);
export const deleteAdminPlace = (id) => api.delete(`/admin/places/${id}`);

// Admin - Categories
export const createAdminCategory = (data) =>
  api.post("/admin/categories", data);
export const deleteAdminCategory = (id) =>
  api.delete(`/admin/categories/${id}`);

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
