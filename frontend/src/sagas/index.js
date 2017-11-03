
import { fork } from 'redux-saga/effects'

import {appGeneric, createGeneric} from './app'

function *mainSaga() {
    yield fork(appGeneric)
    yield fork(createGeneric)
}

export default mainSaga
