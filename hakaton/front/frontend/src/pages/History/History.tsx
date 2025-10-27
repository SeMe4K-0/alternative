import React, { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useCalculationStore } from '../../stores/calculationStore'
import { useObservationStore } from '../../stores/observationStore'
import { useUIStore } from '../../stores/uiStore'
import { apiService } from '../../services/api'
import PageContainer from '../../components/Common/PageContainer'
import Button from '../../components/Common/Button'
import LoadingSpinner from '../../components/Common/LoadingSpinner'
import { format } from 'date-fns'
import { ru } from 'date-fns/locale'

const History: React.FC = () => {
  const navigate = useNavigate()
  const { calculations, deleteCalculation, loadCalculations } = useCalculationStore()
  const { observations } = useObservationStore()
  const { theme } = useUIStore()
  const [filterStatus, setFilterStatus] = useState<'all' | 'completed' | 'pending'>('all')
  const [searchTerm, setSearchTerm] = useState('')
  const [calculatingIds, setCalculatingIds] = useState<Set<string>>(new Set())

  // Получаем первую комету из наблюдений расчёта
  const getCometName = (calculation: any) => {
    // Если название кометы уже есть в расчёте, используем его
    if (calculation.cometName) {
      return calculation.cometName
    }
    
    // Иначе пробуем получить из наблюдений
    const cometNames = calculation.observationIds.map((id: string) => 
      observations.find(obs => obs.id === id)?.cometName
    ).filter(Boolean) as string[]
    return [...new Set(cometNames)][0] || 'Неизвестная комета'
  }

  const filteredCalculations = calculations.filter(calc => {
    // Фильтр по статусу
    if (filterStatus !== 'all') {
      const hasValidApproach = calc.earthApproach && 
        calc.earthApproach.date !== undefined && 
        calc.earthApproach.date !== null &&
        new Date(calc.earthApproach.date).getFullYear() !== 1
      
      const displayStatus = hasValidApproach ? 'completed' : 'pending'
      if (displayStatus !== filterStatus) return false
    }
    
    // Фильтр по поиску
    if (searchTerm) {
      const cometName = getCometName(calc)
      if (!cometName.toLowerCase().includes(searchTerm.toLowerCase())) {
        return false
      }
    }
    
    return true
  })

  const sortedCalculations = [...filteredCalculations].sort((a, b) => {
    // Сортируем по дате создания (новые сначала)
    return new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()
  })

  const handleViewDetails = (calculationId: string) => {
    // Переходим на страницу расчётов с этим ID
    navigate(`/calculations?calculationId=${calculationId}`)
  }

  const handleCalculateApproach = async (calculationId: string) => {
    const calculation = calculations.find(calc => calc.id === calculationId)
    if (!calculation) return
    
    setCalculatingIds(prev => new Set(prev).add(calculationId))
    
    try {
      const earthApproach = await apiService.calculateApproach(calculationId, calculation.orbitParameters)
      
      // Перезагружаем расчеты, чтобы получить обновленные данные
      await loadCalculations()
    } catch (error) {
      console.error('Ошибка при расчете сближения:', error)
      alert('Не удалось рассчитать сближение с Землёй')
    } finally {
      setCalculatingIds(prev => {
        const newSet = new Set(prev)
        newSet.delete(calculationId)
        return newSet
      })
    }
  }

  const handleDelete = async (calculationId: string) => {
    if (window.confirm('Вы уверены, что хотите удалить этот расчёт?')) {
      await deleteCalculation(calculationId)
    }
  }

  return (
    <div className={`min-h-screen bg-gradient-to-br ${theme === 'dark' ? 'from-gray-900 to-gray-800' : 'from-blue-50 to-purple-50'}`}>
      <PageContainer>
        {/* Header */}
        <div className="mb-8">
          <div className={`rounded-2xl shadow-xl border overflow-hidden ${theme === 'dark' ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-100'}`}>
            <div className="bg-gradient-to-r from-purple-600 to-indigo-600 px-8 py-6">
              <h1 className="text-3xl font-bold text-white mb-2">История расчётов</h1>
              <p className="text-purple-100 text-lg">Просмотр всех выполненных расчётов орбит комет</p>
            </div>
          </div>
        </div>

        {/* Filters and Controls */}
        <div className={`rounded-2xl shadow-lg border p-6 mb-8 ${theme === 'dark' ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-100'}`}>
          <div className="flex flex-col sm:flex-row gap-4 items-center justify-between">
            <div className="flex flex-col sm:flex-row gap-4 flex-1">
              <div className="flex-1">
                <label className={`block text-sm font-medium mb-2 ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>
                  Поиск по названию кометы
                </label>
                <input
                  type="text"
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  placeholder="Введите название кометы..."
                  className={`w-full px-4 py-2 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 ${
                    theme === 'dark' ? 'bg-gray-700 border-gray-600 text-white placeholder-gray-400' : 'border border-gray-300'
                  }`}
                />
              </div>
              <div>
                <label className={`block text-sm font-medium mb-2 ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>
                  Фильтр по статусу
                </label>
                <select
                  value={filterStatus}
                  onChange={(e) => setFilterStatus(e.target.value as 'all' | 'completed' | 'pending')}
                  className={`px-4 py-2 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 ${
                    theme === 'dark' ? 'bg-gray-700 border-gray-600 text-white' : 'border border-gray-300'
                  }`}
                >
                  <option value="all">Все расчёты</option>
                  <option value="pending">Незавершённые</option>
                  <option value="completed">Завершённые</option>
                </select>
              </div>
            </div>
            <div className={`text-sm ${theme === 'dark' ? 'text-gray-400' : 'text-gray-600'}`}>
              Всего расчётов: <span className="font-semibold">{filteredCalculations.length}</span>
            </div>
          </div>
        </div>

        {/* Calculations List */}
        {sortedCalculations.length === 0 ? (
          <div className={`rounded-2xl shadow-lg border p-12 text-center ${theme === 'dark' ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-100'}`}>
            <div className={`text-4xl font-bold mb-4 ${theme === 'dark' ? 'text-gray-600' : 'text-gray-300'}`}>Статистика</div>
            <h3 className={`text-2xl font-bold mb-4 ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>Расчётов пока нет</h3>
            <p className={`mb-8 max-w-md mx-auto ${theme === 'dark' ? 'text-gray-400' : 'text-gray-600'}`}>
              Выполните первый расчёт орбиты кометы, чтобы увидеть результаты здесь
            </p>
            <Button className="bg-indigo-600 hover:bg-indigo-700 text-white font-semibold py-3 px-8 rounded-xl">
              Перейти к расчётам
            </Button>
          </div>
        ) : (
          <div className="space-y-6">
            {sortedCalculations.map((calculation) => (
              <div key={calculation.id} className={`rounded-2xl shadow-lg border overflow-hidden ${theme === 'dark' ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-100'}`}>
                {/* Header */}
                <div className={`bg-gradient-to-r px-6 py-4 border-b ${theme === 'dark' ? 'from-gray-700 to-gray-800 border-gray-600' : 'from-gray-50 to-gray-100 border-gray-200'}`}>
                  <div className="flex justify-between items-center">
                    <div>
                      <h3 className={`text-xl font-bold ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>
                        {getCometName(calculation)}
                      </h3>
                      <p className={`mt-1 ${theme === 'dark' ? 'text-gray-400' : 'text-gray-600'}`}>
                        {format(new Date(calculation.createdAt), 'dd MMMM yyyy, HH:mm', { locale: ru })}
                      </p>
                    </div>
                    <div className="flex items-center space-x-3">
                      {(() => {
                        // Проверяем, есть ли данные о сближении с Землёй
                        const hasValidApproach = calculation.earthApproach && 
                          calculation.earthApproach.date !== undefined && 
                          calculation.earthApproach.date !== null &&
                          new Date(calculation.earthApproach.date).getFullYear() !== 1
                        
                        const displayStatus = hasValidApproach ? 'completed' : 'pending'
                        
                        return (
                          <span className={`inline-flex items-center px-3 py-1 rounded-full text-sm font-medium ${
                            displayStatus === 'completed' 
                              ? theme === 'dark'
                                ? 'bg-green-900 text-green-300'
                                : 'bg-green-100 text-green-800'
                              : theme === 'dark'
                                ? 'bg-yellow-900 text-yellow-300'
                                : 'bg-yellow-100 text-yellow-800'
                          }`}>
                            {displayStatus === 'completed' ? 'Завершён' : 'Незавершён'}
                          </span>
                        )
                      })()}
                      <div className="flex space-x-2">
                        <Button
                          onClick={() => handleViewDetails(calculation.id)}
                          variant="secondary"
                          size="sm"
                          className="text-indigo-600 hover:text-indigo-800"
                        >
                          Подробнее
                        </Button>
                        <Button
                          onClick={() => handleDelete(calculation.id)}
                          variant="secondary"
                          size="sm"
                          className="text-red-600 hover:text-red-800"
                        >
                          Удалить
                        </Button>
                      </div>
                    </div>
                  </div>
                </div>

                {/* Content */}
                <div className="p-6">
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    {/* Orbit Parameters */}
                    <div>
                      <h4 className={`text-lg font-semibold mb-4 ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>Параметры орбиты</h4>
                      <div className="space-y-2 text-sm">
                        <div className="flex justify-between">
                          <span className={theme === 'dark' ? 'text-gray-400' : 'text-gray-600'}>Большая полуось (a):</span>
                          <span className={`font-medium ${theme === 'dark' ? 'text-white' : ''}`}>{calculation.orbitParameters.semiMajorAxis.toFixed(6)} AU</span>
                        </div>
                        <div className="flex justify-between">
                          <span className={theme === 'dark' ? 'text-gray-400' : 'text-gray-600'}>Эксцентриситет (e):</span>
                          <span className={`font-medium ${theme === 'dark' ? 'text-white' : ''}`}>{calculation.orbitParameters.eccentricity.toFixed(6)}</span>
                        </div>
                        <div className="flex justify-between">
                          <span className={theme === 'dark' ? 'text-gray-400' : 'text-gray-600'}>Наклонение (i):</span>
                          <span className={`font-medium ${theme === 'dark' ? 'text-white' : ''}`}>{calculation.orbitParameters.inclination.toFixed(3)}°</span>
                        </div>
                        <div className="flex justify-between">
                          <span className={theme === 'dark' ? 'text-gray-400' : 'text-gray-600'}>Долгота восходящего узла (Ω):</span>
                          <span className={`font-medium ${theme === 'dark' ? 'text-white' : ''}`}>{calculation.orbitParameters.longitudeOfAscendingNode.toFixed(3)}°</span>
                        </div>
                        <div className="flex justify-between">
                          <span className={theme === 'dark' ? 'text-gray-400' : 'text-gray-600'}>Аргумент перицентра (ω):</span>
                          <span className={`font-medium ${theme === 'dark' ? 'text-white' : ''}`}>{calculation.orbitParameters.argumentOfPeriapsis.toFixed(3)}°</span>
                        </div>
                        {calculation.orbitParameters.timePerihelion && (
                          <div className="flex justify-between">
                            <span className={theme === 'dark' ? 'text-gray-400' : 'text-gray-600'}>Время прохождения перигелия (T):</span>
                            <span className={`font-medium ${theme === 'dark' ? 'text-white' : ''}`}>
                              {format(new Date(calculation.orbitParameters.timePerihelion), 'dd.MM.yyyy HH:mm:ss', { locale: ru })}
                            </span>
                          </div>
                        )}
                      </div>
                    </div>

                    {/* Earth Approach */}
                    <div>
                      <h4 className={`text-lg font-semibold mb-4 ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>Сближение с Землёй</h4>
                      {(() => {
                        const hasValidApproach = calculation.earthApproach && 
                          calculation.earthApproach.date !== undefined && 
                          calculation.earthApproach.date !== null &&
                          new Date(calculation.earthApproach.date).getFullYear() !== 1
                        
                        return hasValidApproach ? (
                        <div className={`rounded-xl p-4 space-y-2 text-sm ${theme === 'dark' ? 'bg-blue-900' : 'bg-blue-50'}`}>
                          <div className="flex justify-between">
                            <span className={theme === 'dark' ? 'text-gray-400' : 'text-gray-600'}>Дата сближения:</span>
                            <span className={`font-medium ${theme === 'dark' ? 'text-white' : ''}`}>
                              {format(new Date(calculation.earthApproach.date), 'dd MMMM yyyy', { locale: ru })}
                            </span>
                          </div>
                          <div className="flex justify-between">
                            <span className={theme === 'dark' ? 'text-gray-400' : 'text-gray-600'}>Минимальное расстояние:</span>
                            <span className={`font-medium ${theme === 'dark' ? 'text-white' : ''}`}>{calculation.earthApproach.distance.toFixed(6)} AU</span>
                          </div>
                        </div>
                        ) : (
                        <div className={`rounded-xl p-6 text-center ${theme === 'dark' ? 'bg-yellow-900 border-2 border-yellow-700' : 'bg-yellow-50 border-2 border-yellow-200'}`}>
                          <p className={`text-sm mb-4 ${theme === 'dark' ? 'text-yellow-200' : 'text-yellow-800'}`}>
                            Данные о сближении с Землёй ещё не рассчитаны
                          </p>
                          <Button
                            onClick={() => handleCalculateApproach(calculation.id)}
                            disabled={calculatingIds.has(calculation.id)}
                            className={`bg-gradient-to-r from-yellow-600 to-orange-600 text-white font-semibold px-6 py-2 rounded-xl hover:from-yellow-700 hover:to-orange-700 transition-all duration-200 shadow-lg hover:shadow-xl disabled:opacity-50 disabled:cursor-not-allowed`}
                          >
                            {calculatingIds.has(calculation.id) ? (
                              <div className="flex items-center space-x-2">
                                <LoadingSpinner size="sm" color="white" />
                                <span>Расчёт...</span>
                              </div>
                            ) : (
                              'Рассчитать сближение'
                            )}
                          </Button>
                        </div>
                      )
                      })()}
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </PageContainer>
    </div>
  )
}

export default History
