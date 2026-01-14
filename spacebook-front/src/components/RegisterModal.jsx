import { useState } from "react";
import { useAuth } from "../context/AuthContext";
import { register as apiRegister } from "../api/api";

export default function RegisterModal({ onClose, onSwitchToLogin }) {
  const { login } = useAuth();
  const [form, setForm] = useState({
    username: "",
    email: "",
    phone: "",
    password: "",
  });
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      const response = await apiRegister(form);
      login(response.data.user, response.data.token);
      onClose();
    } catch (err) {
      setError(err.response?.data?.error || "Erreur lors de l'inscription");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal" onClick={(e) => e.stopPropagation()}>
        <button className="modal-close" onClick={onClose}>×</button>

        <h2>Inscription</h2>

        {error && (
          <p style={{ color: "var(--danger-red)", textAlign: "center" }}>{error}</p>
        )}

        <form onSubmit={handleSubmit}>
          <input
            type="text"
            className="form-input"
            placeholder="Nom d'utilisateur"
            value={form.username}
            onChange={(e) => setForm({ ...form, username: e.target.value })}
            required
          />

          <input
            type="email"
            className="form-input"
            placeholder="Email"
            value={form.email}
            onChange={(e) => setForm({ ...form, email: e.target.value })}
            required
          />

          <input
            type="tel"
            className="form-input"
            placeholder="Téléphone"
            value={form.phone}
            onChange={(e) => setForm({ ...form, phone: e.target.value })}
          />

          <input
            type="password"
            className="form-input"
            placeholder="Mot de passe"
            value={form.password}
            onChange={(e) => setForm({ ...form, password: e.target.value })}
            required
          />

          <div style={{ textAlign: "center", marginTop: "24px" }}>
            <button
              type="submit"
              className="btn btn-danger"
              disabled={loading}
            >
              {loading ? "Inscription..." : "Inscription"}
            </button>
          </div>
        </form>

        <p style={{ textAlign: "center", marginTop: "16px" }}>
          Déjà un compte ?{" "}
          <button
            onClick={onSwitchToLogin}
            style={{
              background: "none",
              border: "none",
              color: "var(--primary-blue)",
              cursor: "pointer",
              textDecoration: "underline"
            }}
          >
            Se connecter
          </button>
        </p>
      </div>
    </div>
  );
}
