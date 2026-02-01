import { api } from './axios';

export interface LoginRequest {
  username: string;
  password: string;
}

export interface RegisterRequest {
  first_name: string;
  last_name: string;
  email: string;
  username: string;
  password: string;
  role: 'Candidate' | 'Company';
}

export const loginApi = (data: LoginRequest) =>
  api.post('/auth/login', data);

export const registerApi = (data: RegisterRequest) =>
  api.post('/auth/register', data);

export const logoutApi = () =>
  api.post('/auth/logout');

export const refreshApi = () =>
  api.post('/auth/refresh');
