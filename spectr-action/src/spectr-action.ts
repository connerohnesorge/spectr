import * as path from "node:path";
import * as core from "@actions/core";
import * as exec from "@actions/exec";
import {
  downloadVersion,
  resolveVersion,
  tryGetFromToolCache,
} from "./download/download-version";
import type { ValidationOutput } from "./types/spectr";
import { hasReport } from "./types/spectr";
import type { Architecture, Platform } from "./utils/platforms";
import { getArch, getPlatform } from "./utils/platforms";

/**
 * Main entry point for the GitHub Action
 */
async function run(): Promise<void> {
  try {
    // 1. Get inputs
    const version = core.getInput("version");
    const githubToken = core.getInput("github-token");
    const strict = core.getBooleanInput("strict");

    core.info(`Starting spectr validation (strict: ${strict})`);

    // 2. Setup platform and architecture
    const platform = getPlatform();
    const arch = getArch();

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
    } else {
      core.info("Spectr validation passed");
    }
  } catch (error) {
    core.setFailed(`Action failed: ${(error as Error).message}`);
  }
}

/**
 * Setup spectr binary (download/cache)
 * @returns Path to spectr executable
 */
async function setupSpectr(
  platform: Platform,
  arch: Architecture,
  versionInput: string,
  githubToken: string,
): Promise<string> {
  // Resolve version (handle 'latest', semver ranges, etc.)
  const resolvedVersion = await resolveVersion(
    versionInput || "latest",
    githubToken,
  );
  core.info(`Resolved version: ${resolvedVersion}`);
  core.setOutput("spectr-version", resolvedVersion);

  // Try to get from tool cache first
  const toolCacheResult = tryGetFromToolCache(arch, resolvedVersion);
  if (toolCacheResult.installedPath) {
    core.info(
      `Found spectr in tool-cache for version ${toolCacheResult.version}`,
    );
    const executableName =
      platform === "pc-windows-msvc" ? "spectr.exe" : "spectr";
    return path.join(toolCacheResult.installedPath, executableName);
  }

  // Download and cache the binary
  core.info(`Downloading spectr version ${resolvedVersion}...`);
  const downloadResult = await downloadVersion(
    platform,
    arch,
    resolvedVersion,
    githubToken,
  );

  const executableName =
    platform === "pc-windows-msvc" ? "spectr.exe" : "spectr";
  return path.join(downloadResult.cachedToolDir, executableName);
}

/**
 * Run spectr validation and return parsed JSON output
 */
async function runSpectrValidation(
  spectrPath: string,
  strict: boolean,
): Promise<ValidationOutput> {
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

  const options: exec.ExecOptions = {
    cwd: workspacePath,
    ignoreReturnCode: true, // Don't throw on non-zero exit (we'll handle it)
    listeners: {
      stderr: (data: Buffer) => {
        stderr += data.toString();
      },
      stdout: (data: Buffer) => {
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
    throw new Error(
      `No JSON output received from spectr validate. Exit code: ${exitCode}`,
    );
  }

  let validationOutput: ValidationOutput;
  try {
    validationOutput = JSON.parse(stdout) as ValidationOutput;
  } catch (parseError) {
    core.error(`Failed to parse JSON output from spectr validate`);
    core.error(`Raw output: ${stdout.substring(0, 500)}`);
    throw new Error(
      `Invalid JSON from spectr: ${(parseError as Error).message}`,
    );
  }

  core.info(`Validation completed: ${validationOutput.length} items validated`);
  return validationOutput;
}

/**
 * Process validation results and create GitHub annotations
 * @returns true if any errors were found
 */
async function processValidationResults(
  validationOutput: ValidationOutput,
): Promise<boolean> {
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
      core.error(
        `Failed to validate ${result.type} "${result.name}": ${result.error}`,
      );
      hasErrors = true;
      continue;
    }

    // Process results with validation report
    if (!hasReport(result)) {
      continue;
    }

    const { report } = result;
    const itemTitle = `${result.type}: ${result.name}`;

    // Log summary for this item
    if (report.summary.errors > 0) {
      core.error(
        `✗ ${itemTitle} - ${report.summary.errors} errors, ${report.summary.warnings} warnings`,
      );
      hasErrors = true;
    } else if (report.summary.warnings > 0) {
      core.warning(`⚠ ${itemTitle} - ${report.summary.warnings} warnings`);
    } else {
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
      } else if (issue.level === "WARNING") {
        core.warning(issue.message, annotationProps);
        totalWarnings++;
      } else {
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
