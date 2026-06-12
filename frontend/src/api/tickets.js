import client from './client'

export const buyTicket = (data) => client.post('/tickets', data)

export const getMyTickets = () => client.get('/tickets/mine')

export const cancelTicket = (id) => client.delete(`/tickets/${id}`)

export const transferTicket = (id, data) =>
  client.put(`/tickets/${id}/transfer`, data)
