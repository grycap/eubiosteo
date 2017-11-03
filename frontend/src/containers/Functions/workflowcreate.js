
import React from 'react'
import { connect } from 'react-redux'

import CodeMirror from 'react-codemirror'
import {Text} from '../../components/Inputs'

import {actions as CreateActions} from '../../actions/create'

require('codemirror/mode/javascript/javascript');

const options = {
    lineNumbers: true,
}

class WorkflowCreate extends React.Component {
    constructor(props) {
        super(props)

        this.state = {
            code: ''
        }

        this.send = this.send.bind(this)
        this.onChange = this.onChange.bind(this)
    }
    
    send() {
        let {code} = this.state;
        this.props.send(code)
    }

    onChange(e) {
        this.setState({
            code: e
        })
    }

    render() {
        let {variables, steps} = this.state;

        return (
            <div>
                <h1>{'Workflow create'}</h1>
                <div>
                    <CodeMirror 
                        value={variables} 
                        onChange={this.onChange}
                        options={options} 
                    />
                </div>
                <a
                    href="javascript:void(0)"
                    onClick={() => {
                        this.send()
                    }}
                >
                    {'Crear'}
                </a>
            </div>
        )
    }
}


const mapStateToProps = (state, ownProps) => {
    return ownProps
}

const mapDispathToProps = (dispatch) => {
    return {
        send: (content) => {
            dispatch(CreateActions.postWorkflow(content))
        }
    }
}

export default connect(
    mapStateToProps, 
    mapDispathToProps
)(WorkflowCreate)

