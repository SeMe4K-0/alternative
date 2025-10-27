import React, { useState, useMemo } from 'react'
import { Link } from 'react-router-dom'
import { format } from 'date-fns'
import { ru } from 'date-fns/locale'
import { useObservationStore } from '../../stores/observationStore'
import { useCometStore } from '../../stores/cometStore'
import { useUIStore } from '../../stores/uiStore'
import PageContainer from '../../components/Common/PageContainer'

const Observations: React.FC = () => {
  const { observations, deleteObservation, isLoading } = useObservationStore()
  const { comets } = useCometStore()
  const { theme } = useUIStore()
  const [searchTerm, setSearchTerm] = useState('')

  // Группировка наблюдений по кометам
  const groupedObservations = useMemo(() => {
    // Фильтруем наблюдения
    let filtered = observations.filter(obs =>
      obs.cometName.toLowerCase().includes(searchTerm.toLowerCase()) ||
      obs.notes?.toLowerCase().includes(searchTerm.toLowerCase())
    )

    // Группируем по названию кометы
    const grouped = filtered.reduce((acc, obs) => {
      if (!acc[obs.cometName]) {
        acc[obs.cometName] = []
      }
      acc[obs.cometName].push(obs)
      return acc
    }, {} as Record<string, typeof observations>)

    // Сортируем группы по названию кометы
    const sortedGroups = Object.keys(grouped).sort()
    
    return sortedGroups.map(cometName => ({
      cometName,
      observations: grouped[cometName],
      cometInfo: comets.find(c => c.name === cometName)
    }))
  }, [observations, comets, searchTerm])

  const handleDelete = (id: string) => {
    if (window.confirm('Вы уверены, что хотите удалить это наблюдение?')) {
      deleteObservation(id)
    }
  }

  if (isLoading) {
    return (
      <div className={`min-h-screen bg-gradient-to-br ${theme === 'dark' ? 'from-gray-900 to-gray-800' : 'from-slate-50 to-blue-50'} flex items-center justify-center`}>
        <div className="text-center">
          <div className="animate-spin rounded-full h-16 w-16 border-b-4 border-blue-600 mx-auto mb-4"></div>
          <h2 className={`text-xl font-semibold mb-2 ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>Загрузка наблюдений</h2>
          <p className={theme === 'dark' ? 'text-gray-400' : 'text-gray-600'}>Пожалуйста, подождите...</p>
        </div>
      </div>
    )
  }

  return (
    <div className={`min-h-screen bg-gradient-to-br ${theme === 'dark' ? 'from-gray-900 to-gray-800' : 'from-slate-50 to-blue-50'}`}>
      <PageContainer>
        {/* Header */}
        <div className="mb-8">
          <div className="flex justify-between items-center">
            <div>
              <h1 className={`text-3xl font-bold mb-2 ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>Наблюдения комет</h1>
              <p className={`text-lg ${theme === 'dark' ? 'text-gray-400' : 'text-gray-600'}`}>Управление вашими наблюдениями комет</p>
            </div>
            <Link
              to="/comets/select"
              className="inline-flex items-center px-6 py-3 bg-gradient-to-r from-blue-600 to-purple-600 text-white font-semibold rounded-xl hover:from-blue-700 hover:to-purple-700 transition-all duration-200 shadow-lg hover:shadow-xl transform hover:scale-105"
            >
              Добавить наблюдение
            </Link>
          </div>
        </div>

        {/* Filters */}
        <div className={`rounded-2xl shadow-lg border p-6 mb-8 ${theme === 'dark' ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-100'}`}>
          <h2 className={`text-lg font-semibold mb-4 ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>Поиск</h2>
          <div className="grid grid-cols-1 gap-4">
            <div>
              <input
                type="text"
                id="search"
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className={`w-full px-4 py-3 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200 ${
                  theme === 'dark' ? 'bg-gray-700 border-gray-600 text-white' : 'border border-gray-300'
                }`}
                placeholder="По названию кометы или заметкам..."
              />
            </div>
          </div>
        </div>

        {/* Observations Grid */}
        {groupedObservations.length === 0 ? (
          <div className={`rounded-2xl shadow-lg border p-12 text-center ${theme === 'dark' ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-100'}`}>
            <div className={`text-4xl font-bold mb-6 ${theme === 'dark' ? 'text-gray-600' : 'text-gray-300'}`}>Наблюдений пока нет</div>
            {!searchTerm && (
              <Link
                to="/comets/select"
                className="inline-flex items-center px-6 py-3 bg-gradient-to-r from-blue-600 to-purple-600 text-white font-semibold rounded-xl hover:from-blue-700 hover:to-purple-700 transition-all duration-200 shadow-lg hover:shadow-xl transform hover:scale-105"
              >
                Добавить наблюдение
              </Link>
            )}
          </div>
        ) : (
          <div className="space-y-6">
            {groupedObservations.map((group) => (
              <div key={group.cometName} className={`rounded-2xl shadow-lg border overflow-hidden ${theme === 'dark' ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-100'}`}>
                {/* Header кометы */}
                <div className="bg-gradient-to-r from-blue-600 to-purple-600 px-6 py-4">
                  <div className="flex items-center justify-between">
                    <div className="flex-1">
                      <h2 className="text-2xl font-bold text-white">{group.cometName}</h2>
                      {group.cometInfo?.notes && (
                        <p className="text-blue-100 text-sm mt-1">{group.cometInfo.notes}</p>
                      )}
                    </div>
                    <div className="flex items-center space-x-4">
                      <div className="text-right">
                        <div className="text-white text-sm">Наблюдений:</div>
                        <div className="text-2xl font-bold text-white">{group.observations.length}</div>
                      </div>
                      {group.cometInfo && (
                        <Link
                          to={`/comets/${group.cometInfo.id}/edit`}
                          className="px-4 py-2 bg-white bg-opacity-20 hover:bg-opacity-30 text-white font-medium rounded-lg transition-all duration-200"
                          title="Редактировать комету"
                        >
                          Редактировать
                        </Link>
                      )}
                    </div>
                  </div>
                </div>

                {/* Изображение кометы */}
                {group.cometInfo?.imageUrl && (
                  <div className="relative h-48 overflow-hidden">
                    <img
                      src={group.cometInfo.imageUrl}
                      alt={`Фото кометы ${group.cometName}`}
                      className="w-full h-full object-cover"
                    />
                  </div>
                )}

                {/* Список наблюдений */}
                <div className="p-6">
                  <div className="space-y-4">
                    {group.observations.map((observation) => (
                      <div key={observation.id} className={`p-4 rounded-lg border ${theme === 'dark' ? 'bg-gray-700 border-gray-600' : 'bg-gray-50 border-gray-200'}`}>
                        <div className="flex justify-between items-start mb-3">
                          <div className="flex-1">
                            <div className={`text-sm font-medium ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>
                              {format(new Date(observation.observationDate), 'dd.MM.yyyy HH:mm:ss', { locale: ru })}
                            </div>
                          </div>
                          <div className="flex space-x-2">
                            <Link
                              to={`/observations/${observation.id}`}
                              className={`px-3 py-1 rounded-lg text-sm font-medium transition-colors ${
                                theme === 'dark' 
                                  ? 'text-blue-400 hover:bg-blue-900' 
                                  : 'text-blue-600 hover:bg-blue-50'
                              }`}
                              title="Подробнее"
                            >
                              Подробнее
                            </Link>
                            <button
                              onClick={() => handleDelete(observation.id)}
                              className={`px-3 py-1 rounded-lg text-sm font-medium transition-colors ${
                                theme === 'dark'
                                  ? 'text-red-400 hover:text-red-300 hover:bg-red-900'
                                  : 'text-red-600 hover:text-red-700 hover:bg-red-50'
                              }`}
                              title="Удалить"
                            >
                              Удалить
                            </button>
                          </div>
                        </div>
                        
                        <div className={`grid grid-cols-2 gap-4 text-sm ${theme === 'dark' ? 'text-gray-400' : 'text-gray-600'}`}>
                          <div>
                            <span className="font-medium">RA:</span>
                            <span className="ml-2 font-mono">{observation.coordinates.ra.toFixed(3)}h</span>
                          </div>
                          <div>
                            <span className="font-medium">Dec:</span>
                            <span className="ml-2 font-mono">{observation.coordinates.dec.toFixed(3)}°</span>
                          </div>
                        </div>
                        
                        {observation.notes && (
                          <div className="mt-3">
                            <p className={`text-xs p-2 rounded-lg ${theme === 'dark' ? 'text-gray-300 bg-gray-600' : 'text-gray-500 bg-gray-100'}`}>
                              {observation.notes}
                            </p>
                          </div>
                        )}
                      </div>
                    ))}
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}

        {/* Statistics */}
        {observations.length > 0 && (
          <div className={`mt-8 rounded-2xl shadow-lg border p-6 ${theme === 'dark' ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-100'}`}>
            <h3 className={`text-xl font-bold mb-6 ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>Статистика</h3>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              <div className={`text-center p-4 rounded-xl ${theme === 'dark' ? 'bg-blue-900' : 'bg-blue-50'}`}>
                <div className={`text-3xl font-bold mb-2 ${theme === 'dark' ? 'text-blue-300' : 'text-blue-600'}`}>{observations.length}</div>
                <div className={`text-sm font-medium ${theme === 'dark' ? 'text-blue-200' : 'text-blue-800'}`}>Всего наблюдений</div>
              </div>
              <div className={`text-center p-4 rounded-xl ${theme === 'dark' ? 'bg-green-900' : 'bg-green-50'}`}>
                <div className={`text-3xl font-bold mb-2 ${theme === 'dark' ? 'text-green-300' : 'text-green-600'}`}>
                  {groupedObservations.filter(group => group.observations.length >= 5).length}
                </div>
                <div className={`text-sm font-medium ${theme === 'dark' ? 'text-green-200' : 'text-green-800'}`}>Готово для расчёта</div>
              </div>
              <div className={`text-center p-4 rounded-xl ${theme === 'dark' ? 'bg-purple-900' : 'bg-purple-50'}`}>
                <div className={`text-3xl font-bold mb-2 ${theme === 'dark' ? 'text-purple-300' : 'text-purple-600'}`}>
                  {groupedObservations.length}
                </div>
                <div className={`text-sm font-medium ${theme === 'dark' ? 'text-purple-200' : 'text-purple-800'}`}>Комет</div>
              </div>
            </div>
          </div>
        )}
      </PageContainer>
    </div>
  )
}

export default Observations