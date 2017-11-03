
export const types = {
    ERROR_SHOW: "ERROR_SHOW",
    ERROR_HIDE: "ERROR_HIDE"
}

export const actions = {
    showError: (msg) => ({type: types.ERROR_SHOW, msg}),
    hideError: () => ({type: types.ERROR_HIDE})
}
