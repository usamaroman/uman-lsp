"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.deactivate = exports.activate = void 0;
const vscode = require("vscode");
const node_1 = require("vscode-languageclient/node");
let client;
let myTerminal;
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
    const runCodeCmd = vscode.commands.registerCommand("vs-uman.runCode", () => {
        vscode.window.showInformationMessage('Running code...');
        // Replace the line below with your actual code execution logic
        executeCode();
    });
    context.subscriptions.push(runCodeCmd);
    // Create a button in the toolbar
    const runCodeButton = vscode.window.createStatusBarItem(vscode.StatusBarAlignment.Left);
    runCodeButton.text = "$(play) Run Code"; // You can customize the icon and label
    runCodeButton.tooltip = "Run Code";
    runCodeButton.command = "vs-uman.runCode";
    context.subscriptions.push(runCodeButton);
    runCodeButton.show();
    client.start();
}
exports.activate = activate;
function executeCode() {
    myTerminal?.show(true);
    // Get the currently active text editor
    const editor = vscode.window.activeTextEditor;
    if (editor) {
        // Get the file path of the currently open file
        const filePath = editor.document.fileName;
        // Construct the command to run
        const command = `uman ${filePath}`;
        // Run the command in a child process
        if (!myTerminal) {
            myTerminal = vscode.window.createTerminal('My Terminal');
        }
        // Run the command in the terminal
        myTerminal.show(true); // Show the terminal and clear its content
        myTerminal.sendText(command);
    }
    else {
        vscode.window.showErrorMessage("No active text editor.");
    }
}
function deactivate() {
    if (!client) {
        return undefined;
    }
    return client.stop();
}
exports.deactivate = deactivate;
//# sourceMappingURL=extension.js.map