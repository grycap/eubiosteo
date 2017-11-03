
import {setItem, clearItem, getItem} from './storage'

//const API_ROOT = "http://158.42.105.15:10001"

const GET       = "GET"
const POST      = "POST"
const PUT       = "PUT"
const DELETE    = "DELETE"

const OK    = "OK"
const BAD   = "BAD"

// --- storage de los tokens en ./storage

export const storage = {
    setItem,
    clearItem,
    getItem
}

const noAuth = () => {
    return ({
        config: {

        }
    })
}

const authWithToken = (tokenname) => {
    let jwt = getItem(tokenname)
    if (jwt == undefined) {
        return ({error: 'No existe token jwt de nombre ' + tokenname})
    }

    return ({
        config :{
            'headers': {
                'Authorization': 'Bearer ' + jwt
            },
        }
    })
}

const projectAuth = (endpoint) => {
    return authWithToken('project')
}

const userAuth = (endpoint) => {
    return authWithToken('user')
}

const authCred = {
    project: projectAuth,
    user: userAuth,
    empty: noAuth,
}

const Api = (API_ROOT, endpoint, method, auth, data=undefined) => {
    let {config, error} = auth()
    if (!config) {
        // error al obtener la configuracion del auth
        Promise.reject(error)
    }
    
    config = Object.assign({}, config, {method: method})
    
    if (data != undefined) {
        let body = new FormData()
        for (var i in data) {
            body.append(i, data[i])
        }
        
        config = Object.assign({}, config, {body: body})
    }

    const fullUrl = (endpoint.indexOf(API_ROOT) === -1) ? API_ROOT + endpoint : endpoint
    
    console.log("- algo ? -")
    
    return fetch(fullUrl, config)
        .then(response => response.json().then(json => {

            console.log("-- have response --")
            console.log(json)

            if (!response.ok) {
                return Promise.reject(json)
            }
            
            let {Status, Message, Result} = json;

            console.log("-- fetch --")
            console.log(json)
            console.log(Result)

            switch (Status) {
                case OK:
                    console.log("OK")
                    return ({result: Result})
                case BAD:
                    console.log("BAD")
                    return ({error: Message})
                default:
                    console.log(json)
                    console.error("El resultado de la llamada deberia ser OK o BAD pero encontrado " + Status)
            }
            
        }))
}

// El problema es que todas las llamadas no se pueden gestionar desde aqui
// estan esas llamadas que usaban el middleware... O que hagan uso de esta
// libreria tambien

// esto solo son las llamadas api. Sin conexion
// con los reducers

export default {
    defaultProject: (root, endpoint, mode, fields=undefined) => Api(root, endpoint, mode, noAuth, fields),
    defaultGet: (root, endpoint) => Api(root, endpoint, GET, noAuth)
}
