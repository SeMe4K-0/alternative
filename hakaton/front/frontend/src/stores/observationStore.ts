import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { Observation, ObservationState } from '../types'
import { showSuccess, showError } from './uiStore'
import { apiService } from '../services/api'

interface ObservationStore extends ObservationState {
  // Actions
  addObservation: (observation: Omit<Observation, 'id' | 'createdAt' | 'updatedAt'>) => Promise<void>
  updateObservation: (id: string, updates: Partial<Observation>) => Promise<void>
  deleteObservation: (id: string) => Promise<void>
  loadObservations: () => Promise<void>
  setCurrentObservation: (observation: Observation | null) => void
  clearError: () => void
  setLoading: (loading: boolean) => void
  setError: (error: string) => void
}

export const useObservationStore = create<ObservationStore>()(
  persist(
    (set, get) => ({
      observations: [],
      currentObservation: null,
      isLoading: false,
      error: null,

      addObservation: async (observationData) => {
        set({ isLoading: true, error: null })
        
        try {
          const newObservation = await apiService.createObservation(observationData)

          set((state) => ({
            observations: [...state.observations, newObservation],
            isLoading: false,
          }))

          showSuccess('Наблюдение добавлено', `Наблюдение для кометы "${observationData.cometName}" успешно добавлено`)
        } catch (error) {
          set({ isLoading: false, error: 'Ошибка при добавлении наблюдения' })
          showError('Ошибка', 'Не удалось добавить наблюдение')
        }
      },

      updateObservation: async (id, updates) => {
        set({ isLoading: true, error: null })
        
        try {
          const updatedObservation = await apiService.updateObservation(id, updates)
          
          set((state) => ({
            observations: state.observations.map((obs) =>
              obs.id === id ? updatedObservation : obs
            ),
            currentObservation: state.currentObservation?.id === id
              ? updatedObservation
              : state.currentObservation,
            isLoading: false,
          }))

          showSuccess('Наблюдение обновлено', 'Изменения сохранены успешно')
        } catch (error) {
          set({ isLoading: false, error: 'Ошибка при обновлении наблюдения' })
          showError('Ошибка', 'Не удалось обновить наблюдение')
        }
      },

      deleteObservation: async (id) => {
        set({ isLoading: true, error: null })
        
        try {
          const observation = get().observations.find(obs => obs.id === id)
          await apiService.deleteObservation(id)
          
          set((state) => ({
            observations: state.observations.filter((obs) => obs.id !== id),
            currentObservation: state.currentObservation?.id === id ? null : state.currentObservation,
            isLoading: false,
          }))

          showSuccess('Наблюдение удалено', `Комета "${observation?.cometName}" удалена`)
        } catch (error) {
          set({ isLoading: false, error: 'Ошибка при удалении наблюдения' })
          showError('Ошибка', 'Не удалось удалить наблюдение')
        }
      },

      loadObservations: async () => {
        set({ isLoading: true, error: null })
        
        try {
          const observations = await apiService.getObservations()
          set({ observations, isLoading: false })
        } catch (error) {
          set({ isLoading: false, error: 'Ошибка при загрузке наблюдений' })
          showError('Ошибка', 'Не удалось загрузить наблюдения')
        }
      },

      setCurrentObservation: (observation) => {
        set({ currentObservation: observation })
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
      name: 'observation-storage',
      partialize: (state) => ({
        observations: state.observations,
        currentObservation: state.currentObservation,
      }),
    }
  )
)