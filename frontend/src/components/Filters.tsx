import {
  Box,
  Paper,
  TextField,
  MenuItem,
  Typography,
} from '@mui/material';
import type { FilterConfig } from '../types/filters';

interface Props {
  filters: FilterConfig[];
  values: Record<string, string>;
  onChange: (key: string, value: string) => void;
}

export default function Filters({ filters, values, onChange }: Props) {
  return (
    <Paper variant="outlined" sx={{ p: 2 }}>
      <Typography mb={2}>Фильтры</Typography>

      <Box display="flex" flexDirection="column" gap={2}>
        {filters.map((filter) => {
          if (filter.type === 'select') {
            return (
              <TextField
                key={filter.key}
                select
                label={filter.label}
                value={values[filter.key] || ''}
                onChange={(e) => onChange(filter.key, e.target.value)}
              >
                <MenuItem value="">Любой</MenuItem>
                {filter.options?.map((opt) => (
                  <MenuItem key={opt.value} value={opt.value}>
                    {opt.label}
                  </MenuItem>
                ))}
              </TextField>
            );
          }

          return (
            <TextField
              key={filter.key}
              label={filter.label}
              value={values[filter.key] || ''}
              onChange={(e) => onChange(filter.key, e.target.value)}
            />
          );
        })}
      </Box>
    </Paper>
  );
}
