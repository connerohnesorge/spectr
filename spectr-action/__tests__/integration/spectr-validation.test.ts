import assert from "node:assert/strict";
import { describe, it } from "node:test";
import { exec } from "@actions/exec";
import { getFixturePath } from "../helpers/test-utils";

/**
 * Validation result for a single spec or change
 */
interface ValidationItem {
  name: string;
  type: "spec" | "change";
  valid: boolean;
  report: {
    valid: boolean;
    issues: ValidationIssue[];
    summary: {
      errors: number;
      warnings: number;
      info: number;
    };
  };
}

/**
 * Validation issue/error details
 */
interface ValidationIssue {
  level: "ERROR" | "WARNING" | "INFO";
  path: string;
  line: number;
  message: string;
}

/**
 * Execute spectr validate command and parse JSON output
 * @param fixturePath - Absolute path to the fixture directory
 * @returns Array of validation items
 */
async function runSpectrValidation(
  fixturePath: string,
): Promise<ValidationItem[]> {
  let stdout = "";
  let stderr = "";

  const _exitCode = await exec(
    "spectr",
    ["validate", "--all", "--strict", "--json"],
    {
      cwd: fixturePath,
      ignoreReturnCode: true,
      listeners: {
        stderr: (data: Buffer) => {
          stderr += data.toString();
        },
        stdout: (data: Buffer) => {
          stdout += data.toString();
        },
      },
      silent: true,
    },
  );

  // Extract JSON from output (stderr may contain error messages after JSON)
  const jsonMatch = stdout.match(/^\[[\s\S]*\]$/m);
  if (!jsonMatch) {
    throw new Error(
      `Failed to parse JSON from spectr output. stdout: ${stdout}, stderr: ${stderr}`,
    );
  }

  const items = JSON.parse(jsonMatch[0]) as ValidationItem[];
  return items;
}

/**
 * Execute spectr validate command and expect it to fail
 * @param fixturePath - Absolute path to the fixture directory
 * @returns Parsed validation items and exit code
 */
async function runSpectrValidationExpectingFailure(
  fixturePath: string,
): Promise<{ items: ValidationItem[]; exitCode: number }> {
  let stdout = "";
  let stderr = "";

  const exitCode = await exec(
    "spectr",
    ["validate", "--all", "--strict", "--json"],
    {
      cwd: fixturePath,
      ignoreReturnCode: true,
      listeners: {
        stderr: (data: Buffer) => {
          stderr += data.toString();
        },
        stdout: (data: Buffer) => {
          stdout += data.toString();
        },
      },
      silent: true,
    },
  );

  // For failures, JSON is in stdout but error message is in stderr
  const jsonMatch = stdout.match(/^\[[\s\S]*\]$/m);
  if (!jsonMatch) {
    throw new Error(
      `Failed to parse JSON from spectr output. stdout: ${stdout}, stderr: ${stderr}`,
    );
  }

  const items = JSON.parse(jsonMatch[0]) as ValidationItem[];
  return { exitCode, items };
}

