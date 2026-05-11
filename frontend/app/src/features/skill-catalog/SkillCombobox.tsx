import { useMemo } from 'react'
import { Combobox } from '@/shared/ui/Combobox'
import { useGetSkillList } from '@/api/generated/candidate/skill/skill'

interface SkillComboboxProps {
  value?: string
  onChange?: (value: string) => void
  placeholder?: string
}

export function SkillCombobox({
  value,
  onChange,
  placeholder = 'Выберите навык...',
}: SkillComboboxProps) {
  const { data, isLoading } = useGetSkillList()

  const options = useMemo(
    () =>
      (data ?? [])
        .filter((s) => s.id !== undefined && s.name !== undefined)
        .map((s) => ({ value: String(s.id), label: String(s.name) })),
    [data],
  )

  return (
    <Combobox
      options={options}
      value={value}
      onChange={onChange ?? (() => {})}
      placeholder={placeholder}
      searchPlaceholder="Поиск навыка..."
      emptyText="Навыки не найдены"
      isLoading={isLoading}
    />
  )
}
