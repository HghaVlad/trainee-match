import * as React from 'react'
import { Check, ChevronsUpDown, Loader2, X } from 'lucide-react'
import { cn } from '@/shared/lib/utils'
import { Button } from '@/shared/ui/button'
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from '@/shared/ui/command'
import { Popover, PopoverContent, PopoverTrigger } from '@/shared/ui/popover'
import { Badge } from '@/shared/ui/badge'
import { useGetSkillList } from '@/api/generated/candidate/skill/skill'

interface SkillsMultiSelectProps {
  value: string[]
  onChange: (value: string[]) => void
  placeholder?: string
  className?: string
}

export function SkillsMultiSelect({
  value,
  onChange,
  placeholder = 'Выберите навыки…',
  className,
}: SkillsMultiSelectProps) {
  const [open, setOpen] = React.useState(false)
  const { data, isLoading } = useGetSkillList()

  const options = React.useMemo(
    () =>
      (data ?? [])
        .filter((s) => s.id !== undefined && s.name !== undefined)
        .map((s) => ({ value: String(s.id), label: String(s.name) })),
    [data],
  )

  const selectedSet = React.useMemo(() => new Set(value), [value])
  const selectedOptions = React.useMemo(
    () => options.filter((o) => selectedSet.has(o.value)),
    [options, selectedSet],
  )

  function toggle(v: string) {
    if (selectedSet.has(v)) {
      onChange(value.filter((x) => x !== v))
    } else {
      onChange([...value, v])
    }
  }

  function remove(v: string) {
    onChange(value.filter((x) => x !== v))
  }

  return (
    <div className={cn('space-y-2', className)}>
      <Popover open={open} onOpenChange={setOpen}>
        <PopoverTrigger asChild>
          <Button
            type="button"
            variant="outline"
            role="combobox"
            aria-expanded={open}
            className="w-full justify-between"
          >
            <span className="text-muted-foreground">
              {value.length > 0
                ? `Выбрано: ${value.length}`
                : placeholder}
            </span>
            <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
          </Button>
        </PopoverTrigger>
        <PopoverContent className="w-[--radix-popover-trigger-width] p-0" align="start">
          <Command>
            <CommandInput placeholder="Поиск навыка…" />
            <CommandList>
              <CommandEmpty>
                {isLoading ? (
                  <div className="flex items-center justify-center py-6">
                    <Loader2 className="h-4 w-4 animate-spin text-muted-foreground" />
                  </div>
                ) : (
                  'Навыки не найдены'
                )}
              </CommandEmpty>
              <CommandGroup>
                {options.map((option) => {
                  const checked = selectedSet.has(option.value)
                  return (
                    <CommandItem
                      key={option.value}
                      value={option.label}
                      onSelect={() => toggle(option.value)}
                    >
                      <Check
                        className={cn(
                          'mr-2 h-4 w-4',
                          checked ? 'opacity-100' : 'opacity-0',
                        )}
                      />
                      {option.label}
                    </CommandItem>
                  )
                })}
              </CommandGroup>
            </CommandList>
          </Command>
        </PopoverContent>
      </Popover>
      {selectedOptions.length > 0 && (
        <ul className="flex flex-wrap gap-1">
          {selectedOptions.map((o) => (
            <li key={o.value}>
              <Badge variant="secondary" className="gap-1 pr-1">
                {o.label}
                <button
                  type="button"
                  aria-label={`Удалить ${o.label}`}
                  className="rounded-sm hover:bg-muted-foreground/20"
                  onClick={() => remove(o.value)}
                >
                  <X className="h-3 w-3" />
                </button>
              </Badge>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
