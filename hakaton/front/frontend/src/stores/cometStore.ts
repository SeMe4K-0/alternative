import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { Comet } from '../types'
import { showSuccess, showError } from './uiStore'
import { apiService } from '../services/api'

interface CometState {
  comets: Comet[]
  currentComet: Comet | null
  isLoading: boolean
  error: string | null
}

interface CometStore extends CometState {
  // Actions
  createComet: (comet: Omit<Comet, 'id' | 'createdAt' | 'updatedAt'>) => Promise<void>
  updateComet: (id: string, updates: Partial<Comet>) => Promise<void>
  deleteComet: (id: string) => Promise<void>
  loadComets: () => Promise<void>
  setCurrentComet: (comet: Comet | null) => void
  clearError: () => void
  setLoading: (loading: boolean) => void
  setError: (error: string) => void
}

export const useCometStore = create<CometStore>()(
  persist(
    (set, get) => ({
      comets: [],
      currentComet: null,
      isLoading: false,
      error: null,

      createComet: async (cometData) => {
        set({ isLoading: true, error: null })
        
        try {
          const newComet = await apiService.createComet(cometData)

          set((state) => ({
            comets: [...state.comets, newComet],
            isLoading: false,
          }))

          showSuccess('Комета создана', `Комета "${cometData.name}" успешно создана`)
        } catch (error) {
          set({ isLoading: false, error: 'Ошибка при создании кометы' })
          showError('Ошибка', 'Не удалось создать комету')
        }
      },

      updateComet: async (id, updates) => {
        set({ isLoading: true, error: null })
        
        try {
          const updatedComet = await apiService.updateComet(id, updates)
          
          set((state) => ({
            comets: state.comets.map((comet) =>
              comet.id === id ? updatedComet : comet
            ),
            currentComet: state.currentComet?.id === id
              ? updatedComet
              : state.currentComet,
            isLoading: false,
          }))

          showSuccess('Комета обновлена', 'Изменения сохранены успешно')
        } catch (error) {
          set({ isLoading: false, error: 'Ошибка при обновлении кометы' })
          showError('Ошибка', 'Не удалось обновить комету')
        }
      },

      deleteComet: async (id) => {
        set({ isLoading: true, error: null })
        
        try {
          const comet = get().comets.find(comet => comet.id === id)
          await apiService.deleteComet(id)
          
          set((state) => ({
            comets: state.comets.filter((comet) => comet.id !== id),
            currentComet: state.currentComet?.id === id ? null : state.currentComet,
            isLoading: false,
          }))

          showSuccess('Комета удалена', `Комета "${comet?.name}" удалена`)
        } catch (error) {
          set({ isLoading: false, error: 'Ошибка при удалении кометы' })
          showError('Ошибка', 'Не удалось удалить комету')
        }
      },

      loadComets: async () => {
        set({ isLoading: true, error: null })
        
        try {
          const comets = await apiService.getComets()
          set({ comets, isLoading: false })
        } catch (error) {
          set({ isLoading: false, error: 'Ошибка при загрузке комет' })
          showError('Ошибка', 'Не удалось загрузить кометы')
        }
      },

      setCurrentComet: (comet) => {
        set({ currentComet: comet })
      },

      clearError: () => {
        set({ error: null })
      },

      setLoading: (loading) => {
        set({ isLoading: loading })
      },

      setError: (error) => {
        set({ error, isLoading: false })
      },
    }),
    {
      name: 'comet-storage',
      partialize: (state) => ({
        comets: state.comets,
        currentComet: state.currentComet,
      }),
    }
  )
)
