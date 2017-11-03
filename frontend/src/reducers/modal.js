
import { combineReducers } from 'redux'

import {types as ModalTypes} from '../actions/modal'

// -- main info reducer --
let initialState = {
    isOpen: false,
    name: undefined,
    data: undefined
}

const modal = (state = initialState, action) => {

    switch(action.type) {
        case ModalTypes.MODAL_SHOW:
            let {name, data} = action;
            return {isOpen: true, name, data}
        case ModalTypes.MODAL_HIDE:
            return {isOpen: false}
    }

    return state
}

const rootReducer = combineReducers({
    ui: modal
})

export default rootReducer
