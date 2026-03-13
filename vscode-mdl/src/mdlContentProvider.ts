// SPDX-License-Identifier: Apache-2.0

import * as vscode from 'vscode';
import * as cp from 'child_process';
import { resolvedMxcliPath } from './extension';

export class MdlContentProvider implements vscode.TextDocumentContentProvider {
	private mxcliPath: string;
	private mprPath: string | undefined;

	constructor() {
		const config = vscode.workspace.getConfiguration('mdl');
		this.mxcliPath = resolvedMxcliPath();
		const configured = config.get<string>('mprPath', '');
		this.mprPath = configured || undefined;
	}

	updateConfig(): void {
		const config = vscode.workspace.getConfiguration('mdl');
		this.mxcliPath = resolvedMxcliPath();
		const configured = config.get<string>('mprPath', '');
		this.mprPath = configured || undefined;
	}

	async provideTextDocumentContent(uri: vscode.Uri): Promise<string> {
		// URI format: mendix-mdl://describe/<type>/<qualifiedName>
		// uri.authority = "describe", uri.path = "/<type>/<qualifiedName>"
		const parts = uri.path.split('/').filter(Boolean);
		if (parts.length < 2 || uri.authority !== 'describe') {
			return '-- Invalid URI format';
		}

		const elementType = parts[0];
		const qualifiedName = parts.slice(1).join('/');

		const mprFile = await this.findMprFile();
		if (!mprFile) {
			return '-- No .mpr file found. Set mdl.mprPath in settings.';
		}

		return new Promise<string>((resolve) => {
			const args = ['describe', '-p', mprFile, elementType, qualifiedName];
			const env = { ...process.env, MXCLI_QUIET: '1' };

			cp.execFile(this.mxcliPath, args, { env, maxBuffer: 5 * 1024 * 1024 }, (err, stdout, stderr) => {
				if (err) {
					resolve(`-- Error describing ${elementType} ${qualifiedName}:\n-- ${stderr || err.message}`);
					return;
				}
				// Strip "Connected to:" line from output
				const lines = stdout.split('\n');
				const filtered = lines.filter(line => !line.startsWith('Connected to:'));
				resolve(filtered.join('\n').trimStart());
			});
		});
	}

	private async findMprFile(): Promise<string | undefined> {
		if (this.mprPath) {
			return this.mprPath;
		}

		const files = await vscode.workspace.findFiles('**/*.mpr', '**/node_modules/**', 5);
		if (files.length === 0) {
			return undefined;
		}
		return files[0].fsPath;
	}
}
