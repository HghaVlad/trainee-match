import {
  Box,
  Button,
  TextField,
  Typography,
  Paper,
} from '@mui/material';
import { useForm } from 'react-hook-form';
import { loginApi } from '../../api/auth.api';
import { useAuthStore } from '../../store/auth.store';
import { useNavigate } from 'react-router-dom';

interface LoginForm {
  username: string;
  password: string;
}

export default function LoginPage() {
  const { register, handleSubmit } = useForm<LoginForm>();
  const setAuth = useAuthStore((s) => s.setAuth);
  const navigate = useNavigate();

const onSubmit = async (data: LoginForm) => {
  try {
    await loginApi({
      username: data.username,
      password: data.password,
    });

    // cookies уже установлены сервером
    setAuth(null); // роль получим позже
    navigate('/');
  } catch {
    alert('Неверный логин или пароль');
  }
};

  return (
    <Box
      height="100vh"
      display="flex"
      alignItems="center"
      justifyContent="center"
      sx={{ background: 'radial-gradient(circle, #111, #000)' }}
    >
      <Paper sx={{ p: 4, width: 360 }}>
        <Typography variant="h6" textAlign="center" mb={2}>
          Вход
        </Typography>

        <Box component="form" onSubmit={handleSubmit(onSubmit)} display="flex" flexDirection="column" gap={2}>
          <TextField label="Имя пользователя" {...register('username', { required: true })} />
          <TextField type="password" label="Пароль" {...register('password', { required: true })} />

          <Button variant="outlined" type="submit">
            Войти
          </Button>

          <Button variant="text" onClick={() => navigate('/register')}>
            Нет аккаунта? Регистрация
          </Button>
        </Box>
      </Paper>
    </Box>
  );
}
