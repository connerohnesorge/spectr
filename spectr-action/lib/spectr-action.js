"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || (function () {
    var ownKeys = function(o) {
        ownKeys = Object.getOwnPropertyNames || function (o) {
            var ar = [];
            for (var k in o) if (Object.prototype.hasOwnProperty.call(o, k)) ar[ar.length] = k;
            return ar;
        };
        return ownKeys(o);
    };
    return function (mod) {
        if (mod && mod.__esModule) return mod;
        var result = {};
        if (mod != null) for (var k = ownKeys(mod), i = 0; i < k.length; i++) if (k[i] !== "default") __createBinding(result, mod, k[i]);
        __setModuleDefault(result, mod);
        return result;
    };
})();
Object.defineProperty(exports, "__esModule", { value: true });
const path = __importStar(require("node:path"));
const core = __importStar(require("@actions/core"));
const exec = __importStar(require("@actions/exec"));
const download_version_1 = require("./download/download-version");
const spectr_1 = require("./types/spectr");
const platforms_1 = require("./utils/platforms");
/**
 * Main entry point for the GitHub Action
 */
async function run() {
    try {
        // 1. Get inputs
        const version = core.getInput("version");
        const githubToken = core.getInput("github-token");
        const strict = core.getBooleanInput("strict");
        core.info(`Starting spectr validation (strict: ${strict})`);
        // 2. Setup platform and architecture
        const platform = (0, platforms_1.getPlatform)();
        const arch = (0, platforms_1.getArch)();
        if (platform === undefined) {
            throw new Error(`Unsupported platform: ${process.platform}`);
        }
        if (arch === undefined) {
            throw new Error(`Unsupported architecture: ${process.arch}`);
        }
        // 3. Setup spectr binary
        const spectrPath = await setupSpectr(platform, arch, version, githubToken);
        core.info(`Successfully installed spectr at ${spectrPath}`);
        // 4. Run spectr validation
        const validationOutput = await runSpectrValidation(spectrPath, strict);
        // 5. Process results and create annotations
        const hasErrors = await processValidationResults(validationOutput);
        // 6. Set action status
        if (hasErrors) {
            core.setFailed("Spectr validation failed with errors");
        }
        else {
            core.info("Spectr validation passed");
        }
    }
    catch (error) {
        core.setFailed(`Action failed: ${error.message}`);
    }
}
/**
 * Setup spectr binary (download/cache)
 * @returns Path to spectr executable
 */
async function setupSpectr(platform, arch, versionInput, githubToken) {
    // Resolve version (handle 'latest', semver ranges, etc.)
    const resolvedVersion = await (0, download_version_1.resolveVersion)(versionInput || "latest", githubToken);
    core.info(`Resolved version: ${resolvedVersion}`);
    core.setOutput("spectr-version", resolvedVersion);
    // Try to get from tool cache first
    const toolCacheResult = (0, download_version_1.tryGetFromToolCache)(arch, resolvedVersion);
    if (toolCacheResult.installedPath) {
        core.info(`Found spectr in tool-cache for version ${toolCacheResult.version}`);
        const executableName = platform === "pc-windows-msvc" ? "spectr.exe" : "spectr";
        return path.join(toolCacheResult.installedPath, executableName);
    }
    // Download and cache the binary
    core.info(`Downloading spectr version ${resolvedVersion}...`);
    const downloadResult = await (0, download_version_1.downloadVersion)(platform, arch, resolvedVersion, githubToken);
    const executableName = platform === "pc-windows-msvc" ? "spectr.exe" : "spectr";
    return path.join(downloadResult.cachedToolDir, executableName);
}
/**
 * Run spectr validation and return parsed JSON output
 */
async function runSpectrValidation(spectrPath, strict) {
    const workspacePath = process.env.GITHUB_WORKSPACE;
    if (!workspacePath) {
        throw new Error("GITHUB_WORKSPACE environment variable is not set");
    }
    // Build command arguments
    const args = ["validate", "--all", "--json"];
    if (strict) {
        args.push("--strict");
    }
    core.info(`Running: ${spectrPath} ${args.join(" ")}`);
    core.info(`Working directory: ${workspacePath}`);
    // Capture stdout
    let stdout = "";
    let stderr = "";
    const options = {
        cwd: workspacePath,
        ignoreReturnCode: true, // Don't throw on non-zero exit (we'll handle it)
        listeners: {
            stderr: (data) => {
                stderr += data.toString();
            },
            stdout: (data) => {
                stdout += data.toString();
            },
        },
    };
    const exitCode = await exec.exec(spectrPath, args, options);
    // Log stderr if present (warnings, debug info)
    if (stderr) {
        core.debug(`spectr stderr: ${stderr}`);
    }
    // Parse JSON output
    if (!stdout.trim()) {
        throw new Error(`No JSON output received from spectr validate. Exit code: ${exitCode}`);
    }
    let validationOutput;
    try {
        validationOutput = JSON.parse(stdout);
    }
    catch (parseError) {
        core.error(`Failed to parse JSON output from spectr validate`);
        core.error(`Raw output: ${stdout.substring(0, 500)}`);
        throw new Error(`Invalid JSON from spectr: ${parseError.message}`);
    }
    core.info(`Validation completed: ${validationOutput.length} items validated`);
    return validationOutput;
}
/**
 * Process validation results and create GitHub annotations
 * @returns true if any errors were found
 */
async function processValidationResults(validationOutput) {
    const workspacePath = process.env.GITHUB_WORKSPACE || ".";
    let hasErrors = false;
    let totalErrors = 0;
    let totalWarnings = 0;
    let totalInfo = 0;
    for (const result of validationOutput) {
        // Skip valid results with no issues
        if (result.valid && (!result.report || result.report.issues.length === 0)) {
            core.info(`✓ ${result.type}: ${result.name} - valid`);
            continue;
        }
        // Handle results with error field (validation couldn't run)
        if (result.error) {
            core.error(`Failed to validate ${result.type} "${result.name}": ${result.error}`);
            hasErrors = true;
            continue;
        }
        // Process results with validation report
        if (!(0, spectr_1.hasReport)(result)) {
            continue;
        }
        const { report } = result;
        const itemTitle = `${result.type}: ${result.name}`;
        // Log summary for this item
        if (report.summary.errors > 0) {
            core.error(`✗ ${itemTitle} - ${report.summary.errors} errors, ${report.summary.warnings} warnings`);
            hasErrors = true;
        }
        else if (report.summary.warnings > 0) {
            core.warning(`⚠ ${itemTitle} - ${report.summary.warnings} warnings`);
        }
        else {
            core.info(`ℹ ${itemTitle} - ${report.summary.info} info messages`);
        }
        // Create annotations for each issue
        for (const issue of report.issues) {
            const relativePath = path.relative(workspacePath, issue.path);
            const annotationTitle = itemTitle;
            const annotationProps = {
                file: relativePath,
                startLine: issue.line || 1,
                title: annotationTitle,
            };
            if (issue.level === "ERROR") {
                core.error(issue.message, annotationProps);
                totalErrors++;
            }
            else if (issue.level === "WARNING") {
                core.warning(issue.message, annotationProps);
                totalWarnings++;
            }
            else {
                core.notice(issue.message, annotationProps);
                totalInfo++;
            }
        }
    }
    // Log overall summary
    core.info("");
    core.info("=== Validation Summary ===");
    core.info(`Total items validated: ${validationOutput.length}`);
    core.info(`Errors: ${totalErrors}`);
    core.info(`Warnings: ${totalWarnings}`);
    core.info(`Info: ${totalInfo}`);
    return hasErrors;
}
// Execute the action
run();
