import { useEffect, useState } from "react";
import { getResources, createReservation } from "../api/api";
import {
  Container, Typography, Select, MenuItem,
  TextField, Button, FormControl, InputLabel
} from "@mui/material";

export default function Reservations() {
  const [resources, setResources] = useState([]);
  const [form, setForm] = useState({
    userID: "9b8e5a0a-7c0f-4f3c-9c3a-0e6e3e5f8b91",
    resourceID: "",
    startAt: "",
    endAt: "",
  });

  useEffect(() => {
    getResources().then(res => setResources(res.data));
  }, []);

  const submit = () => {
    createReservation(form)
      .then(() => alert("Reservation created"))
      .catch(err => alert(err.response?.data?.error));
  };

  return (
    <Container sx={{ mt: 4 }}>
      <Typography variant="h4">New Reservation</Typography>

      <FormControl fullWidth sx={{ mt: 2 }}>
        <InputLabel>Resource</InputLabel>
        <Select
          value={form.resourceID}
          label="Resource"
          onChange={e => setForm({ ...form, resourceID: e.target.value })}
        >
          {resources.map(r => (
            <MenuItem key={r.ID} value={r.ID}>{r.Name}</MenuItem>
          ))}
        </Select>
      </FormControl>

      <TextField
        type="datetime-local"
        fullWidth
        sx={{ mt: 2 }}
        onChange={e => setForm({ ...form, startAt: e.target.value })}
      />

      <TextField
        type="datetime-local"
        fullWidth
        sx={{ mt: 2 }}
        onChange={e => setForm({ ...form, endAt: e.target.value })}
      />

      <Button variant="contained" sx={{ mt: 3 }} onClick={submit}>
        Create reservation
      </Button>
    </Container>
  );
}
