
// para crear o modificar cosas.

// -- tipos genericos --

export const types = {
    REQUEST_CREATE: 'REQUEST_CREATE',
    SUCCESS_CREATE: 'SUCCESS_CREATE',
    FAILURE_CREATE: 'FAILURE_CREATE',
}

// TODO. Hay que poner aqui importante que parametro se debera actualizar en caso de que sea necesario
// O si se tiene que redirigir a algun sitio

const call = (mode, endpoint, fields, name, authCredentials='noAuth') => ({
    type: types.REQUEST_CREATE,
    mode,
    name,
    endpoint,
    fields,
    authCredentials
})

// put
const put = (endpoint, fields, name) => call('PUT', endpoint, fields, name)

// post
const post = (endpoint, fields, name) => call('POST', endpoint, fields, name)

// delete
const del = (endpoint, name) => call('DELETE', endpoint, {}, name)

export const actions = {
    deleteGeneric: (endpoint, name) => del(endpoint, name),
    postImage: (name, format, image) => post('/images', {name, format, file: image}, 'images'),
    postJob: (name, image, input, output, file, checkfiles) => post('/jobs', {name, image, input, output, file: file, checkfiles}, 'jobs'),
    
    postAlloc: (job, attrs) => post("/allocs", {job, attrs}, 'allocs'),
    postWorkflowAlloc: (id, attrs) => post("/workallocs", {id, attrs}, 'workallocs'),

    downloadAlloc: (alloc) => post(`/allocs/${alloc}/download`, {}),
    postWorkflow: (content) => post('/workflows', {content}, 'workflow')
}
