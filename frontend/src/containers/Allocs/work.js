
import React from 'react'

import Table from '../table'
import {actions as AppActions} from '../../actions/app'

import {connect} from 'react-redux'

const downloadAction = "download"

const DownloadResults = ({item, onAction, onClick}) => {
    return (
        <div>
            <a
                href="javascript:void(0)"
                onClick={() => {
                    onAction('expand', item.ID)
                }}
            >
                {'Download'}
            </a>
        </div>
    )
}

const status = {
    0: ['SSS', 'green'],
    1: ['PENDING', 'yellow'],
    2: ['RUNNING', 'yellow'],
    3: ['COMPLETED', 'green'],
    4: ['CANCELLED', 'yellow'],
    5: ['FAILED', 'red'],
    6: ['UNKNOWN', 'black'],
    7: ['INPUT', 'orange'],
    8: ['OUPUT', 'orange'],
}

const AllocStatus = ({item}) => {
    let statusAttr = status[item.Status]
    let name = statusAttr[0]
    let color = statusAttr[1]

    return (
        <div style={{backgroundColor: color}}>
            {name}
        </div>
    )
}

const AllocView = ({item}) => {
    console.log("-- alloc item --")
    console.log(item)

    let {Initialtime, Finaltime, Elapsedtime} = item;
    let {Logoutput, Logerror} = item;
    let {Input, Output} = item;
    
    return (
        <div>
            {Logoutput}
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
        'width': 150,
        'link': true
    },
    {
        'name': 'Download',
        'field': false,
        'width': 150,
        'component': DownloadResults
    }
]

class AllocsView extends React.Component {
    render() {
        const {query} = this.props;

        return (
            <div>
                <Table 
                    fetchfn={AppActions.fetchWorkAllocs}
                    
                    query={query}
                    endpoint={'workallocs'} 
                    header={allocsHeader}
                    onClick={(field, id) => {

                    }}
                    onQueryChange={(query, str) => {

                    }}
                    onAction={(name, item) => {
                        if (name === downloadAction) {
                            this.props.downloadResults(item)
                        }
                    }}
                    expComp={AllocView}
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
        downloadResults: (id) => {
            
        }
    }
}

export default connect(mapStateToProps, mapDispatchToProps)(AllocsView)
