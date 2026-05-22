import { useEffect, useState } from "react";
import { useAuth } from "../context/AuthContext";
import { getUserNotifications, getAdminNotifications, markNotificationRead } from "../api/api";

export default function Notifications() {
  const { user, isAdmin } = useAuth();
  const [notifications, setNotifications] = useState([]);
  const [loading, setLoading] = useState(true);

  const load = async () => {
    try {
      let res;
      if (isAdmin) {
        res = await getAdminNotifications();
      } else {
        res = await getUserNotifications();
      }
      setNotifications(res.data || []);
    } catch (err) {
      console.error(err);
      setNotifications([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (user) {
      load();
    }
  }, [user, isAdmin]);

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

  const getTypeBadge = (type) => {
    switch (type) {
      case "reservation":
        return <span className="badge badge-pending">Reservation</span>;
      case "resource":
        return <span className="badge badge-available">Ressource</span>;
      default:
        return <span className="badge badge-pending">Info</span>;
    }
  };

  const handleMarkAsRead = async (id) => {
    try {
      await markNotificationRead(id);
      load();
    } catch (err) {
      console.error(err);
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
      <h1 className="page-title">Notifications</h1>

      {notifications.length === 0 ? (
        <p style={{ textAlign: "center", color: "var(--text-gray)" }}>
          Aucune notification
        </p>
      ) : (
        <div className="card" style={{ overflow: "hidden" }}>
          <table style={{ width: "100%", borderCollapse: "collapse" }}>
            <thead>
              <tr style={{ backgroundColor: "var(--primary-blue)", color: "white" }}>
                <th style={{ padding: "12px 16px", textAlign: "left", fontWeight: 600, whiteSpace: "nowrap" }}>
                  Type
                </th>
                <th style={{ padding: "12px 16px", textAlign: "left", fontWeight: 600, width: "100%" }}>
                  Message
                </th>
                <th style={{ padding: "12px 16px", textAlign: "left", fontWeight: 600, whiteSpace: "nowrap" }}>
                  Date
                </th>
                <th style={{ padding: "12px 16px", textAlign: "center", fontWeight: 600, whiteSpace: "nowrap" }}>
                  Statut
                </th>
                {isAdmin && (
                  <th style={{ padding: "12px 16px", textAlign: "center", fontWeight: 600, whiteSpace: "nowrap" }}>
                    Action
                  </th>
                )}
              </tr>
            </thead>
            <tbody>
              {notifications.map((notification) => (
                <tr
                  key={notification.id}
                  style={{
                    borderBottom: "1px solid #e5e5e5",
                    backgroundColor: notification.is_read ? "transparent" : "#f0f7ff",
                  }}
                >
                  <td style={{ padding: "12px 16px", whiteSpace: "nowrap" }}>
                    {getTypeBadge(notification.type)}
                  </td>
                  <td style={{ padding: "12px 16px", fontWeight: notification.is_read ? 400 : 600 }}>
                    {notification.message}
                  </td>
                  <td style={{ padding: "12px 16px", fontSize: "13px", color: "var(--text-gray)", whiteSpace: "nowrap" }}>
                    {formatDate(notification.created_at)}
                  </td>
                  <td style={{ padding: "12px 16px", textAlign: "center", whiteSpace: "nowrap" }}>
                    {notification.is_read ? (
                      <span className="badge badge-available">Lu</span>
                    ) : (
                      <span className="badge badge-pending">Non lu</span>
                    )}
                  </td>
                  {isAdmin && (
                    <td style={{ padding: "12px 16px", textAlign: "center", whiteSpace: "nowrap" }}>
                      {!notification.is_read && (
                        <button
                          className="btn btn-primary"
                          style={{ padding: "6px 12px", fontSize: "13px" }}
                          onClick={() => handleMarkAsRead(notification.id)}
                        >
                          Marquer lu
                        </button>
                      )}
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
