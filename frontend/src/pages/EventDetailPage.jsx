import { useEffect, useState } from 'react'
import { useParams, useNavigate, Link } from 'react-router-dom'
import { getEvent } from '../api/events'
import { buyTicket } from '../api/tickets'
import { joinWaitlist } from '../api/waitlist'
import { useAuth } from '../context/AuthContext'

export default function EventDetailPage() {
  const { id } = useParams()
  const { user } = useAuth()
  const navigate = useNavigate()

  const [event, setEvent] = useState(null)
  const [loading, setLoading] = useState(true)
  const [buying, setBuying] = useState(false)
  const [success, setSuccess] = useState(false)
  const [error, setError] = useState('')
  const [joiningWaitlist, setJoiningWaitlist] = useState(false)
  const [waitlistJoined, setWaitlistJoined] = useState(false)

  useEffect(() => {
    getEvent(id)
      .then((res) => setEvent(res.data.data))
      .catch(() => setError('Evento no encontrado'))
      .finally(() => setLoading(false))
  }, [id])

  async function handleBuy() {
    if (!user) {
      navigate('/login')
      return
    }
    setBuying(true)
    setError('')
    try {
      await buyTicket({ event_id: Number(id) })
      setSuccess(true)
      setEvent((prev) => ({ ...prev, entradas_vendidas: prev.entradas_vendidas + 1 }))
    } catch (err) {
      setError(err.response?.data?.error || 'Error al comprar la entrada')
    } finally {
      setBuying(false)
    }
  }

  async function handleJoinWaitlist() {
    if (!user) {
      navigate('/login')
      return
    }
    setJoiningWaitlist(true)
    setError('')
    try {
      await joinWaitlist(Number(id))
      setWaitlistJoined(true)
    } catch (err) {
      setError(err.response?.data?.error || 'Error al anotarse en la lista de espera')
    } finally {
      setJoiningWaitlist(false)
    }
  }

  if (loading) return <p style={styles.center}>Cargando...</p>
  if (!event) return <p style={styles.center}>Evento no encontrado.</p>

  const disponibles = event.capacidad_total - event.entradas_vendidas
  const sinCupo = disponibles <= 0
  const agotado = sinCupo || event.cancelado

  return (
    <div style={styles.page}>
      <div style={styles.container}>
        <Link to="/" style={styles.back}>← Volver al catálogo</Link>

        {event.foto && (
          <img src={event.foto} alt={event.titulo} style={styles.hero} />
        )}

        <div style={styles.content}>
          <span style={styles.categoria}>{event.categoria}</span>
          {event.cancelado && <span style={styles.canceladoBadge}>CANCELADO</span>}
          <h1 style={styles.title}>{event.titulo}</h1>

          <p style={styles.meta}>
            📅 {new Date(event.fecha_hora).toLocaleString('es-AR', {
              weekday: 'long', day: 'numeric', month: 'long', year: 'numeric',
              hour: '2-digit', minute: '2-digit',
            })}
          </p>
          {event.duracion > 0 && (
            <p style={styles.meta}>⏱ Duración: {event.duracion} min</p>
          )}
          <p style={styles.meta}>
            🎟 {disponibles > 0 ? `${disponibles} entradas disponibles` : 'Sin entradas disponibles'}
          </p>

          {event.descripcion && (
            <p style={styles.description}>{event.descripcion}</p>
          )}

          {success ? (
            <div style={styles.successBox}>
              <h3>🎉 ¡Compra exitosa!</h3>
              <p>Tu entrada fue reservada. Podés verla en{' '}
                <Link to="/mis-entradas">Mis Entradas</Link>.
              </p>
            </div>
          ) : waitlistJoined ? (
            <div style={styles.successBox}>
              <h3>📝 ¡Estás en la lista de espera!</h3>
              <p>Si alguien cancela su entrada, se te asignará automáticamente.
                Vas a poder verla en <Link to="/mis-entradas">Mis Entradas</Link>.
              </p>
            </div>
          ) : sinCupo && !event.cancelado ? (
            <>
              {error && <p style={styles.error}>{error}</p>}
              <p style={styles.meta}>
                Este evento está agotado. Anotate en la lista de espera y te avisamos
                si se libera un lugar.
              </p>
              <button
                style={{ ...styles.waitlistBtn, ...(joiningWaitlist ? styles.buyBtnDisabled : {}) }}
                onClick={handleJoinWaitlist}
                disabled={joiningWaitlist}
              >
                {joiningWaitlist ? 'Procesando...' : 'Anotarme en lista de espera'}
              </button>
              {!user && (
                <p style={styles.loginNote}>
                  Necesitás <Link to="/login">iniciar sesión</Link> para anotarte.
                </p>
              )}
            </>
          ) : (
            <>
              {error && <p style={styles.error}>{error}</p>}
              <button
                style={{ ...styles.buyBtn, ...(agotado ? styles.buyBtnDisabled : {}) }}
                onClick={handleBuy}
                disabled={agotado || buying}
              >
                {buying ? 'Procesando...' : agotado ? 'Sin entradas' : 'Comprar entrada'}
              </button>
              {!user && !agotado && (
                <p style={styles.loginNote}>
                  Necesitás <Link to="/login">iniciar sesión</Link> para comprar.
                </p>
              )}
            </>
          )}
        </div>
      </div>
    </div>
  )
}

