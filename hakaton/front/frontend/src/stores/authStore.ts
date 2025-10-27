import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { apiService } from '../services/api'
import type { AuthState, User } from '../types'
import { showSuccess, showError } from './uiStore'

interface AuthStore extends AuthState {
  login: (email: string, password: string) => Promise<void>
  register: (userData: { name: string; email: string; password: string }) => Promise<void>
  logout: () => void
  checkAuth: () => Promise<void>
  updateUser: (user: User) => void
}

export const useAuthStore = create<AuthStore>()(
  persist(
    (set, get) => ({
      user: null,
      token: null,
      isAuthenticated: false,
      isLoading: false,

      login: async (email: string, password: string) => {
        set({ isLoading: true })
        try {
          const { user, token } = await apiService.login(email, password)
          console.log('Initial user from login:', user)
          console.log('User avatar in store:', user.avatar)
          
          // Загружаем актуальные данные профиля с сервера
          try {
            const currentUser = await apiService.getCurrentUser()
            console.log('Current user from getCurrentUser:', currentUser)
            console.log('Current user avatar:', currentUser.avatar)
            
            set({ user: currentUser, token, isAuthenticated: true, isLoading: false })
            showSuccess('Вход выполнен', `Добро пожаловать, ${currentUser.name}!`)
          } catch (getUserError) {
            console.warn('getCurrentUser failed, using login data:', getUserError)
            set({ user, token, isAuthenticated: true, isLoading: false })
            showSuccess('Вход выполнен', `Добро пожаловать, ${user.name}!`)
          }
        } catch (error) {
          set({ isLoading: false })
          showError('Ошибка входа', 'Неверные данные для входа')
          throw error
        }
      },

      register: async (userData) => {
        set({ isLoading: true })
        try {
          const { user, token } = await apiService.register(userData)
          set({ user, token, isAuthenticated: true, isLoading: false })
          showSuccess('Регистрация успешна', `Добро пожаловать, ${user.name}!`)
        } catch (error) {
          set({ isLoading: false })
          showError('Ошибка регистрации', 'Не удалось создать аккаунт')
          throw error
        }
      },

      logout: async () => {
        try {
          await apiService.logout()
        } catch (error) {
          console.error('Logout error:', error)
        } finally {
          set({ user: null, token: null, isAuthenticated: false, isLoading: false })
          showSuccess('Выход выполнен', 'До свидания!')
        }
      },

      checkAuth: async () => {
        // Zustand persist восстанавливает состояние автоматически
        // Просто убеждаемся, что isLoading = false
        set({ isLoading: false })
      },

      updateUser: (user) => {
        set({ user })
      },
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({
        user: state.user,
        token: state.token,
        isAuthenticated: state.isAuthenticated,
      }),
    }
  )
)
