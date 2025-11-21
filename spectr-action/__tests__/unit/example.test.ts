import assert from "node:assert/strict";
import { describe, it } from "node:test";

describe("Example Test Suite", () => {
  it("should pass a basic assertion", () => {
    const result = 2 + 2;
    assert.equal(result, 4);
  });

  it("should handle string comparisons", () => {
    const greeting = "Hello, World!";
    assert.equal(greeting, "Hello, World!");
  });

  it("should verify object equality", () => {
    const obj1 = { name: "test", value: 42 };
    const obj2 = { name: "test", value: 42 };
    assert.deepEqual(obj1, obj2);
  });
});

describe("Async Operations", () => {
  it("should handle async/await", async () => {
    const promise = Promise.resolve("success");
    const result = await promise;
    assert.equal(result, "success");
  });

  it("should handle promise rejections", async () => {
    await assert.rejects(
      async () => {
        throw new Error("Expected error");
      },
      {
        message: "Expected error",
        name: "Error",
      },
    );
  });
});