const styles = {
  page: { minHeight: '100vh', background: '#f8f8f8', fontFamily: 'sans-serif', padding: '2rem 1rem' },
  container: { maxWidth: '720px', margin: '0 auto', background: '#fff', borderRadius: '8px', boxShadow: '0 2px 12px rgba(0,0,0,0.08)', overflow: 'hidden' },
  back: { display: 'block', padding: '1rem 1.5rem', color: '#4f46e5', textDecoration: 'none', fontSize: '0.9rem' },
  hero: { width: '100%', maxHeight: '360px', objectFit: 'cover' },
  content: { padding: '1.5rem' },
  categoria: { fontSize: '0.75rem', background: '#e0e7ff', color: '#3730a3', padding: '0.2rem 0.6rem', borderRadius: '12px', fontWeight: '600' },
  canceladoBadge: { marginLeft: '0.5rem', fontSize: '0.75rem', background: '#fee2e2', color: '#dc2626', padding: '0.2rem 0.6rem', borderRadius: '12px', fontWeight: '700' },
  title: { margin: '0.75rem 0 0.5rem', fontSize: '1.8rem' },
  meta: { color: '#6b7280', fontSize: '0.95rem', margin: '0.25rem 0' },
  description: { margin: '1rem 0', lineHeight: '1.6', color: '#374151' },
  buyBtn: { marginTop: '1.5rem', display: 'block', width: '100%', padding: '0.9rem', background: '#4f46e5', color: '#fff', border: 'none', borderRadius: '6px', fontSize: '1.1rem', fontWeight: '600', cursor: 'pointer' },
  waitlistBtn: { marginTop: '0.75rem', display: 'block', width: '100%', padding: '0.9rem', background: '#0891b2', color: '#fff', border: 'none', borderRadius: '6px', fontSize: '1.1rem', fontWeight: '600', cursor: 'pointer' },
  buyBtnDisabled: { background: '#9ca3af', cursor: 'not-allowed' },
  successBox: { marginTop: '1.5rem', background: '#ecfdf5', border: '1px solid #6ee7b7', borderRadius: '6px', padding: '1rem' },
  error: { color: '#dc2626', fontSize: '0.9rem', marginTop: '0.5rem' },
  loginNote: { fontSize: '0.85rem', color: '#6b7280', marginTop: '0.5rem' },
  center: { textAlign: 'center', padding: '3rem', color: '#6b7280' },
}
