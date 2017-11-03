
import React, { Component, PropTypes } from 'react'
import {connect} from 'react-redux'

import {actions as ModalActions} from '../../actions/modal'

const Features = ({feats}) => (
    <div>
        {Object.keys(feats).map(function(key, index) {
            let feat = feats[key]

            return (
                <div key={index}>
                    {key} : {feat}
                </div>
            )
        })}
    </div>
)

class Job extends Component {
    render() {
        let {job} = this.props;
        let {ID, Input, Output, DriverImage} = job;

        return (
            <div>
                <h2>{ID}</h2>
                <div>
                    <div>
                        <h5>{'Input'}</h5>
                        <Features feats={Input} />
                    </div>
                    <div>
                        <h5>{'Output'}</h5>
                        <Features feats={Output} />
                    </div>
                </div>
                <hr />
                <a
                    href="javascript:void(0)"
                    onClick={() => {
                        this.props.delete(ID)
                    }}
                >
                {'Delete'}
                </a>
            </div>
        )
    }
}

const mapStateToProps = (state, ownProps) => {
    let job = state.app.entities.jobs[ownProps.name]

    return {
        job
    }
}

const mapDispathToProps = (dispatch) => {
    return {
        delete: (name) => {
            dispatch(ModalActions.showModal('deleteJob', {name}))
        }
    }
}

export default connect(
    mapStateToProps, 
    mapDispathToProps
)(Job)
