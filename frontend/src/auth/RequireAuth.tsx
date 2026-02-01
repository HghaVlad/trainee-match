import { Navigate } from 'react-router-dom';
import { useAuthStore } from '../store/auth.store';
import type { JSX } from 'react';

export default function RequireAuth({
  children,
  allowedRoles,
}: {
  children: JSX.Element;
  allowedRoles?: Array<'Candidate' | 'Company'>;
}) {
  const { isAuth, role } = useAuthStore();

  if (!isAuth) return <Navigate to="/login" />;

  if (allowedRoles && role && !allowedRoles.includes(role)) {
    return <Navigate to="/403" />;
  }

  return children;
}
