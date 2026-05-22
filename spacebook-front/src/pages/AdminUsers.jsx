import { useEffect, useState } from "react";
import {
  getAdminUsers,
  acceptAdminUser,
  refuseAdminUser,
  promoteAdminUser,
  assignUserPlace,
  deleteAdminUser,
  getAdminPlaces,
} from "../api/api";

function AssignPlaceModal({ user, onClose, onDone }) {
  const [places, setPlaces] = useState([]);
  const [placeId, setPlaceId] = useState("");
  const [transfer, setTransfer] = useState(true);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    getAdminPlaces().then((res) => setPlaces(res.data || [])).catch(() => {});
  }, []);

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!placeId) return;
    setLoading(true);
    setError("");
    try {
      await assignUserPlace(user.id, placeId, transfer);
      onDone();
    } catch (err) {
      setError(err.response?.data?.error || "Erreur lors de l'affectation");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal" onClick={(e) => e.stopPropagation()}>
        <button className="modal-close" onClick={onClose}>×</button>
        <h2>Affecter à un lieu</h2>
        <p style={{ textAlign: "center", color: "var(--text-gray)", marginBottom: "20px" }}>
          {user.username}
        </p>

        {error && (
          <p style={{ color: "var(--danger-red)", textAlign: "center", marginBottom: "12px" }}>{error}</p>
        )}

        <form onSubmit={handleSubmit}>
          <select
            className="form-select"
            value={placeId}
            onChange={(e) => setPlaceId(e.target.value)}
            required
          >
            <option value="">Sélectionner un lieu</option>
            {places.map((p) => (
              <option key={p.id} value={p.id}>{p.name}</option>
            ))}
          </select>

          <div style={{ display: "flex", gap: "16px", marginBottom: "24px" }}>
            <label style={{ display: "flex", alignItems: "center", gap: "8px", cursor: "pointer" }}>
              <input
                type="radio"
                name="mode"
                checked={transfer}
                onChange={() => setTransfer(true)}
              />
              Transférer (remplace les lieux actuels)
            </label>
            <label style={{ display: "flex", alignItems: "center", gap: "8px", cursor: "pointer" }}>
              <input
                type="radio"
                name="mode"
                checked={!transfer}
                onChange={() => setTransfer(false)}
              />
              Ajouter (conserve les lieux actuels)
            </label>
          </div>

          <div style={{ textAlign: "center" }}>
            <button type="submit" className="btn btn-primary" disabled={loading || !placeId}>
              {loading ? "Enregistrement..." : "Confirmer"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

export default function AdminUsers() {
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [assignTarget, setAssignTarget] = useState(null);

  const load = async () => {
    try {
      const res = await getAdminUsers();
      setUsers(res.data || []);
    } catch (err) {
      console.error(err);
      setUsers([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    load();
  }, []);

  const getStatusBadge = (status) => {
    switch (status) {
      case "Actif":
        return <span className="badge badge-approved">Actif</span>;
      case "En attente":
        return <span className="badge badge-pending">En attente</span>;
      default:
        return <span className="badge badge-inactive">{status || "—"}</span>;
    }
  };

  const handleAccept = async (userId, username) => {
    if (!window.confirm(`Valider le compte de "${username}" ?`)) return;
    try {
      await acceptAdminUser(userId);
      load();
    } catch (err) {
      alert(err.response?.data?.error || "Erreur lors de la validation");
    }
  };

  const handleRefuse = async (userId, username) => {
    if (!window.confirm(`Refuser et supprimer le compte de "${username}" ?`)) return;
    try {
      await refuseAdminUser(userId);
      load();
    } catch (err) {
      alert(err.response?.data?.error || "Erreur lors du refus");
    }
  };

  const handlePromote = async (userId, username) => {
    if (!window.confirm(`Promouvoir "${username}" au rang d'administrateur ?`)) return;
    try {
      await promoteAdminUser(userId);
      load();
    } catch (err) {
      alert(err.response?.data?.error || "Erreur lors de la promotion");
    }
  };

  const handleDelete = async (userId, username) => {
    if (!window.confirm(`Supprimer l'utilisateur "${username}" ?`)) return;
    try {
      await deleteAdminUser(userId);
      load();
    } catch (err) {
      alert(err.response?.data?.error || "Erreur lors de la suppression");
    }
  };

  if (loading) {
    return (
      <div className="page-container">
        <p style={{ textAlign: "center" }}>Chargement...</p>
      </div>
    );
  }

  return (
    <div className="page-container">
      <h1 className="page-title">Utilisateurs</h1>

      {assignTarget && (
        <AssignPlaceModal
          user={assignTarget}
          onClose={() => setAssignTarget(null)}
          onDone={() => { setAssignTarget(null); load(); }}
        />
      )}

      {users.length === 0 ? (
        <p style={{ textAlign: "center", color: "var(--text-gray)" }}>
          Aucun utilisateur
        </p>
      ) : (
        <div className="card" style={{ overflow: "hidden" }}>
          <table style={{ width: "100%", borderCollapse: "collapse" }}>
            <thead>
              <tr style={{ backgroundColor: "var(--primary-blue)", color: "white" }}>
                <th style={{ padding: "16px 24px", textAlign: "left", fontWeight: 600 }}>Nom</th>
                <th style={{ padding: "16px 24px", textAlign: "left", fontWeight: 600 }}>Email</th>
                <th style={{ padding: "16px 24px", textAlign: "center", fontWeight: 600 }}>Rôle</th>
                <th style={{ padding: "16px 24px", textAlign: "center", fontWeight: 600 }}>Statut</th>
                <th style={{ padding: "16px 24px", textAlign: "center", fontWeight: 600 }}>Actions</th>
              </tr>
            </thead>
            <tbody>
              {users.map((user) => (
                <tr key={user.id} style={{ borderBottom: "1px solid #e5e5e5" }}>
                  <td style={{ padding: "16px 24px", fontWeight: 600 }}>{user.username}</td>
                  <td style={{ padding: "16px 24px" }}>{user.email}</td>
                  <td style={{ padding: "16px 24px", textAlign: "center" }}>
                    <span className={`badge ${user.role === "admin" ? "badge-approved" : "badge-pending"}`}>
                      {user.role === "admin" ? "Administrateur" : "Utilisateur"}
                    </span>
                  </td>
                  <td style={{ padding: "16px 24px", textAlign: "center" }}>
                    {getStatusBadge(user.status)}
                  </td>
                  <td style={{ padding: "16px 24px", textAlign: "center" }}>
                    {user.status === "En attente" ? (
                      <div style={{ display: "flex", gap: "8px", justifyContent: "center" }}>
                        <button
                          className="btn btn-primary"
                          style={{ padding: "8px 16px", fontSize: "14px" }}
                          onClick={() => handleAccept(user.id, user.username)}
                        >
                          Accepter
                        </button>
                        <button
                          className="btn btn-danger"
                          style={{ padding: "8px 16px", fontSize: "14px" }}
                          onClick={() => handleRefuse(user.id, user.username)}
                        >
                          Refuser
                        </button>
                      </div>
                    ) : (
                      <div style={{ display: "flex", gap: "8px", justifyContent: "center", flexWrap: "wrap" }}>
                        <button
                          className="btn btn-outline"
                          style={{ padding: "8px 16px", fontSize: "14px" }}
                          onClick={() => setAssignTarget(user)}
                        >
                          Lieu
                        </button>
                        <button
                          className="btn btn-primary"
                          style={{ padding: "8px 16px", fontSize: "14px" }}
                          onClick={() => handlePromote(user.id, user.username)}
                        >
                          Promouvoir
                        </button>
                        <button
                          className="btn btn-danger"
                          style={{ padding: "8px 16px", fontSize: "14px" }}
                          onClick={() => handleDelete(user.id, user.username)}
                        >
                          Supprimer
                        </button>
                      </div>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
