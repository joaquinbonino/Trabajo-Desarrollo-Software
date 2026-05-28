import client from './client'

export const buyTicket = (data) => client.post('/tickets', data)

export const getMyTickets = () => client.get('/tickets/me')

export const cancelTicket = (id) => client.patch(`/tickets/${id}/cancel`)

export const transferTicket = (id, data) =>
  client.patch(`/tickets/${id}/transfer`, data)
