import React, { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { useObservationStore } from '../../stores/observationStore'
import { useCometStore } from '../../stores/cometStore'
import { useUIStore } from '../../stores/uiStore'
import PageContainer from '../../components/Common/PageContainer'
import LoadingSpinner from '../../components/Common/LoadingSpinner'

const AddObservation: React.FC = () => {
  const navigate = useNavigate()
  const { addObservation, isLoading } = useObservationStore()
  const { currentComet } = useCometStore()
  const { setLoading, isLoading: uiLoading, theme } = useUIStore()
  
  const [observations, setObservations] = useState([
    {
      id: 1,
      observationDate: new Date().toISOString().slice(0, 19), // YYYY-MM-DDTHH:MM:SS
      coordinates: {
        ra: '',
        dec: ''
      },
      notes: '',
      imageFile: null as File | null
    }
  ])
  
  const [errors, setErrors] = useState<Record<string, string>>({})

  const isFormLoading = isLoading || uiLoading('addObservation')

  // Проверяем, выбрана ли комета
  useEffect(() => {
    if (!currentComet) {
      navigate('/comets/select')
    }
  }, [currentComet, navigate])

  // Функции для управления строками таблицы
  const addObservationRow = () => {
    const newId = Math.max(...observations.map(obs => obs.id)) + 1
    setObservations([...observations, {
      id: newId,
      observationDate: new Date().toISOString().slice(0, 19),
      coordinates: { ra: '', dec: '' },
      notes: '',
      imageFile: null
    }])
  }

  const removeObservationRow = (id: number) => {
    if (observations.length > 1) {
      setObservations(observations.filter(obs => obs.id !== id))
    }
  }

  const updateObservation = (id: number, field: string, value: string | File | null) => {
    setObservations(observations.map(obs => {
      if (obs.id === id) {
        if (field === 'observationDate') {
          return { ...obs, observationDate: value as string }
        } else if (field === 'ra' || field === 'dec') {
          return { ...obs, coordinates: { ...obs.coordinates, [field]: value as string } }
        } else if (field === 'notes') {
          return { ...obs, notes: value as string }
        } else if (field === 'imageFile') {
          return { ...obs, imageFile: value as File | null }
        }
      }
      return obs
    }))
  }

  // Валидация формы
  const validateForm = () => {
    const newErrors: Record<string, string> = {}

    observations.forEach((obs) => {
      if (!obs.coordinates.ra.trim()) {
        newErrors[`ra_${obs.id}`] = 'Прямое восхождение обязательно'
      } else if (isNaN(Number(obs.coordinates.ra)) || Number(obs.coordinates.ra) < 0 || Number(obs.coordinates.ra) >= 24) {
        newErrors[`ra_${obs.id}`] = 'Прямое восхождение должно быть от 0 до 24 часов'
      }

      if (!obs.coordinates.dec.trim()) {
        newErrors[`dec_${obs.id}`] = 'Склонение обязательно'
      } else if (isNaN(Number(obs.coordinates.dec)) || Number(obs.coordinates.dec) < -90 || Number(obs.coordinates.dec) > 90) {
        newErrors[`dec_${obs.id}`] = 'Склонение должно быть от -90 до +90 градусов'
      }
    })

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  // Отправка формы
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!validateForm()) {
      return
    }

    try {
      setLoading('addObservation', true)
      
      // Создаем наблюдения для каждой строки таблицы
      for (const obs of observations) {
        const observationData: any = {
          userId: '1',
          cometId: currentComet!.id,
          cometName: currentComet!.name,
          observationDate: obs.observationDate,
          coordinates: {
            ra: Number(obs.coordinates.ra),
            dec: Number(obs.coordinates.dec)
          },
          imageUrl: '',
          annotatedPoints: [],
          notes: obs.notes.trim() || undefined
        }

        if (obs.imageFile) {
          observationData.imageFile = obs.imageFile
        }

        await addObservation(observationData)
      }

      navigate('/observations')
    } catch (error) {
      console.error('Ошибка при добавлении наблюдения:', error)
      // Показываем пользователю более детальную ошибку
      if (error instanceof Error) {
        alert(`Ошибка: ${error.message}`)
      } else {
        alert('Произошла неизвестная ошибка при добавлении наблюдения')
      }
    } finally {
      setLoading('addObservation', false)
    }
  }

  return (
    <div className={`min-h-screen bg-gradient-to-br ${theme === 'dark' ? 'from-gray-900 to-gray-800' : 'from-slate-50 to-blue-50'}`}>
      <PageContainer className="max-w-4xl">
        {/* Header */}
        <div className="mb-8">
          <div className={`rounded-2xl shadow-lg border overflow-hidden ${theme === 'dark' ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-100'}`}>
            <div className="bg-gradient-to-r from-emerald-600 to-teal-600 px-8 py-6">
              <h1 className="text-3xl font-bold text-white mb-2">Добавить наблюдение</h1>
              <p className="text-emerald-100 text-lg">
                Наблюдение для кометы: <span className="font-semibold">{currentComet?.name}</span>
              </p>
              {currentComet?.notes && (
                <p className="text-emerald-200 text-sm mt-1">Заметки: {currentComet.notes}</p>
              )}
            </div>
          </div>
        </div>

        <form onSubmit={handleSubmit} className="space-y-8">
          {/* Observations Table */}
          <div className={`rounded-2xl shadow-lg border p-8 ${theme === 'dark' ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-100'}`}>
            <div className="flex justify-between items-center mb-6">
              <h2 className={`text-2xl font-bold ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>Координаты наблюдений</h2>
              <button
                type="button"
                onClick={addObservationRow}
                className="px-4 py-2 bg-green-600 hover:bg-green-700 text-white rounded-lg font-medium transition-colors"
              >
                + Добавить строку
              </button>
            </div>

            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className={`border-b ${theme === 'dark' ? 'border-gray-700' : 'border-gray-200'}`}>
                    <th className={`text-left py-3 px-4 font-semibold ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>
                      Дата и время наблюдения
                    </th>
                    <th className={`text-left py-3 px-4 font-semibold ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>
                      Прямое восхождение (RA)
                    </th>
                    <th className={`text-left py-3 px-4 font-semibold ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>
                      Склонение (Dec)
                    </th>
                    <th className={`text-left py-3 px-4 font-semibold ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>
                      Комментарий
                    </th>
                    <th className={`text-left py-3 px-4 font-semibold ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>
                      Изображение
                    </th>
                    <th className={`text-center py-3 px-4 font-semibold ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>
                      Действия
                    </th>
                  </tr>
                </thead>
                <tbody>
                  {observations.map((obs) => (
                    <tr key={obs.id} className={`border-b ${theme === 'dark' ? 'border-gray-700' : 'border-gray-200'}`}>
                      <td className="py-3 px-4">
                        <input
                          type="datetime-local"
                          step="1"
                          value={obs.observationDate}
                          onChange={(e) => updateObservation(obs.id, 'observationDate', e.target.value)}
                          className={`w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200 ${
                            theme === 'dark' ? 'border-gray-600 bg-gray-700 text-white' : 'border-gray-300'
                          }`}
                        />
                      </td>
                      <td className="py-3 px-4">
                        <input
                          type="number"
                          step="0.001"
                          min="0"
                          max="23.999"
                          value={obs.coordinates.ra}
                          onChange={(e) => updateObservation(obs.id, 'ra', e.target.value)}
                          className={`w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200 ${
                            errors[`ra_${obs.id}`] 
                              ? theme === 'dark' 
                                ? 'border-red-500 bg-red-900 text-white' 
                                : 'border-red-500 bg-red-50'
                              : theme === 'dark'
                                ? 'border-gray-600 bg-gray-700 text-white'
                                : 'border-gray-300'
                          }`}
                          placeholder="12.345"
                        />
                        {errors[`ra_${obs.id}`] && (
                          <p className="mt-1 text-xs text-red-600">{errors[`ra_${obs.id}`]}</p>
                        )}
                      </td>
                      <td className="py-3 px-4">
                        <input
                          type="number"
                          step="0.001"
                          min="-90"
                          max="90"
                          value={obs.coordinates.dec}
                          onChange={(e) => updateObservation(obs.id, 'dec', e.target.value)}
                          className={`w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200 ${
                            errors[`dec_${obs.id}`] 
                              ? theme === 'dark' 
                                ? 'border-red-500 bg-red-900 text-white' 
                                : 'border-red-500 bg-red-50'
                              : theme === 'dark'
                                ? 'border-gray-600 bg-gray-700 text-white'
                                : 'border-gray-300'
                          }`}
                          placeholder="45.678"
                        />
                        {errors[`dec_${obs.id}`] && (
                          <p className="mt-1 text-xs text-red-600">{errors[`dec_${obs.id}`]}</p>
                        )}
                      </td>
                      <td className="py-3 px-4">
                        <textarea
                          rows={2}
                          value={obs.notes}
                          onChange={(e) => updateObservation(obs.id, 'notes', e.target.value)}
                          className={`w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200 resize-none ${
                            theme === 'dark' ? 'border-gray-600 bg-gray-700 text-white' : 'border-gray-300'
                          }`}
                          placeholder="Комментарий..."
                        />
                      </td>
                      <td className="py-3 px-4">
                        <input
                          type="file"
                          accept="image/*"
                          onChange={(e) => {
                            const file = e.target.files?.[0] || null
                            updateObservation(obs.id, 'imageFile', file)
                          }}
                          className={`w-full px-3 py-2 text-sm border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200 ${
                            theme === 'dark' ? 'border-gray-600 bg-gray-700' : 'border-gray-300'
                          }`}
                        />
                        {obs.imageFile && (
                          <p className="mt-1 text-xs text-green-600">✓ {obs.imageFile.name}</p>
                        )}
                      </td>
                      <td className="py-3 px-4 text-center">
                        <button
                          type="button"
                          onClick={() => removeObservationRow(obs.id)}
                          disabled={observations.length === 1}
                          className={`px-3 py-1 rounded-lg text-sm font-medium transition-colors ${
                            observations.length === 1
                              ? 'text-gray-400 cursor-not-allowed'
                              : theme === 'dark'
                                ? 'text-red-400 hover:text-red-300 hover:bg-red-900'
                                : 'text-red-600 hover:text-red-700 hover:bg-red-50'
                          }`}
                          title={observations.length === 1 ? 'Нельзя удалить единственную строку' : 'Удалить строку'}
                        >
                          Удалить
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            <div className="mt-4 text-sm text-gray-500">
              <p>* Минимум одна строка обязательна. Прямое восхождение: 0-23.999 часов, Склонение: -90 до +90 градусов. Комментарий и изображение опциональны.</p>
            </div>
          </div>

          {/* Action Buttons */}
          <div className="flex justify-end space-x-4">
            <button
              type="button"
              onClick={() => navigate('/observations')}
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
                'Сохранить наблюдения'
              )}
            </button>
          </div>
        </form>
      </PageContainer>
    </div>
  )
}

export default AddObservation