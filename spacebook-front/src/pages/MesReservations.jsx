import { useEffect, useState } from "react";
import { useAuth } from "../context/AuthContext";
import { getUserReservations } from "../api/api";

export default function MesReservations() {
  const { user } = useAuth();
  const [reservations, setReservations] = useState([]);
  const [loading, setLoading] = useState(true);

  const load = async () => {
    if (!user?.id) return;

    try {
      const res = await getUserReservations(user.id);
      setReservations(res.data || []);
    } catch (err) {
      console.error(err);
      setReservations([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    load();
  }, [user]);

  const formatDate = (dateString) => {
    const date = new Date(dateString);
    return date.toLocaleDateString("fr-FR", {
      day: "2-digit",
      month: "2-digit",
      year: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  const getStatusBadge = (status) => {
    switch (status) {
      case "approved":
        return <span className="badge badge-approved">Approuvee</span>;
      case "rejected":
        return <span className="badge badge-rejected">Refusee</span>;
      default:
        return <span className="badge badge-pending">En attente</span>;
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
      <h1 className="page-title">Mes Reservations</h1>

      {reservations.length === 0 ? (
        <p style={{ textAlign: "center", color: "var(--text-gray)" }}>
          Aucune reservation
        </p>
      ) : (
        <div className="card" style={{ overflow: "hidden" }}>
          <table style={{ width: "100%", borderCollapse: "collapse" }}>
            <thead>
              <tr style={{ backgroundColor: "var(--primary-blue)", color: "white" }}>
                <th style={{ padding: "16px 24px", textAlign: "left", fontWeight: 600 }}>
                  Ressource
                </th>
                <th style={{ padding: "16px 24px", textAlign: "left", fontWeight: 600 }}>
                  Type
                </th>
                <th style={{ padding: "16px 24px", textAlign: "left", fontWeight: 600 }}>
                  Horaire
                </th>
                <th style={{ padding: "16px 24px", textAlign: "center", fontWeight: 600 }}>
                  Statut
                </th>
              </tr>
            </thead>
            <tbody>
              {reservations.map((reservation) => (
                <tr
                  key={reservation.id}
                  style={{ borderBottom: "1px solid #e5e5e5" }}
                >
                  <td style={{ padding: "16px 24px", fontWeight: 600 }}>
                    {reservation.resource?.Name || "Ressource"}
                  </td>
                  <td style={{ padding: "16px 24px" }}>
                    {reservation.resource?.Type === "room" ? "Salle" : "Equipement"}
                  </td>
                  <td style={{ padding: "16px 24px", fontSize: "14px" }}>
                    <div>{formatDate(reservation.start_at)}</div>
                    <div style={{ color: "var(--text-gray)" }}>au {formatDate(reservation.end_at)}</div>
                  </td>
                  <td style={{ padding: "16px 24px", textAlign: "center" }}>
                    {getStatusBadge(reservation.status)}
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
