import client from './client'

export const getEvents = (categoria) =>
  client.get('/events', { params: { categoria } })

export const getEvent = (id) => client.get(`/events/${id}`)

export const createEvent = (data) => client.post('/events', data)

export const cancelEvent = (id) => client.patch(`/events/${id}/cancel`)

export const getAllEventsAdmin = () => client.get('/admin/events')

export const updateEvent = (id, data) => client.put(`/events/${id}`, data)
