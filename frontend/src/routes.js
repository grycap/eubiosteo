
// -- crear las rutas --

import React from 'react'
import {Route} from 'react-router'

import App from './containers/app'

import {    
    Functions, 
    Jobs, 
    JobCreate,
    Workflow,
    WorkflowCreate
} from './containers/Functions'

import {
    Allocs, 
    Jobs as JobsAllocs, 
    Workflows as WorkflowsAllocs
} from './containers/Allocs'

import {
    Files, 
    List as FilesList,
} from './containers/Files'

export default (
    <Route path="/" component={App}>
        <Route path="functions" component={Functions} >
            <Route path="jobs/create" component={JobCreate} />
            <Route path="workflows/create" component={WorkflowCreate} />
            <Route path="jobs" component={Jobs} />
            <Route path="workflows" component={Workflow} />
        </Route>
        <Route path="allocs" component={Allocs} >
            <Route path="jobs" component={JobsAllocs} />
            <Route path="workflows" component={WorkflowsAllocs} />
        </Route>
        <Route path="files" component={Files} >
            <Route path="list" component={FilesList} />
        </Route>
    </Route>
)
