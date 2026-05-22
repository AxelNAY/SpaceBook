import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { getAdminPlaces, deleteAdminPlace } from "../api/api";

const statusBadge = (status) => {
  switch (status) {
    case "Actif": return <span className="badge badge-approved">Actif</span>;
    case "En attente": return <span className="badge badge-pending">En attente</span>;
    default: return <span className="badge badge-inactive">{status || "—"}</span>;
  }
};

const resourceStatusBadge = (status) => {
  switch (status) {
    case "available": return <span className="badge badge-approved">Disponible</span>;
    case "unavailable": return <span className="badge badge-rejected">Indisponible</span>;
    default: return <span className="badge badge-pending">{status}</span>;
  }
};

export default function AdminPlaces() {
  const [places, setPlaces] = useState([]);
  const [loading, setLoading] = useState(true);
  const [expanded, setExpanded] = useState({});

  const load = async () => {
    try {
      const res = await getAdminPlaces();
      setPlaces(res.data || []);
    } catch (err) {
      console.error(err);
      setPlaces([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { load(); }, []);

  const toggle = (id) => setExpanded((prev) => ({ ...prev, [id]: !prev[id] }));

  const handleDelete = async (id, name) => {
    if (!window.confirm(`Supprimer le lieu "${name}" ?`)) return;
    try {
      await deleteAdminPlace(id);
      load();
    } catch (err) {
      alert(err.response?.data?.error || "Erreur lors de la suppression");
    }
  };

  if (loading) {
    return <div className="page-container"><p style={{ textAlign: "center" }}>Chargement...</p></div>;
  }

  return (
    <div className="page-container">
      <div style={{ display: "flex", alignItems: "center", justifyContent: "space-between", marginBottom: "32px" }}>
        <h1 className="page-title" style={{ margin: 0 }}>Lieux</h1>
        <Link to="/admin/places/create" className="btn btn-primary" style={{ textDecoration: "none" }}>
          + Nouveau lieu
        </Link>
      </div>

      {places.length === 0 ? (
        <p style={{ textAlign: "center", color: "var(--text-gray)" }}>Aucun lieu</p>
      ) : (
        <div style={{ display: "flex", flexDirection: "column", gap: "16px" }}>
          {places.map((place) => (
            <div key={place.id} className="card">
              {/* En-tête du lieu */}
              <div
                className="card-header"
                style={{ cursor: "pointer" }}
                onClick={() => toggle(place.id)}
              >
                <div>
                  <span style={{ fontSize: "20px", fontWeight: 700 }}>
                    {expanded[place.id] ? "▾" : "▸"} {place.name}
                  </span>
                  {place.description && (
                    <p style={{ margin: "4px 0 0", color: "var(--text-gray)", fontSize: "14px" }}>
                      {place.description}
                    </p>
                  )}
                </div>
                <div style={{ display: "flex", gap: "8px", alignItems: "center" }} onClick={(e) => e.stopPropagation()}>
                  <span style={{ color: "var(--text-gray)", fontSize: "14px" }}>
                    {place.users?.length || 0} utilisateur{(place.users?.length || 0) !== 1 ? "s" : ""}
                    {" · "}
                    {place.resources?.length || 0} ressource{(place.resources?.length || 0) !== 1 ? "s" : ""}
                  </span>
                  <button
                    className="btn btn-danger"
                    style={{ padding: "8px 16px", fontSize: "14px" }}
                    onClick={() => handleDelete(place.id, place.name)}
                  >
                    Supprimer
                  </button>
                </div>
              </div>

              {/* Détails dépliables */}
              {expanded[place.id] && (
                <div style={{ padding: "24px", display: "grid", gridTemplateColumns: "1fr 1fr", gap: "32px" }}>
                  {/* Utilisateurs */}
                  <div>
                    <h3 style={{ margin: "0 0 16px", fontSize: "16px", fontWeight: 700 }}>Utilisateurs</h3>
                    {!place.users || place.users.length === 0 ? (
                      <p style={{ color: "var(--text-gray)", fontSize: "14px" }}>Aucun utilisateur</p>
                    ) : (
                      <table style={{ width: "100%", borderCollapse: "collapse" }}>
                        <tbody>
                          {place.users.map((user) => (
                            <tr key={user.id} style={{ borderBottom: "1px solid #e5e5e5" }}>
                              <td style={{ padding: "10px 0", fontWeight: 600 }}>{user.username}</td>
                              <td style={{ padding: "10px 0", color: "var(--text-gray)", fontSize: "13px" }}>
                                {user.email}
                              </td>
                              <td style={{ padding: "10px 0", textAlign: "center" }}>
                                <span className={`badge ${user.role === "admin" ? "badge-approved" : "badge-pending"}`}>
                                  {user.role === "admin" ? "Admin" : "Utilisateur"}
                                </span>
                              </td>
                              <td style={{ padding: "10px 0", textAlign: "right" }}>
                                {statusBadge(user.status)}
                              </td>
                            </tr>
                          ))}
                        </tbody>
                      </table>
                    )}
                  </div>

                  {/* Ressources */}
                  <div>
                    <h3 style={{ margin: "0 0 16px", fontSize: "16px", fontWeight: 700 }}>Ressources</h3>
                    {!place.resources || place.resources.length === 0 ? (
                      <p style={{ color: "var(--text-gray)", fontSize: "14px" }}>Aucune ressource</p>
                    ) : (
                      <table style={{ width: "100%", borderCollapse: "collapse" }}>
                        <tbody>
                          {place.resources.map((resource) => (
                            <tr key={resource.id} style={{ borderBottom: "1px solid #e5e5e5" }}>
                              <td style={{ padding: "10px 0", fontWeight: 600 }}>{resource.name}</td>
                              <td style={{ padding: "10px 0", textAlign: "right" }}>
                                {resourceStatusBadge(resource.status)}
                              </td>
                            </tr>
                          ))}
                        </tbody>
                      </table>
                    )}
                  </div>
                </div>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
