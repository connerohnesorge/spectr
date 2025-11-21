import assert from "node:assert/strict";
import { describe, it } from "node:test";
import {
  allValid,
  type BulkResult,
  formatAllIssues,
  formatIssue,
  getAllErrors,
  getAllInfo,
  getAllWarnings,
  getFailedResults,
  getResultsWithErrors,
  getTotalErrorCount,
  getTotalInfoCount,
  getTotalIssueCount,
  getTotalWarningCount,
  hasAnyErrors,
  hasError,
  hasReport,
  isValid,
  type ValidationIssue,
  type ValidationOutput,
} from "../../../src/types/spectr";

describe("Type Guards", () => {
  describe("hasReport", () => {
    it("should return true when report is present", () => {
      const result: BulkResult = {
        name: "test-change",
        report: {
          issues: [],
          summary: { errors: 0, info: 0, warnings: 0 },
          valid: true,
        },
        type: "change",
        valid: true,
      };
      assert.equal(hasReport(result), true);
    });

    it("should return false when report is undefined", () => {
      const result: BulkResult = {
        error: "Validation failed",
        name: "test-change",
        type: "change",
        valid: false,
      };
      assert.equal(hasReport(result), false);
    });

    it("should narrow type to include report", () => {
      const result: BulkResult = {
        name: "test-change",
        report: {
          issues: [],
          summary: { errors: 0, info: 0, warnings: 0 },
          valid: true,
        },
        type: "change",
        valid: true,
      };

      if (hasReport(result)) {
        // Type should be narrowed to include report
        assert.ok(result.report);
        assert.equal(result.report.valid, true);
      } else {
        assert.fail("hasReport should return true");
      }
    });
  });

  describe("hasError", () => {
    it("should return true when error is present", () => {
      const result: BulkResult = {
        error: "Failed to parse",
        name: "test-change",
        type: "change",
        valid: false,
      };
      assert.equal(hasError(result), true);
    });

    it("should return false when error is undefined", () => {
      const result: BulkResult = {
        name: "test-change",
        report: {
          issues: [],
          summary: { errors: 0, info: 0, warnings: 0 },
          valid: true,
        },
        type: "change",
        valid: true,
      };
      assert.equal(hasError(result), false);
    });

    it("should narrow type to include error string", () => {
      const result: BulkResult = {
        error: "Parse error",
        name: "test-change",
        type: "change",
        valid: false,
      };

      if (hasError(result)) {
        // Type should be narrowed to include error
        assert.equal(result.error, "Parse error");
      } else {
        assert.fail("hasError should return true");
      }
    });
  });

  describe("isValid", () => {
    it("should return true when result is valid", () => {
      const result: BulkResult = {
        name: "test-change",
        report: {
          issues: [],
          summary: { errors: 0, info: 0, warnings: 0 },
          valid: true,
        },
        type: "change",
        valid: true,
      };
      assert.equal(isValid(result), true);
    });

    it("should return false when result is invalid", () => {
      const result: BulkResult = {
        error: "Validation failed",
        name: "test-change",
        type: "change",
        valid: false,
      };
      assert.equal(isValid(result), false);
    });
  });
});

