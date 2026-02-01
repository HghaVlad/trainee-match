import { Paper, Box, Typography } from '@mui/material';

export default function Filters() {
  return (
    <Box display="flex" flexDirection="column" gap={2}>
      <Paper variant="outlined" sx={{ p: 2, height: 140 }}>
        <Typography>Первый фильтр</Typography>
      </Paper>

      <Paper variant="outlined" sx={{ p: 2, height: 140 }}>
        <Typography>Второй фильтр</Typography>
      </Paper>
    </Box>
  );
}
