import React from 'react'
import { useUIStore } from '../../stores/uiStore'

interface TextareaProps {
  value: string
  onChange: (value: string) => void
  placeholder?: string
  disabled?: boolean
  error?: string
  label?: string
  rows?: number
  maxLength?: number
  className?: string
  resize?: 'none' | 'vertical' | 'horizontal' | 'both'
}

const Textarea: React.FC<TextareaProps> = ({
  value,
  onChange,
  placeholder,
  disabled = false,
  error,
  label,
  rows = 4,
  maxLength,
  className = '',
  resize = 'vertical'
}) => {
  const { theme } = useUIStore()
  const resizeClasses = {
    none: 'resize-none',
    vertical: 'resize-y',
    horizontal: 'resize-x',
    both: 'resize'
  }

  const isDark = theme === 'dark'
  const baseClasses = `
    block w-full rounded-xl border px-4 py-3 text-sm
    focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500
    transition-colors duration-200
    ${resizeClasses[resize]}
    ${disabled ? 'bg-gray-100 text-gray-500 cursor-not-allowed' : isDark ? 'bg-gray-700 text-white' : 'bg-white text-gray-900'}
    ${error ? 'border-red-300 focus:ring-red-500 focus:border-red-500' : isDark ? 'border-gray-600' : 'border-gray-300'}
    ${className}
  `.trim()

  return (
    <div className="space-y-1">
      {label && (
        <label className={`block text-sm font-medium ${isDark ? 'text-gray-300' : 'text-gray-700'}`}>
          {label}
          {maxLength && (
            <span className={`ml-2 text-xs ${isDark ? 'text-gray-400' : 'text-gray-500'}`}>
              ({value.length}/{maxLength})
            </span>
          )}
        </label>
      )}
      <textarea
        value={value}
        onChange={(e) => onChange(e.target.value)}
        placeholder={placeholder}
        disabled={disabled}
        rows={rows}
        maxLength={maxLength}
        className={baseClasses}
      />
      {error && (
        <p className="text-sm text-red-600">{error}</p>
      )}
    </div>
  )
}

export default Textarea
