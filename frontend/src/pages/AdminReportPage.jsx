import { useEffect, useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { getEventReport } from '../api/events'

export default function AdminReportPage() {
  const { id } = useParams()
  const [report, setReport] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    getEventReport(id)
      .then((res) => setReport(res.data.data))
      .catch(() => setError('No se pudo cargar el reporte'))
      .finally(() => setLoading(false))
  }, [id])

  if (loading) return <p style={styles.center}>Cargando reporte...</p>
  if (error) return <p style={{ ...styles.center, color: '#dc2626' }}>{error}</p>

  const ocupacion = report.capacidad_total > 0
    ? Math.round((report.entradas_vendidas / report.capacidad_total) * 100)
    : 0

  return (
    <div style={styles.page}>
      <header style={styles.header}>
        <h1 style={styles.logo}>🎟 TicketApp — Panel Admin</h1>
        <Link to="/admin" style={styles.navLink}>← Volver al panel</Link>
      </header>

      <main style={styles.main}>
        <h2 style={styles.title}>Reporte: {report.titulo}</h2>

        {/* Métricas */}
        <div style={styles.cards}>
          <div style={styles.card}>
            <span style={styles.cardLabel}>Capacidad total</span>
            <span style={styles.cardValue}>{report.capacidad_total}</span>
          </div>
          <div style={{ ...styles.card, ...styles.cardSold }}>
            <span style={styles.cardLabel}>Entradas vendidas</span>
            <span style={styles.cardValue}>{report.entradas_vendidas}</span>
          </div>
          <div style={{ ...styles.card, ...styles.cardAvail }}>
            <span style={styles.cardLabel}>Entradas disponibles</span>
            <span style={styles.cardValue}>{report.entradas_disponibles}</span>
          </div>
        </div>

        {/* Barra de ocupación */}
        <div style={styles.progressSection}>
          <div style={styles.progressHeader}>
            <span style={styles.progressLabel}>Ocupación</span>
            <span style={styles.progressPct}>{ocupacion}%</span>
          </div>
          <div style={styles.progressTrack}>
            <div style={{ ...styles.progressFill, width: `${ocupacion}%`, background: ocupacion >= 90 ? '#dc2626' : ocupacion >= 60 ? '#f59e0b' : '#10b981' }} />
          </div>
          <p style={styles.progressSub}>
            {report.entradas_vendidas} de {report.capacidad_total} entradas vendidas
          </p>
        </div>

        {/* Tabla de compradores */}
        <h3 style={styles.subtitle}>Compradores</h3>
        {report.compradores.length === 0 ? (
          <p style={styles.empty}>Este evento no tiene compradores aún.</p>
        ) : (
          <div style={styles.tableWrapper}>
            <table style={styles.table}>
              <thead>
                <tr>
                  {['Nombre', 'Email', 'Estado', 'Fecha de compra'].map((h) => (
                    <th key={h} style={styles.th}>{h}</th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {report.compradores.map((c) => (
                  <tr key={`${c.user_id}-${c.fecha_compra}`} style={styles.tr}>
                    <td style={styles.td}>{c.nombre}</td>
                    <td style={styles.td}>{c.email}</td>
                    <td style={styles.td}>
                      <span style={{ ...styles.badge, ...badgeColor(c.estado) }}>
                        {c.estado.toUpperCase()}
                      </span>
                    </td>
                    <td style={styles.td}>
                      {new Date(c.fecha_compra).toLocaleDateString('es-AR', {
                        day: 'numeric', month: 'short', year: 'numeric',
                        hour: '2-digit', minute: '2-digit',
                      })}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </main>
    </div>
  )
}

function badgeColor(estado) {
  if (estado === 'activo') return { background: '#d1fae5', color: '#065f46' }
  if (estado === 'cancelado') return { background: '#fee2e2', color: '#991b1b' }
  return { background: '#e0e7ff', color: '#3730a3' }
}

const styles = {
  page: { minHeight: '100vh', background: '#f8f8f8', fontFamily: 'sans-serif' },
  header: { background: '#1e1b4b', color: '#fff', padding: '1rem 2rem', display: 'flex', justifyContent: 'space-between', alignItems: 'center' },
  logo: { margin: 0, fontSize: '1.3rem' },
  navLink: { color: '#a5b4fc', textDecoration: 'none', fontWeight: '500' },
  main: { maxWidth: '900px', margin: '0 auto', padding: '2rem 1rem' },
  title: { marginBottom: '1.5rem', fontSize: '1.5rem' },
  cards: { display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: '1rem', marginBottom: '1.5rem' },
  card: { background: '#fff', borderRadius: '8px', padding: '1.25rem', boxShadow: '0 2px 8px rgba(0,0,0,0.07)', display: 'flex', flexDirection: 'column', gap: '0.4rem' },
  cardSold: { borderTop: '3px solid #4f46e5' },
  cardAvail: { borderTop: '3px solid #10b981' },
  cardLabel: { fontSize: '0.8rem', color: '#6b7280', fontWeight: '600', textTransform: 'uppercase' },
  cardValue: { fontSize: '2rem', fontWeight: '700', color: '#111827' },
  progressSection: { background: '#fff', borderRadius: '8px', padding: '1.25rem', boxShadow: '0 2px 8px rgba(0,0,0,0.07)', marginBottom: '2rem' },
  progressHeader: { display: 'flex', justifyContent: 'space-between', marginBottom: '0.5rem' },
  progressLabel: { fontWeight: '600', color: '#374151' },
  progressPct: { fontWeight: '700', fontSize: '1.1rem', color: '#111827' },
  progressTrack: { height: '16px', background: '#e5e7eb', borderRadius: '8px', overflow: 'hidden' },
  progressFill: { height: '100%', borderRadius: '8px', transition: 'width 0.4s ease' },
  progressSub: { margin: '0.5rem 0 0', fontSize: '0.85rem', color: '#6b7280' },
  subtitle: { fontSize: '1.2rem', margin: '0 0 1rem' },
  tableWrapper: { overflowX: 'auto' },
  table: { width: '100%', borderCollapse: 'collapse', background: '#fff', borderRadius: '8px', overflow: 'hidden', boxShadow: '0 2px 8px rgba(0,0,0,0.07)' },
  th: { background: '#1e1b4b', color: '#fff', padding: '0.75rem 1rem', textAlign: 'left', fontSize: '0.85rem', fontWeight: '600' },
  tr: { borderBottom: '1px solid #f0f0f0' },
  td: { padding: '0.75rem 1rem', fontSize: '0.9rem' },
  badge: { fontSize: '0.7rem', fontWeight: '700', padding: '0.2rem 0.6rem', borderRadius: '12px' },
  empty: { color: '#6b7280', textAlign: 'center', padding: '2rem', background: '#fff', borderRadius: '8px' },
  center: { textAlign: 'center', padding: '3rem', color: '#6b7280' },
}
