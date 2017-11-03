
import { Schemas } from './schema'

// -- tipos genericos --

export const types = {
    REQUEST_PAGE: 'REQUEST_PAGE',
    SUCCESS_PAGE: 'SUCCESS_PAGE',
    FAILURE_PAGE: 'FAILURE_PAGE',
    CANCELED_PAGE: 'CANCELED_PAGE',
    CHANGE_API_ROOT: 'CHANGE_API_ROOT',
}

const emptyQuery = {
    fields: '?',
    page: 0
}

const fetch = (schemaName, query, name, endpoint, paginate=true) => ({
    type: types.REQUEST_PAGE,
    schemaName,
    query,
    name,
    endpoint,
    paginate,
})

export const actions = {
    fetchImages: (query) => fetch('IMAGE_ARRAY', query, 'images', '/images'),
    
    fetchAllocs: (query) => fetch('ALLOC_ARRAY', query, 'allocs', '/allocs'),
    fetchWorkAllocs: (query) => fetch('WORKALLOC_ARRAY', query, 'workallocs', '/workallocs'),
    
    fetchJobs: (query) => fetch('JOB_ARRAY', query, 'jobs', '/jobs'),
    fetchWorkflows: (query) =>  fetch('WORKFLOW_ARRAY', query, 'workflows', '/workflows'),
    
    fetchSuccess: (name, data, query, paginate, where, hasnext) => ({type: types.SUCCESS_PAGE, name, data, query, paginate, where, hasnext}),

    changeApiHost: (host) => ({type: types.CHANGE_API_ROOT, host})
}
