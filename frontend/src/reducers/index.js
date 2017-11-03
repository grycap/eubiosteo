
import { combineReducers } from 'redux'

import { routerReducer as routing } from 'react-router-redux'

import app from './app'
import modal from './modal'
import error from './error'
import config from './config'

const rootReducer = combineReducers({
    routing,
    modal,
    config,
    error,
    app,
})

export default rootReducer
