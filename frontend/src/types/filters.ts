export type FilterType = 'select' | 'text';

export interface FilterConfig {
  key: string;
  label: string;
  type: FilterType;
  options?: { label: string; value: string }[];
}
