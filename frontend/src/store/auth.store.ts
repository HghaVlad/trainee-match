import { create } from 'zustand';

interface AuthState {
  isAuth: boolean;
  role: 'Candidate' | 'Company' | null;
  setAuth: (role: AuthState['role']) => void;
  logout: () => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  isAuth: false,
  role: null,
  setAuth: (role) => set({ isAuth: true, role }),
  logout: () => set({ isAuth: false, role: null }),
}));
