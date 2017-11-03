
import React from 'react'

import Table from '../table'
import {actions as AppActions} from '../../actions/app'
import {actions as CreateActions} from '../../actions/create'
import {actions as ModalActions} from '../../actions/modal'

import {connect} from 'react-redux'

import { browserHistory } from 'react-router';

const applyAction = 'apply'

const ApplyJob = ({item, onAction, onClick}) => {
    return (
        <div>
            <a
                href="javascript:void(0)"
                onClick={() => {
                    onAction(applyAction, item)
                }}
            >
                {'Apply'}
            </a>
        </div>
    )
}

let workflowsHeader = [
    {
        'name': 'ID',
        'label': 'ID',
        'width': 150,
        'link': true
    },
    {
        'name': 'Apply',
        'field': false,
        'width': 150,
        'component': ApplyJob
    }
]

class WorkflowsView extends React.Component {
    render() {
        const {query} = this.props;
        
        return (
            <div>
                <a
                    href="javascript:void(0)"
                    onClick={() => {
                        browserHistory.push('/functions/workflows/create')
                    }}
                >
                    {'Create'}
                </a>
                <Table
                    fetchfn={AppActions.fetchWorkflows}

                    query={query}
                    endpoint={'workflows'} 
                    header={workflowsHeader}
                    onClick={(field, id) => {

                    }}
                    onQueryChange={(query, str) => {

                    }}
                    onAction={(action, item) => {
                        if (action == applyAction) {
                            this.props.applyWorkflow(item)
                        }
                    }}
                />
            </div>
        )
    }
}

const mapStateToProps = (state, ownProps) => {
    return {
        query: ownProps.location.query
    }
}

const mapDispatchToProps = (dispatch, getState) => {
    return {
        applyWorkflow: (item) => {
            // sacar de este item
            console.log("--- workflow ejecutar ---")
            console.log(item)

            let fetchfn = CreateActions.postWorkflowAlloc;
            dispatch(ModalActions.apply(item.ID, fetchfn, item.Entry))
        }
    }
}

export default connect(mapStateToProps, mapDispatchToProps)(WorkflowsView)
