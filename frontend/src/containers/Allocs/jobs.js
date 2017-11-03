
import React from 'react'

import Table from '../table'

import {actions as ModalActions} from '../../actions/modal'
import {actions as AppActions} from '../../actions/app'
import {actions as CreateActions} from '../../actions/create'

import {API_ROOT} from '../../api'

import {connect} from 'react-redux'

import AllocStatus from '../../components/status'
import { browserHistory } from 'react-router';

const viewAction = 'view'

const DownloadResults = ({item, onAction, onClick}) => {
    if (item.Status != 3) {
        return (<div></div>)
    }

    let endpoint = API_ROOT + '/allocs/' + item.ID + '/download'

    return (
        <div>
            <a
                href={endpoint}
            >
                {'Download'}
            </a>
        </div>
    )
}

const ViewResults = ({item, onAction, onClick}) => {
    return (
        <div>
            <a
                href="javascript:void(0)"
                onClick={() => {
                    onAction(viewAction, item.ID)
                }}
            >
                {'View'}
            </a>
        </div>
    )
}

function millisToMinutesAndSeconds(millis) {
    var minutes = Math.floor(millis / 60000);
    var seconds = ((millis % 60000) / 1000).toFixed(0);
    return minutes + ":" + (seconds < 10 ? '0' : '') + seconds;
}

const Timer = ({item}) => {
    let time = item.Elapsedtime

    if (time == 0) {
        return (<div></div>)
    }

    return (
        <div>
            {millisToMinutesAndSeconds(time)}
        </div>
    )
}

let allocsHeader = [
    {
        'name': 'Status',
        'field': false,
        'width': 5,
        'component': AllocStatus
    },
    {
        'name': 'ID',
        'label': 'ID',
        'width': 150
    },
    {
        'name': 'JobID',
        'label': 'JobID',
        'width': 150
    },
    {
        'name': 'Download',
        'field': false,
        'width': 50,
        'component': DownloadResults
    },
    {
        'name': 'View',
        'field': false,
        'width': 50,
        'component': ViewResults,
    },
    {
        'name': 'Time',
        'field': false,
        'width': 20,
        'component': Timer
    }
]

class AllocsView extends React.Component {
    render() {
        const {query} = this.props;

        return (
            <div>
                <Table 
                    fetchfn={AppActions.fetchAllocs}
                    
                    query={query}
                    endpoint={'allocs'} 
                    header={allocsHeader}
                    onClick={(field, id) => {

                    }}
                    onQueryChange={(query, str) => {
                        browserHistory.push('/allocs/jobs' + str)
                    }}
                    onAction={(name, item) => {
                        if (name === viewAction) {
                            this.props.showResults(item)
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
        showResults: (id) => {
            dispatch(ModalActions.infoAlloc(id))
        }
    }
}

export default connect(mapStateToProps, mapDispatchToProps)(AllocsView)
