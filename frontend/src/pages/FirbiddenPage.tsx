import { Typography, Box } from '@mui/material';

export default function ForbiddenPage() {
  return (
    <Box height="100vh" display="flex" alignItems="center" justifyContent="center">
      <Typography variant="h4">403 — Нет доступа</Typography>
    </Box>
  );
}
