import { useEffect, useState } from "react";
import {
  getAdminNotifications,
  markNotificationRead
} from "../api/api";
import {
  Container, Typography, List, ListItem,
  ListItemText, Button
} from "@mui/material";

export default function AdminNotifications() {
  const [notes, setNotes] = useState([]);

  const load = () =>
    getAdminNotifications().then(res => setNotes(res.data));

  useEffect(() => { load(); }, []);

  return (
    <Container sx={{ mt: 4 }}>
      <Typography variant="h4">Admin Notifications</Typography>
      <List>
        {notes.map(n => (
          <ListItem key={n.ID} divider>
            <ListItemText
              primary={n.Message}
              secondary={n.IsRead ? "Read" : "Unread"}
            />
            {!n.IsRead && (
              <Button onClick={() => markNotificationRead(n.ID).then(load)}>
                Mark as read
              </Button>
            )}
          </ListItem>
        ))}
      </List>
    </Container>
  );
}
