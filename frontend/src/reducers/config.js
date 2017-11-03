
import { combineReducers } from 'redux'
import { types as AppTypes } from '../actions/app'

// -- main info reducer --
let initialState = {
    host: 'http://localhost:10001'
}

const config = (state = initialState, action) => {
    switch(action.type) {
        case AppTypes.CHANGE_API_ROOT:
            return {host: action.host}
    }

    return state
}

export default config
