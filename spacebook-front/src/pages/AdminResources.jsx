import { useEffect, useState } from "react";
import { getResources, deleteAdminResource } from "../api/api";
import {
  Container,
  Typography,
  List,
  ListItem,
  ListItemText,
  Button,
  Box,
} from "@mui/material";

export default function AdminResources() {
  const [resources, setResources] = useState([]);

  const load = () => {
    getResources().then((res) => setResources(res.data));
  };

  useEffect(() => {
    load();
  }, []);

  const remove = (id) => {
    if (!window.confirm("Delete this resource?")) return;

    deleteAdminResource(id)
      .then(load)
      .catch((e) =>
        alert(e.response?.data?.error || "Failed to delete resource")
      );
  };


  return (
    <Container sx={{ mt: 4 }}>
      <Typography variant="h4" gutterBottom>
        All Resources
      </Typography>

      <List>
        {resources.map((r) => (
          <ListItem
            key={r.ID}
            divider
            secondaryAction={
              <Button
                color="error"
                variant="outlined"
                onClick={() => remove(r.ID)}
              >
                Delete
              </Button>
            }
          >
            <ListItemText
              primary={`${r.Name} (${r.Type})`}
              secondary={`Category: ${r.Category} | Capacity: ${r.Capacity}`}
            />
          </ListItem>
        ))}
      </List>
    </Container>
  );
}
