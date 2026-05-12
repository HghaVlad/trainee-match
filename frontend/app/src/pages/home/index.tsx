import { Link } from 'react-router'
import { Button } from '@/shared/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/card'
import { Badge } from '@/shared/ui/badge'

const stats = [
  { label: 'Компаний', value: '50+' },
  { label: 'Вакансий', value: '200+' },
  { label: 'Кандидатов', value: '1000+' },
  { label: 'Успешных стажировок', value: '300+' },
]

const howItWorks = [
  {
    step: 1,
    title: 'Создайте профиль',
    description:
      'Зарегистрируйтесь как соискатель или компания и заполните информацию о себе.',
    audience: 'candidate',
  },
  {
    step: 2,
    title: 'Найдите подходящую стажировку',
    description:
      'Просматривайте вакансии, фильтруйте по городу, зарплате и формату работы.',
    audience: 'candidate',
  },
  {
    step: 3,
    title: 'Откликнитесь и получите оффер',
    description:
      'Отправляйте отклики с резюме, общайтесь с компаниями и начинайте карьеру.',
    audience: 'candidate',
  },
  {
    step: 1,
    title: 'Опубликуйте вакансию',
    description:
      'Создайте компанию, разместите стажировку и укажите требования к кандидатам.',
    audience: 'company',
  },
  {
    step: 2,
    title: 'Изучайте отклики',
    description:
      'Просматривайте анкеты кандидатов, отслеживайте статусы в удобной воронке.',
    audience: 'company',
  },
  {
    step: 3,
    title: 'Наймите лучших',
    description:
      'Выбирайте подходящих стажёров, управляйте командой и анализируйте результаты.',
    audience: 'company',
  },
]

const features = [
  {
    title: 'Удобный поиск',
    description: 'Фильтруйте вакансии по городу, зарплате, графику и формату работы.',
  },
  {
    title: 'Воронка откликов',
    description: 'Отслеживайте статусы заявок: от нового отклика до оффера.',
  },
  {
    title: 'Аналитика',
    description: 'Смотрите динамику и статистику по откликам на вакансии.',
  },
  {
    title: 'Командная работа',
    description: 'Добавляйте коллег в компанию и управляйте вакансиями вместе.',
  },
  {
    title: 'Резюме и портфолио',
    description: 'Создавайте резюме, чтобы компании могли узнать вас лучше.',
  },
  {
    title: 'Безопасность',
    description: 'Авторизация через cookie и разделение ролей: кандидат или компания.',
  },
]

