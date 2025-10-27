import React, { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { useCometStore } from '../../stores/cometStore'
import { useUIStore } from '../../stores/uiStore'
import PageContainer from '../../components/Common/PageContainer'
import LoadingSpinner from '../../components/Common/LoadingSpinner'
import Button from '../../components/Common/Button'
import Input from '../../components/Common/Input'
import Textarea from '../../components/Common/Textarea'
import Card from '../../components/Common/Card'
import type { Comet } from '../../types'

const CometSelection: React.FC = () => {
  const navigate = useNavigate()
  const { comets, currentComet, createComet, loadComets, setCurrentComet, isLoading } = useCometStore()
  const { theme } = useUIStore()
  
  const [showCreateForm, setShowCreateForm] = useState(false)
  const [formData, setFormData] = useState({
    name: '',
    notes: ''
  })
  const [errors, setErrors] = useState<Record<string, string>>({})

  useEffect(() => {
    loadComets()
  }, [loadComets])

  const validateForm = () => {
    const newErrors: Record<string, string> = {}

    if (!formData.name.trim()) {
      newErrors.name = 'Название кометы обязательно'
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleCreateComet = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!validateForm()) {
      return
    }

    try {
      await createComet({
        name: formData.name.trim(),
        notes: formData.notes.trim() || undefined
      })
      
      setFormData({ name: '', notes: '' })
      setShowCreateForm(false)
    } catch (error) {
      console.error('Ошибка при создании кометы:', error)
    }
  }

  const handleSelectComet = (comet: Comet) => {
    setCurrentComet(comet)
    navigate('/observations/new')
  }

  if (isLoading) {
    return (
      <div className={`min-h-screen bg-gradient-to-br ${theme === 'dark' ? 'from-gray-900 to-gray-800' : 'from-slate-50 to-blue-50'} flex items-center justify-center`}>
        <LoadingSpinner size="xl" text="Загрузка комет..." />
      </div>
    )
  }

  return (
    <div className={`min-h-screen bg-gradient-to-br ${theme === 'dark' ? 'from-gray-900 to-gray-800' : 'from-slate-50 to-blue-50'}`}>
      <PageContainer className="max-w-6xl">
        {/* Header */}
        <div className="mb-8">
          <div className={`rounded-2xl shadow-lg border overflow-hidden ${theme === 'dark' ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-100'}`}>
            <div className="bg-gradient-to-r from-blue-600 to-purple-600 px-8 py-6">
              <h1 className="text-3xl font-bold text-white mb-2">Выбор кометы</h1>
              <p className="text-blue-100 text-lg">Выберите комету для наблюдения или создайте новую</p>
            </div>
          </div>
        </div>

        {/* Create Comet Form */}
        {showCreateForm && (
          <Card className="mb-8">
            <div className="p-6">
              <h2 className={`text-2xl font-bold mb-6 ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>
                Создать новую комету
              </h2>
              
              <form onSubmit={handleCreateComet} className="space-y-6">
                <div>
                  <Input
                    label="Название кометы"
                    value={formData.name}
                    onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                    error={errors.name}
                    placeholder="Например: Комета Галлея"
                    required
                  />
                </div>
                
                <div>
                  <Textarea
                    label="Заметки"
                    value={formData.notes}
                    onChange={(e) => setFormData({ ...formData, notes: e })}
                    placeholder="Дополнительные заметки о комете..."
                    rows={3}
                  />
                </div>
                
                <div className="flex space-x-4">
                  <Button
                    type="submit"
                    variant="primary"
                    disabled={isLoading}
                  >
                    {isLoading ? 'Создание...' : 'Создать комету'}
                  </Button>
                  
                  <Button
                    type="button"
                    variant="secondary"
                    onClick={() => setShowCreateForm(false)}
                  >
                    Отмена
                  </Button>
                </div>
              </form>
            </div>
          </Card>
        )}

        {/* Comets List */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {comets.map((comet) => (
            <Card
              key={comet.id}
              className={`cursor-pointer transition-all duration-300 transform hover:scale-105 overflow-hidden ${
                currentComet?.id === comet.id
                  ? theme === 'dark'
                    ? 'border-blue-500 bg-blue-900'
                    : 'border-blue-500 bg-blue-50'
                  : ''
              }`}
              onClick={() => handleSelectComet(comet)}
            >
              {/* Изображение кометы */}
              {comet.imageUrl && (
                <div className="aspect-w-16 aspect-h-9 relative">
                  <img
                    src={comet.imageUrl}
                    alt={`Фото кометы ${comet.name}`}
                    className="w-full h-48 object-cover"
                  />
                </div>
              )}
              
              <div className="p-6">
                <div className="flex items-center justify-between mb-4">
                  <h3 className={`text-xl font-bold ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>
                    {comet.name}
                  </h3>
                  {currentComet?.id === comet.id && (
                    <span className="text-blue-600 text-xl">✓</span>
                  )}
                </div>
                
                {comet.notes && (
                  <div className={`text-xs p-2 rounded-lg mb-3 ${theme === 'dark' ? 'text-gray-300 bg-gray-700' : 'text-gray-500 bg-gray-50'}`}>
                    <strong>Заметки:</strong> {comet.notes}
                  </div>
                )}
                
                <div className={`text-xs ${theme === 'dark' ? 'text-gray-400' : 'text-gray-500'}`}>
                  Создана: {new Date(comet.createdAt).toLocaleDateString('ru-RU')}
                </div>
              </div>
            </Card>
          ))}
        </div>

        {/* Create New Comet Button */}
        {!showCreateForm && (
          <div className="mt-8 text-center">
            <Button
              variant="primary"
              onClick={() => setShowCreateForm(true)}
              className="px-8 py-4 text-lg"
            >
              + Создать новую комету
            </Button>
          </div>
        )}
      </PageContainer>
    </div>
  )
}

export default CometSelection
