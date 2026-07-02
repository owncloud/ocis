import { AxiosPromise, AxiosResponse } from 'axios'
import { mock } from 'vitest-mock-extended'

export const mockAxiosResolve = <T>(data: T = {} as any): AxiosResponse<T> => {
  const response = mock<AxiosResponse>({ data })
  return response
}

export const mockAxiosReject = <T>(message = ''): AxiosPromise<T> => {
  return Promise.reject(new Error(message))
}
