
import React from 'react'
import {connect} from 'react-redux'

import Table from '../table'

import {actions as AppActions} from '../../actions/app'
import {actions as ModalActions} from '../../actions/modal'
import {API_ROOT} from '../../api'

import { browserHistory } from 'react-router';

const DownloadFile = ({item, onAction, onClick}) => {
    let endpoint = API_ROOT + '/images/' + item.ID + '/download'
    
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

let imagesHeader = [
    {
        'name': 'Name',
        'label': 'Name',
        'width': 150
    },
    {
        'name': 'Format',
        'width': 150
    },
    {
        'name': 'Download',
        'field': false,
        'width': 50,
        'component': DownloadFile
    },
]

class ImagesView extends React.Component {
    render() {
        const {query} = this.props;

        return (
            <div>
                <a
                    href="javascript:void(0)"
                    onClick={this.props.uploadImage}
                >
                    {'Upload'}
                </a>
                <Table 
                    fetchfn={AppActions.fetchImages}

                    query={query}
                    endpoint={'images'} 
                    header={imagesHeader}
                    onClick={(field, id) => {

                    }}
                    onQueryChange={(query, str) => {
                        browserHistory.push('/files/list' + str)
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
        uploadImage: () => {
            dispatch(ModalActions.postImage())
        }
    }
}

export default connect(mapStateToProps, mapDispatchToProps)(ImagesView)
