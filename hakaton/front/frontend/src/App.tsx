import React, { useEffect } from 'react'
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom'
import Layout from './components/Layout/Layout'
import ProtectedRoute from './components/Auth/ProtectedRoute'
import Login from './pages/Auth/Login'
import Dashboard from './pages/Dashboard/Dashboard'
import CometSelection from './pages/Comets/CometSelection'
import EditComet from './pages/Comets/EditComet'
import AddObservation from './pages/Observations/AddObservation'
import EditObservation from './pages/Observations/EditObservation'
import Observations from './pages/Observations/Observations'
import ObservationDetails from './pages/Observations/ObservationDetails'
import Calculations from './pages/Calculations/Calculations'
import History from './pages/History/History'
import Profile from './pages/Profile/Profile'
import { useUIStore } from './stores/uiStore'
import { useAuthStore } from './stores/authStore'
import { useObservationStore } from './stores/observationStore'
import { useCalculationStore } from './stores/calculationStore'
import { useCometStore } from './stores/cometStore'

function App() {
  const { theme } = useUIStore()
  const { isAuthenticated, checkAuth } = useAuthStore()
  const { loadObservations } = useObservationStore()
  const { loadCalculations } = useCalculationStore()
  const { loadComets } = useCometStore()

  useEffect(() => {
    // Применяем тему к корневому элементу документа
    document.documentElement.classList.toggle('dark', theme === 'dark')
  }, [theme])

  useEffect(() => {
    // Проверяем аутентификацию при загрузке приложения
    checkAuth()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  useEffect(() => {
    // Загружаем данные при успешной аутентификации
    if (isAuthenticated) {
      loadObservations()
      loadCalculations()
      loadComets()
    }
  }, [isAuthenticated, loadObservations, loadCalculations, loadComets])

  return (
    <Router>
      <Routes>
        {/* Public routes */}
        <Route path="/login" element={<Login />} />
        
        {/* Protected routes */}
        <Route path="/" element={
          <ProtectedRoute>
            <Layout />
          </ProtectedRoute>
        }>
          <Route index element={<Navigate to="/dashboard" replace />} />
          <Route path="dashboard" element={<Dashboard />} />
          <Route path="comets/select" element={<CometSelection />} />
          <Route path="comets/:id/edit" element={<EditComet />} />
          <Route path="observations" element={<Observations />} />
          <Route path="observations/new" element={<AddObservation />} />
          <Route path="observations/:id" element={<ObservationDetails />} />
          <Route path="observations/:id/edit" element={<EditObservation />} />
          <Route path="calculations" element={<Calculations />} />
          <Route path="history" element={<History />} />
          <Route path="profile" element={<Profile />} />
        </Route>
        
        {/* Redirect to login for unknown routes */}
        <Route path="*" element={<Navigate to="/login" replace />} />
      </Routes>
    </Router>
  )
}

export default App