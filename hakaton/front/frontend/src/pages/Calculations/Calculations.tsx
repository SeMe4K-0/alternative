import React, { useState, useMemo, useEffect } from 'react'
import { Link, useNavigate, useLocation } from 'react-router-dom'
import { format } from 'date-fns'
import { ru } from 'date-fns/locale'
import { useObservationStore } from '../../stores/observationStore'
import { useCalculationStore } from '../../stores/calculationStore'
import { useCometStore } from '../../stores/cometStore'
import { useUIStore } from '../../stores/uiStore'
import { apiService } from '../../services/api'
import PageContainer from '../../components/Common/PageContainer'
import LoadingSpinner from '../../components/Common/LoadingSpinner'

const Calculations: React.FC = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const { observations } = useObservationStore()
  const { createCalculation, calculations, isLoading, error } = useCalculationStore()
  const { comets } = useCometStore()
  const { setLoading, isLoading: uiLoading, theme } = useUIStore()
  
  const [selectedComet, setSelectedComet] = useState<string>('')
  const [selectedObservations, setSelectedObservations] = useState<string[]>([])
  const [showResults, setShowResults] = useState(false)
  const [currentCalculation, setCurrentCalculation] = useState<any>(null)
  const [isCalculatingApproach, setIsCalculatingApproach] = useState(false)
  const [isLoadingCalculation, setIsLoadingCalculation] = useState(false)

  const isCalculating = isLoading || uiLoading('calculateOrbit')

  // Загрузка конкретного расчёта по ID из URL
  useEffect(() => {
    const params = new URLSearchParams(location.search)
    const calculationId = params.get('calculationId')
    
    console.log('Checking calculationId from URL:', calculationId)
    
    if (calculationId) {
      setIsLoadingCalculation(true)
      const calculation = calculations.find(calc => calc.id === calculationId)
      console.log('Found in store:', !!calculation)
      
      if (calculation) {
        setCurrentCalculation(calculation)
        setShowResults(true)
        setIsLoadingCalculation(false)
      } else {
        // Попробуем загрузить с сервера
        apiService.getCalculation(calculationId)
          .then(calc => {
            console.log('Loaded from server:', calc)
            setCurrentCalculation(calc)
            setShowResults(true)
          })
          .catch(err => {
            console.error('Failed to load calculation:', err)
          })
          .finally(() => {
            setIsLoadingCalculation(false)
          })
      }
    } else {
      // Если нет calculationId, показываем обычный интерфейс
      setShowResults(false)
    }
  }, [location.search, calculations])

  if (isLoadingCalculation) {
    return (
      <div className={`min-h-screen bg-gradient-to-br ${theme === 'dark' ? 'from-gray-900 to-gray-800' : 'from-slate-50 to-purple-50'} flex items-center justify-center`}>
        <LoadingSpinner size="xl" text="Загрузка расчёта..." />
      </div>
    )
  }

  // Группируем наблюдения по кометам
  const observationsByComet = useMemo(() => {
    const grouped = observations.reduce((acc, obs) => {
      if (!acc[obs.cometName]) {
        acc[obs.cometName] = []
      }
      acc[obs.cometName].push(obs)
      return acc
    }, {} as Record<string, typeof observations>)

    return Object.keys(grouped).map(cometName => ({
      cometName,
      observations: grouped[cometName],
      cometInfo: comets.find(c => c.name === cometName)
    }))
  }, [observations, comets])

  // Получаем наблюдения выбранной кометы
  const availableObservations = useMemo(() => {
    if (!selectedComet) return []
    const cometGroup = observationsByComet.find(group => group.cometName === selectedComet)
    return cometGroup ? cometGroup.observations : []
  }, [selectedComet, observationsByComet])

  const handleCometSelect = (cometName: string) => {
    setSelectedComet(cometName)
    setSelectedObservations([]) // Сбрасываем выбранные наблюдения при смене кометы
  }

  const handleObservationToggle = (observationId: string) => {
    setSelectedObservations(prev => {
      if (prev.includes(observationId)) {
        return prev.filter(id => id !== observationId)
      } else {
        return [...prev, observationId]
      }
    })
  }

  const handleCalculate = async () => {
    if (selectedObservations.length < 5) {
      alert('Необходимо выбрать минимум 5 наблюдений для расчета орбиты')
      return
    }

    try {
      setLoading('calculateOrbit', true)
      
      // Получаем cometId из выбранной кометы
      const cometGroup = observationsByComet.find(group => group.cometName === selectedComet)
      if (!cometGroup || !cometGroup.cometInfo) {
        throw new Error('Информация о комете не найдена')
      }
      
      // Этап 1: Расчет параметров орбиты
      await createCalculation(selectedObservations, parseInt(cometGroup.cometInfo.id))
      setShowResults(true)
      
      const latestCalculation = useCalculationStore.getState().calculations[0]
      setCurrentCalculation(latestCalculation)
    } catch (err) {
      console.error('Ошибка при расчете:', err)
      alert(err instanceof Error ? err.message : 'Ошибка при расчете параметров орбиты')
    } finally {
      setLoading('calculateOrbit', false)
    }
  }
  
  // Этап 2: Расчет сближения с Землей
  const handleCalculateApproach = async () => {
    if (!currentCalculation) return
    
    try {
      setIsCalculatingApproach(true)
      
      const response = await fetch(`/api/calculations/${currentCalculation.id}/approach`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          semi_major_axis: currentCalculation.orbitParameters.semiMajorAxis,
          eccentricity: currentCalculation.orbitParameters.eccentricity,
          inclination: currentCalculation.orbitParameters.inclination,
          lon_ascending_node: currentCalculation.orbitParameters.longitudeOfAscendingNode,
          arg_periapsis: currentCalculation.orbitParameters.argumentOfPeriapsis,
          time_perihelion: new Date().toISOString(),
        }),
      })
      
      if (!response.ok) {
        throw new Error('Ошибка при расчете сближения')
      }
      
      const data = await response.json()
      
      // Обновляем расчет с данными сближения
      setCurrentCalculation({
        ...currentCalculation,
        earthApproach: {
          date: data.approach_date,
          distance: data.distance_au,
          velocity: 0,
        },
      })
    } catch (err) {
      console.error('Ошибка при расчете сближения:', err)
      alert(err instanceof Error ? err.message : 'Ошибка при расчете сближения')
    } finally {
      setIsCalculatingApproach(false)
    }
  }

  const resetCalculation = () => {
    setSelectedComet('')
    setSelectedObservations([])
    setShowResults(false)
    setCurrentCalculation(null)
  }

  if (isCalculating) {
    return (
      <div className={`min-h-screen bg-gradient-to-br ${theme === 'dark' ? 'from-gray-900 to-gray-800' : 'from-slate-50 to-purple-50'} flex items-center justify-center`}>
        <LoadingSpinner size="xl" text="Выполняется расчет орбиты..." />
      </div>
    )
  }

  return (
    <div className={`min-h-screen bg-gradient-to-br ${theme === 'dark' ? 'from-gray-900 to-gray-800' : 'from-slate-50 to-purple-50'}`}>
      <PageContainer className="max-w-6xl">
        {/* Header */}
        <div className="mb-8">
          <div className={`rounded-2xl shadow-lg border overflow-hidden ${theme === 'dark' ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-100'}`}>
            <div className="bg-gradient-to-r from-purple-600 to-pink-600 px-8 py-6">
              <h1 className="text-3xl font-bold text-white mb-2">Расчёты орбит</h1>
              <p className="text-purple-100 text-lg">Выберите наблюдения для расчета орбитальных параметров кометы</p>
            </div>
          </div>
        </div>

        {error && (
          <div className={`mb-8 px-6 py-4 rounded-xl border ${
            theme === 'dark'
              ? 'bg-red-900 border-red-700 text-red-300'
              : 'bg-red-50 border-red-200 text-red-600'
          }`}>
            <span className="font-medium">{error}</span>
          </div>
        )}

        {!showResults ? (
          <div className="space-y-8">
            {/* Выбор кометы */}
            <div className={`rounded-2xl shadow-lg border p-8 ${theme === 'dark' ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-100'}`}>
              <h2 className={`text-2xl font-bold mb-6 ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>Выберите комету для расчёта орбиты</h2>
              
              {observationsByComet.length === 0 ? (
                <div className={`text-center py-12 ${theme === 'dark' ? 'text-gray-400' : 'text-gray-600'}`}>
                  <div className={`text-4xl font-bold mb-6 ${theme === 'dark' ? 'text-gray-600' : 'text-gray-300'}`}>Нет наблюдений</div>
                  <p className="mb-8 max-w-md mx-auto">
                    Для расчёта орбиты необходимо сначала добавить наблюдения комет
                  </p>
                  <Link
                    to="/comets/select"
                    className="inline-flex items-center px-6 py-3 bg-gradient-to-r from-blue-600 to-purple-600 text-white font-semibold rounded-xl hover:from-blue-700 hover:to-purple-700 transition-all duration-200 shadow-lg hover:shadow-xl transform hover:scale-105"
                  >
                    Добавить наблюдение
                  </Link>
                </div>
              ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                  {observationsByComet.map((group) => (
                    <div
                      key={group.cometName}
                      className={`border-2 rounded-2xl p-6 cursor-pointer transition-all duration-300 transform hover:scale-105 ${
                        selectedComet === group.cometName
                          ? theme === 'dark' 
                            ? 'border-blue-500 bg-blue-900 shadow-lg'
                            : 'border-blue-500 bg-blue-50 shadow-lg'
                          : theme === 'dark'
                            ? 'border-gray-600 hover:border-gray-500 hover:shadow-md bg-gray-700'
                            : 'border-gray-200 hover:border-gray-300 hover:shadow-md'
                      }`}
                      onClick={() => handleCometSelect(group.cometName)}
                    >
                      <div className="flex items-center justify-between mb-4">
                        <h3 className={`text-lg font-bold ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>
                          {group.cometName}
                        </h3>
                        <div className="flex items-center space-x-2">
                          <span className="text-xs bg-green-100 text-green-800 px-3 py-1 rounded-full font-medium">
                            {group.observations.length} наблюдений
                          </span>
                          {selectedComet === group.cometName && (
                            <span className="text-blue-600 text-xl">✓</span>
                          )}
                        </div>
                      </div>
                      
                      {group.cometInfo?.notes && (
                        <p className={`text-sm ${theme === 'dark' ? 'text-gray-400' : 'text-gray-600'} mb-3`}>
                          {group.cometInfo.notes}
                        </p>
                      )}
                      
                      <div className={`text-xs ${theme === 'dark' ? 'text-gray-500' : 'text-gray-500'}`}>
                        {group.observations.length >= 5 ? 'Готово для расчёта' : `Нужно ещё ${5 - group.observations.length} наблюдений`}
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>

            {/* Выбор наблюдений */}
            {selectedComet && (
              <div className={`rounded-2xl shadow-lg border p-8 ${theme === 'dark' ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-100'}`}>
                <h2 className={`text-2xl font-bold mb-6 ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>
                  Выберите наблюдения кометы "{selectedComet}"
                </h2>
                
                <div className={`mb-6 p-4 rounded-xl ${theme === 'dark' ? 'bg-gray-700' : 'bg-gray-50'}`}>
                  <div className="flex items-center justify-between">
                    <div>
                      <p className={`text-lg font-medium ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>
                        Выбрано: {selectedObservations.length} из {availableObservations.length} наблюдений
                      </p>
                      {selectedObservations.length < 5 && (
                        <p className="text-red-600 mt-1">
                          Минимум 5 наблюдений необходимо для расчёта орбиты
                        </p>
                      )}
                    </div>
                    <div className="text-right">
                      <div className={`text-sm font-bold px-3 py-1 rounded-full ${
                        selectedObservations.length >= 5
                          ? theme === 'dark'
                            ? 'bg-green-900 text-green-300'
                            : 'bg-green-100 text-green-800'
                          : theme === 'dark'
                            ? 'bg-red-900 text-red-300'
                            : 'bg-red-100 text-red-800'
                      }`}>
                        {selectedObservations.length >= 5 ? 'Готово' : 'Недостаточно'}
                      </div>
                    </div>
                  </div>
                </div>
                
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                  {availableObservations.map((observation) => (
                    <div
                      key={observation.id}
                      className={`border-2 rounded-2xl p-6 cursor-pointer transition-all duration-300 transform hover:scale-105 ${
                        selectedObservations.includes(observation.id)
                          ? theme === 'dark' 
                            ? 'border-purple-500 bg-purple-900 shadow-lg'
                            : 'border-purple-500 bg-purple-50 shadow-lg'
                          : theme === 'dark'
                            ? 'border-gray-600 hover:border-gray-500 hover:shadow-md bg-gray-700'
                            : 'border-gray-200 hover:border-gray-300 hover:shadow-md'
                      }`}
                      onClick={() => handleObservationToggle(observation.id)}
                    >
                      <div className="flex items-center justify-between mb-4">
                        <div className={`text-sm font-medium ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>
                          {format(new Date(observation.observationDate), 'dd.MM.yyyy HH:mm:ss', { locale: ru })}
                        </div>
                        {selectedObservations.includes(observation.id) && (
                          <span className="text-purple-600 text-xl">✓</span>
                        )}
                      </div>
                      
                      <div className={`space-y-2 text-sm ${theme === 'dark' ? 'text-gray-400' : 'text-gray-600'}`}>
                        <div className="flex justify-between">
                          <span className="font-medium">RA:</span>
                          <span className="font-mono">{observation.coordinates.ra.toFixed(3)}h</span>
                        </div>
                        <div className="flex justify-between">
                          <span className="font-medium">Dec:</span>
                          <span className="font-mono">{observation.coordinates.dec.toFixed(3)}°</span>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
                
                <div className="mt-8 flex justify-center">
                  <button
                    type="button"
                    onClick={handleCalculate}
                    disabled={selectedObservations.length < 5 || isCalculating}
                    className="px-8 py-4 w-full sm:w-auto bg-gradient-to-r from-purple-600 to-pink-600 text-white font-bold rounded-xl hover:from-purple-700 hover:to-pink-700 focus:outline-none focus:ring-2 focus:ring-purple-500 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200 shadow-lg hover:shadow-xl"
                  >
                    {isCalculating ? (
                      <div className="flex items-center space-x-2">
                        <LoadingSpinner size="sm" color="white" />
                        <span>Расчёт...</span>
                      </div>
                    ) : selectedObservations.length >= 5 ? (
                      'Рассчитать орбиту'
                    ) : (
                      `Нужно ещё ${5 - selectedObservations.length} наблюдений`
                    )}
                  </button>
                </div>
              </div>
            )}
          </div>
        ) : (
          <div className="space-y-8">
            <div className="flex justify-between items-center">
              <h2 className={`text-3xl font-bold ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>Результаты расчета орбиты</h2>
              <button
                onClick={resetCalculation}
                className={`px-6 py-3 border rounded-xl focus:outline-none focus:ring-2 focus:ring-purple-500 transition-all duration-200 font-medium ${
                  theme === 'dark' 
                    ? 'border-gray-600 text-white hover:bg-gray-700 bg-gray-800' 
                    : 'border-gray-300 text-gray-700 hover:bg-gray-50'
                }`}
              >
                Новый расчет
              </button>
            </div>

            {currentCalculation && (
              <>
                <div className={`rounded-2xl shadow-lg border p-8 ${theme === 'dark' ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-100'}`}>
                  <div className="mb-6">
                    <h3 className={`text-2xl font-bold ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>Параметры орбиты</h3>
                  </div>
                  <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                    <div className={`p-6 rounded-xl border ${
                      theme === 'dark' 
                        ? 'bg-gradient-to-br from-blue-900 to-blue-800 border-blue-700' 
                        : 'bg-gradient-to-br from-blue-50 to-blue-100 border-blue-200'
                    }`}>
                      <div className={`text-sm font-medium mb-2 ${theme === 'dark' ? 'text-blue-300' : 'text-blue-600'}`}>Большая полуось</div>
                      <div className={`text-2xl font-bold ${theme === 'dark' ? 'text-white' : 'text-blue-900'}`}>
                        {currentCalculation.orbitParameters.semiMajorAxis.toFixed(6)} AU
                      </div>
                    </div>
                    
                    <div className={`p-6 rounded-xl border ${
                      theme === 'dark' 
                        ? 'bg-gradient-to-br from-green-900 to-green-800 border-green-700' 
                        : 'bg-gradient-to-br from-green-50 to-green-100 border-green-200'
                    }`}>
                      <div className={`text-sm font-medium mb-2 ${theme === 'dark' ? 'text-green-300' : 'text-green-600'}`}>Эксцентриситет</div>
                      <div className={`text-2xl font-bold ${theme === 'dark' ? 'text-white' : 'text-green-900'}`}>
                        {currentCalculation.orbitParameters.eccentricity.toFixed(6)}
                      </div>
                    </div>
                    
                    <div className={`p-6 rounded-xl border ${
                      theme === 'dark' 
                        ? 'bg-gradient-to-br from-purple-900 to-purple-800 border-purple-700' 
                        : 'bg-gradient-to-br from-purple-50 to-purple-100 border-purple-200'
                    }`}>
                      <div className={`text-sm font-medium mb-2 ${theme === 'dark' ? 'text-purple-300' : 'text-purple-600'}`}>Наклонение</div>
                      <div className={`text-2xl font-bold ${theme === 'dark' ? 'text-white' : 'text-purple-900'}`}>
                        {currentCalculation.orbitParameters.inclination.toFixed(3)}°
                      </div>
                    </div>
                    
                    <div className={`p-6 rounded-xl border ${
                      theme === 'dark' 
                        ? 'bg-gradient-to-br from-orange-900 to-orange-800 border-orange-700' 
                        : 'bg-gradient-to-br from-orange-50 to-orange-100 border-orange-200'
                    }`}>
                      <div className={`text-sm font-medium mb-2 ${theme === 'dark' ? 'text-orange-300' : 'text-orange-600'}`}>Долгота восходящего узла</div>
                      <div className={`text-2xl font-bold ${theme === 'dark' ? 'text-white' : 'text-orange-900'}`}>
                        {currentCalculation.orbitParameters.longitudeOfAscendingNode.toFixed(3)}°
                      </div>
                    </div>
                    
                    <div className={`p-6 rounded-xl border ${
                      theme === 'dark' 
                        ? 'bg-gradient-to-br from-pink-900 to-pink-800 border-pink-700' 
                        : 'bg-gradient-to-br from-pink-50 to-pink-100 border-pink-200'
                    }`}>
                      <div className={`text-sm font-medium mb-2 ${theme === 'dark' ? 'text-pink-300' : 'text-pink-600'}`}>Аргумент перигелия</div>
                      <div className={`text-2xl font-bold ${theme === 'dark' ? 'text-white' : 'text-pink-900'}`}>
                        {currentCalculation.orbitParameters.argumentOfPeriapsis.toFixed(3)}°
                      </div>
                    </div>
                    
                    <div className={`p-6 rounded-xl border ${
                      theme === 'dark' 
                        ? 'bg-gradient-to-br from-teal-900 to-teal-800 border-teal-700' 
                        : 'bg-gradient-to-br from-teal-50 to-teal-100 border-teal-200'
                    }`}>
                      <div className={`text-sm font-medium mb-2 ${theme === 'dark' ? 'text-teal-300' : 'text-teal-600'}`}>Время прохождения перигелия</div>
                      <div className={`text-2xl font-bold ${theme === 'dark' ? 'text-white' : 'text-teal-900'}`}>
                        {currentCalculation.orbitParameters.timePerihelion 
                          ? format(new Date(currentCalculation.orbitParameters.timePerihelion), 'dd.MM.yyyy HH:mm:ss', { locale: ru })
                          : 'Не указано'}
                      </div>
                    </div>
                  </div>
                  
                  {/* Кнопка для расчета сближения */}
                  {(() => {
                    // Проверяем, есть ли валидные данные о сближении
                    const hasValidApproach = currentCalculation.earthApproach && 
                      currentCalculation.earthApproach.date !== undefined && 
                      currentCalculation.earthApproach.date !== null &&
                      new Date(currentCalculation.earthApproach.date).getFullYear() !== 1
                    
                    return !hasValidApproach && (
                      <div className="mt-8 text-center">
                        <button
                          onClick={handleCalculateApproach}
                          disabled={isCalculatingApproach}
                          className={`px-8 py-4 bg-gradient-to-r from-yellow-600 to-orange-600 text-white font-bold rounded-xl hover:from-yellow-700 hover:to-orange-700 focus:outline-none focus:ring-2 focus:ring-yellow-500 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200 shadow-lg hover:shadow-xl transform hover:scale-105 ${
                            isCalculatingApproach ? 'animate-pulse' : ''
                          }`}
                        >
                          {isCalculatingApproach ? (
                            <div className="flex items-center space-x-2">
                              <LoadingSpinner size="sm" color="white" />
                              <span>Расчет сближения с Землей...</span>
                            </div>
                          ) : (
                            'Рассчитать сближение с Землей'
                          )}
                        </button>
                      </div>
                    )
                  })()}
                </div>

                {(() => {
                  // Проверяем, есть ли валидные данные о сближении
                  const hasValidApproach = currentCalculation.earthApproach && 
                    currentCalculation.earthApproach.date !== undefined && 
                    currentCalculation.earthApproach.date !== null &&
                    new Date(currentCalculation.earthApproach.date).getFullYear() !== 1
                  
                  return hasValidApproach && (
                  <div className={`rounded-2xl shadow-lg border p-8 ${theme === 'dark' ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-100'}`}>
                    <div className="mb-6">
                      <h3 className={`text-2xl font-bold ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>Сближение с Землей</h3>
                    </div>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                      <div className={`p-6 rounded-xl border ${
                        theme === 'dark' 
                          ? 'bg-gradient-to-br from-blue-900 to-blue-800 border-blue-700' 
                          : 'bg-gradient-to-br from-blue-50 to-blue-100 border-blue-200'
                      }`}>
                        <div className={`text-sm font-medium mb-2 ${theme === 'dark' ? 'text-blue-300' : 'text-blue-600'}`}>Дата сближения</div>
                        <div className={`text-xl font-bold ${theme === 'dark' ? 'text-white' : 'text-blue-900'}`}>
                          {format(new Date(currentCalculation.earthApproach.date), 'dd.MM.yyyy', { locale: ru })}
                        </div>
                      </div>
                      
                      <div className={`p-6 rounded-xl border ${
                        theme === 'dark' 
                          ? 'bg-gradient-to-br from-green-900 to-green-800 border-green-700' 
                          : 'bg-gradient-to-br from-green-50 to-green-100 border-green-200'
                      }`}>
                        <div className={`text-sm font-medium mb-2 ${theme === 'dark' ? 'text-green-300' : 'text-green-600'}`}>Минимальное расстояние</div>
                        <div className={`text-xl font-bold ${theme === 'dark' ? 'text-white' : 'text-green-900'}`}>
                          {currentCalculation.earthApproach.distance.toFixed(3)} AU
                        </div>
                      </div>
                    </div>
                  </div>
                  )
                })()}

                <div className={`rounded-2xl shadow-lg border p-8 ${theme === 'dark' ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-100'}`}>
                  <div className="mb-6">
                    <h3 className={`text-2xl font-bold ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>Информация о расчете</h3>
                  </div>
                  <div className="space-y-6">
                    <div className={`p-4 rounded-xl ${theme === 'dark' ? 'bg-gray-700' : 'bg-gray-50'}`}>
                      <div className={`text-sm font-medium mb-2 ${theme === 'dark' ? 'text-gray-400' : 'text-gray-600'}`}>Дата расчета</div>
                      <div className={`text-lg font-bold ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>
                        {format(new Date(currentCalculation.createdAt), 'dd.MM.yyyy HH:mm', { locale: ru })}
                      </div>
                    </div>
                    
                    <div>
                      <div className={`text-sm font-medium mb-3 ${theme === 'dark' ? 'text-gray-400' : 'text-gray-600'}`}>
                        Использованные наблюдения ({currentCalculation.observationIds.length})
                      </div>
                      <div className="flex flex-wrap gap-2">
                        {currentCalculation.observationIds.map((obsId: string) => {
                          const observation = observations.find(obs => obs.id === obsId)
                          return (
                            <button
                              key={obsId}
                              onClick={() => navigate(`/observations/${obsId}`)}
                              className={`inline-flex items-center px-3 py-1 rounded-full text-sm transition-colors ${
                                theme === 'dark' 
                                  ? 'bg-gray-700 text-gray-200 hover:bg-gray-600' 
                                  : 'bg-gray-100 text-gray-800 hover:bg-gray-200'
                              }`}
                            >
                              {observation 
                                ? `${format(new Date(observation.observationDate), 'dd.MM.yyyy HH:mm:ss', { locale: ru })}`
                                : `ID: ${obsId}`}
                            </button>
                          )
                        })}
                      </div>
                    </div>
                  </div>
                </div>
              </>
            )}
          </div>
        )}
      </PageContainer>
    </div>
  )
}

export default Calculations