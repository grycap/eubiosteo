
import React, { Component, PropTypes } from 'react'
import { connect } from 'react-redux'

import {Text, Link} from '../../../components/Inputs'

import {actions as CreateActions} from '../../../actions/create'
import {actions as AppActions} from '../../../actions/app'

class PostApiHost extends Component {
    constructor(props) {
        super(props)

        this.state = {
            host: ''
        }

        this.onChange = this.onChange.bind(this)
    }

    onChange(val) {
        this.setState({
            host: val
        })
    }

    render() {
        let {host} = this.state
        return (
            <div>
                {'Post API Host'}
                <div>
                    <Text 
                        onChange={this.onChange} 
                        value={host} 
                    />
                    <br />
                    <a
                        href="javascript:void(0)" 
                        onClick={() => this.props.setHost(host)}
                    >
                        {'SetHost'}
                    </a>
                </div>
            </div>
        )
    }
}

const mapStateToProps = (state, ownProps) => {
    return {

    }
}

const mapDispathToProps = (dispatch) => {
    return {
        setHost: (newhost) => {
            dispatch(AppActions.changeApiHost(newhost))
        }
    }
}

export default connect(
    mapStateToProps, 
    mapDispathToProps
)(PostApiHost)
