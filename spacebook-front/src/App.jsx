import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { AuthProvider, useAuth } from "./context/AuthContext";
import Navbar from "./components/Navbar";
import Welcome from "./pages/Welcome";
import Home from "./pages/Home";
import Reservations from "./pages/Reservations";
import MesReservations from "./pages/MesReservations";
import Notifications from "./pages/Notifications";
import AdminCreateResource from "./pages/AdminCreateResource";
import AdminResources from "./pages/AdminResources";
import AdminReservations from "./pages/AdminReservations";
import AdminUsers from "./pages/AdminUsers";

function ProtectedRoute({ children, adminOnly = false }) {
  const { isAuthenticated, isAdmin } = useAuth();

  if (!isAuthenticated) {
    return <Navigate to="/" replace />;
  }

  if (adminOnly && !isAdmin) {
    return <Navigate to="/ressources" replace />;
  }

  return children;
}

function AppRoutes() {
  const { isAuthenticated } = useAuth();

  return (
    <Routes>
      {/* Page d'accueil - Welcome si non connecte, Ressources si connecte */}
      <Route
        path="/"
        element={isAuthenticated ? <Home /> : <Welcome />}
      />

      {/* Pages accessibles uniquement si connecte */}
      <Route
        path="/ressources"
        element={
          <ProtectedRoute>
            <Home />
          </ProtectedRoute>
        }
      />

      <Route
        path="/reserver"
        element={
          <ProtectedRoute>
            <Reservations />
          </ProtectedRoute>
        }
      />

      <Route
        path="/mes-reservations"
        element={
          <ProtectedRoute>
            <MesReservations />
          </ProtectedRoute>
        }
      />

      <Route
        path="/notifications"
        element={
          <ProtectedRoute>
            <Notifications />
          </ProtectedRoute>
        }
      />

      {/* Pages Admin */}
      <Route
        path="/admin/resources"
        element={
          <ProtectedRoute adminOnly>
            <AdminResources />
          </ProtectedRoute>
        }
      />

      <Route
        path="/admin/resources/create"
        element={
          <ProtectedRoute adminOnly>
            <AdminCreateResource />
          </ProtectedRoute>
        }
      />

      <Route
        path="/admin/reservations"
        element={
          <ProtectedRoute adminOnly>
            <AdminReservations />
          </ProtectedRoute>
        }
      />

      <Route
        path="/admin/users"
        element={
          <ProtectedRoute adminOnly>
            <AdminUsers />
          </ProtectedRoute>
        }
      />

      {/* Redirection par defaut */}
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}

export default function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Navbar />
        <AppRoutes />
      </BrowserRouter>
    </AuthProvider>
  );
}
