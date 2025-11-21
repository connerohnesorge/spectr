import assert from "node:assert/strict";
import { describe, it } from "node:test";
import { getArch, getPlatform } from "../../../src/utils/platforms";

describe("Platform Utilities", () => {
  describe("getArch()", () => {
    it("should map arm64 to aarch64", () => {
      // Mock process.arch
      const originalArch = process.arch;
      Object.defineProperty(process, "arch", {
        configurable: true,
        value: "arm64",
      });

      const result = getArch();
      assert.equal(result, "aarch64");

      // Restore original
      Object.defineProperty(process, "arch", {
        configurable: true,
        value: originalArch,
      });
    });

    it("should map ia32 to i686", () => {
      const originalArch = process.arch;
      Object.defineProperty(process, "arch", {
        configurable: true,
        value: "ia32",
      });

      const result = getArch();
      assert.equal(result, "i686");

      Object.defineProperty(process, "arch", {
        configurable: true,
        value: originalArch,
      });
    });

    it("should map x64 to x86_64", () => {
      const originalArch = process.arch;
      Object.defineProperty(process, "arch", {
        configurable: true,
        value: "x64",
      });

      const result = getArch();
      assert.equal(result, "x86_64");

      Object.defineProperty(process, "arch", {
        configurable: true,
        value: originalArch,
      });
    });

    it("should return undefined for unsupported architecture", () => {
      const originalArch = process.arch;
      Object.defineProperty(process, "arch", {
        configurable: true,
        value: "unsupported",
      });

      const result = getArch();
      assert.equal(result, undefined);

      Object.defineProperty(process, "arch", {
        configurable: true,
        value: originalArch,
      });
    });
  });

  describe("getPlatform()", () => {
    it("should map darwin to apple-darwin", () => {
      const originalPlatform = process.platform;
      Object.defineProperty(process, "platform", {
        configurable: true,
        value: "darwin",
      });

      const result = getPlatform();
      assert.equal(result, "apple-darwin");

      Object.defineProperty(process, "platform", {
        configurable: true,
        value: originalPlatform,
      });
    });

    it("should map linux to unknown-linux-gnu", () => {
      const originalPlatform = process.platform;
      Object.defineProperty(process, "platform", {
        configurable: true,
        value: "linux",
      });

      const result = getPlatform();
      assert.equal(result, "unknown-linux-gnu");

      Object.defineProperty(process, "platform", {
        configurable: true,
        value: originalPlatform,
      });
    });

    it("should map win32 to pc-windows-msvc", () => {
      const originalPlatform = process.platform;
      Object.defineProperty(process, "platform", {
        configurable: true,
        value: "win32",
      });

      const result = getPlatform();
      assert.equal(result, "pc-windows-msvc");

      Object.defineProperty(process, "platform", {
        configurable: true,
        value: originalPlatform,
      });
    });

    it("should return undefined for unsupported platform", () => {
      const originalPlatform = process.platform;
      Object.defineProperty(process, "platform", {
        configurable: true,
        value: "unsupported",
      });

      const result = getPlatform();
      assert.equal(result, undefined);

      Object.defineProperty(process, "platform", {
        configurable: true,
        value: originalPlatform,
      });
    });
  });
});
