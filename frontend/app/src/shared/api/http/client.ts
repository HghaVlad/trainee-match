import axios, { type AxiosRequestConfig } from 'axios'

export const httpClient = axios.create({
  withCredentials: true,
})

export const mutatorFn = <T>(config: AxiosRequestConfig): Promise<T> => {
  return httpClient.request<T>(config).then((r) => r.data as T)
}

export default mutatorFn
