import { Routes, Route } from 'react-router-dom';
import LoginPage from './pages/auth/LoginPage';
import RegisterPage from './pages/auth/RegisterPage';
import RequireAuth from './auth/RequireAuth';
import ForbiddenPage from './pages/FirbiddenPage';

export default function App() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route path="/register" element={<RegisterPage />} />

      <Route
        path="/"
        element={
          <RequireAuth allowedRoles={['Candidate']}>
            <div>Список стажировок</div>
          </RequireAuth>
        }
      />

      <Route
        path="/profile"
        element={
          <RequireAuth allowedRoles={['Candidate']}>
            <div>Профиль кандидата</div>
          </RequireAuth>
        }
      />

      <Route path="/403" element={<ForbiddenPage />} />
    </Routes>
  );
}
