import { api, siapkanCsrf } from './client.js';

export async function login({ email, password, ingat }) {
  await siapkanCsrf();
  return (await api.post('/auth/login', { email, password, ingat })).data;
}

export const logout = () => api.post('/auth/logout');

export const saya = () => api.get('/auth/me').then((r) => r.data);
