
import React from 'react'
import {connect} from 'react-redux'

import {actions as ErrorActions} from '../actions/error'

class Error extends React.Component {
    render() {
        let error = this.props.error;
        
        return (
            <div
                className={'errorPanel'}
            >
                {error.msg}
                <a 
                    href="javascript:void(0)" 
                    onClick={this.props.close}
                >
                    {'Close'}
                </a>
            </div>
        )
    }
}

const mapStateToProps = (state, ownProps) => {
    return {
        error: state.error
    }   
}

const mapDispatchToProps = (dispatch, getState) => {
    return {
        close: () => {
            dispatch(ErrorActions.hideError())
        }
    }
}

export default connect(mapStateToProps, mapDispatchToProps)(Error)
