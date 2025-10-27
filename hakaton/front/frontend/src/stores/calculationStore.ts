import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { Calculation, CalculationState } from '../types'
import { showSuccess, showError } from './uiStore'
import { apiService } from '../services/api'

interface CalculationStore extends CalculationState {
  // Actions
  createCalculation: (observationIds: string[], cometId?: number) => Promise<void>
  getCalculation: (id: string) => Calculation | undefined
  deleteCalculation: (id: string) => Promise<void>
  loadCalculations: () => Promise<void>
  setCurrentCalculation: (calculation: Calculation | null) => void
  clearError: () => void
  setLoading: (loading: boolean) => void
  setError: (error: string) => void
}

export const useCalculationStore = create<CalculationStore>()(
  persist(
    (set, get) => ({
      calculations: [],
      currentCalculation: null,
      isLoading: false,
      error: null,

      createCalculation: async (observationIds, cometId) => {
        set({ isLoading: true, error: null })
        
        try {
          const newCalculation = await apiService.createCalculation(observationIds, cometId)

          set((state) => ({
            calculations: [newCalculation, ...state.calculations],
            currentCalculation: newCalculation,
            isLoading: false,
          }))

          showSuccess('Расчёт завершён', 'Параметры орбиты успешно рассчитаны')
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Ошибка при создании расчета',
            isLoading: false 
          })
          showError('Ошибка расчёта', 'Не удалось рассчитать параметры орбиты')
        }
      },

      getCalculation: (id) => {
        const { calculations } = get()
        return calculations.find(calc => calc.id === id)
      },

      deleteCalculation: async (id) => {
        set({ isLoading: true, error: null })
        
        try {
          const calculation = get().calculations.find(calc => calc.id === id)
          await apiService.deleteCalculation(id)
          
          set((state) => ({
            calculations: state.calculations.filter((calc) => calc.id !== id),
            currentCalculation: state.currentCalculation?.id === id ? null : state.currentCalculation,
            isLoading: false,
          }))

          showSuccess('Расчёт удалён', 'Расчёт успешно удалён')
        } catch (error) {
          set({ isLoading: false, error: 'Ошибка при удалении расчёта' })
          showError('Ошибка', 'Не удалось удалить расчёт')
        }
      },

      loadCalculations: async () => {
        set({ isLoading: true, error: null })
        
        try {
          const calculations = await apiService.getCalculations()
          set({ calculations, isLoading: false })
        } catch (error) {
          set({ isLoading: false, error: 'Ошибка при загрузке расчётов' })
          showError('Ошибка', 'Не удалось загрузить расчёты')
        }
      },

      setCurrentCalculation: (calculation) => {
        set({ currentCalculation: calculation })
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
      name: 'calculation-storage',
      partialize: (state) => ({
        calculations: state.calculations,
        currentCalculation: state.currentCalculation,
      }),
    }
  )
)
