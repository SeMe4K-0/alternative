import React from 'react'
import { useUIStore } from '../../stores/uiStore'

interface CardProps {
  children: React.ReactNode
  className?: string
  padding?: 'sm' | 'md' | 'lg'
  shadow?: 'sm' | 'md' | 'lg' | 'xl'
  rounded?: 'sm' | 'md' | 'lg' | 'xl' | '2xl'
  border?: boolean
  hover?: boolean
  onClick?: () => void
}

const Card: React.FC<CardProps> = ({ 
  children, 
  className = '',
  padding = 'md',
  shadow = 'lg',
  rounded = '2xl',
  border = true,
  hover = false,
  onClick
}) => {
  const { theme } = useUIStore()
  
  const paddingClasses = {
    sm: 'p-4',
    md: 'p-6',
    lg: 'p-8'
  }

  const shadowClasses = {
    sm: 'shadow-sm',
    md: 'shadow-md',
    lg: 'shadow-lg',
    xl: 'shadow-xl'
  }

  const roundedClasses = {
    sm: 'rounded-sm',
    md: 'rounded-md',
    lg: 'rounded-lg',
    xl: 'rounded-xl',
    '2xl': 'rounded-2xl'
  }

  const baseClasses = `
    ${theme === 'dark' ? 'bg-gray-800' : 'bg-white'}
    ${paddingClasses[padding]}
    ${shadowClasses[shadow]}
    ${roundedClasses[rounded]}
    ${border ? (theme === 'dark' ? 'border border-gray-700' : 'border border-gray-100') : ''}
    ${hover ? 'hover:shadow-xl transition-shadow duration-200' : ''}
    ${className}
  `.trim()

  return (
    <div className={baseClasses} onClick={onClick}>
      {children}
    </div>
  )
}

export default Card
