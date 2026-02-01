import {
  Box,
  Typography,
  TextField,
  Pagination,
} from '@mui/material';
import { useEffect, useState } from 'react';
import Filters from '../../components/Filters';
import InternshipCard from '../../components/VacancyCard';
import { internshipFilters } from '../../config/vacanciesFilters';
import { applyToVacancy, getVacancies, withdrawApplication } from '../../api/vacancies.api';

export default function VacanciesListPage() {
  const [vacancies, setVacancies] = useState<any[]>([]);
  const [filters, setFilters] = useState<Record<string, string>>({});
  const [search, setSearch] = useState('');
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);

  const loadVacancies = async () => {
    const res = await getVacancies({
      ...filters,
      search,
      page: page,
      size: 5,
    });

    setVacancies(res.data.content);
    setTotalPages(res.data.totalPages);
  };

  useEffect(() => {
    loadVacancies();
  }, [filters, search, page]);

  const handleFilterChange = (key: string, value: string) => {
    setPage(1);
    setFilters((prev) => ({ ...prev, [key]: value }));
  };

  const handleApply = async (id: number) => {
    await applyToVacancy(id);
    loadVacancies();
  };

  const handleWithdraw = async (id: number) => {
    await withdrawApplication(id);
    loadVacancies();
  };

  return (
    <Box p={4}>
      <Typography variant="h5" mb={3}>
        Список вакансий
      </Typography>

      {/* Поиск */}
      <TextField
        fullWidth
        placeholder="Поиск по ключевым словам"
        value={search}
        onChange={(e) => {
          setPage(1);
          setSearch(e.target.value);
        }}
        sx={{ mb: 3 }}
      />

      <Box display="flex" gap={4}>
        {/* Фильтры */}
        <Box width={280}>
          <Filters
            filters={internshipFilters}
            values={filters}
            onChange={handleFilterChange}
          />
        </Box>

        {/* Список */}
        <Box flex={1}>
          {vacancies.map((i) => (
            <InternshipCard
              key={i.id}
              title={i.title}
              company={i.companyName}
              status={i.status}
              onApply={() => handleApply(i.id)}
              onWithdraw={() => handleWithdraw(i.id)}
            />
          ))}
          {/* Пагинация */}
          <Box display="flex" justifyContent="center" mt={3}>
            <Pagination
              count={totalPages}
              page={page}
              onChange={(_, value) => setPage(value)}
            />
          </Box>
        </Box>
      </Box>
    </Box>
  );
}
