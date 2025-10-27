import React, { useState, useEffect } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { useCometStore } from '../../stores/cometStore'
import { useUIStore } from '../../stores/uiStore'
import PageContainer from '../../components/Common/PageContainer'
import LoadingSpinner from '../../components/Common/LoadingSpinner'
import Button from '../../components/Common/Button'
import Input from '../../components/Common/Input'
import Textarea from '../../components/Common/Textarea'
import Card from '../../components/Common/Card'

const EditComet: React.FC = () => {
  const navigate = useNavigate()
  const { id } = useParams<{ id: string }>()
  const { comets, updateComet, loadComets, isLoading } = useCometStore()
  const { theme } = useUIStore()
  
  const [formData, setFormData] = useState({
    name: '',
    notes: ''
  })
  const [errors, setErrors] = useState<Record<string, string>>({})

  useEffect(() => {
    if (id) {
      loadComets()
    }
  }, [id, loadComets])

  useEffect(() => {
    if (id && comets.length > 0) {
      const comet = comets.find(c => c.id === id)
      if (comet) {
        setFormData({
          name: comet.name,
          notes: comet.notes || ''
        })
      }
    }
  }, [id, comets])

  const validateForm = () => {
    const newErrors: Record<string, string> = {}

    if (!formData.name.trim()) {
      newErrors.name = 'Название кометы обязательно'
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!validateForm()) {
      return
    }

    if (!id) {
      console.error('Comet ID is missing')
      return
    }

    try {
      await updateComet(id, {
        name: formData.name.trim(),
        notes: formData.notes.trim() || undefined
      })
      
      navigate('/observations')
    } catch (error) {
      console.error('Ошибка при обновлении кометы:', error)
    }
  }

  const handleCancel = () => {
    navigate('/observations')
  }

  if (isLoading && !formData.name) {
    return (
      <div className={`min-h-screen bg-gradient-to-br ${theme === 'dark' ? 'from-gray-900 to-gray-800' : 'from-slate-50 to-blue-50'} flex items-center justify-center`}>
        <LoadingSpinner size="xl" text="Загрузка данных кометы..." />
      </div>
    )
  }

  const comet = comets.find(c => c.id === id)

  if (!comet && !isLoading) {
    return (
      <div className={`min-h-screen bg-gradient-to-br ${theme === 'dark' ? 'from-gray-900 to-gray-800' : 'from-slate-50 to-blue-50'}`}>
        <PageContainer className="max-w-2xl">
          <Card>
            <div className="p-6 text-center">
              <h2 className={`text-2xl font-bold mb-4 ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>Комета не найдена</h2>
              <p className={`mb-6 ${theme === 'dark' ? 'text-gray-400' : 'text-gray-600'}`}>Комета с ID "{id}" не найдена.</p>
              <Button onClick={handleCancel}>Вернуться к наблюдениям</Button>
            </div>
          </Card>
        </PageContainer>
      </div>
    )
  }

  return (
    <div className={`min-h-screen bg-gradient-to-br ${theme === 'dark' ? 'from-gray-900 to-gray-800' : 'from-slate-50 to-blue-50'}`}>
      <PageContainer className="max-w-3xl">
        {/* Header */}
        <div className="mb-8">
          <div className={`rounded-2xl shadow-lg border overflow-hidden ${theme === 'dark' ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-100'}`}>
            <div className="bg-gradient-to-r from-blue-600 to-purple-600 px-8 py-6">
              <h1 className="text-3xl font-bold text-white mb-2">Редактировать комету</h1>
              <p className="text-blue-100 text-lg">Измените информацию о комете</p>
            </div>
          </div>
        </div>

        {/* Edit Form */}
        <Card>
          <div className={`p-6 ${theme === 'dark' ? 'bg-gray-800' : 'bg-white'}`}>
            <form onSubmit={handleSubmit} className="space-y-6">
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
                  rows={5}
                />
              </div>

              {comet?.imageUrl && (
                <div>
                  <label className={`block text-sm font-medium mb-2 ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>
                    Изображение кометы
                  </label>
                  <div className={`relative rounded-lg overflow-hidden border-2 ${theme === 'dark' ? 'border-gray-600' : 'border-gray-300'}`}>
                    <img
                      src={comet.imageUrl}
                      alt={`Фото кометы ${comet.name}`}
                      className="w-full h-64 object-cover"
                    />
                  </div>
                </div>
              )}
              
              <div className="flex space-x-4">
                <Button
                  type="submit"
                  variant="primary"
                  disabled={isLoading}
                  className="flex-1"
                >
                  {isLoading ? 'Сохранение...' : 'Сохранить изменения'}
                </Button>
                
                <Button
                  type="button"
                  variant="secondary"
                  onClick={handleCancel}
                  className="flex-1"
                >
                  Отмена
                </Button>
              </div>
            </form>
          </div>
        </Card>
      </PageContainer>
    </div>
  )
}

export default EditComet

