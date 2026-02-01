import { createTheme } from '@mui/material/styles';

export const darkTheme = createTheme({
  palette: {
    mode: 'dark',
    background: {
      default: '#000',
      paper: '#111',
    },
    primary: {
      main: '#ffffff',
    },
  },
});
