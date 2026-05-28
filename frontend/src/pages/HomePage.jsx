import { useEffect, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { getEvents } from '../api/events'
import { useAuth } from '../context/AuthContext'

const CATEGORIAS = ['', 'Música', 'Teatro', 'Deportes', 'Cine', 'Arte', 'Otro']

export default function HomePage() {
  const [events, setEvents] = useState([])
  const [categoria, setCategoria] = useState('')
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const { user, logout } = useAuth()
  const navigate = useNavigate()

  useEffect(() => {
    fetchEvents()
  }, [categoria])

  async function fetchEvents() {
    setLoading(true)
    setError('')
    try {
      const res = await getEvents(categoria)
      setEvents(res.data.data || [])
    } catch {
      setError('Error al cargar los eventos')
    } finally {
      setLoading(false)
    }
  }

  function handleLogout() {
    logout()
    navigate('/login')
  }

  return (
    <div style={styles.page}>
      <header style={styles.header}>
        <h1 style={styles.logo}>🎟 TicketApp</h1>
        <nav style={styles.nav}>
          {user ? (
            <>
              <Link to="/mis-entradas" style={styles.navLink}>Mis Entradas</Link>
              <span style={styles.navUser}>Hola, {user.nombre}</span>
              <button onClick={handleLogout} style={styles.navBtn}>Cerrar sesión</button>
            </>
          ) : (
            <>
              <Link to="/login" style={styles.navLink}>Iniciar sesión</Link>
              <Link to="/register" style={styles.navLink}>Registrarse</Link>
            </>
          )}
        </nav>
      </header>

      <main style={styles.main}>
        <div style={styles.filterBar}>
          <h2 style={{ margin: 0 }}>Catálogo de Eventos</h2>
          <select
            style={styles.select}
            value={categoria}
            onChange={(e) => setCategoria(e.target.value)}
          >
            {CATEGORIAS.map((c) => (
              <option key={c} value={c}>{c === '' ? 'Todas las categorías' : c}</option>
            ))}
          </select>
        </div>

        {loading && <p style={styles.info}>Cargando eventos...</p>}
        {error && <p style={styles.error}>{error}</p>}
        {!loading && events.length === 0 && (
          <p style={styles.info}>No hay eventos disponibles.</p>
        )}

        <div style={styles.grid}>
          {events.map((event) => (
            <Link key={event.id} to={`/eventos/${event.id}`} style={styles.cardLink}>
              <div style={styles.card}>
                {event.foto && (
                  <img src={event.foto} alt={event.titulo} style={styles.cardImg} />
                )}
                <div style={styles.cardBody}>
                  <span style={styles.categoria}>{event.categoria}</span>
                  <h3 style={styles.cardTitle}>{event.titulo}</h3>
                  <p style={styles.cardDate}>
                    {new Date(event.fecha_hora).toLocaleDateString('es-AR', {
                      weekday: 'short', day: 'numeric', month: 'long', year: 'numeric',
                    })}
                  </p>
                  <p style={styles.cardCapacity}>
                    {event.entradas_vendidas}/{event.capacidad_total} entradas vendidas
                  </p>
                </div>
              </div>
            </Link>
          ))}
        </div>
      </main>
    </div>
  )
}

const styles = {
  page: { minHeight: '100vh', background: '#f8f8f8', fontFamily: 'sans-serif' },
  header: { background: '#4f46e5', color: '#fff', padding: '1rem 2rem', display: 'flex', justifyContent: 'space-between', alignItems: 'center' },
  logo: { margin: 0, fontSize: '1.5rem' },
  nav: { display: 'flex', gap: '1rem', alignItems: 'center' },
  navLink: { color: '#fff', textDecoration: 'none', fontWeight: '500' },
  navUser: { fontSize: '0.9rem' },
  navBtn: { background: 'transparent', border: '1px solid #fff', color: '#fff', padding: '0.3rem 0.8rem', borderRadius: '4px', cursor: 'pointer' },
  main: { maxWidth: '1100px', margin: '0 auto', padding: '2rem 1rem' },
  filterBar: { display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1.5rem', flexWrap: 'wrap', gap: '1rem' },
  select: { padding: '0.5rem 1rem', borderRadius: '4px', border: '1px solid #ccc', fontSize: '1rem' },
  grid: { display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(280px, 1fr))', gap: '1.5rem' },
  cardLink: { textDecoration: 'none', color: 'inherit' },
  card: { background: '#fff', borderRadius: '8px', boxShadow: '0 2px 8px rgba(0,0,0,0.08)', overflow: 'hidden', transition: 'transform 0.15s', cursor: 'pointer' },
  cardImg: { width: '100%', height: '160px', objectFit: 'cover' },
  cardBody: { padding: '1rem' },
  categoria: { fontSize: '0.75rem', background: '#e0e7ff', color: '#3730a3', padding: '0.2rem 0.6rem', borderRadius: '12px', fontWeight: '600' },
  cardTitle: { margin: '0.5rem 0 0.25rem' },
  cardDate: { color: '#6b7280', fontSize: '0.85rem', margin: '0 0 0.25rem' },
  cardCapacity: { color: '#9ca3af', fontSize: '0.8rem', margin: 0 },
  info: { textAlign: 'center', color: '#6b7280', padding: '2rem' },
  error: { textAlign: 'center', color: '#dc2626', padding: '1rem' },
}
