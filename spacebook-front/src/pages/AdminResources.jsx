import { useEffect, useState } from "react";
import { getAdminResources, deleteAdminResource } from "../api/api";

export default function AdminResources() {
  const [resources, setResources] = useState([]);
  const [loading, setLoading] = useState(true);

  const load = async () => {
    try {
      const res = await getAdminResources();
      setResources(res.data || []);
    } catch (err) {
      console.error(err);
      setResources([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    load();
  }, []);

  const handleDelete = async (resourceId, resourceName) => {
    if (!window.confirm(`Supprimer la ressource "${resourceName}" ?`)) return;

    try {
      await deleteAdminResource(resourceId);
      load();
    } catch (err) {
      alert(err.response?.data?.error || "Échec de la suppression");
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

      {resources.length === 0 ? (
        <p style={{ textAlign: "center", color: "var(--text-gray)" }}>
          Aucune ressource
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
                <th style={{ padding: "16px 24px", textAlign: "center", fontWeight: 600 }}>Action</th>
              </tr>
            </thead>
            <tbody>
              {resources.map((resource) => (
                <tr key={resource.id} style={{ borderBottom: "1px solid #e5e5e5" }}>
                  <td style={{ padding: "16px 24px", fontWeight: 600 }}>{resource.name}</td>
                  <td style={{ padding: "16px 24px" }}>{resource.place?.name || "-"}</td>
                  <td style={{ padding: "16px 24px" }}>{resource.category?.name || "-"}</td>
                  <td style={{ padding: "16px 24px", textAlign: "center" }}>{resource.status}</td>
                  <td style={{ padding: "16px 24px", textAlign: "center" }}>
                    <button
                      className="btn btn-danger"
                      style={{ padding: "8px 20px", fontSize: "14px" }}
                      onClick={() => handleDelete(resource.id, resource.name)}
                    >
                      Supprimer
                    </button>
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
