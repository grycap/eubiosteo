
import Jobs from './jobs'
import JobCreate from './jobcreate'

import Workflow from './workflows'
import WorkflowCreate from './workflowcreate'

import View from '../view'

const functionsView = {
    name: 'Functions',
    back: '/functions/jobs',
    routes: {   // rutas especiales para el back
        '/functions/jobs/create': {
            back: '/functions/jobs',
            name: 'Create Job'
        },
        '/functions/workflows/create': {
            back: '/functions/workflows',
            name: 'Create Workflow'
        }
    },
    views: [
        {
            name: 'Jobs',
            endpoint: '/functions/jobs'
        },
        {
            name: 'Workflows',
            endpoint: '/functions/workflows'
        }
    ]
}

const Functions = View(functionsView)

export {
    Functions, 
    Jobs, 
    JobCreate,
    Workflow,
    WorkflowCreate
}
