import { useEffect, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { getResources, createReservation } from "../api/api";
import { useAuth } from "../context/AuthContext";

export default function Reservations() {
  const { user } = useAuth();
  const location = useLocation();
  const navigate = useNavigate();
  const [resources, setResources] = useState([]);
  const [form, setForm] = useState({
    resource_id: location.state?.resource?.ID || "",
    user_id: user?.id || "",
    start_at: "",
    end_at: "",
  });
  const [error, setError] = useState("");
  const [success, setSuccess] = useState(false);

  useEffect(() => {
    getResources().then((res) => setResources(res.data || []));
  }, []);

  useEffect(() => {
    if (user) {
      setForm((f) => ({ ...f, user_id: user.id }));
    }
  }, [user]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setSuccess(false);

    try {
      await createReservation(form);
      setSuccess(true);
      setTimeout(() => navigate("/mes-reservations"), 1500);
    } catch (err) {
      setError(err.response?.data?.error || "Erreur lors de la reservation");
    }
  };

  return (
    <div className="page-container">
      <h1 className="page-title">Nouvelle Reservation</h1>

      <div className="card" style={{ maxWidth: "500px", margin: "0 auto", padding: "32px" }}>
        {error && (
          <p style={{ color: "var(--danger-red)", textAlign: "center", marginBottom: "16px" }}>
            {error}
          </p>
        )}

        {success && (
          <p style={{ color: "var(--success-green)", textAlign: "center", marginBottom: "16px" }}>
            Reservation creee avec succes !
          </p>
        )}

        <form onSubmit={handleSubmit}>
          <label style={{ display: "block", marginBottom: "8px", fontWeight: 600 }}>
            Ressource
          </label>
          <select
            className="form-select"
            value={form.resource_id}
            onChange={(e) => setForm({ ...form, resource_id: e.target.value })}
            required
          >
            <option value="">Selectionner une ressource</option>
            {resources.map((r) => (
              <option key={r.ID} value={r.ID}>
                {r.Name} ({r.Type === "room" ? "Salle" : "Equipement"})
              </option>
            ))}
          </select>

          <label style={{ display: "block", marginBottom: "8px", fontWeight: 600 }}>
            Date et heure de debut
          </label>
          <input
            type="datetime-local"
            className="form-input"
            value={form.start_at}
            onChange={(e) => setForm({ ...form, start_at: e.target.value })}
            required
          />

          <label style={{ display: "block", marginBottom: "8px", fontWeight: 600 }}>
            Date et heure de fin
          </label>
          <input
            type="datetime-local"
            className="form-input"
            value={form.end_at}
            onChange={(e) => setForm({ ...form, end_at: e.target.value })}
            required
          />

          <div style={{ textAlign: "center", marginTop: "24px" }}>
            <button type="submit" className="btn btn-primary">
              Creer la reservation
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
