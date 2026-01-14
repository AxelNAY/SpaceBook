export default function ResourceCard({ resource, onReserve, showStatus = true }) {
  const getStatusBadge = (status) => {
    switch (status) {
      case "available":
        return <span className="badge badge-available">Disponible</span>;
      case "unavailable":
        return <span className="badge badge-inactive">Indisponible</span>;
      default:
        return <span className="badge badge-available">Disponible</span>;
    }
  };

  return (
    <div className="card">
      <div className="card-header">
        <span style={{ fontWeight: 600, fontSize: "18px" }}>
          {resource.Name}
        </span>

        <div style={{ display: "flex", alignItems: "center", gap: "16px" }}>
          <button className="btn btn-primary" style={{ padding: "8px 24px" }}>
            Description
          </button>

          {showStatus && getStatusBadge(resource.Status)}

          {onReserve && (
            <button
              className="btn btn-primary"
              onClick={onReserve}
              style={{ padding: "8px 24px" }}
            >
              Réserver
            </button>
          )}
        </div>
      </div>

      <div className="card-image">
        {resource.ImageUrl ? (
          <img
            src={resource.ImageUrl}
            alt={resource.Name}
            style={{ maxWidth: "100%", maxHeight: "100%" }}
          />
        ) : (
          <span style={{ color: "var(--text-gray)" }}>Image</span>
        )}
      </div>
    </div>
  );
}
