
import { combineReducers } from 'redux'
import {createPaginationReducer, createEntities} from './pagination'

import merge from 'lodash/merge'

const jobsPagination = createPaginationReducer('jobs', '/jobs/')
const allocsPagination = createPaginationReducer('allocs', '/allocs/')
const imagesPagination = createPaginationReducer('images', '/images')
const workflowsPagination = createPaginationReducer('workflows', '/workflows/')
const workAllocsPagination = createPaginationReducer('workallocs', '/workallocs/')

const pagination = combineReducers({
    'jobs': jobsPagination,
    'allocs': allocsPagination,
    'images': imagesPagination,
    'workflows': workflowsPagination,
    'workallocs': workAllocsPagination,
})

const initialState = {
    jobs: {},
}

const entities = createEntities('project', initialState)

const rootReducer = combineReducers({
    entities,
    pagination,
})

export default rootReducer
