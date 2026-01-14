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
        res = await getUserNotifications(user?.id);
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
                <th style={{ padding: "16px 24px", textAlign: "left", fontWeight: 600 }}>
                  Type
                </th>
                <th style={{ padding: "16px 24px", textAlign: "left", fontWeight: 600 }}>
                  Message
                </th>
                <th style={{ padding: "16px 24px", textAlign: "left", fontWeight: 600 }}>
                  Date
                </th>
                <th style={{ padding: "16px 24px", textAlign: "center", fontWeight: 600 }}>
                  Statut
                </th>
                {isAdmin && (
                  <th style={{ padding: "16px 24px", textAlign: "center", fontWeight: 600 }}>
                    Action
                  </th>
                )}
              </tr>
            </thead>
            <tbody>
              {notifications.map((notification) => (
                <tr
                  key={notification.ID}
                  style={{
                    borderBottom: "1px solid #e5e5e5",
                    backgroundColor: notification.is_read ? "transparent" : "#f0f7ff",
                  }}
                >
                  <td style={{ padding: "16px 24px" }}>
                    {getTypeBadge(notification.type)}
                  </td>
                  <td style={{ padding: "16px 24px", fontWeight: notification.is_read ? 400 : 600 }}>
                    {notification.message}
                  </td>
                  <td style={{ padding: "16px 24px", fontSize: "14px", color: "var(--text-gray)" }}>
                    {formatDate(notification.created_at)}
                  </td>
                  <td style={{ padding: "16px 24px", textAlign: "center" }}>
                    {notification.is_read ? (
                      <span className="badge badge-available">Lu</span>
                    ) : (
                      <span className="badge badge-pending">Non lu</span>
                    )}
                  </td>
                  {isAdmin && (
                    <td style={{ padding: "16px 24px", textAlign: "center" }}>
                      {!notification.is_read && (
                        <button
                          className="btn btn-primary"
                          style={{ padding: "8px 16px", fontSize: "14px" }}
                          onClick={() => handleMarkAsRead(notification.ID)}
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
