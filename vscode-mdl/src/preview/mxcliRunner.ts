// SPDX-License-Identifier: Apache-2.0

import * as vscode from 'vscode';
import * as cp from 'child_process';

export async function generateElk(
	mxcliPath: string,
	mprFile: string,
	elementType?: string,
	qualifiedName?: string
): Promise<string | undefined> {
	return new Promise<string | undefined>((resolve) => {
		const elkType = elementType || 'systemoverview';
		const elkName = qualifiedName || 'SystemOverview';
		const args = ['describe', '-p', mprFile, '--format', 'elk', elkType, elkName];
		const env = { ...process.env, MXCLI_QUIET: '1' };

		cp.execFile(mxcliPath, args, { env, maxBuffer: 5 * 1024 * 1024 }, (err, stdout, stderr) => {
			if (err) {
				vscode.window.showErrorMessage(
					`Failed to generate system overview: ${stderr || err.message}`
				);
				resolve(undefined);
				return;
			}
			// Extract JSON object from output (skip any status messages)
			const jsonStart = stdout.indexOf('{');
			if (jsonStart < 0) {
				vscode.window.showErrorMessage('No JSON data in system overview output');
				resolve(undefined);
				return;
			}
			resolve(stdout.substring(jsonStart).trim());
		});
	});
}

export async function generateMermaid(
	mxcliPath: string,
	elementType: string,
	qualifiedName: string,
	mprFile: string
): Promise<string | undefined> {
	return new Promise<string | undefined>((resolve) => {
		const args = ['describe', '-p', mprFile, '--format', 'mermaid', elementType, qualifiedName];
		const env = { ...process.env, MXCLI_QUIET: '1' };

		cp.execFile(mxcliPath, args, { env, maxBuffer: 5 * 1024 * 1024 }, (err, stdout, stderr) => {
			if (err) {
				vscode.window.showErrorMessage(
					`Failed to generate diagram: ${stderr || err.message}`
				);
				resolve(undefined);
				return;
			}
			// Strip "Connected to:" line from output
			const lines = stdout.split('\n');
			const filtered = lines.filter(line => !line.startsWith('Connected to:'));
			resolve(filtered.join('\n').trim());
		});
	});
}

export async function findMprFile(mprPath: string | undefined): Promise<string | undefined> {
	if (mprPath) {
		return mprPath;
	}
	const files = await vscode.workspace.findFiles('**/*.mpr', '**/node_modules/**', 5);
	if (files.length === 0) {
		return undefined;
	}
	return files[0].fsPath;
}

export function escapeHtml(text: string): string {
	return text
		.replace(/&/g, '&amp;')
		.replace(/</g, '&lt;')
		.replace(/>/g, '&gt;')
		.replace(/"/g, '&quot;');
}

export function escapeForTemplate(str: string): string {
	return str
		.replace(/\\/g, '\\\\')
		.replace(/`/g, '\\`')
		.replace(/\$/g, '\\$');
}
