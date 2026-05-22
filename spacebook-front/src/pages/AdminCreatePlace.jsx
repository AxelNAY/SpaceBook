import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { createAdminPlace } from "../api/api";

export default function AdminCreatePlace() {
  const navigate = useNavigate();
  const [form, setForm] = useState({ name: "", description: "" });
  const [error, setError] = useState("");
  const [success, setSuccess] = useState(false);
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setLoading(true);
    try {
      await createAdminPlace(form);
      setSuccess(true);
      setTimeout(() => navigate("/admin/places"), 1500);
    } catch (err) {
      setError(err.response?.data?.error || "Erreur lors de la création");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="page-container">
      <h1 className="page-title">Créer un lieu</h1>

      <div className="card" style={{ maxWidth: "500px", margin: "0 auto", padding: "32px" }}>
        {error && (
          <p style={{ color: "var(--danger-red)", textAlign: "center", marginBottom: "16px" }}>{error}</p>
        )}
        {success && (
          <p style={{ color: "var(--success-green)", textAlign: "center", marginBottom: "16px" }}>
            Lieu créé avec succès !
          </p>
        )}

        <form onSubmit={handleSubmit}>
          <label style={{ display: "block", marginBottom: "8px", fontWeight: 600 }}>Nom</label>
          <input
            type="text"
            className="form-input"
            placeholder="Nom du lieu"
            value={form.name}
            onChange={(e) => setForm({ ...form, name: e.target.value })}
            required
          />

          <label style={{ display: "block", marginBottom: "8px", fontWeight: 600 }}>Description</label>
          <textarea
            className="form-textarea"
            placeholder="Description (optionnelle)"
            value={form.description}
            onChange={(e) => setForm({ ...form, description: e.target.value })}
          />

          <div style={{ display: "flex", gap: "12px", justifyContent: "center", marginTop: "24px" }}>
            <button
              type="button"
              className="btn btn-outline"
              onClick={() => navigate("/admin/places")}
            >
              Annuler
            </button>
            <button type="submit" className="btn btn-primary" disabled={loading}>
              {loading ? "Création..." : "Créer le lieu"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
