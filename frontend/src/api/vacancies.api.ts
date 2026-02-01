import { api } from './axios';

export const getVacancies = (params?: {
  search?: string;
  city?: string;
  format?: string;
}) => api.get('/vacancies', { params });

export const applyToVacancy = (id: number) =>
  api.post(`/applications/${id}/apply`);

export const withdrawApplication = (id: number) =>
  api.delete(`/applications/${id}/withdraw`);
