import axios from 'axios'
import { useAuthStore } from '@/store'

export const request = axios.create(
  {
    baseURL: import.meta.env.VITE_BASE_API,
    timeout: 12000,
  },
)

request.interceptors.request.use(
  // Request success interceptor
  (config) => {
    if (config.noNeedToken) {
      return config
    }

    const { token } = useAuthStore()
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  // Request failure interceptor
  (error) => {
    return Promise.reject(error)
  },
)

request.interceptors.response.use(
  // Response success interceptor
  (response) => {
    // Business information
    const responseData = response.data
    const { code, message, data } = responseData
    if (code !== 0) { // ! Business status code agreed with backend
      if (data && message !== data) {
        window.$message.error(`${message} ${data}`)
      }
      else {
        window.$message.error(message)
      }
      console.error(responseData) // Console output error information

      const authStore = useAuthStore()
      if (code === 1201) { // Token has issues
        authStore.toLogin()
        return
      }
      // 1202-Token expired
      if (code === 1202 || code === 1203 || code === 1207) {
        authStore.forceOffline()
        return
      }
      return Promise.reject(responseData)
    }
    return Promise.resolve(responseData)
  },
  // Response failure interceptor
  (error) => {
    // Mainly use business status codes to determine status, generally do not operate based on HTTP status codes
    const responseData = error.response?.data
    const { message, data } = responseData
    if (error.response.status === 500) {
      if (message && data) {
        window.$message.error(`${message} ${data}`)
      }
      else {
        window.$message.error('Server error')
      }
    }
    return Promise.reject(error)
  },
)
