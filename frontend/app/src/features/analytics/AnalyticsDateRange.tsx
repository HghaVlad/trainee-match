import { Input } from '@/shared/ui/input'
import { Button } from '@/shared/ui/button'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/ui/select'
import { IntervalQueryParameter } from '@/api/generated/application/schemas'

export interface AnalyticsRangeValue {
  createdFrom?: string
  createdTo?: string
  interval?: IntervalQueryParameter
}

interface Props {
  value: AnalyticsRangeValue
  onChange: (next: AnalyticsRangeValue) => void
  showInterval?: boolean
}

const INTERVAL_LABEL: Record<IntervalQueryParameter, string> = {
  day: 'По дням',
  week: 'По неделям',
  month: 'По месяцам',
}

export function AnalyticsDateRange({ value, onChange, showInterval }: Props) {
  function reset() {
    onChange({
      createdFrom: undefined,
      createdTo: undefined,
      interval: value.interval,
    })
  }

  return (
    <div className="flex flex-wrap items-end gap-3">
      <div className="flex flex-col gap-1">
        <label className="text-xs text-muted-foreground" htmlFor="analytics-from">
          Создано с
        </label>
        <Input
          id="analytics-from"
          type="date"
          value={value.createdFrom ?? ''}
          onChange={(e) =>
            onChange({ ...value, createdFrom: e.target.value || undefined })
          }
        />
      </div>
      <div className="flex flex-col gap-1">
        <label className="text-xs text-muted-foreground" htmlFor="analytics-to">
          Создано до
        </label>
        <Input
          id="analytics-to"
          type="date"
          value={value.createdTo ?? ''}
          onChange={(e) =>
            onChange({ ...value, createdTo: e.target.value || undefined })
          }
        />
      </div>
      {showInterval && (
        <div className="flex flex-col gap-1">
          <span className="text-xs text-muted-foreground">Группировка</span>
          <Select
            value={value.interval ?? IntervalQueryParameter.day}
            onValueChange={(v) =>
              onChange({ ...value, interval: v as IntervalQueryParameter })
            }
          >
            <SelectTrigger className="w-44">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              {(
                Object.values(IntervalQueryParameter) as IntervalQueryParameter[]
              ).map((i) => (
                <SelectItem key={i} value={i}>
                  {INTERVAL_LABEL[i]}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      )}
      <Button type="button" variant="ghost" onClick={reset}>
        Сбросить
      </Button>
    </div>
  )
}
