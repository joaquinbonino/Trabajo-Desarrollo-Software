import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { register as registerApi } from '../api/auth'
import { useAuth } from '../context/AuthContext'

export default function RegisterPage() {
  const [form, setForm] = useState({ nombre: '', email: '', password: '' })
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const { login } = useAuth()
  const navigate = useNavigate()

  function handleChange(e) {
    setForm({ ...form, [e.target.name]: e.target.value })
  }

  async function handleSubmit(e) {
    e.preventDefault()
    setError('')
    if (form.password.length < 6) {
      setError('La contraseña debe tener al menos 6 caracteres')
      return
    }
    setLoading(true)
    try {
      const res = await registerApi(form)
      login(res.data.data.token, res.data.data.user)
      navigate('/')
    } catch (err) {
      setError(err.response?.data?.error || 'Error al registrarse')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div style={styles.container}>
      <div style={styles.card}>
        <h2 style={styles.title}>Crear cuenta</h2>
        <form onSubmit={handleSubmit} style={styles.form}>
          <label style={styles.label}>Nombre</label>
          <input
            style={styles.input}
            type="text"
            name="nombre"
            value={form.nombre}
            onChange={handleChange}
            required
          />
          <label style={styles.label}>Email</label>
          <input
            style={styles.input}
            type="email"
            name="email"
            value={form.email}
            onChange={handleChange}
            required
          />
          <label style={styles.label}>Contraseña</label>
          <input
            style={styles.input}
            type="password"
            name="password"
            value={form.password}
            onChange={handleChange}
            required
          />
          {error && <p style={styles.error}>{error}</p>}
          <button style={styles.button} type="submit" disabled={loading}>
            {loading ? 'Registrando...' : 'Registrarse'}
          </button>
        </form>
        <p style={styles.footer}>
          ¿Ya tenés cuenta?{' '}
          <Link to="/login">Iniciá sesión</Link>
        </p>
      </div>
    </div>
  )
}

const styles = {
  container: { display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '80vh' },
  card: { background: '#fff', padding: '2rem', borderRadius: '8px', boxShadow: '0 2px 12px rgba(0,0,0,0.1)', width: '100%', maxWidth: '400px' },
  title: { marginBottom: '1.5rem', textAlign: 'center' },
  form: { display: 'flex', flexDirection: 'column', gap: '0.5rem' },
  label: { fontWeight: '600', fontSize: '0.9rem' },
  input: { padding: '0.6rem', borderRadius: '4px', border: '1px solid #ccc', fontSize: '1rem' },
  button: { marginTop: '1rem', padding: '0.75rem', background: '#4f46e5', color: '#fff', border: 'none', borderRadius: '4px', fontSize: '1rem', cursor: 'pointer' },
  error: { color: '#dc2626', fontSize: '0.9rem', margin: '0' },
  footer: { textAlign: 'center', marginTop: '1rem', fontSize: '0.9rem' },
}
