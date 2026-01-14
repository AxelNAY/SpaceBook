import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { createAdminResource } from "../api/api";

export default function AdminCreateResource() {
  const navigate = useNavigate();
  const [form, setForm] = useState({
    name: "",
    type: "room",
    category: "",
    capacity: 1,
  });
  const [error, setError] = useState("");
  const [success, setSuccess] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setSuccess(false);

    try {
      await createAdminResource(form);
      setSuccess(true);
      setTimeout(() => navigate("/admin/resources"), 1500);
    } catch (err) {
      setError(err.response?.data?.error || "Erreur lors de la creation");
    }
  };

  return (
    <div className="page-container">
      <h1 className="page-title">Creation de ressource</h1>

      <div className="card" style={{ maxWidth: "500px", margin: "0 auto", padding: "32px" }}>
        {error && (
          <p style={{ color: "var(--danger-red)", textAlign: "center", marginBottom: "16px" }}>
            {error}
          </p>
        )}

        {success && (
          <p style={{ color: "var(--success-green)", textAlign: "center", marginBottom: "16px" }}>
            Ressource creee avec succes !
          </p>
        )}

        <form onSubmit={handleSubmit}>
          <label style={{ display: "block", marginBottom: "8px", fontWeight: 600 }}>
            Nom
          </label>
          <input
            type="text"
            className="form-input"
            placeholder="Nom de la ressource"
            value={form.name}
            onChange={(e) => setForm({ ...form, name: e.target.value })}
            required
          />

          <div style={{ display: "flex", gap: "16px" }}>
            <div style={{ flex: 1 }}>
              <label style={{ display: "block", marginBottom: "8px", fontWeight: 600 }}>
                Type
              </label>
              <select
                className="form-select"
                value={form.type}
                onChange={(e) => setForm({ ...form, type: e.target.value })}
              >
                <option value="room">Salle</option>
                <option value="equipment">Equipement</option>
              </select>
            </div>

            <div style={{ flex: 1 }}>
              <label style={{ display: "block", marginBottom: "8px", fontWeight: 600 }}>
                Capacite
              </label>
              <input
                type="number"
                className="form-input"
                value={form.capacity}
                onChange={(e) => setForm({ ...form, capacity: parseInt(e.target.value) || 1 })}
                min="1"
              />
            </div>
          </div>

          <label style={{ display: "block", marginBottom: "8px", fontWeight: 600 }}>
            Categorie
          </label>
          <select
            className="form-select"
            value={form.category}
            onChange={(e) => setForm({ ...form, category: e.target.value })}
          >
            <option value="">Selectionner une categorie</option>
            <option value="printer">Imprimante</option>
            <option value="projector">Projecteur</option>
            <option value="computer">Ordinateur</option>
            <option value="conference">Conference</option>
            <option value="meeting">Reunion</option>
            <option value="other">Autre</option>
          </select>

          <div style={{ textAlign: "center", marginTop: "24px" }}>
            <button type="submit" className="btn btn-primary">
              Creer la ressource
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
