
import React, { Component, PropTypes } from 'react'
import { connect } from 'react-redux'

import {actions as CreateActions} from '../../actions/create'
import {Text, Select} from '../../components/Inputs'

function makeid(){
    var text = "";
    var possible = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";

    for( var i=0; i < 5; i++ )
        text += possible.charAt(Math.floor(Math.random() * possible.length));

    return text;
}

const selections = [
	"bool",
	"number",
	"slice.number",
	"image.jpeg",
	"image.png",
	"image.nifti",
    "file.other"
]

const baseAttr = {
    name: '',
    value: ''
}

const Attr = ({id, name, value, onChange}) => (
    <tr>
        <td>
            <Text
                value={name}
                onChange={(val) => {
                    onChange({name: val, value: value})
                }}
            />
        </td>
        <td>
            <Select
                value={value}
                values={selections}
                onChange={(val) => {
                    onChange({name: name, value: val})
                }}
            />
        </td>
    </tr>
)

const AttrList = ({title, attrs, changeAttr, addAttr, hasNext}) => {
    return (
        <table>
        <tbody>
            <tr>
                <td>
                    {title}
                </td>
            </tr>
            <tr>
                <td>
                    <table>
                        <tbody>
                            {Object.keys(attrs).map((key, index) => {
                                let {name, value} = attrs[key]
                                
                                return (
                                    <Attr 
                                        key={index} 
                                        onChange={(val) => {changeAttr(key, val)}}
                                        {...attrs[key]} 
                                    />
                                )
                            })}
                        </tbody>
                    </table>
                    {hasNext &&
                        <a
                            href="javascript:void(0)"
                            onClick={addAttr}
                        >
                            {'+'}
                        </a>
                    }
                </td>
            </tr>
        </tbody>
        </table>
    )
}

function convertAttrs(attrs) {
    let newattrs = {}
    for (var i in attrs) {
        let {name, value} = attrs[i]
        newattrs[name] = value
    }

    return newattrs
}

class PostJob extends Component {
    constructor(props) {
        super(props)

        this.file = undefined;

        this.state = {
            name: '',
            image: '',
            input: {},
            output: {},
            checkfiles: false,
        }

        this.send = this.send.bind(this)
        this.hasNext = this.hasNext.bind(this)
        this.addAttr = this.addAttr.bind(this)
        this.changeAttr = this.changeAttr.bind(this)
        this.changeValue = this.changeValue.bind(this)
        this.onFileChange = this.onFileChange.bind(this)
    }

    onFileChange(e) {
        const file = e.target.files[0]
        this.file = file;
    }

    send() {
        let {name, image, attrs, input, output, checkfiles} = this.state;

        console.log("##########")
        console.log("checkfiles")
        console.log(checkfiles)

        let newinput = convertAttrs(input)
        let newoutput = convertAttrs(output)

        this.props.send(name, image, newinput, newoutput, this.file, checkfiles)
    }
    
    addAttr(field) {
        let id = makeid()

        this.setState({
            [field]: Object.assign({}, this.state[field], {[id]: baseAttr})
        })
    }

    changeValue(key, value) {
        this.setState(Object.assign({}, this.state, {[key]: value}))
    }

    changeAttr(field, id, value) {
        this.setState({
            [field]: Object.assign({}, this.state[field], {[id]: value})
        })
    }
    
    hasNext(field) {
        let attrs = this.state[field];
        for (var i in attrs) {
            let {name, value} = attrs[i]

            if (name == '' || value == '') {
                return false
            }
        }
        return true
    }

    render() {
        const {name, image} = this.state;
        let {input, output, checkfiles} = this.state;

        console.log("-- add attr --")
        console.log(this.addAttr)

        return (
            <div>
                {'Post Job'}
                <div>
                    <table>
                        <tbody>
                            <tr>
                                <td>
                                    {'Name'}
                                </td>
                                <td>
                                    <Text
                                        value={name}
                                        onChange={(e) => {this.changeValue('name', e)}}
                                    />
                                </td>
                            </tr>
                            <tr>
                                <td>
                                    {'Image'}
                                </td>
                                <td>
                                    <Text
                                        value={image}
                                        onChange={(e) => {this.changeValue('image', e)}}
                                    />
                                </td>
                            </tr>
                            <tr>
                                <td>
                                    {'Check files'}
                                </td>
                                <td>
                                    <input
                                        type="checkbox"
                                        checked={checkfiles}
                                        onChange={(e) => {this.changeValue('checkfiles', e.target.checked)}}
                                    />
                                </td>
                            </tr>
                        </tbody>
                    </table>
                    <AttrList
                        title={'Input'}
                        attrs={input}
                        addAttr={() => {this.addAttr('input')}} 
                        changeAttr={(id, value) => {this.changeAttr('input', id, value)}}
                        hasNext={this.hasNext('input')}
                    />
                    <AttrList
                        title={'Output'}
                        attrs={output}
                        addAttr={() => {this.addAttr('output')}} 
                        changeAttr={(id, value) => {this.changeAttr('output', id, value)}}
                        hasNext={this.hasNext('output')}
                    />
                    <input
                        type="file"
                        onChange={this.onFileChange}
                    />
                </div>
                <a
                    href="javascript:void(0)"
                    onClick={this.send}
                >
                    {'Create'}
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
        send: (name, image, input, output, file, checkfiles) => {
            console.log("- send -")
            console.log(name)
            console.log(image)
            console.log(input)
            console.log(output)
            console.log(checkfiles)

            dispatch(CreateActions.postJob(name, image, JSON.stringify(input), JSON.stringify(output), file, checkfiles))
        }
    }
}

export default connect(
    mapStateToProps, 
    mapDispathToProps
)(PostJob)
