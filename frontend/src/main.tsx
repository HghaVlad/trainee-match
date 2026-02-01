import ReactDOM from 'react-dom/client';
import { ThemeProvider, CssBaseline } from '@mui/material';
import { darkTheme } from './theme';
import { BrowserRouter } from 'react-router-dom';
import App from './App';

import { worker } from './mocks/browser';

// Запуск mock-сервера только в dev
if (import.meta.env.DEV) {
  console.log("Start in development.")
  await worker.start({
    onUnhandledRequest: 'warn', // предупреждает о незамоканных запросах
  });
}

ReactDOM.createRoot(document.getElementById('root')!).render(
  <ThemeProvider theme={darkTheme}>
    <CssBaseline />
    <BrowserRouter>
      <App />
    </BrowserRouter>
  </ThemeProvider>
);
