import React from 'react'
import { useParams, useNavigate, Link } from 'react-router-dom'
import { useObservationStore } from '../../stores/observationStore'
import { useUIStore } from '../../stores/uiStore'
import PageContainer from '../../components/Common/PageContainer'
import Button from '../../components/Common/Button'
import Card from '../../components/Common/Card'
import Modal from '../../components/Common/Modal'
import LoadingSpinner from '../../components/Common/LoadingSpinner'
import { format } from 'date-fns'
import { ru } from 'date-fns/locale'
import { ArrowLeftIcon, PencilIcon, TrashIcon } from '@heroicons/react/24/outline'

const ObservationDetails: React.FC = () => {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const { observations, deleteObservation, isLoading } = useObservationStore()
  const { openModal, closeModal, isModalOpen, theme } = useUIStore()

  const observation = observations.find(obs => obs.id === id)

  const handleDelete = async () => {
    if (!observation) return
    
    try {
      await deleteObservation(observation.id)
      closeModal('deleteObservation')
      navigate('/observations')
    } catch (error) {
      console.error('Ошибка при удалении наблюдения:', error)
    }
  }

  const handleEdit = () => {
    navigate(`/observations/${id}/edit`)
  }

  if (isLoading) {
    return (
      <div className={`min-h-screen bg-gradient-to-br ${theme === 'dark' ? 'from-gray-900 to-gray-800' : 'from-blue-50 to-purple-50'} flex items-center justify-center`}>
        <LoadingSpinner size="lg" text="Загрузка наблюдения..." />
      </div>
    )
  }

  if (!observation) {
    return (
      <div className={`min-h-screen bg-gradient-to-br ${theme === 'dark' ? 'from-gray-900 to-gray-800' : 'from-blue-50 to-purple-50'}`}>
        <PageContainer>
          <div className="text-center py-12">
            <h1 className={`text-3xl font-bold mb-4 ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>Наблюдение не найдено</h1>
            <p className={`mb-8 ${theme === 'dark' ? 'text-gray-400' : 'text-gray-600'}`}>Возможно, это наблюдение было удалено или не существует.</p>
            <Link to="/observations">
              <Button className="bg-indigo-600 hover:bg-indigo-700 text-white px-6 py-3 rounded-xl">
                Вернуться к списку
              </Button>
            </Link>
          </div>
        </PageContainer>
      </div>
    )
  }

  return (
    <div className={`min-h-screen bg-gradient-to-br ${theme === 'dark' ? 'from-gray-900 to-gray-800' : 'from-blue-50 to-purple-50'}`}>
      <PageContainer>
        {/* Header */}
        <div className="mb-8">
          <div className={`rounded-2xl shadow-xl border overflow-hidden ${theme === 'dark' ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-100'}`}>
            <div className="bg-gradient-to-r from-emerald-600 to-teal-600 px-8 py-6">
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-4">
                  <Link to="/observations">
                    <Button variant="secondary" className="flex items-center space-x-2">
                      <ArrowLeftIcon className="h-5 w-5" />
                      <span>Назад</span>
                    </Button>
                  </Link>
                  <div>
                    <h1 className="text-3xl font-bold text-white mb-2">{observation.cometName}</h1>
                    <p className="text-emerald-100 text-lg">Детали наблюдения кометы</p>
                  </div>
                </div>
                <div className="flex flex-wrap gap-3">
                  <Button
                    onClick={handleEdit}
                    className="bg-emerald-600 hover:bg-emerald-700 text-white font-bold shadow-lg px-4 py-2 rounded-xl flex items-center space-x-2 transition-all border-2 border-white"
                  >
                    <PencilIcon className="h-5 w-5" />
                    <span>Редактировать</span>
                  </Button>
                  <Button
                    onClick={() => openModal('deleteObservation')}
                    className="bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded-xl flex items-center space-x-2"
                  >
                    <TrashIcon className="h-5 w-5" />
                    <span>Удалить</span>
                  </Button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          {/* Image Section */}
          {observation.imageUrl && (
            <Card>
              <h2 className={`text-2xl font-bold mb-6 ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>Изображение наблюдения</h2>
              <div className="space-y-4">
                <div className="relative w-full rounded-xl overflow-hidden" style={{ aspectRatio: '4/3' }}>
                  <img
                    src={observation.imageUrl}
                    alt={`Изображение наблюдения ${observation.cometName}`}
                    className="w-full h-full object-cover"
                  />
                </div>
              </div>
            </Card>
          )}

          {/* Details Section */}
          <div className="space-y-6">
            {/* Basic Info */}
            <Card>
              <h2 className={`text-2xl font-bold mb-6 ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>Основная информация</h2>
              <div className="space-y-4">
                <div>
                  <label className={`block text-sm font-medium mb-1 ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>Название кометы</label>
                  <div className={`px-4 py-3 rounded-xl font-semibold ${theme === 'dark' ? 'bg-gray-700 text-white' : 'bg-gray-50 text-gray-900'}`}>
                    {observation.cometName}
                  </div>
                </div>
                <div>
                  <label className={`block text-sm font-medium mb-1 ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>Дата и время наблюдения</label>
                  <div className={`px-4 py-3 rounded-xl ${theme === 'dark' ? 'bg-gray-700 text-white' : 'bg-gray-50 text-gray-900'}`}>
                    {format(new Date(observation.observationDate), 'dd MMMM yyyy, HH:mm', { locale: ru })}
                  </div>
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className={`block text-sm font-medium mb-1 ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>Прямое восхождение (RA)</label>
                    <div className={`px-4 py-3 rounded-xl ${theme === 'dark' ? 'bg-gray-700 text-white' : 'bg-gray-50 text-gray-900'}`}>
                      {observation.coordinates.ra} ч
                    </div>
                  </div>
                  <div>
                    <label className={`block text-sm font-medium mb-1 ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>Склонение (Dec)</label>
                    <div className={`px-4 py-3 rounded-xl ${theme === 'dark' ? 'bg-gray-700 text-white' : 'bg-gray-50 text-gray-900'}`}>
                      {observation.coordinates.dec}°
                    </div>
                  </div>
                </div>
              </div>
            </Card>

            {/* Notes */}
            {observation.notes && (
              <Card>
                <h2 className={`text-2xl font-bold mb-6 ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>Заметки</h2>
                <div className={`px-4 py-3 rounded-xl whitespace-pre-wrap ${theme === 'dark' ? 'bg-gray-700 text-white' : 'bg-gray-50 text-gray-900'}`}>
                  {observation.notes}
                </div>
              </Card>
            )}

            {/* Metadata */}
            <Card>
              <h2 className={`text-2xl font-bold mb-6 ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>Метаданные</h2>
              <div className="space-y-3">
                <div className="flex justify-between">
                  <span className={theme === 'dark' ? 'text-gray-400' : 'text-gray-600'}>Создано:</span>
                  <span className={`font-medium ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>
                    {format(new Date(observation.createdAt), 'dd.MM.yyyy HH:mm', { locale: ru })}
                  </span>
                </div>
                <div className="flex justify-between">
                  <span className={theme === 'dark' ? 'text-gray-400' : 'text-gray-600'}>Обновлено:</span>
                  <span className={`font-medium ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>
                    {format(new Date(observation.updatedAt), 'dd.MM.yyyy HH:mm', { locale: ru })}
                  </span>
                </div>
              </div>
            </Card>
          </div>
        </div>

        {/* Delete Confirmation Modal */}
        <Modal
          isOpen={isModalOpen('deleteObservation')}
          onClose={() => closeModal('deleteObservation')}
          title="Подтверждение удаления"
        >
          <p className={`mb-6 ${theme === 'dark' ? 'text-gray-400' : 'text-gray-600'}`}>
            Вы уверены, что хотите удалить наблюдение "{observation?.cometName}"? 
            Это действие нельзя отменить.
          </p>
          <div className="flex space-x-4">
            <Button
              onClick={() => closeModal('deleteObservation')}
              variant="secondary"
              className="flex-1"
            >
              Отмена
            </Button>
            <Button
              onClick={handleDelete}
              className="flex-1 bg-red-600 hover:bg-red-700 text-white"
            >
              Удалить
            </Button>
          </div>
        </Modal>

      </PageContainer>
    </div>
  )
}

export default ObservationDetails
