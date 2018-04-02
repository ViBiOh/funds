const path = require('path');
const ExtractTextPlugin = require('extract-text-webpack-plugin');

const config = {
  context: path.join(__dirname, 'ui'),
  entry: ['babel-polyfill', './index.jsx', './index.css'],

  resolve: {
    modules: ['node_modules', 'src'],
    extensions: ['.js', '.jsx'],
  },

  module: {
    rules: [
      {
        test: /\.jsx?$/,
        enforce: 'pre',
        exclude: /node_modules/,
        use: 'eslint-loader',
      },
      {
        test: /\.jsx?$/,
        exclude: /node_modules/,
        use: 'babel-loader',
      },
      {
        test: /\.css$/,
        use: ExtractTextPlugin.extract({
          fallback: 'style-loader',
          use: 'css-loader',
        }),
      },
      {
        test: /\.less$/,
        use: ExtractTextPlugin.extract({
          use:
            'css-loader?modules&importLoaders=1&localIdentName=[name]_[local]_[hash:base64:5]!less-loader',
        }),
      },
    ],
  },

  plugins: [
    new ExtractTextPlugin({
      filename: 'main.css',
      allChunks: true,
    }),
  ],

  output: {
    filename: 'index.js',
    path: path.join(__dirname, 'ui/dist/static'),
  },

  devServer: {
    setup: (ui) => {
      ui.get('/env', (req, res) => {
        res.json({ API_URL: process.env.API_URL });
      });
    },
  },
};

module.exports = config;
