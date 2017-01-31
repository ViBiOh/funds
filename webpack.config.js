const path = require('path');
const webpack = require('webpack');
const ExtractTextPlugin = require('extract-text-webpack-plugin');

const config = {
  context: path.join(__dirname, 'app'),
  entry: ['./index.jsx', './style.scss'],

  resolve: {
    modules: ['node_modules', 'src'],
    extensions: ['.js', '.jsx'],
  },

  module: {
    rules: [{
      test: /\.jsx?$/,
      enforce: 'pre',
      exclude: /node_modules/,
      loader: 'eslint-loader',
    }, {
      test: /\.jsx?$/,
      exclude: /node_modules/,
      use: 'babel-loader',
    }, {
      test: /\.scss$/,
      use: ExtractTextPlugin.extract({
        fallbackLoader: 'style-loader',
        loader: 'css-loader!sass-loader',
      }),
    }, {
      test: /\.css$/,
      use: ExtractTextPlugin.extract({
        loader: 'css-loader?modules&importLoaders=1&localIdentName=[name]_[local]_[hash:base64:5]!sass-loader',
      }),
    }],
  },

  plugins: [
    new ExtractTextPlugin({
      filename: 'app.css',
      allChunks: true,
    }),
  ],

  output: {
    filename: 'app.js',
    path: path.join(__dirname, 'dist/static'),
  },
};


if (process.env.PRODUCTION) {
  config.plugins.push(new webpack.DefinePlugin({
    'process.env': {
      NODE_ENV: JSON.stringify('production'),
    },
  }));
  config.plugins.push(new webpack.optimize.UglifyJsPlugin({
    compress: {
      warnings: false,
    },
  }));
}

module.exports = config;
