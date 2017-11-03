
// -- storage basado en localStorage

export const setItem = (item, value) => {
    localStorage.setItem(item, value)
}

export const clearItem = (item) => {
    localStorage.removeItem(item)
}

export const getItem = (item) => {
    return localStorage.getItem(item)
}
