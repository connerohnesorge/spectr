import assert from "node:assert/strict";
import { describe, it } from "node:test";
import type {
  ValidationIssue,
  ValidationOutput,
} from "../../../src/types/spectr";

/**
 * Mock implementation of the validation output processing logic
 * This simulates the behavior from src/spectr-action.ts
 */

interface AnnotationProps {
  file: string;
  startLine: number;
  title: string;
}

interface ProcessingResult {
  hasErrors: boolean;
  totalErrors: number;
  totalWarnings: number;
  totalInfo: number;
  annotations: Array<{
    level: "ERROR" | "WARNING" | "INFO";
    message: string;
    props: AnnotationProps;
  }>;
  summaries: Array<{
    itemTitle: string;
    logLevel: "error" | "warning" | "info";
    message: string;
  }>;
}

/**
 * Process validation output and create annotations
 * This is a testable version of the logic from processValidationResults
 */
function processValidationOutput(
  validationOutput: ValidationOutput,
  workspacePath = ".",
): ProcessingResult {
  let hasErrors = false;
  let totalErrors = 0;
  let totalWarnings = 0;
  let totalInfo = 0;
  const annotations: ProcessingResult["annotations"] = [];
  const summaries: ProcessingResult["summaries"] = [];

  for (const result of validationOutput) {
    // Skip valid results with no issues
    if (result.valid && (!result.report || result.report.issues.length === 0)) {
      summaries.push({
        itemTitle: `${result.type}: ${result.name}`,
        logLevel: "info",
        message: "valid",
      });
      continue;
    }

    // Handle results with error field
    if (result.error) {
      hasErrors = true;
      summaries.push({
        itemTitle: `${result.type}: ${result.name}`,
        logLevel: "error",
        message: result.error,
      });
      continue;
    }

    // Process results with validation report
    if (!result.report) {
      continue;
    }

    const { report } = result;
    const itemTitle = `${result.type}: ${result.name}`;

    // Log summary for this item
    if (report.summary.errors > 0) {
      hasErrors = true;
      summaries.push({
        itemTitle,
        logLevel: "error",
        message: `${report.summary.errors} errors, ${report.summary.warnings} warnings`,
      });
    } else if (report.summary.warnings > 0) {
      summaries.push({
        itemTitle,
        logLevel: "warning",
        message: `${report.summary.warnings} warnings`,
      });
    } else {
      summaries.push({
        itemTitle,
        logLevel: "info",
        message: `${report.summary.info} info messages`,
      });
    }

    // Create annotations for each issue
    for (const issue of report.issues) {
      const relativePath = issue.path.startsWith(workspacePath)
        ? issue.path.slice(workspacePath.length + 1)
        : issue.path;

      const annotationProps: AnnotationProps = {
        file: relativePath,
        startLine: issue.line || 1,
        title: itemTitle,
      };

      if (issue.level === "ERROR") {
        annotations.push({
          level: "ERROR",
          message: issue.message,
          props: annotationProps,
        });
        totalErrors++;
      } else if (issue.level === "WARNING") {
        annotations.push({
          level: "WARNING",
          message: issue.message,
          props: annotationProps,
        });
        totalWarnings++;
      } else {
        annotations.push({
          level: "INFO",
          message: issue.message,
          props: annotationProps,
        });
        totalInfo++;
      }
    }
  }

  return {
    annotations,
    hasErrors,
    summaries,
    totalErrors,
    totalInfo,
    totalWarnings,
  };
}

