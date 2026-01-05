import axios from "axios";

const api = axios.create({
  baseURL: "http://localhost:8080",
});

export const getResources = () => api.get("/resources");

export const createReservation = (data) =>
  api.post("/reservations", {
    ...data,
    startAt: new Date(data.startAt).toISOString(),
    endAt: new Date(data.endAt).toISOString(),
  });

export const getAdminNotifications = () =>
  api.get("/admin/notifications", {
    headers: { "X-ROLE": "admin" },
  });

export const markNotificationRead = (id) =>
  api.put(`/admin/notifications/${id}/read`, null, {
    headers: { "X-ROLE": "admin" },
  });

export const createAdminResource = (data) =>
  api.post("/admin/resources", data, {
    headers: { "X-ROLE": "admin" },
  });

export const getReservations = () => api.get("/reservations");

export const approveReservation = (id) =>
  api.put(`/admin/reservations/${id}/approve`, null, {
    headers: { "X-ROLE": "admin" },
  });

export const deleteAdminResource = (id) =>
  api.delete(`/admin/resources/${id}`, {
    headers: { "X-ROLE": "admin" },
  });

export const getAdminReservations = () =>
  axios.get("http://localhost:8080/admin/reservations", {
    headers: { "X-ROLE": "admin" }
  });

export const deleteResource = (id) =>
  axios.delete(`http://localhost:8080/admin/resources/${id}`, {
    headers: { "X-ROLE": "admin" }
  });

export default api;
