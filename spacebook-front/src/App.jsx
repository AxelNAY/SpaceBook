import { BrowserRouter, Routes, Route } from "react-router-dom";
import Navbar from "./components/Navbar";
import Home from "./pages/Home";
import Reservations from "./pages/Reservations";
import AdminNotifications from "./pages/AdminNotifications";
import AdminCreateResource from "./pages/AdminCreateResource";
import AdminResources from "./pages/AdminResources";
import AdminReservations from "./pages/AdminReservations";

export default function App() {
  return (
    <BrowserRouter>
      <Navbar />
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/reserve" element={<Reservations />} />
        <Route path="/admin" element={<AdminNotifications />} />
        <Route path="/admin/resources" element={<AdminResources />} />
        <Route path="/admin/resources/create" element={<AdminCreateResource />} />
        <Route path="/admin/reservations" element={<AdminReservations />} />
      </Routes>
    </BrowserRouter>
  );
}
