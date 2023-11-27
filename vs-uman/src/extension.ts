import * as vscode from "vscode";
import * as child_process from 'child_process';

import {
  LanguageClient,
  LanguageClientOptions,
  ServerOptions,
  TransportKind,
} from "vscode-languageclient/node";

let client: LanguageClient;
let myTerminal: vscode.Terminal | undefined

export function activate(context: vscode.ExtensionContext) {
  // vscode.window.showInformationMessage(`Starting LSP`);
  let serverOptions: ServerOptions = {
    command: "umanlsp",
    transport: TransportKind.stdio,
  };

  let clientOptions: LanguageClientOptions = {
    documentSelector: [{ scheme: "file", language: "um" }],
  };

  client = new LanguageClient("um", "um", serverOptions, clientOptions);

  const runCodeCmd = vscode.commands.registerCommand("vs-uman.runCode", () => {
    vscode.window.showInformationMessage('Running code...');
        // Replace the line below with your actual code execution logic
        executeCode();
  })
  context.subscriptions.push(runCodeCmd)
  
   // Create a button in the toolbar
  const runCodeButton = vscode.window.createStatusBarItem(vscode.StatusBarAlignment.Left);
  runCodeButton.text = "$(play) Run Code"; // You can customize the icon and label
  runCodeButton.tooltip = "Run Code";
  runCodeButton.command = "vs-uman.runCode";
  context.subscriptions.push(runCodeButton)
  runCodeButton.show()
  
  client.start();
}

function executeCode() {
  myTerminal?.show(true)
  
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
  } else {
      vscode.window.showErrorMessage("No active text editor.");
  }
}

export function deactivate(): Thenable<void> | undefined {
  if (!client) {
    return undefined;
  }
  return client.stop();
}