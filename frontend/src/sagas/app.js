
import { delay } from 'redux-saga'
import { call, put, cancelled, fork, take, cancel, select } from 'redux-saga/effects'
import Api, {storage, Api as defaultAPI} from '../api'

import {types as AppTypes, actions as AppActions} from '../actions/app'
import {types as CreateTypes, actions as CreateActions} from '../actions/create'
import {types as ErrorTypes} from '../actions/error'

import {actions as ModalActions} from '../actions/modal'

import { normalize, schema } from 'normalizr';
import {Schemas} from '../actions/schema'

// POST. DELETE. PUT

export const getHost = (state) => state.config.host

function* createGeneric() {
    while(true) {
        const {name, mode, endpoint, fields, authCredentials} = yield take(CreateTypes.REQUEST_CREATE)
        let root_host = yield select(getHost)

        console.log("-- REQUEST CREATE ss --")
        


        /////// se suposa que desde que fas el request aci.... fins a que tens un resultat
        /////// en algun dels trycatch.. es que esta carregant algo...

        /////// Para el modal estaria be que mirarra tots els request_create i que mire si el nom que te
        /////// es el mateix que tenia ell (es a dir), el request sa fet mentres ell estaba actiu
        /////// en ixe momento se dona conter que les dades son per a ell i pot mostrar els missatges de error si fa falta 

        const {result, error} = yield call(Api.defaultProject, root_host, endpoint, mode, fields)
        if (error == undefined) {
            yield put({type: CreateTypes.SUCCESS_CREATE, name})
        } else {
            yield put({type: CreateTypes.FAILURE_CREATE, error})
            yield put({type: ErrorTypes.ERROR_SHOW, msg: error})
        }
        
        // que se fa despres???
        // eliminar la modal...
        console.log("-- cerrar modal --")
        

        yield put(ModalActions.closeModal())

        // fer el fetch per a buscar dades noves...
        // o fer redirect al altre...

    }
}

// GET. single o con pagination

function* appGeneric() {
    while(true) {
        
        const {schemaName, query, name, endpoint, paginate} = yield take(AppTypes.REQUEST_PAGE)
        let root_host = yield select(getHost)

        console.log("-- requested --")
        console.log(schemaName)
        console.log(name)
        console.log(endpoint)
        console.log(query)

        const {fields, page} = query;

        const {result} = yield call(Api.defaultGet, root_host, endpoint + fields + '&page=' + page, 'empty')
        const {Result, HasNext} = result;

        console.log("-- result --")
        console.log(Result)
        
        let schema = Schemas[schemaName]

        console.log("-- schema --")
        console.log(schema)

        let res = normalize(Result, schema)
        
        console.log("-- res --")
        console.log(res)

        yield put(AppActions.fetchSuccess(name, res, query, paginate, 'project', HasNext))
        
        /*

        console.log(name)
        console.log(schema)
        console.log(endpoint)

        const {result} = yield call(Api.defaultProjectGet, endpoint)

        console.log("-- result --")
        console.log(result)

        let res = normalize(result, schema)

        console.log(res)
        yield put(AppActions.fetchSuccess(name, name, res))

        */

    }
}

export {appGeneric, createGeneric}
