const addDays = (days: number): Date => {
  const date = new Date()
  date.setDate(date.getDate() + days)
  return date
}

const addMonth = (noOfMonths: number): Date => {
  const date = new Date()
  date.setMonth(date.getMonth() + noOfMonths)
  return date
}

export const getActualExpiryDate = (
  dateType: 'day' | 'week' | 'month' | 'year',
  dateOfExpiration: string
): Date => {
  switch (dateType) {
    case 'day':
      return addDays(parseInt(dateOfExpiration))
    case 'week':
      return addDays(parseInt(dateOfExpiration) * 7)
    case 'month':
      return addMonth(parseInt(dateOfExpiration))
    case 'year':
      return new Date(new Date().setFullYear(new Date().getFullYear() + parseInt(dateOfExpiration)))
  }
}
