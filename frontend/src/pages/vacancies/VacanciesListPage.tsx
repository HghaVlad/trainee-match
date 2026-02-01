import { Box, Typography } from '@mui/material';
import { useEffect, useState } from 'react';
import Filters from '../../components/Filters';
import VacancyCard from '../../components/VacancyCard';
import {
  getVacancies,
  applyToVacancy,
  withdrawApplication,
} from '../../api/vacancies.api';

interface Vacancy {
  id: number;
  title: string;
  companyName: string;
  status?: 'SENT' | 'REJECTED';
}

export default function VacanciesListPage() {
  const [vacancies, setVacancies] = useState<Vacancy[]>([]);

  const loadVacancies = async () => {
    const res = await getVacancies();
    console.log(res.data);
    setVacancies(res.data);
  };

  useEffect(() => {
    loadVacancies();
  }, []);

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

      <Box display="flex" gap={4}>
        {/* Левая колонка */}
        <Box width={260}>
          <Filters />
        </Box>

        {/* Правая колонка */}
        <Box flex={1}>
          {vacancies.map((i) => (
            <VacancyCard
              key={i.id}
              title={i.title}
              company={i.companyName}
              status={i.status}
              onApply={() => handleApply(i.id)}
              onWithdraw={() => handleWithdraw(i.id)}
            />
          ))}
        </Box>
      </Box>
    </Box>
  );
}