describe("Spectr Validation Integration Tests", () => {
  describe("Valid Fixtures", () => {
    it("should pass validation for valid-spectr-project", async () => {
      const fixturePath = getFixturePath("valid-spectr-project");
      const items = await runSpectrValidation(fixturePath);

      assert.equal(items.length, 1, "Should have exactly 1 spec");
      assert.equal(items[0].name, "example-capability");
      assert.equal(items[0].type, "spec");
      assert.equal(items[0].valid, true);
      assert.equal(items[0].report.valid, true);
      assert.equal(items[0].report.issues.length, 0);
      assert.equal(items[0].report.summary.errors, 0);
      assert.equal(items[0].report.summary.warnings, 0);
    });

    it("should pass validation for valid-with-changes", async () => {
      const fixturePath = getFixturePath("valid-with-changes");
      const items = await runSpectrValidation(fixturePath);

      assert.equal(items.length, 2, "Should have 2 items (1 change + 1 spec)");

      // Find the change and spec
      const change = items.find((item) => item.type === "change");
      const spec = items.find((item) => item.type === "spec");

      assert.ok(change, "Should have a change");
      assert.ok(spec, "Should have a spec");

      // Validate the change
      assert.equal(change.name, "add-two-factor-auth");
      assert.equal(change.valid, true);
      assert.equal(change.report.valid, true);
      assert.equal(change.report.issues.length, 0);
      assert.equal(change.report.summary.errors, 0);

      // Validate the spec
      assert.equal(spec.name, "auth");
      assert.equal(spec.valid, true);
      assert.equal(spec.report.valid, true);
      assert.equal(spec.report.issues.length, 0);
      assert.equal(spec.report.summary.errors, 0);
    });

    it("should pass validation for empty-spectr-project", async () => {
      const fixturePath = getFixturePath("empty-spectr-project");
      const items = await runSpectrValidation(fixturePath);

      assert.equal(items.length, 0, "Should have 0 items for empty project");
    });

    it("should pass validation for multi-capability project", async () => {
      const fixturePath = getFixturePath("multi-capability");
      const items = await runSpectrValidation(fixturePath);

      assert.equal(items.length, 3, "Should have 3 specs");

      // All items should be valid specs
      for (const item of items) {
        assert.equal(item.type, "spec");
        assert.equal(item.valid, true);
        assert.equal(item.report.valid, true);
        assert.equal(item.report.issues.length, 0);
        assert.equal(item.report.summary.errors, 0);
      }

      // Check that we have the expected capabilities
      const names = items.map((item) => item.name).sort();
      assert.deepEqual(names, ["authentication", "authorization", "logging"]);
    });
  });

  describe("Invalid Fixtures", () => {
    it("should fail validation for invalid-missing-scenarios", async () => {
      const fixturePath = getFixturePath("invalid-missing-scenarios");
      const result = await runSpectrValidationExpectingFailure(fixturePath);

      assert.equal(result.exitCode, 1, "Should exit with code 1");
      assert.equal(result.items.length, 1, "Should have 1 spec");

      const spec = result.items[0];
      assert.equal(spec.name, "api");
      assert.equal(spec.type, "spec");
      assert.equal(spec.valid, false);
      assert.equal(spec.report.valid, false);
      assert.equal(spec.report.summary.errors, 3);

      // Verify all three requirements are flagged
      assert.equal(spec.report.issues.length, 3);

      // Check that all errors are about missing scenarios
      for (const issue of spec.report.issues) {
        assert.equal(issue.level, "ERROR");
        assert.ok(
          issue.message.includes("at least one scenario"),
          `Expected scenario error, got: ${issue.message}`,
        );
      }

      // Verify specific requirements with errors
      const requirementNames = spec.report.issues.map((issue) => {
        const match = issue.path.match(/Requirement '([^']+)'/);
        return match ? match[1] : "";
      });

      assert.ok(requirementNames.includes("Authentication"));
      assert.ok(requirementNames.includes("Rate Limiting"));
      assert.ok(requirementNames.includes("Error Response Format"));
    });

    it("should fail validation for invalid-malformed-spec", async () => {
      const fixturePath = getFixturePath("invalid-malformed-spec");
      const result = await runSpectrValidationExpectingFailure(fixturePath);

      assert.equal(result.exitCode, 1, "Should exit with code 1");
      assert.equal(result.items.length, 1, "Should have 1 spec");

      const spec = result.items[0];
      assert.equal(spec.name, "data-storage");
      assert.equal(spec.type, "spec");
      assert.equal(spec.valid, false);
      assert.equal(spec.report.valid, false);
      assert.equal(spec.report.summary.errors, 5);

      // Check for both types of errors
      const issueMessages = spec.report.issues.map((issue) => issue.message);

      const missingScenarioErrors = issueMessages.filter((msg) =>
        msg.includes("at least one scenario"),
      );
      const formatErrors = issueMessages.filter((msg) =>
        msg.includes("'#### Scenario:' format"),
      );

      assert.equal(
        missingScenarioErrors.length,
        3,
        "Should have 3 missing scenario errors",
      );
      assert.equal(formatErrors.length, 2, "Should have 2 format errors");

      // Verify error levels
      for (const issue of spec.report.issues) {
        assert.equal(issue.level, "ERROR");
        assert.ok(issue.line > 0, "Should have valid line number");
      }
    });
  });

  describe("JSON Output Parsing", () => {
    it("should parse validation report structure correctly", async () => {
      const fixturePath = getFixturePath("valid-spectr-project");
      const items = await runSpectrValidation(fixturePath);

      const item = items[0];

      // Verify top-level structure
      assert.ok("name" in item);
      assert.ok("type" in item);
      assert.ok("valid" in item);
      assert.ok("report" in item);

      // Verify report structure
      assert.ok("valid" in item.report);
      assert.ok("issues" in item.report);
      assert.ok("summary" in item.report);

      // Verify summary structure
      assert.ok("errors" in item.report.summary);
      assert.ok("warnings" in item.report.summary);
      assert.ok("info" in item.report.summary);
    });

    it("should parse validation issue structure correctly", async () => {
      const fixturePath = getFixturePath("invalid-missing-scenarios");
      const result = await runSpectrValidationExpectingFailure(fixturePath);

      const spec = result.items[0];
      const issue = spec.report.issues[0];

      // Verify issue structure
      assert.ok("level" in issue);
      assert.ok("path" in issue);
      assert.ok("line" in issue);
      assert.ok("message" in issue);

      // Verify issue types
      assert.equal(typeof issue.level, "string");
      assert.equal(typeof issue.path, "string");
      assert.equal(typeof issue.line, "number");
      assert.equal(typeof issue.message, "string");

      // Verify issue values
      assert.ok(
        ["ERROR", "WARNING", "INFO"].includes(issue.level),
        `Invalid issue level: ${issue.level}`,
      );
      assert.ok(issue.path.length > 0, "Path should not be empty");
      assert.ok(issue.line > 0, "Line should be positive");
      assert.ok(issue.message.length > 0, "Message should not be empty");
    });
  });

  describe("Error Detection and Reporting", () => {
    it("should detect multiple validation errors in a single spec", async () => {
      const fixturePath = getFixturePath("invalid-malformed-spec");
      const result = await runSpectrValidationExpectingFailure(fixturePath);

      const spec = result.items[0];

      // Should have multiple errors
      assert.ok(spec.report.issues.length > 1, "Should have multiple errors");
      assert.equal(
        spec.report.summary.errors,
        spec.report.issues.length,
        "Error count should match issues length",
      );
    });

    it("should provide line numbers for all errors", async () => {
      const fixturePath = getFixturePath("invalid-missing-scenarios");
      const result = await runSpectrValidationExpectingFailure(fixturePath);

      const spec = result.items[0];

      // All issues should have valid line numbers
      for (const issue of spec.report.issues) {
        assert.ok(
          issue.line > 0,
          `Issue should have valid line number: ${JSON.stringify(issue)}`,
        );
      }

      // Line numbers should be in ascending order (generally)
      const lineNumbers = spec.report.issues.map((issue) => issue.line);
      const sortedLineNumbers = [...lineNumbers].sort((a, b) => a - b);
      assert.deepEqual(
        lineNumbers,
        sortedLineNumbers,
        "Line numbers should be in order",
      );
    });

    it("should include file paths in error reports", async () => {
      const fixturePath = getFixturePath("invalid-missing-scenarios");
      const result = await runSpectrValidationExpectingFailure(fixturePath);

      const spec = result.items[0];

      // All issues should have paths
      for (const issue of spec.report.issues) {
        assert.ok(
          issue.path.includes("spec.md"),
          "Path should include spec.md",
        );
        assert.ok(
          issue.path.includes("api"),
          "Path should include capability name",
        );
      }
    });

    it("should differentiate between error types", async () => {
      const fixturePath = getFixturePath("invalid-malformed-spec");
      const result = await runSpectrValidationExpectingFailure(fixturePath);

      const spec = result.items[0];

      // Should have different error messages for different issues
      const uniqueMessages = new Set(
        spec.report.issues.map((issue) => issue.message),
      );
      assert.ok(
        uniqueMessages.size > 1,
        "Should have multiple types of errors",
      );

      // Should have both missing scenario and format errors
      const hasMissingScenarioError = spec.report.issues.some((issue) =>
        issue.message.includes("at least one scenario"),
      );
      const hasFormatError = spec.report.issues.some((issue) =>
        issue.message.includes("'#### Scenario:' format"),
      );

      assert.ok(hasMissingScenarioError, "Should have missing scenario error");
      assert.ok(hasFormatError, "Should have format error");
    });
  });

  describe("Edge Cases", () => {
    it("should handle empty validation results", async () => {
      const fixturePath = getFixturePath("empty-spectr-project");
      const items = await runSpectrValidation(fixturePath);

      assert.ok(Array.isArray(items), "Should return an array");
      assert.equal(items.length, 0, "Should be empty array");
    });

    it("should handle mixed valid and invalid items gracefully", async () => {
      // Note: This test assumes that when one item fails, the command still
      // returns all results. Based on our testing, invalid items are marked
      // as valid: false but still included in the output.
      const fixturePath = getFixturePath("invalid-missing-scenarios");
      const result = await runSpectrValidationExpectingFailure(fixturePath);

      // Should still return results even though validation failed
      assert.ok(result.items.length > 0, "Should have results");
      assert.ok(
        result.items.some((item) => !item.valid),
        "Should have at least one invalid item",
      );
    });
  });

  describe("Summary Statistics", () => {
    it("should accurately count errors in summary", async () => {
      const fixturePath = getFixturePath("invalid-missing-scenarios");
      const result = await runSpectrValidationExpectingFailure(fixturePath);

      const spec = result.items[0];

      assert.equal(
        spec.report.summary.errors,
        spec.report.issues.filter((issue) => issue.level === "ERROR").length,
        "Error count should match ERROR level issues",
      );
    });

    it("should accurately count warnings in summary", async () => {
      const fixturePath = getFixturePath("valid-spectr-project");
      const items = await runSpectrValidation(fixturePath);

      const spec = items[0];

      assert.equal(
        spec.report.summary.warnings,
        spec.report.issues.filter((issue) => issue.level === "WARNING").length,
        "Warning count should match WARNING level issues",
      );
    });

    it("should show zero counts for valid specs", async () => {
      const fixturePath = getFixturePath("valid-spectr-project");
      const items = await runSpectrValidation(fixturePath);

      const spec = items[0];

      assert.equal(spec.report.summary.errors, 0);
      assert.equal(spec.report.summary.warnings, 0);
      assert.equal(spec.report.summary.info, 0);
    });
  });
});
