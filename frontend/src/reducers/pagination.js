
import {types} from '../actions/app'
import {types as ActionTypes} from '../actions/create'
import merge from 'lodash/merge'

const initialPage = {
    items: [],
    loaded: false,
    hasNext: false, // comprubea si hay otro luego
    time: undefined
}

const createEntities = (name, initialStateEntities) => {
    return (state=initialStateEntities, action) => {
        if (action.type == types.SUCCESS_PAGE) {
            if (action.data && action.data.entities) {
                if (action.where == name) {
                    return merge({}, state, action.data.entities)
                }
            }
        }
        
        return state
    }
}

const updatePage = (state=initialPage, action) => {
    switch(action.type) {
        case types.CANCELED_PAGE:
            // se ha cancelado la carga de la pagina. (Volver al estado inicial)
            return initialPage
            
        case types.REQUEST_PAGE:
            return Object.assign({}, state, {isLoading: true})

        case types.SUCCESS_PAGE:
            
            let res = action.data.result;
            return Object.assign({}, state, {items: res, loaded: true, isLoading: false, time: Date.now(), hasNext: action.hasnext})
        case types.FAILURE_PAGE:
        default:
            return state
    }
}

const initialQuery = {
    pages: {},
}

const updateQuery = (query=initialQuery, action, page) => {
    return Object.assign({}, query, {
        pages: Object.assign({}, query.pages, {
            [page]: updatePage(query.pages[page], action)
        })
    })
}

const initialPagination = {}

const createPaginationReducer = (name, endpoint) => {

    return (state=initialPagination, action) => {
        if (action.name != name) {  // no es para este reductor. Algo mes para saber que es un tipo de query de request_page
            return state
        }

        if (action.type == types.CHANGE_API_ROOT) {
            return initialPagination
        }
        
        if (action.type == ActionTypes.SUCCESS_CREATE) {
            return initialPagination
        }

        if (action.type == 'INVALIDATE') {
            return initialPagination
        }
        
        if (!action.paginate) {
            return state
        }

        let {query} = action;
        let {fields, page} = query;

        return Object.assign({}, state, {
            [fields]: updateQuery(state[fields], action, page)
        })
    }
}

export {createPaginationReducer, createEntities}
