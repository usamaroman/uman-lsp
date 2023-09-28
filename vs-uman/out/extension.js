"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.deactivate = exports.activate = void 0;
const node_1 = require("vscode-languageclient/node");
let client;
function activate(context) {
    // vscode.window.showInformationMessage(`Starting LSP`);
    let serverOptions = {
        command: "/home/chechyotka/projects/usamaroman/uman-lsp/umanlsp",
        transport: node_1.TransportKind.stdio,
    };
    let clientOptions = {
        documentSelector: [{ scheme: "file", language: "um" }],
    };
    client = new node_1.LanguageClient("um", "um", serverOptions, clientOptions);
    client.start();
}
exports.activate = activate;
function deactivate() {
    if (!client) {
        return undefined;
    }
    return client.stop();
}
exports.deactivate = deactivate;
//# sourceMappingURL=extension.js.map