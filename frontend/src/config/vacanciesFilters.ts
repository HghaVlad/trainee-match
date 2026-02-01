import type { FilterConfig } from '../types/filters';

export const internshipFilters: FilterConfig[] = [
    {
        key: 'city',
        label: 'Город',
        type: 'select',
        options: [
            { label: 'Москва', value: 'moscow' },
            { label: 'Санкт-Петербург', value: 'spb' },
            { label: 'Казань', value: 'kazan' },
        ],
    },
    {
        key: 'format',
        label: 'Формат',
        type: 'select',
        options: [
            { label: 'Удалённо', value: 'remote' },
            { label: 'Офис', value: 'office' },
        ],
    },
    {
        key: 'duration',
        label: 'Длительность',
        type: 'select',
        options: [
            { label: '1-3 месяца', value: 'short' },
            { label: '3-6 месяцев', value: 'medium' },
        ],
    },
];
