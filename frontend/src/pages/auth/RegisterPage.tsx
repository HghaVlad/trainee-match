import {
    Box,
    Button,
    TextField,
    Typography,
    RadioGroup,
    FormControlLabel,
    Radio,
    Paper,
} from '@mui/material';
import { useForm } from 'react-hook-form';
import { registerApi } from '../../api/auth.api';
import { useNavigate } from 'react-router-dom';

interface RegisterForm {
    firstName: string;
    lastName: string;
    email: string;
    username: string;
    password: string;
    confirmPassword: string;
    role: 'Candidate' | 'Company';
}

export default function RegisterPage() {
    const navigate = useNavigate();
    const { register, handleSubmit } = useForm<RegisterForm>({
        defaultValues: { role: 'Candidate' },
    });

    const onSubmit = async (data: RegisterForm) => {
        if (data.password !== data.confirmPassword) {
            alert('Пароли не совпадают');
            return;
        }

        await registerApi({
            first_name: data.firstName,
            last_name: data.lastName,
            email: data.email,
            username: data.username,
            password: data.password,
            role: data.role, // Candidate | Company
        });

        navigate('/login');
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
                    Регистрация
                </Typography>

                <Box component="form" onSubmit={handleSubmit(onSubmit)} display="flex" flexDirection="column" gap={2}>
                    <TextField label="Имя" {...register('firstName', { required: true })} />
                    <TextField label="Фамилия" {...register('lastName', { required: true })} />
                    <TextField label="Почта" {...register('email', { required: true })} />
                    <TextField label="Имя пользователя" {...register('username', { required: true })} />
                    <TextField type="password" label="Пароль" {...register('password', { required: true })} />
                    <TextField type="password" label="Повторите пароль" {...register('confirmPassword', { required: true })} />

                    <RadioGroup {...register('role')}>
                        <FormControlLabel value="CANDIDATE" control={<Radio />} label="Кандидат" />
                        <FormControlLabel value="HR" control={<Radio />} label="Представитель компании" />
                    </RadioGroup>

                    <Button variant="outlined" type="submit">
                        Зарегистрироваться
                    </Button>

                    <Button variant="text" onClick={() => navigate('/login')}>
                        Уже есть аккаунт? Войти
                    </Button>
                </Box>
            </Paper>
        </Box>
    );
}
