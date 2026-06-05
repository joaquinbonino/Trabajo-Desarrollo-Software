import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { getAllEventsAdmin, cancelEvent, updateEvent } from '../api/events'
import { useAuth } from '../context/AuthContext'

export default function AdminEventsPage() {
  const { logout } = useAuth()
  const [events, setEvents] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [editingId, setEditingId] = useState(null)
  const [editForm, setEditForm] = useState({})
  const [actionError, setActionError] = useState({})

  useEffect(() => {
    fetchEvents()
  }, [])

  async function fetchEvents() {
    setLoading(true)
    setError('')
    try {
      const res = await getAllEventsAdmin()
      setEvents(res.data.data || [])
    } catch {
      setError('Error al cargar los eventos')
    } finally {
      setLoading(false)
    }
  }

  function startEdit(event) {
    setEditingId(event.id)
    setEditForm({
      titulo: event.titulo,
      descripcion: event.descripcion || '',
      categoria: event.categoria || '',
      fecha_hora: event.fecha_hora ? event.fecha_hora.slice(0, 16) : '',
      duracion: event.duracion || 0,
      capacidad_total: event.capacidad_total,
      foto: event.foto || '',
    })
    setActionError((prev) => ({ ...prev, [event.id]: '' }))
  }

  function cancelEdit() {
    setEditingId(null)
    setEditForm({})
  }

  async function handleSaveEdit(id) {
    try {
      const payload = {
        ...editForm,
        duracion: Number(editForm.duracion),
        capacidad_total: Number(editForm.capacidad_total),
        fecha_hora: editForm.fecha_hora ? new Date(editForm.fecha_hora).toISOString() : undefined,
      }
      const res = await updateEvent(id, payload)
      setEvents((prev) => prev.map((e) => (e.id === id ? res.data.data : e)))
      setEditingId(null)
    } catch (err) {
      setActionError((prev) => ({
        ...prev,
        [id]: err.response?.data?.error || 'Error al guardar',
      }))
    }
  }

  async function handleCancel(id) {
    if (!confirm('¿Cancelar este evento? Esta acción no se puede deshacer.')) return
    try {
      await cancelEvent(id)
      setEvents((prev) =>
        prev.map((e) => (e.id === id ? { ...e, cancelado: true } : e))
      )
    } catch (err) {
      setActionError((prev) => ({
        ...prev,
        [id]: err.response?.data?.error || 'Error al cancelar',
      }))
    }
  }

  if (loading) return <p style={styles.center}>Cargando eventos...</p>

  return (
    <div style={styles.page}>
      <header style={styles.header}>
        <h1 style={styles.logo}>🎟 TicketApp — Panel Admin</h1>
        <nav style={styles.nav}>
          <Link to="/" style={styles.navLink}>Catálogo</Link>
          <button onClick={logout} style={styles.navBtn}>Cerrar sesión</button>
        </nav>
      </header>

      <main style={styles.main}>
        <div style={styles.titleRow}>
          <h2 style={{ margin: 0 }}>Gestión de Eventos</h2>
          <Link to="/admin/nuevo" style={styles.newBtn}>+ Nuevo evento</Link>
        </div>

        {error && <p style={styles.error}>{error}</p>}

        {events.length === 0 ? (
          <p style={styles.empty}>No hay eventos registrados.</p>
        ) : (
          <div style={styles.tableWrapper}>
            <table style={styles.table}>
              <thead>
                <tr>
                  {['ID', 'Título', 'Categoría', 'Fecha', 'Capacidad', 'Vendidas', 'Estado', 'Acciones'].map((h) => (
                    <th key={h} style={styles.th}>{h}</th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {events.map((event) => (
                  <>
                    <tr key={event.id} style={styles.tr}>
                      <td style={styles.td}>{event.id}</td>
                      <td style={styles.td}>{event.titulo}</td>
                      <td style={styles.td}>{event.categoria || '—'}</td>
                      <td style={styles.td}>
                        {new Date(event.fecha_hora).toLocaleDateString('es-AR', {
                          day: 'numeric', month: 'short', year: 'numeric',
                        })}
                      </td>
                      <td style={{ ...styles.td, textAlign: 'center' }}>{event.capacidad_total}</td>
                      <td style={{ ...styles.td, textAlign: 'center' }}>{event.entradas_vendidas}</td>
                      <td style={styles.td}>
                        <span style={{ ...styles.badge, ...(event.cancelado ? styles.badgeCancelado : styles.badgeActivo) }}>
                          {event.cancelado ? 'CANCELADO' : 'ACTIVO'}
                        </span>
                      </td>
                      <td style={styles.td}>
                        <div style={styles.actions}>
                          {!event.cancelado && (
                            <>
                              <button
                                style={styles.editBtn}
                                onClick={() => editingId === event.id ? cancelEdit() : startEdit(event)}
                              >
                                {editingId === event.id ? 'Cerrar' : 'Editar'}
                              </button>
                              <button
                                style={styles.cancelBtn}
                                onClick={() => handleCancel(event.id)}
                              >
                                Cancelar
                              </button>
                            </>
                          )}
                          <Link to={`/admin/reportes/${event.id}`} style={styles.reportBtn}>
                            Reporte
                          </Link>
                        </div>
                        {actionError[event.id] && (
                          <p style={styles.inlineError}>{actionError[event.id]}</p>
                        )}
                      </td>
                    </tr>

                    {editingId === event.id && (
                      <tr key={`edit-${event.id}`}>
                        <td colSpan={8} style={styles.editRow}>
                          <div style={styles.editForm}>
                            <h4 style={{ margin: '0 0 1rem' }}>Editar evento #{event.id}</h4>
                            <div style={styles.formGrid}>
                              <label style={styles.label}>
                                Título
                                <input style={styles.input} value={editForm.titulo}
                                  onChange={(e) => setEditForm({ ...editForm, titulo: e.target.value })} />
                              </label>
                              <label style={styles.label}>
                                Categoría
                                <input style={styles.input} value={editForm.categoria}
                                  onChange={(e) => setEditForm({ ...editForm, categoria: e.target.value })} />
                              </label>
                              <label style={styles.label}>
                                Fecha y hora
                                <input type="datetime-local" style={styles.input} value={editForm.fecha_hora}
                                  onChange={(e) => setEditForm({ ...editForm, fecha_hora: e.target.value })} />
                              </label>
                              <label style={styles.label}>
                                Duración (min)
                                <input type="number" style={styles.input} value={editForm.duracion}
                                  onChange={(e) => setEditForm({ ...editForm, duracion: e.target.value })} />
                              </label>
                              <label style={styles.label}>
                                Capacidad total
                                <input type="number" style={styles.input} value={editForm.capacidad_total}
                                  onChange={(e) => setEditForm({ ...editForm, capacidad_total: e.target.value })} />
                              </label>
                              <label style={styles.label}>
                                URL Foto
                                <input style={styles.input} value={editForm.foto}
                                  onChange={(e) => setEditForm({ ...editForm, foto: e.target.value })} />
                              </label>
                              <label style={{ ...styles.label, gridColumn: '1 / -1' }}>
                                Descripción
                                <textarea style={{ ...styles.input, height: '80px', resize: 'vertical' }}
                                  value={editForm.descripcion}
                                  onChange={(e) => setEditForm({ ...editForm, descripcion: e.target.value })} />
                              </label>
                            </div>
                            <div style={styles.editActions}>
                              <button style={styles.saveBtn} onClick={() => handleSaveEdit(event.id)}>
                                Guardar cambios
                              </button>
                              <button style={styles.discardBtn} onClick={cancelEdit}>
                                Descartar
                              </button>
                            </div>
                          </div>
                        </td>
                      </tr>
                    )}
                  </>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </main>
    </div>
  )
}

const styles = {
  page: { minHeight: '100vh', background: '#f8f8f8', fontFamily: 'sans-serif' },
  header: { background: '#1e1b4b', color: '#fff', padding: '1rem 2rem', display: 'flex', justifyContent: 'space-between', alignItems: 'center' },
  logo: { margin: 0, fontSize: '1.3rem' },
  nav: { display: 'flex', gap: '1rem', alignItems: 'center' },
  navLink: { color: '#a5b4fc', textDecoration: 'none', fontWeight: '500' },
  navBtn: { background: 'transparent', border: '1px solid #a5b4fc', color: '#a5b4fc', padding: '0.3rem 0.8rem', borderRadius: '4px', cursor: 'pointer' },
  main: { maxWidth: '1200px', margin: '0 auto', padding: '2rem 1rem' },
  titleRow: { display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1.5rem' },
  newBtn: { background: '#4f46e5', color: '#fff', padding: '0.5rem 1.2rem', borderRadius: '6px', textDecoration: 'none', fontWeight: '600', fontSize: '0.9rem' },
  tableWrapper: { overflowX: 'auto' },
  table: { width: '100%', borderCollapse: 'collapse', background: '#fff', borderRadius: '8px', overflow: 'hidden', boxShadow: '0 2px 8px rgba(0,0,0,0.07)' },
  th: { background: '#1e1b4b', color: '#fff', padding: '0.75rem 1rem', textAlign: 'left', fontSize: '0.85rem', fontWeight: '600' },
  tr: { borderBottom: '1px solid #f0f0f0' },
  td: { padding: '0.75rem 1rem', fontSize: '0.9rem', verticalAlign: 'middle' },
  badge: { fontSize: '0.7rem', fontWeight: '700', padding: '0.2rem 0.6rem', borderRadius: '12px' },
  badgeActivo: { background: '#d1fae5', color: '#065f46' },
  badgeCancelado: { background: '#fee2e2', color: '#991b1b' },
  actions: { display: 'flex', gap: '0.4rem', flexWrap: 'wrap' },
  editBtn: { padding: '0.3rem 0.7rem', background: '#e0e7ff', color: '#3730a3', border: 'none', borderRadius: '4px', cursor: 'pointer', fontSize: '0.8rem', fontWeight: '600' },
  cancelBtn: { padding: '0.3rem 0.7rem', background: '#fee2e2', color: '#dc2626', border: 'none', borderRadius: '4px', cursor: 'pointer', fontSize: '0.8rem', fontWeight: '600' },
  reportBtn: { padding: '0.3rem 0.7rem', background: '#fef3c7', color: '#92400e', borderRadius: '4px', textDecoration: 'none', fontSize: '0.8rem', fontWeight: '600' },
  inlineError: { color: '#dc2626', fontSize: '0.78rem', margin: '0.3rem 0 0' },
  editRow: { background: '#f8f7ff', padding: 0 },
  editForm: { padding: '1.5rem' },
  formGrid: { display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(220px, 1fr))', gap: '1rem', marginBottom: '1rem' },
  label: { display: 'flex', flexDirection: 'column', gap: '0.3rem', fontSize: '0.85rem', fontWeight: '600', color: '#374151' },
  input: { padding: '0.5rem', border: '1px solid #d1d5db', borderRadius: '4px', fontSize: '0.9rem', fontFamily: 'sans-serif' },
  editActions: { display: 'flex', gap: '0.75rem' },
  saveBtn: { padding: '0.5rem 1.2rem', background: '#4f46e5', color: '#fff', border: 'none', borderRadius: '4px', cursor: 'pointer', fontWeight: '600' },
  discardBtn: { padding: '0.5rem 1.2rem', background: '#e5e7eb', color: '#374151', border: 'none', borderRadius: '4px', cursor: 'pointer', fontWeight: '600' },
  error: { color: '#dc2626', textAlign: 'center' },
  empty: { textAlign: 'center', color: '#6b7280', padding: '2rem' },
  center: { textAlign: 'center', padding: '3rem', color: '#6b7280' },
}