describe("Issue Extraction", () => {
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

  describe("getAllErrors", () => {
    it("should extract all errors from validation output", () => {
      const output: ValidationOutput = [
        {
          name: "change-1",
          report: {
            issues: [
              createIssue("ERROR", "spec.md", "Missing requirement", 10),
              createIssue("WARNING", "spec.md", "Consider adding scenario", 15),
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
            issues: [createIssue("ERROR", "design.md", "Invalid format", 5)],
            summary: { errors: 1, info: 0, warnings: 0 },
            valid: false,
          },
          type: "change",
          valid: false,
        },
      ];

      const errors = getAllErrors(output);
      assert.equal(errors.length, 2);
      assert.equal(errors[0].level, "ERROR");
      assert.equal(errors[0].message, "Missing requirement");
      assert.equal(errors[1].level, "ERROR");
      assert.equal(errors[1].message, "Invalid format");
    });

    it("should return empty array when no errors", () => {
      const output: ValidationOutput = [
        {
          name: "change-1",
          report: {
            issues: [createIssue("INFO", "spec.md", "All good", 10)],
            summary: { errors: 0, info: 1, warnings: 0 },
            valid: true,
          },
          type: "change",
          valid: true,
        },
      ];

      const errors = getAllErrors(output);
      assert.equal(errors.length, 0);
    });

    it("should skip results without reports", () => {
      const output: ValidationOutput = [
        {
          error: "Failed to parse",
          name: "change-1",
          type: "change",
          valid: false,
        },
      ];

      const errors = getAllErrors(output);
      assert.equal(errors.length, 0);
    });
  });

  describe("getAllWarnings", () => {
    it("should extract all warnings from validation output", () => {
      const output: ValidationOutput = [
        {
          name: "change-1",
          report: {
            issues: [
              createIssue("WARNING", "spec.md", "Missing scenario", 10),
              createIssue("ERROR", "spec.md", "Required field missing", 15),
            ],
            summary: { errors: 1, info: 0, warnings: 1 },
            valid: true,
          },
          type: "change",
          valid: true,
        },
      ];

      const warnings = getAllWarnings(output);
      assert.equal(warnings.length, 1);
      assert.equal(warnings[0].level, "WARNING");
      assert.equal(warnings[0].message, "Missing scenario");
    });
  });

  describe("getAllInfo", () => {
    it("should extract all info messages from validation output", () => {
      const output: ValidationOutput = [
        {
          name: "change-1",
          report: {
            issues: [
              createIssue("INFO", "spec.md", "Validation passed", 10),
              createIssue("WARNING", "spec.md", "Consider improvement", 15),
            ],
            summary: { errors: 0, info: 1, warnings: 1 },
            valid: true,
          },
          type: "change",
          valid: true,
        },
      ];

      const info = getAllInfo(output);
      assert.equal(info.length, 1);
      assert.equal(info[0].level, "INFO");
      assert.equal(info[0].message, "Validation passed");
    });
  });
});

describe("Count Functions", () => {
  describe("getTotalErrorCount", () => {
    it("should count total errors across all results", () => {
      const output: ValidationOutput = [
        {
          name: "change-1",
          report: {
            issues: [],
            summary: { errors: 2, info: 0, warnings: 1 },
            valid: false,
          },
          type: "change",
          valid: false,
        },
        {
          name: "change-2",
          report: {
            issues: [],
            summary: { errors: 3, info: 0, warnings: 0 },
            valid: false,
          },
          type: "change",
          valid: false,
        },
      ];

      assert.equal(getTotalErrorCount(output), 5);
    });

    it("should handle results without reports", () => {
      const output: ValidationOutput = [
        {
          error: "Parse failed",
          name: "change-1",
          type: "change",
          valid: false,
        },
        {
          name: "change-2",
          report: {
            issues: [],
            summary: { errors: 2, info: 0, warnings: 0 },
            valid: false,
          },
          type: "change",
          valid: false,
        },
      ];

      assert.equal(getTotalErrorCount(output), 2);
    });
  });

  describe("getTotalWarningCount", () => {
    it("should count total warnings across all results", () => {
      const output: ValidationOutput = [
        {
          name: "change-1",
          report: {
            issues: [],
            summary: { errors: 0, info: 0, warnings: 2 },
            valid: true,
          },
          type: "change",
          valid: true,
        },
        {
          name: "change-2",
          report: {
            issues: [],
            summary: { errors: 0, info: 0, warnings: 3 },
            valid: true,
          },
          type: "change",
          valid: true,
        },
      ];

      assert.equal(getTotalWarningCount(output), 5);
    });
  });

  describe("getTotalInfoCount", () => {
    it("should count total info messages across all results", () => {
      const output: ValidationOutput = [
        {
          name: "change-1",
          report: {
            issues: [],
            summary: { errors: 0, info: 2, warnings: 0 },
            valid: true,
          },
          type: "change",
          valid: true,
        },
        {
          name: "change-2",
          report: {
            issues: [],
            summary: { errors: 0, info: 3, warnings: 0 },
            valid: true,
          },
          type: "change",
          valid: true,
        },
      ];

      assert.equal(getTotalInfoCount(output), 5);
    });
  });

  describe("getTotalIssueCount", () => {
    it("should count total issues across all results", () => {
      const output: ValidationOutput = [
        {
          name: "change-1",
          report: {
            issues: [
              { level: "ERROR", message: "Error 1", path: "a.md" },
              { level: "WARNING", message: "Warning 1", path: "b.md" },
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
            issues: [{ level: "INFO", message: "Info 1", path: "c.md" }],
            summary: { errors: 0, info: 1, warnings: 0 },
            valid: true,
          },
          type: "change",
          valid: true,
        },
      ];

      assert.equal(getTotalIssueCount(output), 3);
    });
  });
});

describe("Validation Predicates", () => {
  describe("hasAnyErrors", () => {
    it("should return true when any result has errors", () => {
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
          name: "change-2",
          report: {
            issues: [],
            summary: { errors: 1, info: 0, warnings: 0 },
            valid: false,
          },
          type: "change",
          valid: false,
        },
      ];

      assert.equal(hasAnyErrors(output), true);
    });

    it("should return false when no results have errors", () => {
      const output: ValidationOutput = [
        {
          name: "change-1",
          report: {
            issues: [],
            summary: { errors: 0, info: 0, warnings: 1 },
            valid: true,
          },
          type: "change",
          valid: true,
        },
      ];

      assert.equal(hasAnyErrors(output), false);
    });
  });

  describe("allValid", () => {
    it("should return true when all results are valid", () => {
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
          name: "change-2",
          report: {
            issues: [],
            summary: { errors: 0, info: 0, warnings: 0 },
            valid: true,
          },
          type: "change",
          valid: true,
        },
      ];

      assert.equal(allValid(output), true);
    });

    it("should return false when any result is invalid", () => {
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
      ];

      assert.equal(allValid(output), false);
    });
  });
});

