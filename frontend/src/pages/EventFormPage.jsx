import { useEffect, useState } from 'react'
import { useNavigate, useParams, Link } from 'react-router-dom'
import { getEvent, createEvent, updateEvent } from '../api/events'

const CATEGORIAS = ['Música', 'Teatro', 'Deportes', 'Cine', 'Arte', 'Otro']

const emptyForm = {
  titulo: '',
  descripcion: '',
  categoria: '',
  fecha_hora: '',
  duracion: '',
  capacidad_total: '',
  foto: '',
}

export default function EventFormPage() {
  const { id } = useParams()
  const isEdit = Boolean(id)
  const navigate = useNavigate()

  const [form, setForm] = useState(emptyForm)
  const [errors, setErrors] = useState({})
  const [apiError, setApiError] = useState('')
  const [loading, setLoading] = useState(isEdit)
  const [submitting, setSubmitting] = useState(false)

  useEffect(() => {
    if (!isEdit) return
    getEvent(id)
      .then((res) => {
        const e = res.data.data
        setForm({
          titulo: e.titulo || '',
          descripcion: e.descripcion || '',
          categoria: e.categoria || '',
          fecha_hora: e.fecha_hora ? e.fecha_hora.slice(0, 16) : '',
          duracion: e.duracion ?? '',
          capacidad_total: e.capacidad_total ?? '',
          foto: e.foto || '',
        })
      })
      .catch(() => setApiError('No se pudo cargar el evento'))
      .finally(() => setLoading(false))
  }, [id, isEdit])

  function validate() {
    const errs = {}
    if (!form.titulo.trim()) errs.titulo = 'El título es requerido'
    if (!form.fecha_hora) errs.fecha_hora = 'La fecha y hora son requeridas'
    if (!form.capacidad_total || Number(form.capacidad_total) < 1)
      errs.capacidad_total = 'La capacidad debe ser al menos 1'
    if (form.duracion !== '' && Number(form.duracion) < 0)
      errs.duracion = 'La duración no puede ser negativa'
    return errs
  }

  async function handleSubmit(e) {
    e.preventDefault()
    setApiError('')
    const errs = validate()
    if (Object.keys(errs).length > 0) {
      setErrors(errs)
      return
    }
    setErrors({})
    setSubmitting(true)

    const payload = {
      titulo: form.titulo.trim(),
      descripcion: form.descripcion.trim(),
      categoria: form.categoria,
      fecha_hora: new Date(form.fecha_hora).toISOString(),
      duracion: form.duracion !== '' ? Number(form.duracion) : 0,
      capacidad_total: Number(form.capacidad_total),
      foto: form.foto.trim(),
    }

    try {
      if (isEdit) {
        await updateEvent(id, payload)
      } else {
        await createEvent(payload)
      }
      navigate('/admin')
    } catch (err) {
      setApiError(err.response?.data?.error || 'Error al guardar el evento')
    } finally {
      setSubmitting(false)
    }
  }

  function handleChange(field, value) {
    setForm((prev) => ({ ...prev, [field]: value }))
    if (errors[field]) setErrors((prev) => ({ ...prev, [field]: '' }))
  }

  if (loading) return <p style={styles.center}>Cargando evento...</p>

  return (
    <div style={styles.page}>
      <header style={styles.header}>
        <h1 style={styles.logo}>🎟 TicketApp — Panel Admin</h1>
        <Link to="/admin" style={styles.navLink}>← Volver al panel</Link>
      </header>

      <main style={styles.main}>
        <h2 style={styles.title}>
          {isEdit ? `Editar evento #${id}` : 'Nuevo evento'}
        </h2>

        {apiError && <p style={styles.apiError}>{apiError}</p>}

        <form onSubmit={handleSubmit} style={styles.form} noValidate>
          <div style={styles.grid}>

            <div style={styles.field}>
              <label style={styles.label}>Título *</label>
              <input
                style={{ ...styles.input, ...(errors.titulo ? styles.inputError : {}) }}
                type="text"
                value={form.titulo}
                onChange={(e) => handleChange('titulo', e.target.value)}
                placeholder="Nombre del evento"
              />
              {errors.titulo && <span style={styles.fieldError}>{errors.titulo}</span>}
            </div>

            <div style={styles.field}>
              <label style={styles.label}>Categoría</label>
              <select
                style={styles.input}
                value={form.categoria}
                onChange={(e) => handleChange('categoria', e.target.value)}
              >
                <option value="">Sin categoría</option>
                {CATEGORIAS.map((c) => (
                  <option key={c} value={c}>{c}</option>
                ))}
              </select>
            </div>

            <div style={styles.field}>
              <label style={styles.label}>Fecha y hora *</label>
              <input
                style={{ ...styles.input, ...(errors.fecha_hora ? styles.inputError : {}) }}
                type="datetime-local"
                value={form.fecha_hora}
                onChange={(e) => handleChange('fecha_hora', e.target.value)}
              />
              {errors.fecha_hora && <span style={styles.fieldError}>{errors.fecha_hora}</span>}
            </div>

            <div style={styles.field}>
              <label style={styles.label}>Duración (minutos)</label>
              <input
                style={{ ...styles.input, ...(errors.duracion ? styles.inputError : {}) }}
                type="number"
                min="0"
                value={form.duracion}
                onChange={(e) => handleChange('duracion', e.target.value)}
                placeholder="ej: 120"
              />
              {errors.duracion && <span style={styles.fieldError}>{errors.duracion}</span>}
            </div>

            <div style={styles.field}>
              <label style={styles.label}>Capacidad total *</label>
              <input
                style={{ ...styles.input, ...(errors.capacidad_total ? styles.inputError : {}) }}
                type="number"
                min="1"
                value={form.capacidad_total}
                onChange={(e) => handleChange('capacidad_total', e.target.value)}
                placeholder="ej: 200"
              />
              {errors.capacidad_total && <span style={styles.fieldError}>{errors.capacidad_total}</span>}
            </div>

            <div style={styles.field}>
              <label style={styles.label}>URL de foto</label>
              <input
                style={styles.input}
                type="text"
                value={form.foto}
                onChange={(e) => handleChange('foto', e.target.value)}
                placeholder="https://..."
              />
            </div>

            <div style={{ ...styles.field, gridColumn: '1 / -1' }}>
              <label style={styles.label}>Descripción</label>
              <textarea
                style={{ ...styles.input, height: '100px', resize: 'vertical' }}
                value={form.descripcion}
                onChange={(e) => handleChange('descripcion', e.target.value)}
                placeholder="Descripción del evento..."
              />
            </div>

          </div>

          <div style={styles.formActions}>
            <button
              type="submit"
              style={{ ...styles.submitBtn, ...(submitting ? styles.btnDisabled : {}) }}
              disabled={submitting}
            >
              {submitting ? 'Guardando...' : isEdit ? 'Guardar cambios' : 'Crear evento'}
            </button>
            <Link to="/admin" style={styles.cancelLink}>Cancelar</Link>
          </div>
        </form>
      </main>
    </div>
  )
}

