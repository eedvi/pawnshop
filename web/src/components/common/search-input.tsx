import { useState, useEffect } from 'react'
import { Search, X } from 'lucide-react'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'
import { SEARCH_DEBOUNCE } from '@/lib/constants'

interface SearchInputProps {
  value?: string
  onChange: (value: string) => void
  placeholder?: string
  debounce?: number
  className?: string
}

export function SearchInput({
  value: externalValue,
  onChange,
  placeholder = 'Buscar...',
  debounce = SEARCH_DEBOUNCE,
  className,
}: SearchInputProps) {
  const [internalValue, setInternalValue] = useState(externalValue || '')

  // Sync external value changes
  useEffect(() => {
    if (externalValue !== undefined) {
      setInternalValue(externalValue)
    }
  }, [externalValue])

  // Debounce the onChange callback
  useEffect(() => {
    const timer = setTimeout(() => {
      if (internalValue !== externalValue) {
        onChange(internalValue)
      }
    }, debounce)

    return () => clearTimeout(timer)
  }, [internalValue, debounce, onChange, externalValue])

  return (
    <div className={cn('relative', className)}>
      <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
      <Input
        placeholder={placeholder}
        value={internalValue}
        onChange={(e) => setInternalValue(e.target.value)}
        className="pl-8 pr-8"
      />
      {internalValue && (
        <Button
          variant="ghost"
          size="icon"
          className="absolute right-0 top-0 h-9 w-9"
          onClick={() => {
            setInternalValue('')
            onChange('')
          }}
        >
          <X className="h-4 w-4" />
        </Button>
      )}
    </div>
  )
}
