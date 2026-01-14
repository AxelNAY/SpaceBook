import { useEffect, useState } from "react";
import { getAdminUsers, deleteAdminUser } from "../api/api";

export default function AdminUsers() {
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);

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

  const getRoleBadge = (role) => {
    if (role === "admin") {
      return <span className="badge badge-pending">Admin</span>;
    }
    return <span className="badge badge-available">User</span>;
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

      {users.length === 0 ? (
        <p style={{ textAlign: "center", color: "var(--text-gray)" }}>
          Aucun utilisateur
        </p>
      ) : (
        <div className="card" style={{ overflow: "hidden" }}>
          <table style={{ width: "100%", borderCollapse: "collapse" }}>
            <thead>
              <tr style={{ backgroundColor: "var(--primary-blue)", color: "white" }}>
                <th style={{ padding: "16px 24px", textAlign: "left", fontWeight: 600 }}>
                  Nom
                </th>
                <th style={{ padding: "16px 24px", textAlign: "left", fontWeight: 600 }}>
                  Email
                </th>
                <th style={{ padding: "16px 24px", textAlign: "center", fontWeight: 600 }}>
                  Role
                </th>
                <th style={{ padding: "16px 24px", textAlign: "center", fontWeight: 600 }}>
                  Action
                </th>
              </tr>
            </thead>
            <tbody>
              {users.map((user) => (
                <tr
                  key={user.id}
                  style={{ borderBottom: "1px solid #e5e5e5" }}
                >
                  <td style={{ padding: "16px 24px", fontWeight: 600 }}>
                    {user.username}
                  </td>
                  <td style={{ padding: "16px 24px" }}>
                    {user.email}
                  </td>
                  <td style={{ padding: "16px 24px", textAlign: "center" }}>
                    {getRoleBadge(user.role)}
                  </td>
                  <td style={{ padding: "16px 24px", textAlign: "center" }}>
                    <button
                      className="btn btn-danger"
                      style={{ padding: "8px 20px", fontSize: "14px" }}
                      onClick={() => handleDelete(user.id, user.username)}
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
