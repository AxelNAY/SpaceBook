import { useEffect, useState } from "react";
import { getResources } from "../api/api";
import { Container, Typography, Grid } from "@mui/material";
import ResourceCard from "../components/ResourceCard";

export default function Home() {
  const [resources, setResources] = useState([]);

  useEffect(() => {
    getResources().then(res => setResources(res.data));
  }, []);

  return (
    <Container sx={{ mt: 4 }}>
      <Typography variant="h4" gutterBottom>Resources</Typography>
      <Grid container spacing={2}>
        {resources.map(r => (
          <Grid item xs={12} md={4} key={r.ID}>
            <ResourceCard resource={r} />
          </Grid>
        ))}
      </Grid>
    </Container>
  );
}
