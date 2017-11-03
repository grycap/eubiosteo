
import React, {Component} from 'react'

export default class Background extends Component {
    render() {
        return (
            <div className="modal-background">
                {this.props.children}
            </div>
        )
    }
}
