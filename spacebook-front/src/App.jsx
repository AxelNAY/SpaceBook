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
import AdminPlaces from "./pages/AdminPlaces";
import AdminCreatePlace from "./pages/AdminCreatePlace";

const USER_HOME = "/ressources";
const ADMIN_HOME = "/admin/reservations";

function ProtectedRoute({ children, adminOnly = false, userOnly = false }) {
  const { isAuthenticated, isAdmin } = useAuth();

  if (!isAuthenticated) {
    return <Navigate to="/" replace />;
  }

  if (adminOnly && !isAdmin) {
    return <Navigate to={USER_HOME} replace />;
  }

  if (userOnly && isAdmin) {
    return <Navigate to={ADMIN_HOME} replace />;
  }

  return children;
}

function AppRoutes() {
  const { isAuthenticated, isAdmin } = useAuth();

  return (
    <Routes>
      {/* Page d'accueil - redirige vers la home du rôle si connecté */}
      <Route
        path="/"
        element={
          !isAuthenticated ? <Welcome /> :
          isAdmin ? <Navigate to={ADMIN_HOME} replace /> :
          <Navigate to={USER_HOME} replace />
        }
      />

      {/* Pages utilisateur uniquement */}
      <Route
        path="/ressources"
        element={
          <ProtectedRoute userOnly>
            <Home />
          </ProtectedRoute>
        }
      />

      <Route
        path="/reserver"
        element={
          <ProtectedRoute userOnly>
            <Reservations />
          </ProtectedRoute>
        }
      />

      <Route
        path="/mes-reservations"
        element={
          <ProtectedRoute userOnly>
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

      <Route
        path="/admin/places"
        element={
          <ProtectedRoute adminOnly>
            <AdminPlaces />
          </ProtectedRoute>
        }
      />

      <Route
        path="/admin/places/create"
        element={
          <ProtectedRoute adminOnly>
            <AdminCreatePlace />
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
