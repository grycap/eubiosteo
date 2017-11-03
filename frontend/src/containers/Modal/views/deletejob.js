
import React, { Component, PropTypes } from 'react'
import { connect } from 'react-redux'

import {actions as CreateActions} from '../../../actions/create'

class deleteJob extends Component {
    render() {
        let {name} = this.props;

        return (
            <div>
                <h4>{'Delete'} {name}</h4>
                <a
                    href="javascript:void(0)"
                    onClick={() => {
                        this.props.delete(name)
                    }}
                >
                    {'Confirm'}
                </a>
            </div>
        )
    }
}

const mapStateToProps = (state, ownProps) => {
    return {
        name: state.modal.ui.data.name
    }
}

const mapDispathToProps = (dispatch) => {
    return {
        delete: function(id) {
            dispatch(CreateActions.deleteGeneric(`/jobs/${id}`, 'jobs'))
        }
    }
}

export default connect(
    mapStateToProps, 
    mapDispathToProps
)(deleteJob)
