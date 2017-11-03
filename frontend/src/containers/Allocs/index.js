
import Workflows from './work'
import Jobs from './jobs'

import View from '../view'

const allocsView = {
    name: 'Allocs',
    views: [
        {
            name: 'Jobs',
            endpoint: '/allocs/jobs'
        },
        {
            name: 'Workflows',
            endpoint: '/allocs/workflows'
        }
    ],
    back: '/allocs/jobs',
}

const Allocs = View(allocsView)

export {
    Workflows, 
    Jobs, 
    Allocs
}
