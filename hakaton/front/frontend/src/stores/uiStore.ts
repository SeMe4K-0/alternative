import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { Toast } from '../components/Common/Toast'

interface UIState {
  // Theme
  theme: 'light' | 'dark'
  setTheme: (theme: 'light' | 'dark') => void
  toggleTheme: () => void

  // Sidebar
  sidebarOpen: boolean
  setSidebarOpen: (open: boolean) => void
  toggleSidebar: () => void

  // Notifications
  toasts: Toast[]
  addToast: (toast: Omit<Toast, 'id'>) => void
  removeToast: (id: string) => void
  clearToasts: () => void

  // Loading states
  loadingStates: Record<string, boolean>
  setLoading: (key: string, loading: boolean) => void
  isLoading: (key: string) => boolean

  // Modals
  modals: Record<string, boolean>
  openModal: (key: string) => void
  closeModal: (key: string) => void
  isModalOpen: (key: string) => boolean
}

export const useUIStore = create<UIState>()(
  persist(
    (set, get) => ({
      // Theme
      theme: 'light',
      setTheme: (theme) => set({ theme }),
      toggleTheme: () => set((state) => ({ 
        theme: state.theme === 'light' ? 'dark' : 'light' 
      })),

      // Sidebar
      sidebarOpen: false,
      setSidebarOpen: (open) => set({ sidebarOpen: open }),
      toggleSidebar: () => set((state) => ({ 
        sidebarOpen: !state.sidebarOpen 
      })),

      // Notifications
      toasts: [],
      addToast: (toast) => {
        const id = Math.random().toString(36).substr(2, 9)
        const newToast: Toast = {
          ...toast,
          id,
          duration: toast.duration || 5000
        }
        set((state) => ({
          toasts: [...state.toasts, newToast]
        }))
      },
      removeToast: (id) => set((state) => ({
        toasts: state.toasts.filter(toast => toast.id !== id)
      })),
      clearToasts: () => set({ toasts: [] }),

      // Loading states
      loadingStates: {},
      setLoading: (key, loading) => set((state) => ({
        loadingStates: {
          ...state.loadingStates,
          [key]: loading
        }
      })),
      isLoading: (key) => get().loadingStates[key] || false,

      // Modals
      modals: {},
      openModal: (key) => set((state) => ({
        modals: {
          ...state.modals,
          [key]: true
        }
      })),
      closeModal: (key) => set((state) => ({
        modals: {
          ...state.modals,
          [key]: false
        }
      })),
      isModalOpen: (key) => get().modals[key] || false,
    }),
    {
      name: 'ui-storage',
      partialize: (state) => ({
        theme: state.theme,
        sidebarOpen: state.sidebarOpen,
      }),
    }
  )
)

// Helper functions disabled - no toasts
export const showSuccess = (title: string, message?: string) => {
  // Уведомления отключены
}

export const showError = (title: string, message?: string) => {
  // Уведомления отключены
}

export const showWarning = (title: string, message?: string) => {
  // Уведомления отключены
}

export const showInfo = (title: string, message?: string) => {
  // Уведомления отключены
}
