import client from './client'

export const joinWaitlist = (eventId) =>
  client.post(`/events/${eventId}/waitlist`)

export const getMyWaitlist = () => client.get('/waitlist/mine')
