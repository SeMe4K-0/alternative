import React from 'react'
import { useUIStore } from '../../stores/uiStore'

interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  label?: string
  error?: string
  helperText?: string
}

const Input: React.FC<InputProps> = ({
  label,
  error,
  helperText,
  className = '',
  id,
  ...props
}) => {
  const { theme } = useUIStore()
  const inputId = id || `input-${Math.random().toString(36).substr(2, 9)}`
  
  const baseClasses = 'block w-full px-3 py-2 border rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 transition-colors'
  const errorClasses = error ? 'border-red-300 focus:ring-red-500 focus:border-red-500' : (theme === 'dark' ? 'border-gray-600' : 'border-gray-300')
  const disabledClasses = props.disabled ? 'bg-gray-50 cursor-not-allowed' : (theme === 'dark' ? 'bg-gray-700 text-white' : 'bg-white')

  return (
    <div className="space-y-1">
      {label && (
        <label htmlFor={inputId} className={`block text-sm font-medium ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>
          {label}
        </label>
      )}
      <input
        id={inputId}
        className={`${baseClasses} ${errorClasses} ${disabledClasses} ${theme === 'dark' ? 'text-white' : ''} ${className}`}
        {...props}
      />
      {error && (
        <p className="text-sm text-red-600">{error}</p>
      )}
      {helperText && !error && (
        <p className={`text-sm ${theme === 'dark' ? 'text-gray-400' : 'text-gray-500'}`}>{helperText}</p>
      )}
    </div>
  )
}

export default Input
