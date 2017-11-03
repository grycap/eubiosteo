
export const types = {
    MODAL_SHOW: "MODAL_SHOW",
    MODAL_HIDE: "MODAL_HIDE"
}

export const actions = {
    showModal: (name, data={}) => ({type: types.MODAL_SHOW, name, data}),
    closeModal: () => ({type: types.MODAL_HIDE}),

    apply: (id, fetchfn, inputs) => actions.showModal('apply', {id, fetchfn, inputs}),
    
    // -- specific visuals --
    postImage: () => actions.showModal('postImage'),
    postJob: () => actions.showModal('postJob'),
    applyJob: (jobid) => actions.showModal('applyJob', {jobid}),
    infoAlloc: (allocid) => actions.showModal('infoAlloc', {allocid})
}
