
import React, { Component } from 'react'

const status = {
    0: ['SSS', 'green'],
    1: ['PENDING', 'yellow'],
    2: ['RUNNING', 'yellow'],
    3: ['COMPLETED', 'green'],
    4: ['CANCELLED', 'yellow'],
    5: ['FAILED', 'red'],
    6: ['UNKNOWN', 'black'],
    7: ['INPUT', 'orange'],
    8: ['OUPUT', 'orange'],
}

const AllocStatus = ({item}) => {
    let statusAttr = status[item.Status]
    let name = statusAttr[0]
    let color = statusAttr[1]

    return (
        <div style={{backgroundColor: color}}>
            {name}
        </div>
    )
}

export default AllocStatus
