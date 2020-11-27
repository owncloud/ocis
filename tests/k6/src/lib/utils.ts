export const randomString = (): string => {
    return Math.random().toString(36).slice(2)
}

export const extension = (p: string): string | undefined => {
    return (p.split('/').pop())!.split('.').pop()
}

