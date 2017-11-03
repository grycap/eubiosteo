
import React, {Component} from 'react'
import {connect} from 'react-redux'

import Background from './background'
import {actions as ModalActions} from '../../actions/modal'

// -- importar los views --

import apply from './views/apply'
import postImage from './views/postImage'
import infoAlloc from './views/infoAlloc'
import deleteJob from './views/deletejob'
import postapihost from './views/postApiHost'

const views = {
    apply,
    postImage,
    infoAlloc,
    deleteJob,
    postapihost,
}

class Modal extends Component {
    render() {
        let {name, data} = this.props.modal;
        let View = views[name];

        return (
            <Background>
                <div className="dialog">
                    Ventana dentro
                    <a href="javascript:void(0)" onClick={this.props.close}>Close</a>
                    <View 
                        data={data}
                    />
                </div>
            </Background>
        )
    }
}

const mapStateToProps = (state, ownProps) => {
    return {
        modal: state.modal.ui
    }
}

const mapDispatchToProps = (dispatch, getState) => {
    return {
        close: () => {
            dispatch(ModalActions.closeModal())
        }
    }
}

export default connect(mapStateToProps, mapDispatchToProps)(Modal)

