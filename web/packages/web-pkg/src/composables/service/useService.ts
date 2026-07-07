import { inject } from 'vue'

export const useService = <T>(name: string): T => inject(name)
