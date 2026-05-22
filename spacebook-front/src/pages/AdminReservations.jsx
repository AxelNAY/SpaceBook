import { useEffect, useState } from "react";
import {
  getAdminReservations,
  approveReservation,
  rejectReservation,
  approveGroupResources,
  rejectGroupResources,
} from "../api/api";

export default function AdminReservations() {
  const [reservations, setReservations] = useState([]);
  const [loading, setLoading] = useState(true);
  const [expanded, setExpanded] = useState({});

  const load = async () => {
    try {
      const res = await getAdminReservations();
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
  }, []);

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString("fr-FR", {
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
        return <span className="badge badge-approved">Approuvé</span>;
      case "rejected":
        return <span className="badge badge-rejected">Refusé</span>;
      default:
        return <span className="badge badge-pending">En attente</span>;
    }
  };

  const isGroup = (r) => !r.ressource_id && r.reservation_resources?.length > 0;

  const toggleExpand = (id) =>
    setExpanded((prev) => ({ ...prev, [id]: !prev[id] }));

  const handleApprove = async (id) => {
    try {
      await approveReservation(id);
      load();
    } catch (err) {
      alert(err.response?.data?.error || "Erreur lors de l'approbation");
    }
  };

  const handleReject = async (id) => {
    try {
      await rejectReservation(id);
      load();
    } catch (err) {
      alert(err.response?.data?.error || "Erreur lors du refus");
    }
  };

  const handleApproveResource = async (reservationId, resourceId) => {
    try {
      await approveGroupResources(reservationId, [resourceId]);
      load();
    } catch (err) {
      alert(err.response?.data?.error || "Erreur lors de l'approbation");
    }
  };

  const handleRejectResource = async (reservationId, resourceId) => {
    try {
      await rejectGroupResources(reservationId, [resourceId]);
      load();
    } catch (err) {
      alert(err.response?.data?.error || "Erreur lors du refus");
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
      <h1 className="page-title">Réservations</h1>

      {reservations.length === 0 ? (
        <p style={{ textAlign: "center", color: "var(--text-gray)" }}>
          Aucune réservation
        </p>
      ) : (
        <div className="card" style={{ overflow: "hidden" }}>
          <table style={{ width: "100%", borderCollapse: "collapse" }}>
            <thead>
              <tr style={{ backgroundColor: "var(--primary-blue)", color: "white" }}>
                <th style={{ padding: "12px 16px", textAlign: "left", fontWeight: 600, width: "100%" }}>Ressource</th>
                <th style={{ padding: "12px 16px", textAlign: "left", fontWeight: 600, whiteSpace: "nowrap" }}>Lieu</th>
                <th style={{ padding: "12px 16px", textAlign: "left", fontWeight: 600, whiteSpace: "nowrap" }}>Horaire</th>
                <th style={{ padding: "12px 16px", textAlign: "left", fontWeight: 600, whiteSpace: "nowrap" }}>Demandé par</th>
                <th style={{ padding: "12px 16px", textAlign: "center", fontWeight: 600, whiteSpace: "nowrap" }}>Statut</th>
                <th style={{ padding: "12px 16px", textAlign: "center", fontWeight: 600, whiteSpace: "nowrap" }}>Actions</th>
              </tr>
            </thead>
            <tbody>
              {reservations.map((reservation) =>
                isGroup(reservation) ? (
                  <>
                    {/* Ligne principale groupe */}
                    <tr
                      key={reservation.id}
                      style={{ borderBottom: "1px solid #e5e5e5", cursor: "pointer" }}
                      onClick={() => toggleExpand(reservation.id)}
                    >
                      <td style={{ padding: "12px 16px", fontWeight: 600 }}>
                        <span style={{ marginRight: "8px" }}>
                          {expanded[reservation.id] ? "▾" : "▸"}
                        </span>
                        Groupe ({reservation.reservation_resources.length} ressources)
                      </td>
                      <td style={{ padding: "12px 16px", color: "var(--text-gray)", whiteSpace: "nowrap" }}>—</td>
                      <td style={{ padding: "12px 16px", fontSize: "14px", whiteSpace: "nowrap" }}>
                        <div>{formatDate(reservation.start_datetime)}</div>
                        <div style={{ color: "var(--text-gray)" }}>
                          au {formatDate(reservation.end_datetime)}
                        </div>
                      </td>
                      <td style={{ padding: "12px 16px", whiteSpace: "nowrap" }}>
                        {reservation.user?.username || "Utilisateur"}
                      </td>
                      <td style={{ padding: "12px 16px", textAlign: "center", whiteSpace: "nowrap" }}>
                        {getStatusBadge(reservation.status)}
                      </td>
                      <td style={{ padding: "12px 16px", textAlign: "center", color: "var(--text-gray)", fontSize: "13px", whiteSpace: "nowrap" }}>
                        Gérer par ressource ↓
                      </td>
                    </tr>

                    {/* Sous-lignes par ressource */}
                    {expanded[reservation.id] &&
                      reservation.reservation_resources.map((rr) => (
                        <tr
                          key={rr.id}
                          style={{
                            borderBottom: "1px solid #f0f0f0",
                            backgroundColor: "#fafafa",
                          }}
                        >
                          <td style={{ padding: "10px 16px 10px 40px", color: "var(--text-gray)" }}>
                            ↳ {rr.resource?.name || "Ressource"}
                          </td>
                          <td style={{ padding: "10px 16px", color: "var(--text-gray)", fontSize: "13px", whiteSpace: "nowrap" }}>
                            {rr.resource?.place?.name || "—"}
                          </td>
                          <td colSpan={2} style={{ padding: "10px 16px" }} />
                          <td style={{ padding: "10px 16px", textAlign: "center", whiteSpace: "nowrap" }}>
                            {getStatusBadge(rr.status)}
                          </td>
                          <td style={{ padding: "10px 16px", textAlign: "center", whiteSpace: "nowrap" }}>
                            <div style={{ display: "flex", gap: "8px", justifyContent: "center" }}>
                              <button
                                className="btn btn-primary"
                                style={{ padding: "6px 12px", fontSize: "13px" }}
                                onClick={(e) => {
                                  e.stopPropagation();
                                  handleApproveResource(reservation.id, rr.resource_id);
                                }}
                                disabled={rr.status === "approved"}
                              >
                                Accepter
                              </button>
                              <button
                                className="btn btn-danger"
                                style={{ padding: "6px 12px", fontSize: "13px" }}
                                onClick={(e) => {
                                  e.stopPropagation();
                                  handleRejectResource(reservation.id, rr.resource_id);
                                }}
                                disabled={rr.status === "rejected"}
                              >
                                Refuser
                              </button>
                            </div>
                          </td>
                        </tr>
                      ))}
                  </>
                ) : (
                  /* Ligne réservation simple */
                  <tr key={reservation.id} style={{ borderBottom: "1px solid #e5e5e5" }}>
                    <td style={{ padding: "12px 16px", fontWeight: 600 }}>
                      {reservation.ressource?.name || "Ressource"}
                    </td>
                    <td style={{ padding: "12px 16px", whiteSpace: "nowrap" }}>
                      {reservation.ressource?.place?.name || "—"}
                    </td>
                    <td style={{ padding: "12px 16px", fontSize: "14px", whiteSpace: "nowrap" }}>
                      <div>{formatDate(reservation.start_datetime)}</div>
                      <div style={{ color: "var(--text-gray)" }}>
                        au {formatDate(reservation.end_datetime)}
                      </div>
                    </td>
                    <td style={{ padding: "12px 16px", whiteSpace: "nowrap" }}>
                      {reservation.user?.username || "Utilisateur"}
                    </td>
                    <td style={{ padding: "12px 16px", textAlign: "center", whiteSpace: "nowrap" }}>
                      {getStatusBadge(reservation.status)}
                    </td>
                    <td style={{ padding: "12px 16px", textAlign: "center", whiteSpace: "nowrap" }}>
                      <div style={{ display: "flex", gap: "8px", justifyContent: "center" }}>
                        <button
                          className="btn btn-primary"
                          style={{ padding: "6px 12px", fontSize: "13px" }}
                          onClick={() => handleApprove(reservation.id)}
                          disabled={reservation.status === "approved"}
                        >
                          Accepter
                        </button>
                        <button
                          className="btn btn-danger"
                          style={{ padding: "6px 12px", fontSize: "13px" }}
                          onClick={() => handleReject(reservation.id)}
                          disabled={reservation.status === "rejected"}
                        >
                          Refuser
                        </button>
                      </div>
                    </td>
                  </tr>
                )
              )}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
