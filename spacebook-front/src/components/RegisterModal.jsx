import { useState, useEffect } from "react";
import { register as apiRegister, getCompanies, getPlaces } from "../api/api";

export default function RegisterModal({ onClose, onSwitchToLogin }) {
  const [form, setForm] = useState({
    username: "",
    email: "",
    phone: "",
    password: "",
    place_id: "",
  });
  const [companies, setCompanies] = useState([]);
  const [places, setPlaces] = useState([]);
  const [selectedCompany, setSelectedCompany] = useState("");
  const [error, setError] = useState("");
  const [success, setSuccess] = useState(false);
  const [loading, setLoading] = useState(false);

  const passwordRules = [
    { label: "8 caractères minimum", test: (p) => p.length >= 8 },
    { label: "Une majuscule", test: (p) => /[A-Z]/.test(p) },
    { label: "Une minuscule", test: (p) => /[a-z]/.test(p) },
    { label: "Un chiffre", test: (p) => /[0-9]/.test(p) },
    { label: "Un caractère spécial", test: (p) => /[^A-Za-z0-9]/.test(p) },
  ];
  const passwordValid = passwordRules.every((r) => r.test(form.password));

  useEffect(() => {
    getCompanies().then((res) => setCompanies(res.data || [])).catch(() => {});
  }, []);

  useEffect(() => {
    if (selectedCompany) {
      getPlaces(selectedCompany).then((res) => setPlaces(res.data || [])).catch(() => {});
      setForm((f) => ({ ...f, place_id: "" }));
    } else {
      setPlaces([]);
      setForm((f) => ({ ...f, place_id: "" }));
    }
  }, [selectedCompany]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    if (!passwordValid) {
      setError("Le mot de passe ne respecte pas les critères de sécurité.");
      setLoading(false);
      return;
    }

    try {
      await apiRegister({ ...form, phone: form.phone ? parseInt(form.phone, 10) : 0 });
      setSuccess(true);
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

        {success ? (
          <div style={{ textAlign: "center" }}>
            <p style={{ color: "var(--success-green)", fontSize: "16px", marginBottom: "16px" }}>
              Votre compte a été créé et est en attente de validation par un administrateur.
            </p>
            <button className="btn btn-primary" onClick={onSwitchToLogin}>
              Se connecter
            </button>
          </div>
        ) : (
        <>
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
          {form.password && (
            <ul style={{ fontSize: "12px", paddingLeft: "8px", marginBottom: "10px", listStyle: "none" }}>
              {passwordRules.map((r) => (
                <li key={r.label} style={{ color: r.test(form.password) ? "var(--success-green)" : "var(--danger-red)" }}>
                  {r.test(form.password) ? "✓" : "✗"} {r.label}
                </li>
              ))}
            </ul>
          )}

          <select
            className="form-select"
            value={selectedCompany}
            onChange={(e) => setSelectedCompany(e.target.value)}
          >
            <option value="">Sélectionner une entreprise</option>
            {companies.map((c) => (
              <option key={c.id} value={c.id}>{c.name}</option>
            ))}
          </select>

          <select
            className="form-select"
            value={form.place_id}
            onChange={(e) => setForm({ ...form, place_id: e.target.value })}
            disabled={!selectedCompany}
          >
            <option value="">Sélectionner un lieu</option>
            {places.map((p) => (
              <option key={p.id} value={p.id}>{p.name}</option>
            ))}
          </select>

          <div style={{ textAlign: "center", marginTop: "12px" }}>
            <button
              type="submit"
              className="btn btn-danger"
              disabled={loading}
            >
              {loading ? "Inscription..." : "Inscription"}
            </button>
          </div>
        </form>

        <p style={{ textAlign: "center", marginTop: "10px" }}>
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
        </>
        )}
      </div>
    </div>
  );
}
