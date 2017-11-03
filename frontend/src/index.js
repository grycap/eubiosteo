
// -- react --
import React from 'react'
import { render } from 'react-dom'

import reducer from './reducers'
import mySaga from './sagas'

import routes from './routes'

// -- redux --
import { createStore, applyMiddleware, compose } from 'redux'
import { Provider } from 'react-redux'
import createSagaMiddleware from 'redux-saga'

// -- router --
import { Router, browserHistory } from 'react-router'
import { syncHistoryWithStore } from 'react-router-redux'

// -- devtools --
import { createDevTools } from 'redux-devtools'
import LogMonitor from 'redux-devtools-log-monitor'
import DockMonitor from 'redux-devtools-dock-monitor'

const DevTools = createDevTools(
    <DockMonitor toggleVisibilityKey="ctrl-h" changePositionKey="ctrl-q">
        <LogMonitor theme="tomorrow" preserveScrollTop={false} />
    </DockMonitor>
)

// -- saga middleware --
const sagaMiddleware = createSagaMiddleware()

// -- store --
const store = createStore(
    reducer,
    compose (
        DevTools.instrument()
    ),
    applyMiddleware(
        sagaMiddleware
    )
)

sagaMiddleware.run(mySaga)

const history = syncHistoryWithStore(browserHistory, store)

const Root = () => (
    <Provider store={store}>
        <div id="sub-root">
            <Router history={history} routes={routes} />
        </div>
    </Provider>
)

render(
    <Root />,
    document.getElementById('root')
)
