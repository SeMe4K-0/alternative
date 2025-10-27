import React, { useState, useEffect } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { useObservationStore } from '../../stores/observationStore'
import { useCometStore } from '../../stores/cometStore'
import { useUIStore } from '../../stores/uiStore'
import PageContainer from '../../components/Common/PageContainer'
import LoadingSpinner from '../../components/Common/LoadingSpinner'

const EditObservation: React.FC = () => {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const { observations, updateObservation, isLoading } = useObservationStore()
  const { currentComet } = useCometStore()
  const { setLoading, isLoading: uiLoading, theme } = useUIStore()
  
  const observation = observations.find(obs => obs.id === id)

  const [formData, setFormData] = useState({
    observationDate: '',
    coordinates: {
      ra: '',
      dec: ''
    },
    notes: ''
  })
  
  const [imageFile, setImageFile] = useState<File | null>(null)
  const [previewUrl, setPreviewUrl] = useState<string>('')
  const [errors, setErrors] = useState<Record<string, string>>({})

  const isFormLoading = isLoading || uiLoading('editObservation')

  // Загружаем данные наблюдения
  useEffect(() => {
    if (observation) {
      setFormData({
        observationDate: observation.observationDate.slice(0, 16), // Без секунд для datetime-local
        coordinates: {
          ra: observation.coordinates.ra.toString(),
          dec: observation.coordinates.dec.toString()
        },
        notes: observation.notes || ''
      })
      
      // Устанавливаем превью текущего изображения
      if (observation.imageUrl) {
        setPreviewUrl(observation.imageUrl)
      }
    }
  }, [observation])

  // Проверяем, что наблюдение существует
  useEffect(() => {
    if (!observation) {
      navigate('/observations')
    }
  }, [observation, navigate])

  // Валидация формы
  const validateForm = () => {
    const newErrors: Record<string, string> = {}

    if (!formData.coordinates.ra.trim()) {
      newErrors.ra = 'Прямое восхождение обязательно'
    } else if (isNaN(Number(formData.coordinates.ra)) || Number(formData.coordinates.ra) < 0 || Number(formData.coordinates.ra) >= 24) {
      newErrors.ra = 'Прямое восхождение должно быть от 0 до 24 часов'
    }

    if (!formData.coordinates.dec.trim()) {
      newErrors.dec = 'Склонение обязательно'
    } else if (isNaN(Number(formData.coordinates.dec)) || Number(formData.coordinates.dec) < -90 || Number(formData.coordinates.dec) > 90) {
      newErrors.dec = 'Склонение должно быть от -90 до +90 градусов'
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  // Обработка загрузки файла
  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (file) {
      setImageFile(file)
      const url = URL.createObjectURL(file)
      setPreviewUrl(url)
    }
  }

  // Отправка формы
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!validateForm()) {
      return
    }

    try {
      setLoading('editObservation', true)
      
      const updateData: any = {
        observationDate: formData.observationDate,
        coordinates: {
          ra: Number(formData.coordinates.ra),
          dec: Number(formData.coordinates.dec)
        },
        notes: formData.notes.trim() || undefined
      }

      if (imageFile) {
        updateData.imageFile = imageFile
      }
      
      await updateObservation(id!, updateData)

      navigate(`/observations/${id}`)
    } catch (error) {
      console.error('Ошибка при обновлении наблюдения:', error)
      if (error instanceof Error) {
        alert(`Ошибка: ${error.message}`)
      } else {
        alert('Произошла неизвестная ошибка при обновлении наблюдения')
      }
    } finally {
      setLoading('editObservation', false)
    }
  }

  if (!observation) {
    return null
  }

  return (
    <div className={`min-h-screen bg-gradient-to-br ${theme === 'dark' ? 'from-gray-900 to-gray-800' : 'from-slate-50 to-blue-50'}`}>
      <PageContainer className="max-w-4xl">
        {/* Header */}
        <div className="mb-8">
          <div className={`rounded-2xl shadow-lg border overflow-hidden ${theme === 'dark' ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-100'}`}>
            <div className="bg-gradient-to-r from-emerald-600 to-teal-600 px-8 py-6">
              <h1 className="text-3xl font-bold text-white mb-2">Редактировать наблюдение</h1>
              <p className="text-emerald-100 text-lg">
                Комета: <span className="font-semibold">{observation.cometName}</span>
              </p>
            </div>
          </div>
        </div>

        <form onSubmit={handleSubmit} className="space-y-8">
          {/* Basic Information */}
          <div className={`rounded-2xl shadow-lg border p-8 ${theme === 'dark' ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-100'}`}>
            
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label htmlFor="observationDate" className={`block text-sm font-medium mb-3 ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>
                  Дата и время наблюдения *
                </label>
                <input
                  type="datetime-local"
                  id="observationDate"
                  value={formData.observationDate}
                  onChange={(e) => setFormData(prev => ({ ...prev, observationDate: e.target.value }))}
                  className={`w-full px-4 py-3 border rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200 ${
                    theme === 'dark' ? 'border-gray-600 bg-gray-700 text-white' : 'border-gray-300'
                  }`}
                />
              </div>
            </div>

            <div className="mt-8">
              <h3 className={`text-lg font-semibold mb-4 ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>Координаты</h3>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <label htmlFor="ra" className={`block text-sm font-medium mb-3 ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>
                    Прямое восхождение (RA) *
                  </label>
                  <input
                    type="number"
                    id="ra"
                    step="0.001"
                    min="0"
                    max="23.999"
                    value={formData.coordinates.ra}
                    onChange={(e) => setFormData(prev => ({ 
                      ...prev, 
                      coordinates: { ...prev.coordinates, ra: e.target.value }
                    }))}
                    className={`w-full px-4 py-3 border rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200 ${
                      errors.ra 
                        ? theme === 'dark' 
                          ? 'border-red-500 bg-red-900 text-white' 
                          : 'border-red-500 bg-red-50'
                        : theme === 'dark'
                          ? 'border-gray-600 bg-gray-700 text-white'
                          : 'border-gray-300'
                    }`}
                    placeholder="12.345"
                  />
                  <p className={`mt-2 text-xs ${theme === 'dark' ? 'text-gray-400' : 'text-gray-500'}`}>Часы (0-23.999)</p>
                  {errors.ra && <p className="mt-1 text-sm text-red-600">{errors.ra}</p>}
                </div>

                <div>
                  <label htmlFor="dec" className={`block text-sm font-medium mb-3 ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>
                    Склонение (Dec) *
                  </label>
                  <input
                    type="number"
                    id="dec"
                    step="0.001"
                    min="-90"
                    max="90"
                    value={formData.coordinates.dec}
                    onChange={(e) => setFormData(prev => ({ 
                      ...prev, 
                      coordinates: { ...prev.coordinates, dec: e.target.value }
                    }))}
                    className={`w-full px-4 py-3 border rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200 ${
                      errors.dec 
                        ? theme === 'dark' 
                          ? 'border-red-500 bg-red-900 text-white' 
                          : 'border-red-500 bg-red-50'
                        : theme === 'dark'
                          ? 'border-gray-600 bg-gray-700 text-white'
                          : 'border-gray-300'
                    }`}
                    placeholder="45.678"
                  />
                  <p className={`mt-2 text-xs ${theme === 'dark' ? 'text-gray-400' : 'text-gray-500'}`}>Градусы (-90 до +90)</p>
                  {errors.dec && <p className="mt-1 text-sm text-red-600">{errors.dec}</p>}
                </div>
              </div>
            </div>

            <div className="mt-8">
              <label htmlFor="notes" className={`block text-sm font-medium mb-3 ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>
                Заметки
              </label>
              <textarea
                id="notes"
                rows={4}
                value={formData.notes}
                onChange={(e) => setFormData(prev => ({ ...prev, notes: e.target.value }))}
                className={`w-full px-4 py-3 border rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200 ${
                  theme === 'dark' ? 'border-gray-600 bg-gray-700 text-white' : 'border-gray-300'
                }`}
                placeholder="Дополнительная информация о наблюдении..."
              />
            </div>

            <div className="mt-8">
              <label htmlFor="image" className={`block text-sm font-medium mb-3 ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>
                Изображение
              </label>
              <input
                type="file"
                id="image"
                accept="image/*"
                onChange={handleFileChange}
                className={`w-full px-4 py-3 border rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200 ${
                  theme === 'dark' ? 'border-gray-600 bg-gray-700' : 'border-gray-300'
                }`}
              />
              {previewUrl && (
                <div className="mt-4">
                  <img
                    src={previewUrl}
                    alt="Предварительный просмотр"
                    className="w-full h-48 object-cover rounded-xl"
                  />
                </div>
              )}
            </div>
          </div>

          {/* Action Buttons */}
          <div className="flex justify-end space-x-4">
            <button
              type="button"
              onClick={() => navigate(`/observations/${id}`)}
              className={`px-6 py-3 border rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all duration-200 font-medium ${
                theme === 'dark' 
                  ? 'border-gray-600 text-white hover:bg-gray-700 bg-gray-800' 
                  : 'border-gray-300 text-gray-700 hover:bg-gray-50'
              }`}
            >
              Отмена
            </button>
            <button
              type="submit"
              disabled={isFormLoading}
              className="px-8 py-3 bg-gradient-to-r from-emerald-600 to-teal-600 text-white rounded-xl hover:from-emerald-700 hover:to-teal-700 focus:outline-none focus:ring-2 focus:ring-emerald-500 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200 font-semibold shadow-lg hover:shadow-xl transform hover:scale-105"
            >
              {isFormLoading ? (
                <div className="flex items-center space-x-2">
                  <LoadingSpinner size="sm" color="white" />
                  <span>Сохранение...</span>
                </div>
              ) : (
                'Сохранить изменения'
              )}
            </button>
          </div>
        </form>
      </PageContainer>
    </div>
  )
}

export default EditObservation

