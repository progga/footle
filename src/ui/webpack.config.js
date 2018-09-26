const path = require('path')

module.exports = {
  mode: 'development',
  entry: {
    ui: './scripts/ui.js'
  },
  output: {
    path: path.resolve(__dirname, 'build/scripts'),
    filename: 'ui-bundle.js'
  },
  devtool: 'source-map'
}
