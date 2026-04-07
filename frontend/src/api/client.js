import axios from 'axios'

const API_BASE = import.meta.env.VITE_API_URL || '/api/v1'

const client = axios.create({
  baseURL: API_BASE,
  headers: {
    'Content-Type': 'application/json',
  },
})

client.interceptors.request.use((config) => {
  const token = localStorage.getItem('api_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

client.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('api_token')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export default {
  get: (url) => client.get(url),
  post: (url, data) => {
    if (data instanceof FormData) {
      return axios.post(`${API_BASE}${url}`, data, {
        headers: { 'Content-Type': 'multipart/form-data' },
      })
    }
    return client.post(url, data)
  },
  put: (url, data) => client.put(url, data),
  delete: (url) => client.delete(url),
}
