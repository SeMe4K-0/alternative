import React, { useMemo } from 'react'
import { Link } from 'react-router-dom'
import { useAuthStore } from '../../stores/authStore'
import { useObservationStore } from '../../stores/observationStore'
import { useCalculationStore } from '../../stores/calculationStore'
import { useUIStore } from '../../stores/uiStore'
import Button from '../../components/Common/Button'
import PageContainer from '../../components/Common/PageContainer'

const Dashboard: React.FC = () => {
  const { user } = useAuthStore()
  const { observations, isLoading: observationsLoading } = useObservationStore()
  const { calculations, isLoading: calculationsLoading } = useCalculationStore()
  const { theme } = useUIStore()

  // Группируем наблюдения по кометам
  const groupedObservations = useMemo(() => {
    const grouped = observations.reduce((acc, obs) => {
      if (!acc[obs.cometName]) {
        acc[obs.cometName] = []
      }
      acc[obs.cometName].push(obs)
      return acc
    }, {} as Record<string, typeof observations>)

    return Object.keys(grouped).map(cometName => ({
      cometName,
      observations: grouped[cometName]
    }))
  }, [observations])


  
  // Считаем кометы с 5+ наблюдениями
  const readyForCalculation = groupedObservations.filter(group => group.observations.length >= 5).length
  const uniqueComets = groupedObservations.length

  return (
    <div className={`h-full bg-gradient-to-br ${theme === 'dark' ? 'from-gray-900 to-gray-800' : 'from-blue-50 to-purple-50'}`}>
      <PageContainer>
        {/* Welcome Header */}
        <div className="mb-8">
          <div className={`rounded-2xl shadow-xl border overflow-hidden ${theme === 'dark' ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-100'}`}>
            <div className="bg-gradient-to-r from-blue-600 to-purple-600 px-8 py-6">
              <h1 className="text-3xl font-bold text-white mb-2">
                Добро пожаловать, {user?.name}!
              </h1>
              <p className="text-blue-100 text-lg">
                Управляйте своими наблюдениями комет и рассчитывайте их орбиты
              </p>
            </div>
            <div className="px-8 py-6">
              <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                <div className="text-center">
                  <div className={`text-3xl font-bold mb-1 ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>
                    {observationsLoading ? '...' : observations.length}
                  </div>
                  <div className={theme === 'dark' ? 'text-gray-400 text-sm' : 'text-sm text-gray-600'}>Наблюдений</div>
                </div>
                <div className="text-center">
                  <div className={`text-3xl font-bold mb-1 ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>
                    {calculationsLoading ? '...' : calculations.length}
                  </div>
                  <div className={theme === 'dark' ? 'text-gray-400 text-sm' : 'text-sm text-gray-600'}>Расчётов</div>
                </div>
                <div className="text-center">
                  <div className={`text-3xl font-bold mb-1 ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>
                    {uniqueComets}
                  </div>
                  <div className={theme === 'dark' ? 'text-gray-400 text-sm' : 'text-sm text-gray-600'}>Комет</div>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Quick Actions */}
        <div className="mb-8">
          <h2 className={`text-2xl font-bold mb-6 ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>Быстрые действия</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {/* Add Observation */}
            <div className={`group rounded-2xl shadow-lg hover:shadow-xl transition-all duration-300 border overflow-hidden ${theme === 'dark' ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-100'}`}>
              <div className="bg-gradient-to-br from-blue-500 to-cyan-600 p-6">
                <div>
                  <h3 className="text-xl font-bold text-white mb-1">Новое наблюдение</h3>
                  <p className="text-blue-100">Добавить данные о комете</p>
                </div>
              </div>
              <div className="p-6">
                <Link to="/comets/select">
                  <Button className="w-full bg-blue-600 hover:bg-blue-700 text-white font-semibold py-3 rounded-xl">
                    Добавить наблюдение
                  </Button>
                </Link>
              </div>
            </div>

            {/* Calculate Orbit */}
            <div className={`group rounded-2xl shadow-lg hover:shadow-xl transition-all duration-300 border overflow-hidden ${theme === 'dark' ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-100'}`}>
              <div className="bg-gradient-to-br from-indigo-500 to-purple-600 p-6">
                <div>
                  <h3 className="text-xl font-bold text-white mb-1">Расчёт орбиты</h3>
                  <p className="text-indigo-100">Рассчитать параметры орбиты</p>
                </div>
              </div>
              <div className="p-6">
                <Link to="/calculations">
                  <Button 
                    className="w-full bg-indigo-600 hover:bg-indigo-700 text-white font-semibold py-3 rounded-xl"
                    disabled={readyForCalculation === 0}
                  >
                    {readyForCalculation > 0 ? 'Рассчитать орбиту' : `Нужно ${5} наблюдений для кометы`}
                  </Button>
                </Link>
              </div>
            </div>

            {/* View History */}
            <div className={`group rounded-2xl shadow-lg hover:shadow-xl transition-all duration-300 border overflow-hidden ${theme === 'dark' ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-100'}`}>
              <div className="bg-gradient-to-br from-violet-500 to-purple-600 p-6">
                <div>
                  <h3 className="text-xl font-bold text-white mb-1">История</h3>
                  <p className="text-violet-100">Просмотр всех расчётов</p>
                </div>
              </div>
              <div className="p-6">
                <Link to="/history">
                  <Button className="w-full bg-indigo-600 hover:bg-violet-700 text-white font-semibold py-3 rounded-xl">
                    Просмотреть историю
                  </Button>
                </Link>
              </div>
            </div>
          </div>
        </div>

        {/* Statistics Cards */}
        <div className="mb-8">
          <h2 className={`text-2xl font-bold mb-6 ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>Статистика</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            <div className={`rounded-2xl shadow-lg p-6 border ${theme === 'dark' ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-100'}`}>
              <div className="flex items-center">
                <div className="ml-4">
                  <p className={`text-sm font-medium ${theme === 'dark' ? 'text-gray-400' : 'text-gray-600'}`}>Всего наблюдений</p>
                  <p className={`text-2xl font-bold ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>
                    {observationsLoading ? '...' : observations.length}
                  </p>
                </div>
              </div>
            </div>

            <div className={`rounded-2xl shadow-lg p-6 border ${theme === 'dark' ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-100'}`}>
              <div className="flex items-center">
                <div className="ml-4">
                  <p className={`text-sm font-medium ${theme === 'dark' ? 'text-gray-400' : 'text-gray-600'}`}>Готово для расчёта</p>
                  <p className={`text-2xl font-bold ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>{readyForCalculation}</p>
                </div>
              </div>
            </div>

            <div className={`rounded-2xl shadow-lg p-6 border ${theme === 'dark' ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-100'}`}>
              <div className="flex items-center">
                <div className="ml-4">
                  <p className={`text-sm font-medium ${theme === 'dark' ? 'text-gray-400' : 'text-gray-600'}`}>Уникальных комет</p>
                  <p className={`text-2xl font-bold ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>{uniqueComets}</p>
                </div>
              </div>
            </div>

            <div className={`rounded-2xl shadow-lg p-6 border ${theme === 'dark' ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-100'}`}>
              <div className="flex items-center">
                <div className="ml-4">
                  <p className={`text-sm font-medium ${theme === 'dark' ? 'text-gray-400' : 'text-gray-600'}`}>Выполнено расчётов</p>
                  <p className={`text-2xl font-bold ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>
                    {calculationsLoading ? '...' : calculations.length}
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Empty State */}
        {observations.length === 0 && (
          <div className={`rounded-2xl shadow-lg border p-12 text-center ${theme === 'dark' ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-100'}`}>
            <div className={`text-4xl font-bold mb-4 ${theme === 'dark' ? 'text-gray-600' : 'text-gray-300'}`}>Начните с первого наблюдения</div>
            <h3 className={`text-2xl font-bold mb-4 ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>Добавьте первое наблюдение</h3>
            <p className={`mb-8 max-w-md mx-auto ${theme === 'dark' ? 'text-gray-400' : 'text-gray-600'}`}>
              Добавьте первое наблюдение кометы, загрузите фотографию и отметьте точки для расчёта орбиты
            </p>
            <Link to="/comets/select">
              <Button className="bg-blue-600 hover:bg-blue-700 text-white font-semibold px-8 py-3 rounded-xl">
                Добавить первое наблюдение
              </Button>
            </Link>
          </div>
        )}
      </PageContainer>
    </div>
  )
}

export default Dashboard