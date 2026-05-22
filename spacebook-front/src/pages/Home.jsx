import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { getResources } from "../api/api";
import { useAuth } from "../context/AuthContext";

export default function Home() {
  const { isAuthenticated } = useAuth();
  const navigate = useNavigate();
  const [resources, setResources] = useState([]);
  const [loading, setLoading] = useState(true);
  const [filterStatus, setFilterStatus] = useState("");

  useEffect(() => {
    getResources()
      .then((res) => setResources(res.data || []))
      .finally(() => setLoading(false));
  }, []);

  const filteredResources = resources.filter((r) => {
    if (filterStatus && r.status !== filterStatus) return false;
    return true;
  });

  const handleReserve = (resource) => {
    navigate("/reserver", { state: { resource } });
  };

  const getStatusBadge = (status) => {
    switch (status) {
      case "available":
        return <span className="badge badge-available">Disponible</span>;
      case "unavailable":
        return <span className="badge badge-rejected">Indisponible</span>;
      default:
        return <span className="badge badge-available">Disponible</span>;
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
      <h1 className="page-title">Ressources</h1>

      <div className="filter-bar">
        <select
          className="filter-select"
          value={filterStatus}
          onChange={(e) => setFilterStatus(e.target.value)}
        >
          <option value="">Toutes disponibilites</option>
          <option value="available">Disponible</option>
          <option value="unavailable">Indisponible</option>
        </select>
      </div>

      {filteredResources.length === 0 ? (
        <p style={{ textAlign: "center", color: "var(--text-gray)" }}>
          Aucune ressource trouvee
        </p>
      ) : (
        <div className="card" style={{ overflow: "hidden" }}>
          <table style={{ width: "100%", borderCollapse: "collapse" }}>
            <thead>
              <tr style={{ backgroundColor: "var(--primary-blue)", color: "white" }}>
                <th style={{ padding: "16px 24px", textAlign: "left", fontWeight: 600 }}>Nom</th>
                <th style={{ padding: "16px 24px", textAlign: "left", fontWeight: 600 }}>Lieu</th>
                <th style={{ padding: "16px 24px", textAlign: "left", fontWeight: 600 }}>Catégorie</th>
                <th style={{ padding: "16px 24px", textAlign: "center", fontWeight: 600 }}>Statut</th>
                {isAuthenticated && (
                  <th style={{ padding: "16px 24px", textAlign: "center", fontWeight: 600 }}>Action</th>
                )}
              </tr>
            </thead>
            <tbody>
              {filteredResources.map((resource) => (
                <tr key={resource.id} style={{ borderBottom: "1px solid #e5e5e5" }}>
                  <td style={{ padding: "16px 24px", fontWeight: 600 }}>{resource.name}</td>
                  <td style={{ padding: "16px 24px" }}>{resource.place?.name || "-"}</td>
                  <td style={{ padding: "16px 24px" }}>{resource.category?.name || "-"}</td>
                  <td style={{ padding: "16px 24px", textAlign: "center" }}>
                    {getStatusBadge(resource.status)}
                  </td>
                  {isAuthenticated && (
                    <td style={{ padding: "16px 24px", textAlign: "center" }}>
                      <button
                        className="btn btn-primary"
                        style={{ padding: "8px 20px", fontSize: "14px" }}
                        onClick={() => handleReserve(resource)}
                      >
                        Reserver
                      </button>
                    </td>
                  )}
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