const styles = {
  page: { minHeight: '100vh', background: '#f8f8f8', fontFamily: 'sans-serif' },
  header: { background: '#1e1b4b', color: '#fff', padding: '1rem 2rem', display: 'flex', justifyContent: 'space-between', alignItems: 'center' },
  logo: { margin: 0, fontSize: '1.3rem' },
  navLink: { color: '#a5b4fc', textDecoration: 'none', fontWeight: '500' },
  main: { maxWidth: '720px', margin: '0 auto', padding: '2rem 1rem' },
  title: { marginBottom: '1.5rem', fontSize: '1.5rem' },
  form: { background: '#fff', borderRadius: '8px', boxShadow: '0 2px 8px rgba(0,0,0,0.07)', padding: '2rem' },
  grid: { display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(280px, 1fr))', gap: '1.25rem', marginBottom: '1.5rem' },
  field: { display: 'flex', flexDirection: 'column', gap: '0.3rem' },
  label: { fontSize: '0.85rem', fontWeight: '600', color: '#374151' },
  input: { padding: '0.55rem 0.75rem', border: '1px solid #d1d5db', borderRadius: '4px', fontSize: '0.95rem', fontFamily: 'sans-serif', outline: 'none' },
  inputError: { borderColor: '#dc2626' },
  fieldError: { fontSize: '0.78rem', color: '#dc2626' },
  apiError: { background: '#fef2f2', border: '1px solid #fca5a5', color: '#dc2626', padding: '0.75rem 1rem', borderRadius: '6px', marginBottom: '1rem', fontSize: '0.9rem' },
  formActions: { display: 'flex', alignItems: 'center', gap: '1rem' },
  submitBtn: { padding: '0.65rem 1.8rem', background: '#4f46e5', color: '#fff', border: 'none', borderRadius: '6px', fontSize: '1rem', fontWeight: '600', cursor: 'pointer' },
  btnDisabled: { background: '#9ca3af', cursor: 'not-allowed' },
  cancelLink: { color: '#6b7280', textDecoration: 'none', fontSize: '0.95rem' },
  center: { textAlign: 'center', padding: '3rem', color: '#6b7280' },
}
