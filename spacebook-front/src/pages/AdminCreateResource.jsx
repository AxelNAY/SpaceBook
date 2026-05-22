import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { createAdminResource, getAdminPlaces, getCategories } from "../api/api";

export default function AdminCreateResource() {
  const navigate = useNavigate();
  const [form, setForm] = useState({
    name: "",
    place_id: "",
    category_id: "",
  });
  const [places, setPlaces] = useState([]);
  const [categories, setCategories] = useState([]);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState(false);

  useEffect(() => {
    getAdminPlaces().then((res) => setPlaces(res.data)).catch(() => {});
    getCategories().then((res) => setCategories(res.data)).catch(() => {});
  }, []);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setSuccess(false);

    if (!form.place_id) {
      setError("Veuillez sélectionner un lieu");
      return;
    }

    try {
      await createAdminResource(form);
      setSuccess(true);
      setTimeout(() => navigate("/admin/resources"), 1500);
    } catch (err) {
      setError(err.response?.data?.error || "Erreur lors de la création");
    }
  };

  return (
    <div className="page-container">
      <h1 className="page-title">Création de ressource</h1>

      <div className="card" style={{ maxWidth: "500px", margin: "0 auto", padding: "32px" }}>
        {error && (
          <p style={{ color: "var(--danger-red)", textAlign: "center", marginBottom: "16px" }}>
            {error}
          </p>
        )}

        {success && (
          <p style={{ color: "var(--success-green)", textAlign: "center", marginBottom: "16px" }}>
            Ressource créée avec succès !
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

          <label style={{ display: "block", marginBottom: "8px", fontWeight: 600 }}>
            Lieu
          </label>
          <select
            className="form-select"
            value={form.place_id}
            onChange={(e) => setForm({ ...form, place_id: e.target.value })}
            required
          >
            <option value="">Sélectionner un lieu</option>
            {places.map((p) => (
              <option key={p.id} value={p.id}>{p.name}</option>
            ))}
          </select>

          <label style={{ display: "block", marginBottom: "8px", fontWeight: 600 }}>
            Catégorie
          </label>
          <select
            className="form-select"
            value={form.category_id}
            onChange={(e) => setForm({ ...form, category_id: e.target.value })}
          >
            <option value="">Sélectionner une catégorie</option>
            {categories.map((cat) => (
              <option key={cat.id} value={cat.id}>{cat.name}</option>
            ))}
          </select>

          <div style={{ textAlign: "center", marginTop: "24px" }}>
            <button type="submit" className="btn btn-primary">
              Créer la ressource
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
