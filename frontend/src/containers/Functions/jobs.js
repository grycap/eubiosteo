
import React from 'react'

import Table from '../table'

import {actions as AppActions} from '../../actions/app'
import {actions as ModalActions} from '../../actions/modal'
import {actions as CreateActions} from '../../actions/create'

import {connect} from 'react-redux'

import { browserHistory } from 'react-router';

import Job from './job'

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

let jobsHeader = [
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

class JobsView extends React.Component {
    render() {
        const {view, query} = this.props;

        if (view != undefined) {
            return (
                <div>
                    <a
                        href="javascript:void(0)"
                        onClick={() => {
                            browserHistory.push("/functions/jobs")
                        }}
                    >
                        {'Back'}
                    </a>
                    <Job name={view} />
                </div>
            )
        }

        return (
            <div>
                <a
                    href="javascript:void(0)"
                    onClick={() => {
                        browserHistory.push('/functions/jobs/create')
                    }}
                >
                    {'Create'}
                </a>
                <Table 
                    fetchfn={AppActions.fetchJobs}

                    query={query}
                    endpoint={'jobs'} 
                    header={jobsHeader}
                    onClick={(field, id) => {
                        if (field == "ID") {
                            browserHistory.push("/functions/jobs?view=" + id)
                        }
                    }}
                    onQueryChange={(query, str) => {

                    }}
                    onAction={(action, item) => {
                        if (action == applyAction) {
                            this.props.applyJob(item.ID, item.Input)
                        }
                    }}
                />
            </div>
        )
    }
}

const mapStateToProps = (state, ownProps) => {
    console.log("-- props --")
    console.log(ownProps.location)
    console.log(ownProps.location.query.view)

    return {
        view: ownProps.location.query.view,
        query: ownProps.location.query
    }
}

const mapDispatchToProps = (dispatch, getState) => {
    return {
        uploadJob: function() {
            dispatch(ModalActions.postJob())
        },
        applyJob: function(id, inputs) {
            let fetchfn = CreateActions.postAlloc;
            dispatch(ModalActions.apply(id, fetchfn, inputs))
        }
    }
}

export default connect(mapStateToProps, mapDispatchToProps)(JobsView)
