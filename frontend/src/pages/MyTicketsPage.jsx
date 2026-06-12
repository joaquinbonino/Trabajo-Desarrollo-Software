import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { getMyTickets, cancelTicket, transferTicket } from '../api/tickets'
import { getMyWaitlist } from '../api/waitlist'

export default function MyTicketsPage() {
  const [tickets, setTickets] = useState([])
  const [waitlist, setWaitlist] = useState([])
  const [loading, setLoading] = useState(true)
  const [pageError, setPageError] = useState('')
  const [transferState, setTransferState] = useState({})

  useEffect(() => {
    fetchTickets()
  }, [])

  async function fetchTickets() {
    setLoading(true)
    setPageError('')
    try {
      const [ticketsRes, waitlistRes] = await Promise.all([
        getMyTickets(),
        getMyWaitlist(),
      ])
      setTickets(ticketsRes.data.data || [])
      setWaitlist(waitlistRes.data.data || [])
    } catch {
      setPageError('Error al cargar las entradas')
    } finally {
      setLoading(false)
    }
  }

  async function handleCancel(ticketId) {
    if (!confirm('¿Cancelar esta entrada? Esta acción no se puede deshacer.')) return
    try {
      await cancelTicket(ticketId)
      setTickets((prev) =>
        prev.map((t) => (t.id === ticketId ? { ...t, estado: 'cancelado' } : t))
      )
    } catch (err) {
      setTransferState((prev) => ({
        ...prev,
        [ticketId]: { ...prev[ticketId], cancelError: err.response?.data?.error || 'Error al cancelar' },
      }))
    }
  }

  function handleTransferEmailChange(ticketId, value) {
    setTransferState((prev) => ({
      ...prev,
      [ticketId]: { ...prev[ticketId], email: value, error: '', success: false },
    }))
  }

  function toggleTransferForm(ticketId) {
    setTransferState((prev) => ({
      ...prev,
      [ticketId]: { email: '', open: !prev[ticketId]?.open, error: '', success: false },
    }))
  }

  async function handleTransfer(ticketId) {
    const email = transferState[ticketId]?.email?.trim()
    if (!email) {
      setTransferState((prev) => ({
        ...prev,
        [ticketId]: { ...prev[ticketId], error: 'Ingresá el email del destinatario' },
      }))
      return
    }

    setTransferState((prev) => ({
      ...prev,
      [ticketId]: { ...prev[ticketId], loading: true, error: '' },
    }))

    try {
      await transferTicket(ticketId, { destino_email: email })
      setTickets((prev) => prev.filter((t) => t.id !== ticketId))
    } catch (err) {
      setTransferState((prev) => ({
        ...prev,
        [ticketId]: {
          ...prev[ticketId],
          loading: false,
          error: err.response?.data?.error || 'Error al transferir',
        },
      }))
    }
  }

  if (loading) return <p style={styles.center}>Cargando entradas...</p>

  return (
    <div style={styles.page}>
      <div style={styles.container}>
        <div style={styles.headerRow}>
          <h2 style={{ margin: 0 }}>Mis Entradas</h2>
          <Link to="/" style={styles.back}>← Volver al catálogo</Link>
        </div>

        {pageError && <p style={styles.pageError}>{pageError}</p>}

        {tickets.length === 0 && (
          <p style={styles.empty}>
            No tenés entradas aún. <Link to="/">¡Explorá los eventos!</Link>
          </p>
        )}

        <div style={styles.list}>
          {tickets.map((ticket) => {
            const ts = transferState[ticket.id] || {}
            return (
              <div key={ticket.id} style={styles.card}>
                <div style={styles.cardLeft}>
                  <span style={{ ...styles.badge, ...badgeColor(ticket.estado) }}>
                    {ticket.estado.toUpperCase()}
                  </span>
                  <h3 style={styles.eventTitle}>{ticket.event?.titulo || '—'}</h3>
                  <p style={styles.eventDate}>
                    📅 {ticket.event?.fecha_hora
                      ? new Date(ticket.event.fecha_hora).toLocaleString('es-AR', {
                          day: 'numeric', month: 'long', year: 'numeric',
                          hour: '2-digit', minute: '2-digit',
                        })
                      : '—'}
                  </p>
                  <p style={styles.meta}>
                    Comprada el {new Date(ticket.fecha_compra).toLocaleDateString('es-AR')}
                  </p>
                  {ts.cancelError && <p style={styles.inlineError}>{ts.cancelError}</p>}
                </div>

                {ticket.estado === 'activo' && (
                  <div style={styles.actions}>
                    <button style={styles.cancelBtn} onClick={() => handleCancel(ticket.id)}>
                      Cancelar
                    </button>
                    <button style={styles.transferBtn} onClick={() => toggleTransferForm(ticket.id)}>
                      {ts.open ? 'Cerrar' : 'Transferir'}
                    </button>

                    {ts.open && (
                      <div style={styles.transferForm}>
                        <input
                          style={styles.input}
                          type="email"
                          placeholder="Email del destinatario"
                          value={ts.email || ''}
                          onChange={(e) => handleTransferEmailChange(ticket.id, e.target.value)}
                        />
                        {ts.error && <p style={styles.inlineError}>{ts.error}</p>}
                        <button
                          style={{ ...styles.confirmBtn, ...(ts.loading ? styles.btnDisabled : {}) }}
                          onClick={() => handleTransfer(ticket.id)}
                          disabled={ts.loading}
                        >
                          {ts.loading ? 'Transfiriendo...' : 'Confirmar transferencia'}
                        </button>
                      </div>
                    )}
                  </div>
                )}
              </div>
            )
          })}
        </div>

        {waitlist.length > 0 && (
          <div style={styles.waitlistSection}>
            <h2 style={{ margin: '0 0 1rem' }}>Lista de espera</h2>
            <div style={styles.list}>
              {waitlist.map((w) => (
                <div key={`wl-${w.id}`} style={styles.card}>
                  <div style={styles.cardLeft}>
                    <span style={{ ...styles.badge, ...waitlistBadgeColor(w.estado) }}>
                      {w.estado === 'asignado' ? 'ASIGNADA' : 'EN ESPERA'}
                    </span>
                    <h3 style={styles.eventTitle}>{w.event?.titulo || '—'}</h3>
                    <p style={styles.eventDate}>
                      📅 {w.event?.fecha_hora
                        ? new Date(w.event.fecha_hora).toLocaleString('es-AR', {
                            day: 'numeric', month: 'long', year: 'numeric',
                            hour: '2-digit', minute: '2-digit',
                          })
                        : '—'}
                    </p>
                    {w.estado === 'asignado' ? (
                      <p style={styles.assignedNote}>
                        🎉 ¡Se liberó un lugar y te lo asignamos! Ya tenés tu entrada
                        {w.fecha_asignacion
                          ? ` (${new Date(w.fecha_asignacion).toLocaleDateString('es-AR')})`
                          : ''}.
                      </p>
                    ) : (
                      <p style={styles.meta}>
                        Te avisamos acá si se libera un lugar.
                      </p>
                    )}
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  )
}

function waitlistBadgeColor(estado) {
  if (estado === 'asignado') return { background: '#d1fae5', color: '#065f46' }
  return { background: '#fef3c7', color: '#92400e' }
}

function badgeColor(estado) {
  if (estado === 'activo') return { background: '#d1fae5', color: '#065f46' }
  if (estado === 'cancelado') return { background: '#fee2e2', color: '#991b1b' }
  return { background: '#e0e7ff', color: '#3730a3' }
}

const styles = {
  page: { minHeight: '100vh', background: '#f8f8f8', fontFamily: 'sans-serif', padding: '2rem 1rem' },
  container: { maxWidth: '800px', margin: '0 auto' },
  headerRow: { display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1.5rem' },
  back: { color: '#4f46e5', textDecoration: 'none', fontSize: '0.9rem' },
  list: { display: 'flex', flexDirection: 'column', gap: '1rem' },
  card: { background: '#fff', borderRadius: '8px', boxShadow: '0 2px 8px rgba(0,0,0,0.07)', padding: '1.25rem', display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', gap: '1rem', flexWrap: 'wrap' },
  cardLeft: { flex: 1 },
  badge: { fontSize: '0.7rem', fontWeight: '700', padding: '0.2rem 0.6rem', borderRadius: '12px' },
  eventTitle: { margin: '0.4rem 0 0.2rem', fontSize: '1.1rem' },
  eventDate: { color: '#6b7280', fontSize: '0.85rem', margin: '0.1rem 0' },
  meta: { color: '#9ca3af', fontSize: '0.8rem', margin: '0.1rem 0' },
  actions: { display: 'flex', flexDirection: 'column', gap: '0.5rem', minWidth: '160px' },
  cancelBtn: { padding: '0.5rem', background: '#fee2e2', color: '#dc2626', border: 'none', borderRadius: '4px', cursor: 'pointer', fontWeight: '600' },
  transferBtn: { padding: '0.5rem', background: '#e0e7ff', color: '#4338ca', border: 'none', borderRadius: '4px', cursor: 'pointer', fontWeight: '600' },
  transferForm: { display: 'flex', flexDirection: 'column', gap: '0.4rem' },
  input: { padding: '0.5rem', borderRadius: '4px', border: '1px solid #ccc', fontSize: '0.9rem' },
  confirmBtn: { padding: '0.5rem', background: '#4f46e5', color: '#fff', border: 'none', borderRadius: '4px', cursor: 'pointer', fontWeight: '600' },
  btnDisabled: { background: '#9ca3af', cursor: 'not-allowed' },
  inlineError: { color: '#dc2626', fontSize: '0.82rem', margin: '0', background: '#fef2f2', padding: '0.3rem 0.5rem', borderRadius: '4px' },
  pageError: { color: '#dc2626', textAlign: 'center' },
  empty: { textAlign: 'center', color: '#6b7280', padding: '2rem' },
  center: { textAlign: 'center', padding: '3rem', color: '#6b7280' },
  waitlistSection: { marginTop: '2.5rem' },
  assignedNote: { color: '#065f46', fontSize: '0.85rem', margin: '0.4rem 0 0', background: '#ecfdf5', padding: '0.4rem 0.6rem', borderRadius: '4px' },
}