describe("Validation Output Processing", () => {
  const createIssue = (
    level: "ERROR" | "WARNING" | "INFO",
    path: string,
    message: string,
    line?: number,
  ): ValidationIssue => ({
    level,
    line,
    message,
    path,
  });

  describe("Processing valid results", () => {
    it("should handle valid results with no issues", () => {
      const output: ValidationOutput = [
        {
          name: "test-change",
          report: {
            issues: [],
            summary: { errors: 0, info: 0, warnings: 0 },
            valid: true,
          },
          type: "change",
          valid: true,
        },
      ];

      const result = processValidationOutput(output);

      assert.equal(result.hasErrors, false);
      assert.equal(result.totalErrors, 0);
      assert.equal(result.totalWarnings, 0);
      assert.equal(result.totalInfo, 0);
      assert.equal(result.annotations.length, 0);
      assert.equal(result.summaries.length, 1);
      assert.equal(result.summaries[0].logLevel, "info");
      assert.equal(result.summaries[0].message, "valid");
    });

    it("should skip valid results without report", () => {
      const output: ValidationOutput = [
        {
          name: "test-change",
          type: "change",
          valid: true,
        },
      ];

      const result = processValidationOutput(output);

      assert.equal(result.hasErrors, false);
      assert.equal(result.summaries.length, 1);
      assert.equal(result.summaries[0].message, "valid");
    });
  });

  describe("Processing results with errors", () => {
    it("should handle results with error field", () => {
      const output: ValidationOutput = [
        {
          error: "Failed to parse file",
          name: "test-change",
          type: "change",
          valid: false,
        },
      ];

      const result = processValidationOutput(output);

      assert.equal(result.hasErrors, true);
      assert.equal(result.summaries.length, 1);
      assert.equal(result.summaries[0].logLevel, "error");
      assert.equal(result.summaries[0].message, "Failed to parse file");
      assert.equal(result.annotations.length, 0);
    });

    it("should process validation errors and create annotations", () => {
      const output: ValidationOutput = [
        {
          name: "test-change",
          report: {
            issues: [
              createIssue(
                "ERROR",
                "/workspace/spectr/changes/test/spec.md",
                "Missing requirement header",
                42,
              ),
            ],
            summary: { errors: 1, info: 0, warnings: 0 },
            valid: false,
          },
          type: "change",
          valid: false,
        },
      ];

      const result = processValidationOutput(output, "/workspace");

      assert.equal(result.hasErrors, true);
      assert.equal(result.totalErrors, 1);
      assert.equal(result.annotations.length, 1);

      const annotation = result.annotations[0];
      assert.equal(annotation.level, "ERROR");
      assert.equal(annotation.message, "Missing requirement header");
      assert.equal(annotation.props.file, "spectr/changes/test/spec.md");
      assert.equal(annotation.props.startLine, 42);
      assert.equal(annotation.props.title, "change: test-change");
    });

    it("should set hasErrors flag when errors present", () => {
      const output: ValidationOutput = [
        {
          name: "change-1",
          report: {
            issues: [createIssue("ERROR", "spec.md", "Error 1", 10)],
            summary: { errors: 1, info: 0, warnings: 0 },
            valid: false,
          },
          type: "change",
          valid: false,
        },
      ];

      const result = processValidationOutput(output);

      assert.equal(result.hasErrors, true);
      assert.equal(result.summaries[0].logLevel, "error");
    });
  });

  describe("Processing warnings", () => {
    it("should process warnings and create annotations", () => {
      const output: ValidationOutput = [
        {
          name: "test-spec",
          report: {
            issues: [
              createIssue(
                "WARNING",
                "/workspace/spectr/specs/auth/spec.md",
                "Consider adding more scenarios",
                15,
              ),
            ],
            summary: { errors: 0, info: 0, warnings: 1 },
            valid: true,
          },
          type: "spec",
          valid: true,
        },
      ];

      const result = processValidationOutput(output, "/workspace");

      assert.equal(result.hasErrors, false);
      assert.equal(result.totalWarnings, 1);
      assert.equal(result.annotations.length, 1);

      const annotation = result.annotations[0];
      assert.equal(annotation.level, "WARNING");
      assert.equal(annotation.message, "Consider adding more scenarios");
      assert.equal(annotation.props.file, "spectr/specs/auth/spec.md");
      assert.equal(annotation.props.startLine, 15);
    });

    it("should not set hasErrors flag for warnings only", () => {
      const output: ValidationOutput = [
        {
          name: "test-change",
          report: {
            issues: [createIssue("WARNING", "spec.md", "Warning", 10)],
            summary: { errors: 0, info: 0, warnings: 1 },
            valid: true,
          },
          type: "change",
          valid: true,
        },
      ];

      const result = processValidationOutput(output);

      assert.equal(result.hasErrors, false);
      assert.equal(result.summaries[0].logLevel, "warning");
    });
  });

  describe("Processing info messages", () => {
    it("should process info messages and create annotations", () => {
      const output: ValidationOutput = [
        {
          name: "test-change",
          report: {
            issues: [
              createIssue(
                "INFO",
                "spectr/changes/test/proposal.md",
                "Validation passed",
                1,
              ),
            ],
            summary: { errors: 0, info: 1, warnings: 0 },
            valid: true,
          },
          type: "change",
          valid: true,
        },
      ];

      const result = processValidationOutput(output);

      assert.equal(result.hasErrors, false);
      assert.equal(result.totalInfo, 1);
      assert.equal(result.annotations.length, 1);

      const annotation = result.annotations[0];
      assert.equal(annotation.level, "INFO");
      assert.equal(annotation.message, "Validation passed");
    });
  });

  describe("Mixed severity levels", () => {
    it("should process multiple issues with different severities", () => {
      const output: ValidationOutput = [
        {
          name: "test-change",
          report: {
            issues: [
              createIssue("ERROR", "spec.md", "Error 1", 10),
              createIssue("ERROR", "spec.md", "Error 2", 20),
              createIssue("WARNING", "spec.md", "Warning 1", 30),
              createIssue("INFO", "spec.md", "Info 1", 40),
            ],
            summary: { errors: 2, info: 1, warnings: 1 },
            valid: false,
          },
          type: "change",
          valid: false,
        },
      ];

      const result = processValidationOutput(output);

      assert.equal(result.hasErrors, true);
      assert.equal(result.totalErrors, 2);
      assert.equal(result.totalWarnings, 1);
      assert.equal(result.totalInfo, 1);
      assert.equal(result.annotations.length, 4);

      // Verify annotation levels
      const errorAnnotations = result.annotations.filter(
        (a) => a.level === "ERROR",
      );
      const warningAnnotations = result.annotations.filter(
        (a) => a.level === "WARNING",
      );
      const infoAnnotations = result.annotations.filter(
        (a) => a.level === "INFO",
      );

      assert.equal(errorAnnotations.length, 2);
      assert.equal(warningAnnotations.length, 1);
      assert.equal(infoAnnotations.length, 1);
    });
  });

  describe("Multiple results", () => {
    it("should aggregate statistics across multiple results", () => {
      const output: ValidationOutput = [
        {
          name: "change-1",
          report: {
            issues: [
              createIssue("ERROR", "spec1.md", "Error 1", 10),
              createIssue("WARNING", "spec1.md", "Warning 1", 20),
            ],
            summary: { errors: 1, info: 0, warnings: 1 },
            valid: false,
          },
          type: "change",
          valid: false,
        },
        {
          name: "change-2",
          report: {
            issues: [
              createIssue("ERROR", "spec2.md", "Error 2", 30),
              createIssue("ERROR", "spec2.md", "Error 3", 40),
            ],
            summary: { errors: 2, info: 0, warnings: 0 },
            valid: false,
          },
          type: "change",
          valid: false,
        },
        {
          name: "spec-1",
          report: {
            issues: [createIssue("INFO", "spec3.md", "Info 1", 50)],
            summary: { errors: 0, info: 1, warnings: 0 },
            valid: true,
          },
          type: "spec",
          valid: true,
        },
      ];

      const result = processValidationOutput(output);

      assert.equal(result.hasErrors, true);
      assert.equal(result.totalErrors, 3);
      assert.equal(result.totalWarnings, 1);
      assert.equal(result.totalInfo, 1);
      assert.equal(result.annotations.length, 5);
      assert.equal(result.summaries.length, 3);
    });

    it("should handle mix of valid and invalid results", () => {
      const output: ValidationOutput = [
        {
          name: "change-1",
          report: {
            issues: [],
            summary: { errors: 0, info: 0, warnings: 0 },
            valid: true,
          },
          type: "change",
          valid: true,
        },
        {
          error: "Parse error",
          name: "change-2",
          type: "change",
          valid: false,
        },
        {
          name: "change-3",
          report: {
            issues: [createIssue("ERROR", "spec.md", "Error", 10)],
            summary: { errors: 1, info: 0, warnings: 0 },
            valid: false,
          },
          type: "change",
          valid: false,
        },
      ];

      const result = processValidationOutput(output);

      assert.equal(result.hasErrors, true);
      assert.equal(result.totalErrors, 1);
      assert.equal(result.summaries.length, 3);
    });
  });

  describe("Edge cases", () => {
    it("should handle empty validation output", () => {
      const output: ValidationOutput = [];

      const result = processValidationOutput(output);

      assert.equal(result.hasErrors, false);
      assert.equal(result.totalErrors, 0);
      assert.equal(result.totalWarnings, 0);
      assert.equal(result.totalInfo, 0);
      assert.equal(result.annotations.length, 0);
      assert.equal(result.summaries.length, 0);
    });

    it("should handle issues without line numbers", () => {
      const output: ValidationOutput = [
        {
          name: "test-change",
          report: {
            issues: [createIssue("ERROR", "spec.md", "No line number")],
            summary: { errors: 1, info: 0, warnings: 0 },
            valid: false,
          },
          type: "change",
          valid: false,
        },
      ];

      const result = processValidationOutput(output);

      assert.equal(result.annotations.length, 1);
      assert.equal(result.annotations[0].props.startLine, 1); // Default to 1
    });

    it("should handle results without report or error", () => {
      const output: ValidationOutput = [
        {
          name: "test-change",
          type: "change",
          valid: false,
        },
      ];

      const result = processValidationOutput(output);

      assert.equal(result.hasErrors, false);
      assert.equal(result.annotations.length, 0);
      assert.equal(result.summaries.length, 0);
    });

    it("should correctly format item titles", () => {
      const output: ValidationOutput = [
        {
          name: "my-change",
          report: {
            issues: [createIssue("ERROR", "spec.md", "Error", 10)],
            summary: { errors: 1, info: 0, warnings: 0 },
            valid: false,
          },
          type: "change",
          valid: false,
        },
        {
          name: "auth-spec",
          report: {
            issues: [createIssue("ERROR", "spec.md", "Error", 20)],
            summary: { errors: 1, info: 0, warnings: 0 },
            valid: false,
          },
          type: "spec",
          valid: false,
        },
      ];

      const result = processValidationOutput(output);

      assert.equal(result.annotations[0].props.title, "change: my-change");
      assert.equal(result.annotations[1].props.title, "spec: auth-spec");
    });

    it("should handle relative path conversion", () => {
      const output: ValidationOutput = [
        {
          name: "test-change",
          report: {
            issues: [
              createIssue(
                "ERROR",
                "/home/user/workspace/spectr/changes/test/spec.md",
                "Error",
                10,
              ),
            ],
            summary: { errors: 1, info: 0, warnings: 0 },
            valid: false,
          },
          type: "change",
          valid: false,
        },
      ];

      const result = processValidationOutput(output, "/home/user/workspace");

      assert.equal(
        result.annotations[0].props.file,
        "spectr/changes/test/spec.md",
      );
    });

    it("should handle paths that do not start with workspace", () => {
      const output: ValidationOutput = [
        {
          name: "test-change",
          report: {
            issues: [
              createIssue("ERROR", "relative/path/spec.md", "Error", 10),
            ],
            summary: { errors: 1, info: 0, warnings: 0 },
            valid: false,
          },
          type: "change",
          valid: false,
        },
      ];

      const result = processValidationOutput(output, "/workspace");

      assert.equal(result.annotations[0].props.file, "relative/path/spec.md");
    });
  });

  describe("Summary generation", () => {
    it("should generate error summary for results with errors", () => {
      const output: ValidationOutput = [
        {
          name: "test-change",
          report: {
            issues: [],
            summary: { errors: 5, info: 0, warnings: 2 },
            valid: false,
          },
          type: "change",
          valid: false,
        },
      ];

      const result = processValidationOutput(output);

      assert.equal(result.summaries[0].message, "5 errors, 2 warnings");
      assert.equal(result.summaries[0].logLevel, "error");
    });

    it("should generate warning summary for results with warnings only", () => {
      const output: ValidationOutput = [
        {
          name: "test-change",
          report: {
            issues: [
              createIssue("WARNING", "spec.md", "Warning 1", 10),
              createIssue("WARNING", "spec.md", "Warning 2", 20),
              createIssue("WARNING", "spec.md", "Warning 3", 30),
            ],
            summary: { errors: 0, info: 0, warnings: 3 },
            valid: true,
          },
          type: "change",
          valid: true,
        },
      ];

      const result = processValidationOutput(output);

      assert.equal(result.summaries[0].message, "3 warnings");
      assert.equal(result.summaries[0].logLevel, "warning");
    });

    it("should generate info summary for results with info only", () => {
      const output: ValidationOutput = [
        {
          name: "test-change",
          report: {
            issues: [
              createIssue("INFO", "spec.md", "Info 1", 10),
              createIssue("INFO", "spec.md", "Info 2", 20),
            ],
            summary: { errors: 0, info: 2, warnings: 0 },
            valid: true,
          },
          type: "change",
          valid: true,
        },
      ];

      const result = processValidationOutput(output);

      assert.equal(result.summaries[0].message, "2 info messages");
      assert.equal(result.summaries[0].logLevel, "info");
    });
  });
});
