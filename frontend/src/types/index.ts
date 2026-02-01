export interface Vacancy {
  id: number;
  title: string;
  companyName: string;
  hasApplied: boolean;
  status?: 'SENT' | 'REJECTED';
}
