
import { combineReducers } from 'redux'
import { types as ErrorTypes } from '../actions/error'

// -- main info reducer --
let initialState = {
    msg: '',
    used: false
}

const error = (state = initialState, action) => {
    switch(action.type) {
        case ErrorTypes.ERROR_SHOW:
            return {used: true, msg: action.msg}
        case ErrorTypes.ERROR_HIDE:
            return initialState
    }

    return state
}

export default error