describe("Result Filtering", () => {
  describe("getFailedResults", () => {
    it("should return all invalid results", () => {
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
            issues: [],
            summary: { errors: 1, info: 0, warnings: 0 },
            valid: false,
          },
          type: "change",
          valid: false,
        },
      ];

      const failed = getFailedResults(output);
      assert.equal(failed.length, 2);
      assert.equal(failed[0].name, "change-2");
      assert.equal(failed[1].name, "change-3");
    });
  });

  describe("getResultsWithErrors", () => {
    it("should return results with errors only", () => {
      const output: ValidationOutput = [
        {
          name: "change-1",
          report: {
            issues: [],
            summary: { errors: 0, info: 0, warnings: 1 },
            valid: true,
          },
          type: "change",
          valid: true,
        },
        {
          name: "change-2",
          report: {
            issues: [],
            summary: { errors: 2, info: 0, warnings: 0 },
            valid: false,
          },
          type: "change",
          valid: false,
        },
        {
          error: "No report available",
          name: "change-3",
          type: "change",
          valid: false,
        },
      ];

      const withErrors = getResultsWithErrors(output);
      assert.equal(withErrors.length, 1);
      assert.equal(withErrors[0].name, "change-2");
    });
  });
});

describe("Formatting Functions", () => {
  describe("formatIssue", () => {
    it("should format issue with line number", () => {
      const issue: ValidationIssue = {
        level: "ERROR",
        line: 42,
        message: "Missing requirement header",
        path: "spectr/changes/test/spec.md",
      };

      const formatted = formatIssue(issue);
      assert.equal(
        formatted,
        "[ERROR] spectr/changes/test/spec.md:42: Missing requirement header",
      );
    });

    it("should format issue without line number", () => {
      const issue: ValidationIssue = {
        level: "WARNING",
        message: "Consider adding more detail",
        path: "spectr/changes/test/proposal.md",
      };

      const formatted = formatIssue(issue);
      assert.equal(
        formatted,
        "[WARNING] spectr/changes/test/proposal.md: Consider adding more detail",
      );
    });
  });

  describe("formatAllIssues", () => {
    it("should format all issues from a result with report", () => {
      const result: BulkResult = {
        name: "test-change",
        report: {
          issues: [
            {
              level: "ERROR",
              line: 10,
              message: "Error 1",
              path: "spec.md",
            },
            {
              level: "WARNING",
              line: 20,
              message: "Warning 1",
              path: "spec.md",
            },
          ],
          summary: { errors: 1, info: 0, warnings: 1 },
          valid: false,
        },
        type: "change",
        valid: false,
      };

      const formatted = formatAllIssues(result);
      assert.equal(formatted.length, 2);
      assert.equal(formatted[0], "[ERROR] spec.md:10: Error 1");
      assert.equal(formatted[1], "[WARNING] spec.md:20: Warning 1");
    });

    it("should format error message when result has error", () => {
      const result: BulkResult = {
        error: "Failed to parse file",
        name: "test-change",
        type: "change",
        valid: false,
      };

      const formatted = formatAllIssues(result);
      assert.equal(formatted.length, 1);
      assert.equal(formatted[0], "Error: Failed to parse file");
    });

    it("should return empty array when no report and no error", () => {
      const result: BulkResult = {
        name: "test-change",
        type: "change",
        valid: true,
      };

      const formatted = formatAllIssues(result);
      assert.equal(formatted.length, 0);
    });
  });
});

describe("Edge Cases", () => {
  it("should handle empty validation output", () => {
    const output: ValidationOutput = [];

    assert.equal(getTotalErrorCount(output), 0);
    assert.equal(getTotalWarningCount(output), 0);
    assert.equal(getTotalInfoCount(output), 0);
    assert.equal(getTotalIssueCount(output), 0);
    assert.equal(hasAnyErrors(output), false);
    assert.equal(allValid(output), true);
    assert.equal(getAllErrors(output).length, 0);
  });

  it("should handle results with empty issues array", () => {
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
    ];

    assert.equal(getTotalIssueCount(output), 0);
    assert.equal(getAllErrors(output).length, 0);
  });

  it("should handle mixed results with and without reports", () => {
    const output: ValidationOutput = [
      {
        error: "Parse error",
        name: "change-1",
        type: "change",
        valid: false,
      },
      {
        name: "change-2",
        report: {
          issues: [
            {
              level: "ERROR",
              message: "Missing requirement",
              path: "spec.md",
            },
          ],
          summary: { errors: 1, info: 0, warnings: 0 },
          valid: false,
        },
        type: "change",
        valid: false,
      },
      {
        name: "change-3",
        type: "change",
        valid: true,
      },
    ];

    assert.equal(allValid(output), false);
    assert.equal(getTotalErrorCount(output), 1);
    assert.equal(getAllErrors(output).length, 1);
  });
});
