
import React, { Component, PropTypes } from 'react'
import { connect } from 'react-redux'

import Table from '../../table'
import {actions as AppActions} from '../../../actions/app'
import {actions as CreateActions} from '../../../actions/create'

function sliceToString(val) {
    val = val.filter(i => !isNaN(i))
    let valStr = val.map(i => i.toString())
    return valStr.join(",")
}

function stringToSlice(val) {

    let valSlice = val.split(",")
    let valInt = valSlice.map(i => parseInt(i))

    return valInt
}

const SliceType = ({value=[], onChange}) => {
    return (
        <input
            type="text"
            onChange={(e) => {
                onChange(stringToSlice(e.target.value))
            }}
            value={sliceToString(value)}
        />
    )
}

const TextType = ({value, onChange}) => (
    <input 
        type="text" 
        onChange={(e) => {
            onChange(e.target.value)
        }} 
        value={value} 
    />
)

const ACTION_LOAD_IMAGE = 'loadimage'
const ACTION_IMAGE_LOADED = 'imageloaded'
const ACTION_MOVE_BACK = 'move_back'

const ImageType = ({value, onChange, onAction}) => (
    <div>
        <TextType
            value={value}
        />
        <a
            href="javascript:void(0)"
            onClick={() => {
                onAction(ACTION_LOAD_IMAGE)
            }}
        >
            {'#'}
        </a>
    </div>    
)

const NumberType = ({value, onChange}) => (
    <input 
        type="number" 
        onChange={(e) => {
            onChange(parseInt(e.target.value))
        }} 
        value={value == '' ? 0 : value} 
    />
)

// para annadir mas cosas
function inputTypes(type) {
    if (type.startsWith('image')) {
        return ImageType
    }

    if (type == "file.other") {
        return ImageType
    }

    if (type == 'number') {
        return NumberType
    }

    return TextType
}

function inputValues(job) {
    let inputs = job.Input
    let values = {}

    for (var i in inputs) {
        values[i] = ''
    }

    return values
}

// -- views --

// -- input --

class InputValuesView extends React.Component {
    constructor(props) {
        super(props)
    }

    render() {
        let {onChange, onAction} = this.props;
        let {inputs, values} = this.props;
        
        if (Object.keys(values).length == 0) {
            values = inputValues(inputs)
        }

        return (
            <div>
                <table>
                    <thead>
                        <tr>
                            <td>{'Name'}</td>
                            <td>{'Type'}</td>
                            <td>{'Value'}</td>
                        </tr>
                    </thead>
                    <tbody>
                        {Object.keys(inputs).map((name, index) => {
                            let value = values[name]
                            let type  = inputs[name]

                            let View = inputTypes(type)

                            return (
                                <tr key={index}>
                                    <td>{name}</td>
                                    <td>{type}</td>
                                    <td>
                                        <View 
                                            value={value}
                                            onChange={(e) => {
                                                this.props.onChange(viewValues, Object.assign({}, values, {[name]: e}))
                                            }}
                                            onAction={(action) => {
                                                this.props.onAction(action, {name, type})
                                            }}
                                        />
                                    </td>
                                </tr>
                            )
                        })}
                    </tbody>
                </table>
                <div>
                    <a
                        href="javascript:void(0)"
                        onClick={this.props.send}
                    >
                        {'Apply'}
                    </a>
                </div>
            </div>
        )
    }
}

// -- images --

const ApplyImage = ({item, onAction}) => {
    return (
        <div>
            <a
                href="javascript:void(0)"
                onClick={() => {
                    onAction('select', item.ID)
                }}
            >
                {'Select'}
            </a>
        </div>
    )
}

let imagesHeader = [
    {
        'name': 'Name',
        'label': 'Name',
        'width': 150,
    },
    {
        'name': 'Select',
        'field': false,
        'width': 150,
        'component': ApplyImage
    }
]

class InputImageSelect extends React.Component {
    constructor(props) {
        super(props)

        this.state = {
            query: {}
        }

        this.onQueryChange = this.onQueryChange.bind(this)
    }

    onQueryChange(query) {
        this.setState({
            query: query
        })
    }
    