export default function HomePage() {
  return (
    <div>
      {/* Hero */}
      <section className="relative overflow-hidden border-b border-border bg-gradient-to-b from-primary/5 via-background to-background">
        <div className="mx-auto max-w-5xl px-6 pb-24 pt-16 text-center sm:pb-32 sm:pt-24">
          <Badge variant="outline" className="mb-6 text-sm">
            Платформа для поиска стажировок
          </Badge>
          <h1 className="text-4xl font-bold tracking-tight sm:text-5xl lg:text-6xl">
            Найди стажировку{' '}
            <span className="text-primary/70">своей мечты</span>
          </h1>
          <p className="mx-auto mt-6 max-w-2xl text-lg text-muted-foreground">
            trainee-match соединяет талантливых стажёров с лучшими компаниями.
            Просматривай вакансии, откликайся одной кнопкой и начинай карьеру.
          </p>
          <div className="mt-10 flex flex-col items-center justify-center gap-4 sm:flex-row">
            <Button size="lg" asChild>
              <Link to="/vacancies">Найти стажировку</Link>
            </Button>
            <Button size="lg" variant="outline" asChild>
              <Link to="/register">Зарегистрироваться</Link>
            </Button>
          </div>
          <p className="mt-4 text-xs text-muted-foreground">
            Компаниям —{' '}
            <Link
              to="/register"
              className="font-medium text-primary underline underline-offset-2"
            >
              разместить вакансию
            </Link>
          </p>
        </div>
      </section>

      {/* Stats */}
      <section className="relative border-b border-border">
        <div className="mx-auto max-w-5xl px-6 py-16">
          <div className="grid grid-cols-2 gap-6 sm:grid-cols-4">
            {stats.map((s) => (
              <div key={s.label} className="text-center">
                <div className="text-3xl font-bold text-primary sm:text-4xl">
                  {s.value}
                </div>
                <div className="mt-1 text-sm text-muted-foreground">
                  {s.label}
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Маленькая плашка */}
        <div className="absolute bottom-3 right-3 rounded-md border border-border bg-background/80 px-3 py-1 text-xs text-muted-foreground shadow-sm backdrop-blur-sm">
          *Цифры выдуманы для красивой картинки
        </div>
      </section>

      {/* How it works */}
      <section className="border-b border-border" id="how-it-works">
        <div className="mx-auto max-w-5xl px-6 py-16 sm:py-24">
          <div className="text-center">
            <h2 className="text-3xl font-bold tracking-tight">Как это работает</h2>
            <p className="mt-2 text-muted-foreground">
              Для соискателей и компаний
            </p>
          </div>

          <div className="mt-12 grid gap-8 md:grid-cols-2">
            <div>
              <h3 className="mb-6 text-lg font-semibold text-primary">
                Для соискателей
              </h3>
              <div className="space-y-4">
                {howItWorks
                  .filter((h) => h.audience === 'candidate')
                  .map((h) => (
                    <Card key={`cand-${h.step}`}>
                      <CardContent className="flex items-start gap-4 p-4">
                        <span className="flex size-8 shrink-0 items-center justify-center rounded-full bg-primary/10 text-sm font-bold text-primary">
                          {h.step}
                        </span>
                        <div>
                          <div className="font-medium">{h.title}</div>
                          <p className="mt-1 text-sm text-muted-foreground">
                            {h.description}
                          </p>
                        </div>
                      </CardContent>
                    </Card>
                  ))}
              </div>
            </div>

            <div>
              <h3 className="mb-6 text-lg font-semibold text-primary/70">
                Для компаний
              </h3>
              <div className="space-y-4">
                {howItWorks
                  .filter((h) => h.audience === 'company')
                  .map((h) => (
                    <Card key={`comp-${h.step}`}>
                      <CardContent className="flex items-start gap-4 p-4">
                        <span className="flex size-8 shrink-0 items-center justify-center rounded-full bg-primary/10 text-sm font-bold text-primary">
                          {h.step}
                        </span>
                        <div>
                          <div className="font-medium">{h.title}</div>
                          <p className="mt-1 text-sm text-muted-foreground">
                            {h.description}
                          </p>
                        </div>
                      </CardContent>
                    </Card>
                  ))}
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Features */}
      <section className="border-b border-border">
        <div className="mx-auto max-w-5xl px-6 py-16 sm:py-24">
          <div className="text-center">
            <h2 className="text-3xl font-bold tracking-tight">
              Возможности платформы
            </h2>
            <p className="mt-2 text-muted-foreground">
              Всё необходимое для успешного поиска и найма стажёров
            </p>
          </div>

          <div className="mt-12 grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {features.map((f) => (
              <Card key={f.title}>
                <CardHeader>
                  <CardTitle className="text-base">{f.title}</CardTitle>
                </CardHeader>
                <CardContent>
                  <p className="text-sm text-muted-foreground">{f.description}</p>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      </section>

      {/* CTA */}
      <section className="border-b border-border">
        <div className="mx-auto max-w-3xl px-6 py-16 text-center sm:py-24">
          <h2 className="text-3xl font-bold tracking-tight">
            Готовы начать?
          </h2>
          <p className="mt-2 text-muted-foreground">
            Присоединяйтесь к trainee-match — платформе, где стажёры находят
            работу, а компании — таланты.
          </p>
          <div className="mt-8 flex flex-col items-center justify-center gap-4 sm:flex-row">
            <Button size="lg" asChild>
              <Link to="/register">Создать аккаунт</Link>
            </Button>
            <Button size="lg" variant="outline" asChild>
              <Link to="/vacancies">Смотреть вакансии</Link>
            </Button>
          </div>
        </div>
      </section>
    </div>
  )
}
