import React, { useState, useEffect } from 'react'
import { useAuthStore } from '../../stores/authStore'
import { useObservationStore } from '../../stores/observationStore'
import { useCometStore } from '../../stores/cometStore'
import { useCalculationStore } from '../../stores/calculationStore'
import { useUIStore } from '../../stores/uiStore'
import { apiService } from '../../services/api'
import { showSuccess, showError } from '../../stores/uiStore'
import type { User } from '../../types'
import PageContainer from '../../components/Common/PageContainer'
import Button from '../../components/Common/Button'

const Profile: React.FC = () => {
  const { user, logout, updateUser } = useAuthStore()
  const { observations } = useObservationStore()
  const { comets } = useCometStore()
  const { calculations } = useCalculationStore()
  const { theme } = useUIStore()
  const [isEditing, setIsEditing] = useState(false)
  const [isChangingPassword, setIsChangingPassword] = useState(false)
  const [avatarFile, setAvatarFile] = useState<File | null>(null)
  const [avatarPreview, setAvatarPreview] = useState<string>(user?.avatar || '')
  
  // Обновляем preview аватара при изменении user
  useEffect(() => {
    if (user?.avatar) {
      console.log('Updating avatarPreview from user.avatar:', user.avatar)
      // Используем полный URL из user.avatar
      setAvatarPreview(user.avatar)
    } else {
      console.log('No avatar in user, clearing avatarPreview')
      setAvatarPreview('')
    }
  }, [user])
  
  // Обновляем formData при изменении user
  useEffect(() => {
    if (user) {
      setFormData({
        name: user.name || '',
        email: user.email || '',
      })
    }
  }, [user])
  
  const [isLoading, setIsLoading] = useState(false)
  const [formData, setFormData] = useState({
    name: user?.name || '',
    email: user?.email || '',
  })
  const [passwordData, setPasswordData] = useState({
    currentPassword: '',
    newPassword: '',
    confirmPassword: '',
  })
  const [passwordError, setPasswordError] = useState('')

  const handleAvatarChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (file) {
      setAvatarFile(file)
      const reader = new FileReader()
      reader.onloadend = () => {
        setAvatarPreview(reader.result as string)
      }
      reader.readAsDataURL(file)
    }
  }

  const handleSave = async () => {
    setIsLoading(true)
    try {
      const updateData: any = {}
      
      // Всегда обновляем, если есть изменения
      if (formData.name !== user?.name) {
        updateData.username = formData.name
      }
      
      if (avatarFile) {
        updateData.avatar = avatarFile
      }
      
      // Отправляем запрос, даже если только имя изменилось
      console.log('Отправляем данные на обновление:', updateData)
      console.log('Текущий user:', user)
      console.log('formData.name:', formData.name, 'user?.name:', user?.name)
      console.log('avatarFile:', avatarFile)
      
      if (updateData.username || updateData.avatar) {
        try {
          const updatedUser = await apiService.updateProfile(updateData)
          console.log('Обновленный пользователь:', updatedUser)
          
          // Сразу обновляем данные пользователя в store
          if (updatedUser) {
            console.log('Updated user from server:', updatedUser)
            console.log('Updated user avatar:', updatedUser.avatar)
            
            // Обновляем authStore только с нужными полями
            const updatedUserData: User = {
              ...user!, // Сохраняем существующие данные
              name: updatedUser.name || user!.name, // Обновляем имя если изменилось
              avatar: updatedUser.avatar || user!.avatar, // Обновляем аватар если изменилось
            }
            
            updateUser(updatedUserData)
            
            showSuccess('Профиль обновлён', 'Изменения сохранены успешно')
            
            setIsEditing(false)
            
            // Сбрасываем файл аватара
            setAvatarFile(null)
            
            // Обновляем formData
            setFormData({
              name: updatedUser.name || user!.name,
              email: user!.email
            })
          }
        } catch (updateError) {
          console.error('Ошибка в updateProfile:', updateError)
          throw updateError
        }
      } else {
        console.log('Нет изменений для отправки')
        showSuccess('Профиль обновлён', 'Изменения сохранены успешно')
      }
    } catch (error) {
      console.error('Ошибка при сохранении профиля:', error)
      showError('Ошибка', 'Не удалось сохранить изменения')
      setIsEditing(false)
    } finally {
      setIsLoading(false)
    }
  }

  const handleCancel = () => {
    setFormData({
      name: user?.name || '',
      email: user?.email || '',
    })
    setAvatarFile(null)
    setAvatarPreview(user?.avatar || '')
    setIsEditing(false)
  }

  const handlePasswordChange = async () => {
    setPasswordError('')
    
    if (passwordData.newPassword !== passwordData.confirmPassword) {
      setPasswordError('Новые пароли не совпадают')
      return
    }
    
    if (passwordData.newPassword.length < 6) {
      setPasswordError('Пароль должен содержать минимум 6 символов')
      return
    }
    
    setIsLoading(true)
    try {
      await apiService.updateProfile({
        password: passwordData.newPassword
      })
      
      showSuccess('Пароль изменён', 'Ваш пароль был успешно изменён')
      setPasswordData({
        currentPassword: '',
        newPassword: '',
        confirmPassword: '',
      })
      setIsChangingPassword(false)
    } catch (error) {
      console.error('Ошибка при изменении пароля:', error)
      showError('Ошибка', 'Не удалось изменить пароль')
    } finally {
      setIsLoading(false)
    }
  }

  const handleCancelPasswordChange = () => {
    setPasswordData({
      currentPassword: '',
      newPassword: '',
      confirmPassword: '',
    })
    setPasswordError('')
    setIsChangingPassword(false)
  }

  return (
    <div className={`min-h-screen bg-gradient-to-br ${theme === 'dark' ? 'from-gray-900 to-gray-800' : 'from-blue-50 to-purple-50'}`}>
      <PageContainer>
        {/* Header */}
        <div className="mb-8">
          <div className={`rounded-2xl shadow-xl border overflow-hidden ${theme === 'dark' ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-100'}`}>
            <div className="bg-gradient-to-r from-blue-600 to-purple-600 px-8 py-6">
              <h1 className="text-3xl font-bold text-white mb-2">Профиль пользователя</h1>
              <p className="text-blue-100 text-lg">Управление личными данными и настройками</p>
            </div>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Основная информация */}
          <div className="lg:col-span-2">
            <div className={`rounded-2xl shadow-lg border p-8 ${theme === 'dark' ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-100'}`}>
              <div className="flex justify-between items-center mb-6">
                <h2 className={`text-2xl font-bold ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>Личная информация</h2>
                {!isEditing && (
                  <Button
                    onClick={() => {
                      setIsEditing(true)
                      // Сбрасываем состояние при начале редактирования
                      setAvatarFile(null)
                      // Обновляем preview аватара при начале редактирования с правильным URL
                      if (user?.avatar) {
                        console.log('Setting avatarPreview for edit mode:', user.avatar)
                        setAvatarPreview(user.avatar)
                      } else {
                        setAvatarPreview('')
                      }
                    }}
                    className="bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded-xl"
                  >
                    Редактировать
                  </Button>
                )}
              </div>

              <div className="space-y-6">
                {/* Avatar Upload */}
                <div>
                  <label className={`block text-sm font-medium mb-2 ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>
                    Аватар
                  </label>
                  {isEditing ? (
                    <div className="flex items-center space-x-4">
                      <div className="relative">
                        {avatarPreview ? (
                          <img
                            key={`edit-${avatarPreview}`}
                            src={`${avatarPreview}${avatarPreview.startsWith('http') ? '?t=' + Date.now() : ''}`}
                            alt="Avatar preview"
                            className="w-24 h-24 rounded-full object-cover border-4 border-indigo-500"
                          />
                        ) : (
                          <div className="w-24 h-24 rounded-full bg-gradient-to-r from-blue-500 to-purple-500 flex items-center justify-center">
                            <span className="text-white text-2xl font-bold">
                              {user?.name?.charAt(0).toUpperCase() || 'U'}
                            </span>
                          </div>
                        )}
                        <label className="absolute bottom-0 right-0 bg-indigo-600 hover:bg-indigo-700 text-white p-2 rounded-full cursor-pointer transition-all">
                          <span className="text-xs">📷</span>
                          <input
                            type="file"
                            accept="image/*"
                            onChange={handleAvatarChange}
                            className="hidden"
                          />
                        </label>
                      </div>
                      <div>
                        <p className={`text-sm ${theme === 'dark' ? 'text-gray-400' : 'text-gray-600'}`}>
                          Нажмите на иконку камеры, чтобы изменить аватар
                        </p>
                      </div>
                    </div>
                  ) : (
                    <div className="flex items-center">
                      {user?.avatar ? (
                        <img
                          src={`${user.avatar}?t=${Date.now()}`}
                          alt="Avatar"
                          key={user.avatar}
                          className="w-24 h-24 rounded-full object-cover border-4 border-indigo-500"
                        />
                      ) : (
                        <div className="w-24 h-24 rounded-full bg-gradient-to-r from-blue-500 to-purple-500 flex items-center justify-center">
                          <span className="text-white text-2xl font-bold">
                            {user?.name?.charAt(0).toUpperCase() || 'U'}
                          </span>
                        </div>
                      )}
                    </div>
                  )}
                </div>

                <div>
                  <label className={`block text-sm font-medium mb-2 ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>
                    Имя пользователя
                  </label>
                  {isEditing ? (
                    <input
                      type="text"
                      value={formData.name}
                      onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                      className={`w-full px-4 py-3 rounded-xl focus:ring-2 focus:ring-indigo-500 ${
                        theme === 'dark' 
                          ? 'bg-gray-700 border-gray-600 text-white focus:border-indigo-500' 
                          : 'border border-gray-300 focus:border-indigo-500'
                      }`}
                    />
                  ) : (
                    <div className={`px-4 py-3 rounded-xl ${theme === 'dark' ? 'bg-gray-700 text-white' : 'bg-gray-50 text-gray-900'}`}>
                      {user?.name || 'Не указано'}
                    </div>
                  )}
                </div>

                <div>
                  <label className={`block text-sm font-medium mb-2 ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>
                    Email адрес
                  </label>
                  <div className={`px-4 py-3 rounded-xl ${theme === 'dark' ? 'bg-gray-700 text-white' : 'bg-gray-50 text-gray-900'}`}>
                    {user?.email || 'Не указано'}
                  </div>
                </div>

                <div>
                  <label className={`block text-sm font-medium mb-2 ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>
                    Дата регистрации
                  </label>
                  <div className={`px-4 py-3 rounded-xl ${theme === 'dark' ? 'bg-gray-700 text-white' : 'bg-gray-50 text-gray-900'}`}>
                    {user?.createdAt ? new Date(user.createdAt).toLocaleDateString('ru-RU') : 'Не указано'}
                  </div>
                </div>

                {isEditing && (
                  <div className="flex space-x-4 pt-4">
                    <Button
                      onClick={handleSave}
                      className="bg-green-600 hover:bg-green-700 text-white px-6 py-3 rounded-xl"
                      disabled={isLoading}
                    >
                      {isLoading ? 'Сохранение...' : 'Сохранить'}
                    </Button>
                    <Button
                      onClick={handleCancel}
                      variant="secondary"
                      className="px-6 py-3 rounded-xl"
                    >
                      Отмена
                    </Button>
                  </div>
                )}
              </div>
            </div>

            {/* Password Change Section */}
            <div className={`rounded-2xl shadow-lg border p-8 mt-6 ${theme === 'dark' ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-100'}`}>
              <div className="flex justify-between items-center mb-6">
                <h2 className={`text-2xl font-bold ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>Безопасность</h2>
                {!isChangingPassword && (
                  <Button
                    onClick={() => setIsChangingPassword(true)}
                    className="bg-orange-600 hover:bg-orange-700 text-white px-4 py-2 rounded-xl"
                  >
                    Изменить пароль
                  </Button>
                )}
              </div>

              {isChangingPassword ? (
                <div className="space-y-6">
                  <div>
                    <label className={`block text-sm font-medium mb-2 ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>
                      Текущий пароль
                    </label>
                    <input
                      type="password"
                      value={passwordData.currentPassword}
                      onChange={(e) => setPasswordData({ ...passwordData, currentPassword: e.target.value })}
                      className={`w-full px-4 py-3 rounded-xl focus:ring-2 focus:ring-indigo-500 ${
                        theme === 'dark' 
                          ? 'bg-gray-700 border-gray-600 text-white focus:border-indigo-500' 
                          : 'border border-gray-300 focus:border-indigo-500'
                      }`}
                      placeholder="Введите текущий пароль"
                    />
                  </div>

                  <div>
                    <label className={`block text-sm font-medium mb-2 ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>
                      Новый пароль
                    </label>
                    <input
                      type="password"
                      value={passwordData.newPassword}
                      onChange={(e) => setPasswordData({ ...passwordData, newPassword: e.target.value })}
                      className={`w-full px-4 py-3 rounded-xl focus:ring-2 focus:ring-indigo-500 ${
                        theme === 'dark' 
                          ? 'bg-gray-700 border-gray-600 text-white focus:border-indigo-500' 
                          : 'border border-gray-300 focus:border-indigo-500'
                      }`}
                      placeholder="Введите новый пароль"
                    />
                  </div>

                  <div>
                    <label className={`block text-sm font-medium mb-2 ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>
                      Подтвердите новый пароль
                    </label>
                    <input
                      type="password"
                      value={passwordData.confirmPassword}
                      onChange={(e) => setPasswordData({ ...passwordData, confirmPassword: e.target.value })}
                      className={`w-full px-4 py-3 rounded-xl focus:ring-2 focus:ring-indigo-500 ${
                        theme === 'dark' 
                          ? 'bg-gray-700 border-gray-600 text-white focus:border-indigo-500' 
                          : 'border border-gray-300 focus:border-indigo-500'
                      }`}
                      placeholder="Подтвердите новый пароль"
                    />
                  </div>

                  {passwordError && (
                    <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded-xl">
                      {passwordError}
                    </div>
                  )}

                  <div className="flex space-x-4 pt-4">
                    <Button
                      onClick={handlePasswordChange}
                      className="bg-green-600 hover:bg-green-700 text-white px-6 py-3 rounded-xl"
                      disabled={isLoading}
                    >
                      {isLoading ? 'Изменение...' : 'Изменить пароль'}
                    </Button>
                    <Button
                      onClick={handleCancelPasswordChange}
                      variant="secondary"
                      className="px-6 py-3 rounded-xl"
                    >
                      Отмена
                    </Button>
                  </div>
                </div>
              ) : (
                <div className={theme === 'dark' ? 'text-gray-400' : 'text-gray-600'}>
                  <p>Для изменения пароля нажмите кнопку "Изменить пароль"</p>
                </div>
              )}
            </div>
          </div>

          {/* Боковая панель */}
          <div className="space-y-6">
            {/* Статистика */}
            <div className={`rounded-2xl shadow-lg border p-6 ${theme === 'dark' ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-100'}`}>
              <h3 className={`text-xl font-bold mb-4 ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>Статистика</h3>
              <div className="space-y-3">
                <div className="flex justify-between">
                  <span className={theme === 'dark' ? 'text-gray-400' : 'text-gray-600'}>Наблюдений:</span>
                  <span className={`font-semibold ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>{observations.length}</span>
                </div>
                <div className="flex justify-between">
                  <span className={theme === 'dark' ? 'text-gray-400' : 'text-gray-600'}>Расчётов:</span>
                  <span className={`font-semibold ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>{calculations.length}</span>
                </div>
                <div className="flex justify-between">
                  <span className={theme === 'dark' ? 'text-gray-400' : 'text-gray-600'}>Комет:</span>
                  <span className={`font-semibold ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>{comets.length}</span>
                </div>
              </div>
            </div>

            {/* Действия */}
            <div className={`rounded-2xl shadow-lg border p-6 ${theme === 'dark' ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-100'}`}>
              <h3 className={`text-xl font-bold mb-4 ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>Действия</h3>
              <div className="space-y-3">
                <Button
                  onClick={logout}
                  className="w-full bg-red-600 hover:bg-red-700 text-white py-3 rounded-xl"
                >
                  Выйти из аккаунта
                </Button>
              </div>
            </div>
          </div>
        </div>
      </PageContainer>
    </div>
  )
}

export default Profile