    render() {
        const {values} = this.props;
        const {name, type} = values;

        let {query} = this.state;

        return (
            <div>
                {'load images'}
                <div>
                    {name}
                </div>
                <a
                    href="javascript:void(0)"
                    onClick={() => {
                        this.props.onAction(ACTION_MOVE_BACK)
                    }}
                >
                    {'Back'}
                </a>
                <div>
                    <Table 
                        fetchfn={AppActions.fetchImages}

                        query={query}
                        endpoint={'images'} 
                        header={imagesHeader}
                        onClick={(field, id) => {

                        }}
                        onQueryChange={(query, str) => {
                            this.onQueryChange(query)
                        }}
                        onAction={(old, id) => {
                            this.props.onAction(ACTION_IMAGE_LOADED, {name, id})
                        }}
                    />
                </div>
            </div>
        )
    }
}

const viewValues = 'VALUES'
const viewImages = 'IMAGES'

const viewToClass = {
    [viewValues]: InputValuesView,
    [viewImages]: InputImageSelect
}

const CHANGEVIEW = 'CHANGE_VIEW'

function postProcessSingleValue(value, type) {

    console.log("-- value --")
    console.log(value)
    console.log(type)

    if (type == "slice.number") {
        try {
            let valueSlice = value.split(",")
            return {value: valueSlice.map(i => parseInt(i))}
        } catch (error) {
            return {error}
        }
    }

    return {value}
}

function postProcessValues(values, inputs) {
    let newvalues = {}
    for (var name in values) {
        let inputType = inputs[name]

        let {value, error} = postProcessSingleValue(values[name], inputType)
        if (error != undefined) {
            console.error(error)
        }

        newvalues[name] = value
    }

    return newvalues
}

class PostAlloc extends Component {
    constructor(props) {
        super(props)

        this.state = {
            view: viewValues,
            values: {}
        }

        this.send = this.send.bind(this)
        this.onAction = this.onAction.bind(this)
        this.onChange = this.onChange.bind(this)
    }

    send() {
        let values = this.state.values[viewValues];
        let {id, inputs, fetchfn} = this.props;

        // TODO. mirar que todos los inputs del job estan en el values

        if (Object.keys(values).length == 0) {
            return
        }

        let newvalues  = postProcessValues(values, inputs)
        this.props.send(id, fetchfn, newvalues)
    }
    
    // cambiar el view
    onAction(actioname, info) {
        
        // cargar imagen
        if (actioname == ACTION_LOAD_IMAGE) {
            this.setState({
                view: viewImages,
                values: Object.assign({}, this.state.values, {[viewImages]: info})
            })
        }

        // volver al values
        if (actioname == ACTION_MOVE_BACK) {
            this.setState({
                view: viewValues
            })
        }
        
        // volver al image loaded
        if (actioname == ACTION_IMAGE_LOADED) {
            let {name, id} = info;

            let oldvalues = this.state.values[viewValues]
            let newvalues = Object.assign({}, oldvalues, {[name]: id})

            console.log("-- new values --")
            console.log(newvalues)

            console.log(name)
            console.log(id)
            console.log(info)
            
            this.setState({
                view: viewValues,
                values: Object.assign({}, this.state.values, {[viewValues]: newvalues})
            })
        }
    }
    
    // cambiar valores
    onChange(viewname, newvalues) {
        let {view, values} = this.state;

        this.setState({
            values: Object.assign({}, values, {
                [view]: newvalues
            })
        })
    }

    render() {
        let {view, values} = this.state;
        let {inputs} = this.props;

        let View = viewToClass[view]
        let viewValues = values[view] || {}
        
        return (
            <div>
                <View
                    inputs={inputs}
                    values={viewValues}
                    onChange={this.onChange}
                    onAction={this.onAction}
                    send={this.send}
                />
            </div>
        )
    }
}

const mapStateToProps = (state, ownProps) => {
    let {id, inputs, fetchfn} = ownProps.data;
    
    return {
        id,
        inputs,
        fetchfn,
    }
}

const mapDispathToProps = (dispatch) => {
    return {
        send: (id, fetchfn, attrs) => {
            let newattrs = JSON.stringify(attrs)

            console.log("-- attrs --")
            console.log(id)
            console.log(fetchfn)
            console.log(attrs)
            console.log(newattrs)

            dispatch(fetchfn(id, newattrs))
        }
    }
}

export default connect(
    mapStateToProps, 
    mapDispathToProps
)(PostAlloc)
