
import React from 'react'
import Tools from '../components/tools'

import { browserHistory } from 'react-router';
import {connect} from 'react-redux'

import {actions as ModalActions} from '../actions/modal'

import Header from '../components/header'
import Modal from './Modal'

import Error from './error'

const tools = [
    {
        name: 'Functions',
        endpoint: '/functions',
        dest: '/functions/jobs'
    },
    {
        name: 'Allocs',
        endpoint: '/allocs',
        dest: '/allocs/jobs'
    },
    {
        name: 'Files',
        endpoint: '/files',
        dest: '/files/list'
    }
]

class App extends React.Component {
    render() {
        
        // -- modal --
        let modalUI = this.props.modal
        const {isOpen, name, data} = modalUI;

        // error message
        let error = this.props.error;

        return (
            <div style={{height: '100%'}}>
                {isOpen &&
                    <Modal />
                }
                <div className={"hg"}>
                    <Header
                        onHostChange={this.props.changeHost}
                    />
                    <main className={"hg-body"}>
                        <article className={"hg-content"} style={{position: 'relative'}}>
                            {this.props.children}
                            {error.used &&
                                <Error error={error} />
                            }
                        </article>
                        <nav className={"hg-nav"} style={{'position': 'relative'}}>
                            <Tools
                                pathname={this.props.pathname}
                                tools={tools}
                                onClick={(endpoint) => {
                                    browserHistory.push(endpoint)
                                }}
                            />
                            <img
                                className="logoi3m"
                                src="/static/images/i3m.jpg"
                            />
                        </nav>
                    </main>
                </div>
            </div>
        )
    }
}

const mapStateToProps = (state, ownProps) => {
    return {
        modal: state.modal.ui,
        error: state.error,
        pathname: ownProps.location.pathname
    }
}

const mapDispatchToProps = (dispatch, getState) => {
    return {
        changeHost: () => {
            dispatch(ModalActions.showModal('postapihost', {}))
        }
    }
}

export default connect(mapStateToProps, mapDispatchToProps)(App)
