import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import { login as apiLogin } from "../api/api";

export default function LoginModal({ onClose, onSwitchToRegister }) {
  const { login } = useAuth();
  const navigate = useNavigate();
  const [form, setForm] = useState({
    email: "",
    password: "",
  });
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      const response = await apiLogin(form);
      const user = response.data.user;
      login(user, response.data.token);
      onClose();

      // Rediriger l'admin vers la page de creation
      if (user.role === "admin") {
        navigate("/admin/resources/create");
      }
    } catch (err) {
      setError(err.response?.data?.error || "Erreur de connexion");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal" onClick={(e) => e.stopPropagation()}>
        <button className="modal-close" onClick={onClose}>×</button>

        <h2>Connexion</h2>

        {error && (
          <p style={{ color: "var(--danger-red)", textAlign: "center" }}>{error}</p>
        )}

        <form onSubmit={handleSubmit}>
          <input
            type="email"
            className="form-input"
            placeholder="Email"
            value={form.email}
            onChange={(e) => setForm({ ...form, email: e.target.value })}
            required
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
              {loading ? "Connexion..." : "Connexion"}
            </button>
          </div>
        </form>

        <p style={{ textAlign: "center", marginTop: "16px" }}>
          Pas encore de compte ?{" "}
          <button
            onClick={onSwitchToRegister}
            style={{
              background: "none",
              border: "none",
              color: "var(--primary-blue)",
              cursor: "pointer",
              textDecoration: "underline"
            }}
          >
            S'inscrire
          </button>
        </p>
      </div>
    </div>
  );
}
