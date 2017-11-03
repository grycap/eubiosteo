

import { normalize, schema } from 'normalizr'

// schema en forma de normalizr para convertir los objetos
// de las llamadas api

// -- jobs --

const jobSchema = new schema.Entity('jobs', {}, {
    idAttribute: (job) => job.ID
})

// -- allocs --

const allocSchema = new schema.Entity('allocs', {}, {
    idAttribute: (alloc) => alloc.ID
})

// -- work allocs --

const workAllocSchema = new schema.Entity('workallocs', {}, {
    idAttribute: (alloc) => alloc.ID
})

// -- images --

const imagesSchema = new schema.Entity('images', {}, {
    idAttribute: (image) => image.ID
})

// -- workflows --

const workflowsSchema = new schema.Entity('workflows', {}, {
    idAttribute: (workflow) => workflow.ID
})

export const Schemas = {
    JOB: jobSchema,
    JOB_ARRAY: [jobSchema],

    ALLOC: allocSchema,
    ALLOC_ARRAY: [allocSchema],

    IMAGE: imagesSchema,
    IMAGE_ARRAY: [imagesSchema],

    WORKFLOW: workflowsSchema,
    WORKFLOW_ARRAY: [workflowsSchema],

    WORKALLOC: workAllocSchema,
    WORKALLOC_ARRAY: [workAllocSchema]
}
