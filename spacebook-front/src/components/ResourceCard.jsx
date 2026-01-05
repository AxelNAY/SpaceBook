import { Card, CardContent, Typography } from "@mui/material";

export default function ResourceCard({ resource }) {
  return (
    <Card>
      <CardContent>
        <Typography variant="h6">{resource.Name}</Typography>
        <Typography>Type: {resource.Type}</Typography>
        <Typography>Category: {resource.Category}</Typography>
        <Typography>Capacity: {resource.Capacity}</Typography>
      </CardContent>
    </Card>
  );
}
