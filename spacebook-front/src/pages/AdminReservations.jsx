import { useEffect, useState } from "react";
import { getAdminReservations, approveReservation } from "../api/api.js";
import { Container, Typography, Button, Box, Paper } from "@mui/material";

export default function AdminReservations() {
  const [reservations, setReservations] = useState([]);
  const [loading, setLoading] = useState(true);

  const fetchReservations = async () => {
    try {
      const res = await getAdminReservations();
      setReservations(res.data ?? []);
    } catch (e) {
      console.error(e);
      setReservations([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchReservations();
  }, []);

  const handleApprove = async (id) => {
    await approveReservation(id);
    fetchReservations();
  };

  if (loading) return <Typography>Loading...</Typography>;

  return (
    <Container>
      <Typography variant="h4" gutterBottom>
        Reservations
      </Typography>

      {reservations.length === 0 && (
        <Typography>No reservations found</Typography>
      )}

      {reservations.map((r) => (
        <Paper key={r.id} sx={{ p: 2, mb: 2 }}>
          <Typography>
            <strong>Resource:</strong> {r.resource?.name || "Unknown"}
          </Typography>

          <Typography>
            <strong>User:</strong> {r.user?.name || "Unknown"}
          </Typography>

          <Typography>
            <strong>Status:</strong> {r.status}
          </Typography>

          <Typography>
            <strong>From:</strong>{" "}
            {r.start_at ? new Date(r.start_at).toLocaleString() : "—"}
          </Typography>

          <Typography>
            <strong>To:</strong>{" "}
            {r.end_at ? new Date(r.end_at).toLocaleString() : "—"}
          </Typography>
        </Paper>
      ))}
    </Container>
  );
}
