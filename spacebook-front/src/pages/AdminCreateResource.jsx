import { useState } from "react";
import {
  Container, Typography, TextField, Button,
  Select, MenuItem, FormControl, InputLabel
} from "@mui/material";
import { createAdminResource } from "../api/api";

export default function AdminCreateResource() {
  const [form, setForm] = useState({
    name: "",
    type: "room",
    category: "none",
    capacity: 1,
  });

  const submit = () => {
    createAdminResource(form)
      .then(() => alert("Resource created"))
      .catch(e => alert(e.response?.data?.error));
  };

  return (
    <Container sx={{ mt: 4 }}>
      <Typography variant="h4">Create Resource</Typography>

      <TextField
        label="Name"
        fullWidth sx={{ mt: 2 }}
        onChange={e => setForm({ ...form, name: e.target.value })}
      />

      <FormControl fullWidth sx={{ mt: 2 }}>
        <InputLabel>Type</InputLabel>
        <Select
          value={form.type}
          label="Type"
          onChange={e => setForm({ ...form, type: e.target.value })}
        >
          <MenuItem value="room">Room</MenuItem>
          <MenuItem value="equipment">Equipment</MenuItem>
        </Select>
      </FormControl>

      {form.type === "equipment" && (
        <>
          <TextField
            label="Category"
            fullWidth sx={{ mt: 2 }}
            onChange={e => setForm({ ...form, category: e.target.value })}
          />
          <TextField
            label="Capacity"
            type="number"
            fullWidth sx={{ mt: 2 }}
            onChange={e => setForm({ ...form, capacity: Number(e.target.value) })}
          />
        </>
      )}

      <Button variant="contained" sx={{ mt: 3 }} onClick={submit}>
        Create
      </Button>
    </Container>
  );
}
