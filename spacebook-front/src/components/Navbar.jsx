import { AppBar, Toolbar, Button, Box } from "@mui/material";
import { Link } from "react-router-dom";

export default function Navbar() {
  return (
    <AppBar position="static">
      <Toolbar sx={{ display: "flex", justifyContent: "space-between" }}>
        
        {/* LEFT – USER */}
        <Box>
          <Button color="inherit" component={Link} to="/">
            Resources
          </Button>

          <Button color="inherit" component={Link} to="/reserve">
            Reserve
          </Button>
        </Box>

        {/* RIGHT – ADMIN */}
        <Box>
          <Button color="inherit" component={Link} to="/admin">
            Admin Notifications
          </Button>

          <Button color="inherit" component={Link} to="/admin/resources">
            Admin Resources
          </Button>

          <Button color="inherit" component={Link} to="/admin/resources/create">
            Create Resource
          </Button>

          <Button color="inherit" component={Link} to="/admin/reservations">
            Admin Reservations
          </Button>
        </Box>

      </Toolbar>
    </AppBar>
  );
}
