import axios from "axios";

const api = axios.create({
  baseURL: "http://localhost:8000",
});

// Intercepteur pour ajouter le token JWT à chaque requête
api.interceptors.request.use((config) => {
  const token = sessionStorage.getItem("token");
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
  api.post("/reservation", {
    ressource_id: data.ressource_id,
    start_datetime: new Date(data.start_datetime).toISOString(),
    end_datetime: new Date(data.end_datetime).toISOString(),
  });

export const createGroupReservation = (data) =>
  api.post("/reservation_group", {
    resource_ids: data.resource_ids,
    start_datetime: new Date(data.start_datetime).toISOString(),
    end_datetime: new Date(data.end_datetime).toISOString(),
  });

export const getUserReservations = () => api.get("/reservations");

// Notifications (user)
export const getUserNotifications = () => api.get("/notifications");

// Admin - Notifications
export const getAdminNotifications = () => api.get("/admin/notifications");
export const markNotificationRead = (id) =>
  api.put(`/admin/notifications/${id}/read`);

// Admin - Resources
export const getAdminResources = () => api.get("/admin/resources");

// Admin - Places
export const getAdminPlaces = () => api.get("/admin/places");
export const getAdminPlaceDetails = (id) => api.get(`/admin/places/${id}`);
export const createAdminPlace = (data) => api.post("/admin/places", data);
export const deleteAdminPlace = (id) => api.delete(`/admin/places/${id}`);

// Admin - Categories
export const createAdminCategory = (data) =>
  api.post("/admin/categories", data);
export const deleteAdminCategory = (id) =>
  api.delete(`/admin/categories/${id}`);

export const createAdminResource = (data) => api.post("/admin/resources", data);
export const deleteAdminResource = (id) => api.delete(`/admin/resources/${id}`);

// Admin - Reservations
export const getAdminReservations = () => api.get("/admin/reservations");
export const approveReservation = (id) =>
  api.put(`/admin/reservations/${id}/approve`);
export const rejectReservation = (id) =>
  api.put(`/admin/reservations/${id}/reject`);
export const approveGroupResources = (id, resourceIds) =>
  api.put(`/admin/reservations/${id}/resources/approve`, { resource_ids: resourceIds });
export const rejectGroupResources = (id, resourceIds) =>
  api.put(`/admin/reservations/${id}/resources/reject`, { resource_ids: resourceIds });

// Admin - Users
export const getAdminUsers = () => api.get("/admin/users");
export const acceptAdminUser = (id) => api.put(`/admin/users/accept/${id}`);
export const refuseAdminUser = (id) => api.put(`/admin/users/refuse/${id}`);
export const promoteAdminUser = (id) => api.put(`/admin/users/promote/${id}`);
export const assignUserPlace = (id, placeId, transfer) =>
  api.put(`/admin/users/${id}/place`, { place_id: placeId, transfer });
export const deleteAdminUser = (id) => api.delete(`/admin/user/${id}`);

export default api;
