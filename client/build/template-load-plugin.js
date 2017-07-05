// template-load-plugin.js

function MyPlugin(options) {
    // Configure your plugin with options...
    this.mode = options.mode || 'production'
}

MyPlugin.prototype.apply = function (compiler) {
    var self = this
    compiler.plugin('compilation', function (compilation) {

        compilation.plugin('html-webpack-plugin-before-html-processing', function (htmlPluginData, callback) {
            // do nothing in production mode
            // remove all cdn scripts in development mode
            // use unminified cdn scripts in local server mode
            if (self.mode != 'production') {
                htmlPluginData.html = htmlPluginData.html.replace(/<!--cdn-begin-->([^]*)<!--cdn-end-->/g, function (match) {
                    if (self.mode == 'development') {
                        return ''
                    } else {
                        return match.replace(/\.min\./g, '.')
                    }
                })
            }
            callback(null, htmlPluginData);
        });
    });

};

module.exports = MyPlugin;
