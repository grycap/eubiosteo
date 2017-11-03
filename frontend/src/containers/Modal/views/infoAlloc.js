
import React, { Component, PropTypes } from 'react'
import { connect } from 'react-redux'

import AllocStatus from '../../../components/status'

const tableStyle = {
     'border': "1|0"
}

class InfoAlloc extends Component {
    constructor(props) {
        super(props)
    }

    render() {
        let {allocid, alloc} = this.props;

        return (
            <div>
                <h2>
                    {allocid}
                </h2>
                <div>
                    <table className={'allocview'}>
                        <tbody>
                            <tr>
                                <td>{'Status'}</td>
                                <td>
                                    <AllocStatus item={alloc} />
                                </td>
                            </tr>
                            <tr>
                                <td>{'Workflow'}</td>
                                <td>
                                    {alloc.WorkflowAllocID ? 'Yes' : 'No'}
                                </td>
                            </tr>
                            <tr>
                                <td>{'Error'}</td>
                                <td>
                                    {alloc.Error}
                                </td>
                            </tr>
                            <tr>
                                <td>{'Log Output'}</td>
                                <td>
                                    {alloc.Logoutput}
                                </td>
                            </tr>
                            <tr>
                                <td>{'Log Error'}</td>
                                <td>
                                    {alloc.Logerror}
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        )
    }
}

const mapStateToProps = (state, ownProps) => {
    let {allocid} = ownProps.data;
    let alloc = state.app.entities.allocs[allocid]

    return {
        allocid,
        alloc
    }
}

const mapDispathToProps = (dispatch) => {
    return {

    }
}

export default connect(
    mapStateToProps, 
    mapDispathToProps
)(InfoAlloc)
