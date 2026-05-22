import { useEffect, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { getResources, createReservation, createGroupReservation } from "../api/api";
import { useAuth } from "../context/AuthContext";

export default function Reservations() {
  const { user } = useAuth();
  const location = useLocation();
  const navigate = useNavigate();

  const [resources, setResources] = useState([]);
  const [mode, setMode] = useState("simple");
  const [form, setForm] = useState({
    ressource_id: location.state?.resource?.id || "",
    start_datetime: "",
    end_datetime: "",
  });
  const [selectedIds, setSelectedIds] = useState(
    location.state?.resource?.id ? [location.state.resource.id] : []
  );
  const [error, setError] = useState("");
  const [success, setSuccess] = useState(false);

  useEffect(() => {
    getResources().then((res) => setResources(res.data || []));
  }, []);

  const toggleResource = (id) => {
    setSelectedIds((prev) =>
      prev.includes(id) ? prev.filter((r) => r !== id) : [...prev, id]
    );
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setSuccess(false);

    try {
      if (mode === "simple") {
        await createReservation(form);
      } else {
        if (selectedIds.length < 2) {
          setError("Sélectionnez au moins deux ressources pour une réservation de groupe");
          return;
        }
        await createGroupReservation({
          resource_ids: selectedIds,
          start_datetime: form.start_datetime,
          end_datetime: form.end_datetime,
        });
      }
      setSuccess(true);
      setTimeout(() => navigate("/mes-reservations"), 1500);
    } catch (err) {
      setError(err.response?.data?.error || "Erreur lors de la réservation");
    }
  };

  return (
    <div className="page-container">
      <h1 className="page-title">Nouvelle Réservation</h1>

      <div className="card" style={{ maxWidth: "560px", margin: "0 auto", padding: "32px" }}>
        {/* Toggle Simple / Groupe */}
        <div style={{ display: "flex", gap: "8px", marginBottom: "24px" }}>
          <button
            type="button"
            className={mode === "simple" ? "btn btn-primary" : "btn btn-secondary"}
            style={{ flex: 1 }}
            onClick={() => { setMode("simple"); setError(""); }}
          >
            Simple
          </button>
          <button
            type="button"
            className={mode === "groupe" ? "btn btn-primary" : "btn btn-secondary"}
            style={{ flex: 1 }}
            onClick={() => { setMode("groupe"); setError(""); }}
          >
            Groupe
          </button>
        </div>

        {error && (
          <p style={{ color: "var(--danger-red)", textAlign: "center", marginBottom: "16px" }}>
            {error}
          </p>
        )}
        {success && (
          <p style={{ color: "var(--success-green)", textAlign: "center", marginBottom: "16px" }}>
            Réservation créée avec succès !
          </p>
        )}

        <form onSubmit={handleSubmit}>
          {mode === "simple" ? (
            <>
              <label style={{ display: "block", marginBottom: "8px", fontWeight: 600 }}>
                Ressource
              </label>
              <select
                className="form-select"
                value={form.ressource_id}
                onChange={(e) => setForm({ ...form, ressource_id: e.target.value })}
                required
              >
                <option value="">Sélectionner une ressource</option>
                {resources.map((r) => (
                  <option key={r.id} value={r.id}>
                    {r.name}{r.place ? ` — ${r.place.name}` : ""}
                  </option>
                ))}
              </select>
            </>
          ) : (
            <>
              <label style={{ display: "block", marginBottom: "8px", fontWeight: 600 }}>
                Ressources{" "}
                <span style={{ fontWeight: 400, color: "var(--text-gray)", fontSize: "13px" }}>
                  (minimum 2)
                </span>
              </label>
              <div
                style={{
                  border: "1px solid #e5e5e5",
                  borderRadius: "8px",
                  maxHeight: "200px",
                  overflowY: "auto",
                  marginBottom: "16px",
                }}
              >
                {resources.map((r) => (
                  <label
                    key={r.id}
                    style={{
                      display: "flex",
                      alignItems: "center",
                      gap: "10px",
                      padding: "10px 14px",
                      cursor: "pointer",
                      borderBottom: "1px solid #f0f0f0",
                      backgroundColor: selectedIds.includes(r.id) ? "var(--light-blue, #e8f0fe)" : "transparent",
                    }}
                  >
                    <input
                      type="checkbox"
                      checked={selectedIds.includes(r.id)}
                      onChange={() => toggleResource(r.id)}
                    />
                    <span style={{ fontWeight: 500 }}>{r.name}</span>
                    {r.place && (
                      <span style={{ color: "var(--text-gray)", fontSize: "13px" }}>
                        — {r.place.name}
                      </span>
                    )}
                  </label>
                ))}
              </div>
              {selectedIds.length > 0 && (
                <p style={{ fontSize: "13px", color: "var(--text-gray)", marginBottom: "16px" }}>
                  {selectedIds.length} ressource(s) sélectionnée(s)
                </p>
              )}
            </>
          )}

          <label style={{ display: "block", marginBottom: "8px", fontWeight: 600 }}>
            Date et heure de début
          </label>
          <input
            type="datetime-local"
            className="form-input"
            value={form.start_datetime}
            onChange={(e) => setForm({ ...form, start_datetime: e.target.value })}
            required
          />

          <label style={{ display: "block", marginBottom: "8px", fontWeight: 600 }}>
            Date et heure de fin
          </label>
          <input
            type="datetime-local"
            className="form-input"
            value={form.end_datetime}
            onChange={(e) => setForm({ ...form, end_datetime: e.target.value })}
            required
          />

          <div style={{ textAlign: "center", marginTop: "24px" }}>
            <button type="submit" className="btn btn-primary">
              Créer la réservation
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
