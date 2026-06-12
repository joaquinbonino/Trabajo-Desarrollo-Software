import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider, useAuth } from './context/AuthContext'
import HomePage from './pages/HomePage'
import EventDetailPage from './pages/EventDetailPage'
import MyTicketsPage from './pages/MyTicketsPage'
import LoginPage from './pages/LoginPage'
import RegisterPage from './pages/RegisterPage'
import AdminEventsPage from './pages/AdminEventsPage'
import EventFormPage from './pages/EventFormPage'
import AdminReportPage from './pages/AdminReportPage'

function PrivateRoute({ children }) {
  const { token } = useAuth()
  return token ? children : <Navigate to="/login" replace />
}

function AdminRoute({ children }) {
  const { token, user } = useAuth()
  if (!token) return <Navigate to="/login" replace />
  if (user?.rol !== 'admin') return <Navigate to="/" replace />
  return children
}

export default function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/eventos/:id" element={<EventDetailPage />} />
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />
          <Route
            path="/mis-entradas"
            element={<PrivateRoute><MyTicketsPage /></PrivateRoute>}
          />
          <Route
            path="/admin"
            element={<AdminRoute><AdminEventsPage /></AdminRoute>}
          />
          <Route
            path="/admin/nuevo"
            element={<AdminRoute><EventFormPage /></AdminRoute>}
          />
          <Route
            path="/admin/editar/:id"
            element={<AdminRoute><EventFormPage /></AdminRoute>}
          />
          <Route
            path="/admin/reportes/:id"
            element={<AdminRoute><AdminReportPage /></AdminRoute>}
          />
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  )
}
