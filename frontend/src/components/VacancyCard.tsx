import {
  Box,
  Button,
  Typography,
  Paper,
} from '@mui/material';

interface Props {
  title: string;
  company: string;
  status?: 'SENT' | 'REJECTED';
  onApply: () => void;
  onWithdraw: () => void;
}

export default function VacancyCard({
  title,
  company,
  status,
  onApply,
  onWithdraw,
}: Props) {
  const renderButton = () => {
    if (status === 'SENT') {
      return (<Button variant="outlined" onClick={onWithdraw}>
          Отозвать
        </Button>);
    }

    if (status === 'REJECTED') {
      return (
        <Button variant="outlined" color="error" disabled>
          Отказано
        </Button>
      );
    }

    return (
      <Button variant="outlined" onClick={onApply}>
        Подать заявку
      </Button>
    );
  };

  return (
    <Paper
      variant="outlined"
      sx={{
        p: 2,
        mb: 2,
        borderRadius: 2,
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center',
      }}
    >
      <Box>
        <Typography variant="subtitle1">
          {title}
        </Typography>
        <Typography variant="body2" color="gray">
          {company}
        </Typography>
      </Box>

      {renderButton()}
    </Paper>
  );
}
