var path = require("path");
var webpack = require('webpack');

module.exports = {
	entry: ["babel-polyfill", "./src/index"],
	output: {
		path: path.join(__dirname, "/public"),
		filename: 'bundle.js'
	},
	devServer: {
		contentBase: './public',
		historyApiFallback: true
	},
    resolve: {
        extensions: ['*','.js']
    },
	module: {
		rules: [
			{
				test: /\.json$/,
				loader: 'json-loader'
			},
            { 
				test: /\.css$/, 
				loader: "style-loader!css-loader" 
			},
			{
				test: /\.js?$/,
				exclude: /(node_modules|bower_components)/,
				loader: 'babel-loader',
				query: {
					presets: ['es2015', 'react']
				}
			}
		]
	}
}
